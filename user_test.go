package main

import (
	"testing"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func connect(collection string) (c *mgo.Collection, session *mgo.Session, err error) {
	session, err = mgo.Dial(url)
	c = session.DB(database).C(collection)
	return
}

func Test_startSession(t *testing.T) {
	c, session, err := connect("sessions")
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	count, err := c.Count()
	if err != nil {
		t.Fatal(err)
	}
	sessionId, err := startSession(session, "Nemo")
	if err != nil {
		t.Error(err)
	}
	defer c.Remove(sessionId)
	if n, err := c.Count(); n != count + 1 {
		t.Error("It does not look like a session was started.")
	} else if err != nil {
		t.Fatal(err)
	}
}

func Test_getSession(t *testing.T) {
	c, session, err := connect("sessions")
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	// getSession when sessionId does not exist
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: "Nemo"}
	if result, err := getSession(session, sessionId); err != mgo.ErrNotFound {
		t.Error(err)
	} else if result != nil {
		t.Error("result is not nil")
	}

	// getSession when sessionId exists
	sessionId, err = startSession(session, "Nemo")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Remove(sessionId)
	if result, err := getSession(session, sessionId); err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("result is nil")
	}
}

func Test_endSession(t *testing.T) {
	_, session, err := connect("sessions")
	defer session.Close()
	if err != nil {
		t.Fatal(err)
	}
	// sessionId does not exist
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: "Nemo"}
	if err := endSession(session, sessionId); err != mgo.ErrNotFound {
		t.Error(err)
	}
	// sessionId exists
	sessionId, err = startSession(session, "john")
	if err != nil || sessionId == nil {
		t.Fatal(err)
	}
	err = endSession(session, sessionId)
	if err != nil {
		t.Error(err)
	}
}
