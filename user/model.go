package user

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id				bson.ObjectId		`bson:"_id"`
	Name 			string 				`bson:"name"`
	PassHash 		string 				`bson:"passHash"`
	Goals 			[]bson.ObjectId 	`bson:"goals"`
	CheckerOf		[]bson.ObjectId 	`bson:"checkerOf"`
	CheckeeOf		[]bson.ObjectId		`bson:"checkeeOf"`
	Email 			string				`bson:"email"`
}
