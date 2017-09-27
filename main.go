package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	sharedSession = session
	defer session.Close()

	drugsColl := session.DB("DDI").C("drugs")

	drugs = []Drug{}
	drugsColl.Find(bson.M{}).All(&drugs)

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/drugs", DrugHandler)
	http.HandleFunc("/drug", DrugInfoHandler)

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
