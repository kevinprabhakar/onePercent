package post

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	Id 			bson.ObjectId 		`bson:"_id"`
	Action 		string 				`bson:"action"`
	Feeling 	string 				`bson:"feeling"`
	Learning 	string 				`bson:"string"`
	Owner 		bson.ObjectId 		`bson:"owner"`
	Goal 		bson.ObjectId		`bson:"goal"`
	Created 	bson.ObjectId 		`bson:"created"`
}
