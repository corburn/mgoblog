package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type SessionId struct {
	Id bson.ObjectId "_id,omitempty"
	Username string
}

// startSession will create a new document in the sessions collection and return the _id
func startSession(session *mgo.Session, username string) (*SessionId, error) {
	c := session.DB(database).C(sessions)
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: username}
	if err := c.Insert(sessionId); err != nil {
		return nil, err
	}
	return sessionId, nil
}

// getSession returns the requested SessionId if it exists
// err == mgo.ErrNotFound if the SessionId does not exist
func getSession(session *mgo.Session, sessionId *SessionId) (result bson.M, err error) {
	c := session.DB(database).C(sessions)
	err = c.Find(sessionId).One(&result)
	return
}

// endSession will end a new user session by deleting it from the sessions table
func endSession(session *mgo.Session, sessionId *SessionId) error {
	c := session.DB(database).C(sessions)
	err := c.Remove(sessionId)
	return err
}

type User struct {
	Id string `bson:"_id"`
	Password []byte
	Email string `bson:",omitempty"`
}

// newUser creates a new user in the database
func newUser(session *mgo.Session, username, password, email string) (*User, error) {
	pHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	user := &User{Id: username, Password: pHash, Email: email}
	c := session.DB(database).C(users)
	err = c.Insert(user)
	return user, err
}
