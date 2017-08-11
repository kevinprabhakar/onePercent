package mongo

import (
	"gopkg.in/mgo.v2"
)

const DBname = "onePercent"

func GetMongoSession() (*mgo.Session){
	session, err := mgo.Dial("mongodb://localhost")

	if (err != nil){
		panic(err)
	}

	return session
}

func GetDataBase(session *mgo.Session) *mgo.Database{
	return session.DB(DBname)
}

func GetPostCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Post")
}

func GetGoalCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Goal")
}

func GetUserCollection(db *mgo.Database) *mgo.Collection {
	return db.C("User")
}
