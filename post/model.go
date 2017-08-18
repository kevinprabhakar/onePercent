package post

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	Id 			bson.ObjectId 		`bson:"_id"`
	Action 		string 				`bson:"action"`
	Feeling 	string 				`bson:"feeling"`
	Learning 	string 				`bson:"learning"`
	Owner 		bson.ObjectId 		`bson:"owner"`
	Goal 		bson.ObjectId		`bson:"goal"`
	Created 	time.Time 			`bson:"created"`
}

type PostCreateParams struct {
	Action 		string 				`json:"action"`
	Feeling 	string 				`json:"feeling"`
	Learning 	string 				`json:"learning`
	Owner 		string 				`json:"owner"`
	Goal 		string 				`json:"goal"`
	Created 	int64 				`json:"created"`
}

