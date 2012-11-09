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
	c, session, err := connect(sessions)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	count, err := c.Count()
	if err != nil {
		t.Error(err)
	}
	sessionId, err := startSession(session, "Nemo")
	if err != nil {
		t.Error(err)
	}
	if n, err := c.Count(); n != count+1 {
		t.Error("It does not look like a session was started.")
	} else if err != nil {
		t.Error(err)
	}

	if err = c.Remove(sessionId); err != nil {
		t.Error(err)
	}
}

func Test_getSession(t *testing.T) {
	c, session, err := connect(sessions)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	t.Logf("Test SessionId does not exist.")
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: "Nemo"}
	if result, err := getSession(session, sessionId); err != mgo.ErrNotFound {
		t.Error(err)
	} else if result != nil {
		t.Error("result is not nil")
	}

	t.Logf("Test SessionId exists.")
	sessionId, err = startSession(session, "Nemo")
	if err != nil {
		t.Fatal(err)
	}
	if result, err := getSession(session, sessionId); err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("result is nil")
	}

	if err = c.Remove(sessionId); err != nil {
		t.Error(err)
	}
}

func Test_endSession(t *testing.T) {
	_, session, err := connect(sessions)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	t.Logf("Test SessionId does not exist.")
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: "Nemo"}
	if err := endSession(session, sessionId); err != mgo.ErrNotFound {
		t.Error(err)
	}
	t.Logf("Test SessionId exists.")
	sessionId, err = startSession(session, "john")
	if err != nil || sessionId == nil {
		t.Fatal(err)
	}
	err = endSession(session, sessionId)
	if err != nil {
		t.Error(err)
	}
}

func Test_newUser(t *testing.T) {
	c, session, err := connect(users)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Test without email")
	foo, err := newUser(session, "Foo", "password", "")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Test with email")
	nemo1, err := newUser(session, "Nemo", "password", "nemo@example.com")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Test duplicate user")
	nemo2, err := newUser(session, "Nemo", "password", "nemo@example.com")
	if !mgo.IsDup(err) {
		t.Error(err)
	}

	if err = c.Remove(foo); err != nil {
		t.Error(err)
	}
	if err = c.Remove(nemo1); err != nil {
		t.Error(err)
	}
	if err = c.Remove(nemo2); err != mgo.ErrNotFound {
		t.Error(err)
	}
}

func Test_hashStr(t *testing.T) {
	if s, err := hashStr("test"); s != "1e2379853564a6fd9ff69b0a99cd82d4" || err != nil {
		t.Error("Failed to return a md5 hmac", err)
	}
}
