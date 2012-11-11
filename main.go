package main

import (
	"fmt"
	"net/http"
	"html/template"

	"labix.org/v2/mgo"
)

const (
	url      = "mongodb://localhost"
	database = "blog"
	sessions = "sessions"
	users    = "users"

	tmplDir = "tmpl/"
)

var (
	session *mgo.Session
	templates = template.Must(template.ParseFiles(tmplDir+"signup.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl string, u *User) {
	if err := templates.ExecuteTemplate(w, tmpl, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a place holder for the blog")
}

func signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	user := &User{Id: username, Password: []byte(password), Email: email}
	renderTemplate(w, "signup.html", user)
}

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	http.HandleFunc("/", root)
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}
