package mongo

import (
	"gopkg.in/mgo.v2"
)

const DBname = "blogServer"

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
	return db.C("BlogPost")
}

func GetGoalCollection(db *mgo.Database) *mgo.Collection {
	return db.C("BlogPost")
}

func GetUserCollection(db *mgo.Database) *mgo.Collection {
	return db.C("User")
}