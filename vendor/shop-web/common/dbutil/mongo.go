package dbutil

import "gopkg.in/mgo.v2"

const (
	DBName = "shop"
	Host = "mongodb://root:uqJiSCyj7Yam@35.208.24.119:27017"
)

var session *mgo.Session

func createDBSession()  {

	var err error
	session, err = mgo.Dial(Host)
	if err != nil {
		logger.Error(err)
	}
}

func GetSession()  *mgo.Session{
	if session == nil {
		createDBSession()
	}
	return session
}

func InitMongoDb()  {
	createDBSession()
}