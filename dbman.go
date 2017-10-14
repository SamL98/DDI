package main

import (
	"errors"
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

func OpenConnection(url string) (*mgo.Session, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("session is nil")
	}

	return session, nil
}

func FetchAssociations(class string, rank int, assocs *([]Assoc)) error {
	return session.DB(class).C("assocs").Find(bson.M{}).Limit(rank).Sort("-or").All(assocs)
}

func FetchAssociation(base []string, added []string, assoc *Assoc) error {
	class := "S6"
	if len(base) == 0 {
		class = fmt.Sprintf("S%d", len(added))
	} else if len(base) == 1 {
		class = fmt.Sprintf("S%d", len(added)+3)
	}

	return session.DB(class).C("assocs").Find(bson.M{
		"added": added, "base": base,
	}).One(assoc)
}
