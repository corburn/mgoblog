package main

import (
	"fmt"
	"net/http"

	"labix.org/v2/mgo"
)

const (
	url      = "mongodb://localhost"
	database = "blog"
	sessions = "sessions"
	users    = "users"
)

var session *mgo.Session

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a place holder for the blog")
}

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}
