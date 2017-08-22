package mongo

import (
	"gopkg.in/mgo.v2"
	"os"
)

const DBname = "heroku_8hrctnmg"

func GetMongoSession(debug bool) (*mgo.Session){
	if (debug){
		uri := os.Getenv("MONGODB_URI")
		session, err := mgo.Dial(uri)

		if (err != nil){
			panic(err)
		}
		return session
	}else{
		uri := ("mongodb://localhost:27017/onePercent")
		session, err := mgo.Dial(uri)

		if (err != nil){
			panic(err)
		}
		return session
	}

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
