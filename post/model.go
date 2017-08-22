package post

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	Id 			bson.ObjectId 		`bson:"_id" json:"id"`
	Action 		string 				`bson:"action" json:"action"`
	Feeling 	string 				`bson:"feeling" json:"feeling"`
	Learning 	string 				`bson:"learning" json:"learning"`
	Owner 		bson.ObjectId 		`bson:"owner" json:"owner"`
	Goal 		bson.ObjectId		`bson:"goal" json:"goal"`
	Created 	time.Time 			`bson:"created" json:"created"`
}

type PostCreateParams struct {
	Action 		string 				`json:"action"`
	Feeling 	string 				`json:"feeling"`
	Learning 	string 				`json:"learning`
	Owner 		string 				`json:"owner"`
	Goal 		string 				`json:"goal"`
	Created 	int64 				`json:"created"`
}

