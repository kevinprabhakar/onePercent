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

type UserSignUpParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
	Name 			string 			`json:"name"`
}

type UserSignInParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
}

type GoalCreateParams struct {
	Owner 			string 				`json:"owner"`
	Name 			string 				`json:"name"`
	Description 	string 				`json:"description"`
	Created 		int64 				`json:"created"`
	UpdateBy 		int64 				`json:"updateBy"`
}

type UserIdList struct {
	IdList 			[]string 			`json:"idList"`
}