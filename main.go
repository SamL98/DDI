package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Drug struct {
	Id   int
	Name string
}

var drugs []Drug

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func DrugHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if fmt.Sprintf("%s", p) == "Start" {
			for i := range drugs {
				d, err := json.Marshal(drugs[i])
				if err != nil {
					log.Println(err)
					continue
				}
				conn.WriteMessage(messageType, d)
			}
			if err := conn.WriteMessage(messageType, []byte("End")); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	drugsColl := session.DB("DDI").C("drugs")

	drugs = []Drug{}
	drugsColl.Find(bson.M{}).All(&drugs)

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/drugs", DrugHandler)

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
