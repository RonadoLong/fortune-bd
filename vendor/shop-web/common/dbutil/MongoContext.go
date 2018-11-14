package dbutil

import "gopkg.in/mgo.v2"

type MongoContext struct {
	Session *mgo.Session
}

func (m *MongoContext)Close(){
	m.Session.Clone()
}

func (m *MongoContext)DBCollection(cname string) *mgo.Collection {
	return m.Session.DB(DBName).C(cname)
}

func NewDBContext() *MongoContext {
	session := GetSession()
	m := &MongoContext{}
	m.Session = session
	return m
}