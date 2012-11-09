package main

import (
	"labix.org/v2/mgo"
)

const (
	url = "mongodb://localhost"
	database = "blog"
	sessions = "sessions"
	users = "users"
)

var session *mgo.Session

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
}
