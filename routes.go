package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := templateHandler{filename: "home.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("public", t.filename)))
	})
	t.templ.Execute(w, make(map[string]interface{}))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := templateHandler{filename: "index.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("public", t.filename)))
	})

	class := r.URL.Query().Get("class")
	rankStr := r.URL.Query().Get("rank")

	rank, err := strconv.ParseInt(rankStr, 10, 64)
	if err != nil {
		log.Println("Cannot parse int from ", rankStr)
		return
	}

	data := map[string]interface{}{
		"drugs": GetPerms(class, int(rank)),
	}
	t.templ.Execute(w, data)
}
