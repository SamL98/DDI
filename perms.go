package main

import (
	"fmt"
	"log"
	"strings"
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

func GetPerms(class string, rank int) string {
	firstAssocs := []Assoc{}

	if err := FetchAssociations(class, int(rank), &firstAssocs); err != nil {
		log.Println("Error fetching association for class and rank ", class, rank)
		return err.Error()
	}

	if len(firstAssocs) == 0 {
		log.Println("No association fetched for class and rank ", class, rank)
		return "no assoc found"
	}

	firstAssoc := firstAssocs[len(firstAssocs)-1]

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

	assocStr := ""
	for i, pair := range perms {
		assoc := Assoc{}
		if err := FetchAssociation(pair.Base, pair.Added, &assoc); err != nil {
			log.Println("Error fetching association for base and added ", pair.Base, pair.Added)
			continue
		}

		if len(assoc.Base) > 0 {
			assocStr += strings.Join(assoc.Base, ", ") + "->"
		} else {
			assocStr += "BaseLine->"
		}
		assocStr += strings.Join(assoc.Added, ", ") + ":" + fmt.Sprintf("%.4f", assoc.Or)

		// TODO: Figure out why trailing | is added
		if i < len(perms)-1 {
			assocStr += "|"
		}
	}

	return assocStr
}
