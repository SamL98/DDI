package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

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

var dbname string

func main() {
	env := getEnv()
	dbURL := env["MONGODB_URI"]
	envPort := env["PORT"]
	dbname = env["DBNAME"]

	sess, err := OpenConnection(dbURL)
	if err != nil {
		panic(err)
	}

	session = sess
	defer session.Close()

	addr := ":" + envPort
	if envPort == "" || envPort == "8080" {
		addr = ":8080"
	}

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/show", IndexHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
