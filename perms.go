package main

import (
	"encoding/json"
	"log"
)

type AssocPair struct {
	Base  []string
	Added []string
}

type Assoc struct {
	Base  []string
	Added []string
	Or    float64
}

func Perms2(src []string) []AssocPair {
	return []AssocPair{
		AssocPair{[]string{}, []string{src[0]}},
		AssocPair{
			[]string{src[0]}, []string{src[1]},
		},
		AssocPair{[]string{}, []string{src[1]}},
		AssocPair{
			[]string{src[1]}, []string{src[0]},
		},
		AssocPair{
			[]string{}, []string{src[0], src[1]},
		},
	}
}

func Perms3(src []string) []AssocPair {
	return []AssocPair{
		AssocPair{[]string{}, []string{src[0]}},
		AssocPair{
			[]string{src[0]},
			[]string{src[1]},
		},
		AssocPair{
			[]string{src[0]},
			[]string{src[2]},
		},
		AssocPair{
			[]string{src[0]},
			[]string{src[1], src[2]},
		},
		AssocPair{[]string{}, []string{src[1]}},
		AssocPair{
			[]string{src[1]},
			[]string{src[0]},
		},
		AssocPair{
			[]string{src[1]},
			[]string{src[2]},
		},
		AssocPair{
			[]string{src[1]},
			[]string{src[0], src[2]},
		},
		AssocPair{[]string{}, []string{src[2]}},
		AssocPair{
			[]string{src[2]},
			[]string{src[0]},
		},
		AssocPair{
			[]string{src[2]},
			[]string{src[1]},
		},
		AssocPair{
			[]string{src[2]},
			[]string{src[0], src[1]},
		},
		AssocPair{[]string{}, []string{src[0], src[1]}},
		AssocPair{
			[]string{src[0], src[1]},
			[]string{src[2]},
		},
		AssocPair{[]string{}, []string{src[0], src[2]}},
		AssocPair{
			[]string{src[0], src[2]},
			[]string{src[1]},
		},
		AssocPair{[]string{}, []string{src[1], src[2]}},
		AssocPair{
			[]string{src[1], src[2]},
			[]string{src[0]},
		},
		AssocPair{[]string{}, src},
	}
}

func GetPerms(class string, rank int) []byte {
	firstAssocs := []Assoc{}
	assocs := []Assoc{}

	if err := FetchAssociations(class, int(rank), &firstAssocs); err != nil {
		log.Println("Error fetching association for class and rank ", class, rank)
	}

	firstAssoc := firstAssocs[len(firstAssocs)-1]
	log.Println(firstAssoc)

	if len(firstAssoc.Base) > 0 && firstAssoc.Base[0] == "0" {
		firstAssoc.Base = []string{}
	}

	drugs := firstAssoc.Base
	for _, drug := range firstAssoc.Added {
		drugs = append(drugs, drug)
	}

	perms := []AssocPair{AssocPair{[]string{}, drugs}}
	if len(drugs) == 3 {
		perms = Perms3(drugs)
	} else if len(drugs) == 2 {
		perms = Perms2(drugs)
	}

	for _, pair := range perms {
		assoc := Assoc{}
		if err := FetchAssociation(pair.Base, pair.Added, &assoc); err != nil {
			log.Println("Error fetching association for base and added ", pair.Base, pair.Added)
			continue
		}
		assocs = append(assocs, assoc)
	}

	json, err := json.Marshal(assocs)
	if err != nil {
		log.Println("Error marshalling association json ", err)
	}
	return json
}
