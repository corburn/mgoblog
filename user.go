package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type SessionId struct {
	Id bson.ObjectId "_id,omitempty"
	Username string
}

// startSession will create a new document in the sessions collection and return the _id
func startSession(session *mgo.Session, username string) (*SessionId, error) {
	c := session.DB(database).C("sessions")
	sessionId := &SessionId{Id: bson.NewObjectId(), Username: username}
	if err := c.Insert(sessionId); err != nil {
		return nil, err
	}
	return sessionId, nil
}
