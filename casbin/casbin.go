package main

import (
	"log"

	"github.com/casbin/casbin/v2"
)

func main() {
	// ref RESTful https://casbin.io/docs/supported-models#examples
	// Model  https://github.com/casbin/casbin/blob/master/examples/keymatch_model.conf
	// Policy https://github.com/casbin/casbin/blob/master/examples/keymatch_policy.csv

	enforcer, err := casbin.NewEnforcer("./rest_model.conf", "./rest_policy.csv")
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	sub := "alice"
	obj := "/alice_data/something"
	act := "GET"
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		log.Fatalf("enforce error: %v", err)
	}

	log.Printf("%s %s %s is %v", sub, act, obj, ok)

	sub = "bob"
	ok, err = enforcer.Enforce(sub, obj, act)
	if err != nil {
		log.Fatalf("enforce error: %v", err)
	}

	log.Printf("%s %s %s is %v", sub, act, obj, ok)
}
