package goal

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Goal struct {
	Id 				bson.ObjectId 		`bson:"_id"`
	Owner 			bson.ObjectId		`bson:"owner"`
	Posts 			[]bson.ObjectId		`bson:"posts"`
	Name 			string 				`bson:"name"`
	Description 	string 				`bson:"description"`
	Created 		time.Time 			`bson:"created"`
	UpdateBy 		time.Time 			`bson:"updateBy"`
}

