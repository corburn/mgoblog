package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
	trace = log.New(os.Stdout, "trace:", log.Lshortfile)

	session   *mgo.Session
	templates = template.Must(template.ParseFiles(tmplDir + "signup.html", tmplDir + "login.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl string, u *User) {
	if err := templates.ExecuteTemplate(w, tmpl, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a place holder for the blog")

	trace.Println(r, "\n")
}

func signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	user := &User{Id: username, Password: []byte(password), Email: email}
	renderTemplate(w, "signup.html", user)

	trace.Println(r, "\n")
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := &User{Id: username, Password: []byte(password)}
	renderTemplate(w, "login.html", user)

	trace.Println(r, "\n")
}

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	http.HandleFunc("/", root)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}
