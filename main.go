package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Drug struct {
	Id   int
	Name string
}

type Assoc struct {
	Base   []int
	Added  []int
	Result float64
}

var sharedSession *mgo.Session
var drugs []Drug

func DrugHandler(w http.ResponseWriter, r *http.Request) {
	baseNames := r.URL.Query()["base"]
	baseNums := []int{}
	query := bson.M{"base": []int{}}

	for _, drug := range drugs {
		for _, name := range baseNames {
			if drug.Name == name {
				baseNums = append(baseNums, drug.Id)
				query["base"] = append(query["base"].([]int), drug.Id)
			}
		}
	}

	var assocs []Assoc
	sharedSession.DB("DDI").C("associations").Find(query).All(&assocs)

	drugNames := []string{}
	for _, assoc := range assocs {
		addedNames := []string{}
		for _, added := range assoc.Added {
			var drug Drug
			sharedSession.DB("DDI").C("drugs").Find(bson.M{"id": added}).One(&drug)
			addedNames = append(addedNames, drug.Name)
		}
		drugNames = append(drugNames, strings.Join(addedNames, ","))
	}

	w.Write([]byte(strings.Join(drugNames, ":")))
}

func DrugInfoHandler(w http.ResponseWriter, r *http.Request) {
	var baseNames []string
	baseNames = r.URL.Query()["base"]
	baseNums := []int{}

	var addedNames []string
	addedNames = r.URL.Query()["added"]
	addedNums := []int{}

	for _, drug := range drugs {
		for _, baseName := range baseNames {
			if drug.Name == baseName {
				baseNums = append(baseNums, drug.Id)
			}
		}
		for _, addedName := range addedNames {
			if drug.Name == addedName {
				addedNums = append(addedNums, drug.Id)
			}
		}
	}

	var assoc Assoc
	sharedSession.DB("DDI").C("associations").Find(bson.M{
		"base":  baseNums,
		"added": addedNums,
	}).One(&assoc)
	log.Println(assoc)

	w.Write([]byte(fmt.Sprintf("%f", assoc.Result)))
}

func getEnv() map[string]string {
	envvars := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := splits[1]
		envvars[key] = val
	}
	return envvars
}

func main() {
	env := getEnv()
	dbURL := env["MONGODB_URI"]
	dbName := env["DBNAME"]
	envPort := env["PORT"]

	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	sharedSession = session
	defer session.Close()

	drugsColl := session.DB(dbName).C("drugs")

	drugs = []Drug{}
	drugsColl.Find(bson.M{}).All(&drugs)

	addr := ":" + envPort
	if envPort == "" || envPort == "8080" {
		addr = ":8080"
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/drugs", DrugHandler)
	http.HandleFunc("/drug", DrugInfoHandler)

	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
