package main

import (
	"testing"

	"labix.org/v2/mgo"
)

func Test_startSession(t *testing.T) {
	session, err := mgo.Dial(url)
	defer session.Close()
	if err != nil {
		t.Fatal(err)
	}
	c := session.DB(database).C("sessions")
	count, err := c.Count()
	if err != nil {
		// TODO: may not be fatal
		t.Fatal(err)
	}
	sessionId, err := startSession(session, "john")
	if err != nil {
		t.Error(err)
	}
	defer c.Remove(sessionId)
	if n, err := c.Count(); n != count + 1 {
		t.Error("It does not look like a session was started.")
	} else if err != nil {
		// TODO: may not be fatal
		t.Fatal(err)
	}
}
