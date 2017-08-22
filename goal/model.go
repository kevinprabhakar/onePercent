package goal

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Goal struct {
	Id 				bson.ObjectId 		`bson:"_id" json:"id"`
	Owner 			bson.ObjectId		`bson:"owner" json:"owner"`
	Posts 			[]bson.ObjectId		`bson:"posts" json:"posts"`
	Name 			string 				`bson:"name" json:"name"`
	Description 	string 				`bson:"description" json:"description"`
	Created 		time.Time 			`bson:"created" json:"created"`
	UpdateBy 		time.Time 			`bson:"updateBy" json:"updateBy"`
}

type StreakReturn struct {
	LongW 			int 				`json:"longW"`
	CurrW 			int 				`json:"currW"`
}