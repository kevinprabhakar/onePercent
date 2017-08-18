package user

import (
	"gopkg.in/mgo.v2/bson"
	"onePercent/goal"
	"onePercent/post"
)

type Checker struct{
	Email 			string 				`bson:"email" json:"email"`
	Name 			string 				`bson:"name" json:"name"`
}

type User struct {
	Id				bson.ObjectId		`bson:"_id"`
	Name 			string 				`bson:"name"`
	PassHash 		string 				`bson:"passHash"`
	Goals 			[]bson.ObjectId 	`bson:"goals"`
	CheckerOf		[]bson.ObjectId 	`bson:"checkerOf"`
	CheckeeOf		[]Checker			`bson:"checkeeOf"`
	Email 			string				`bson:"email"`
}

type UserSignUpParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
	PasswordVerify	string 			`json:"passwordVerify"`
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

type CheckerEmailName struct{
	CheckerList		[]Checker			`json:"checkerList"`
}

type UserIdList struct {
	IdList 			[]string 			`json:"idList"`
}

type ChangePasswordParams struct{
	OldPassword 	string 				`json:"oldPassword"`
	NewPassword 	string 				`json:"newPassword"`
}

type GoalPosts struct{
	Goal 			goal.Goal 			`json:"goal"`
	Posts 			[]post.Post			`json:"posts"`
}

type MessageEmails struct{
	FromEmail		string 				`json:"fromEmail"`
	ToEmail 		string 				`json:"toEmail"`
}