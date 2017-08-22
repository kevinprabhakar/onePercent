package user

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"strings"
	"onePercent/mongo"
	"onePercent/util"
	"errors"
	//"fmt"
	"golang.org/x/net/html"
)

type UserController struct {
	Session 		*mgo.Session
	Logger 			*util.Logger
}

func NewUserController(session *mgo.Session, logger *util.Logger)(*UserController) {
	return &UserController{session, logger}
}

func (self *UserController)SignUp(params string)(*User, error){
	decoder := json.NewDecoder(strings.NewReader(params))
	var SignUpParams UserSignUpParams
	decodeErr := decoder.Decode(&SignUpParams)

	if (decodeErr != nil){
		return nil, decodeErr
	}

	if (!util.IsValidEmail(SignUpParams.Email)){
		return nil, errors.New("InvalidEmailAddress")
	}

	if(len(SignUpParams.Password) < 6){
		return nil, errors.New("PasswordTooShort")
	}

	if (SignUpParams.Password != SignUpParams.PasswordVerify){
		return nil, errors.New("PasswordsDontMatch")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	var findUser User

	err := userCollection.Find(bson.M{ "email": SignUpParams.Email}).One(&findUser)

	if (err != mgo.ErrNotFound){
		return nil, err
	}

	passHash, err := util.HashPassword(SignUpParams.Password)
	if (err != nil){
		return nil, err
	}

	newUser := User{
		Id			: bson.NewObjectId(),
		Name		: strings.Replace(html.EscapeString(SignUpParams.Name),"&#39;","'",-1),
		PassHash	: passHash,
		Goals 		: []bson.ObjectId{},
		CheckerOf	: []bson.ObjectId{},
		CheckeeOf   : []Checker{},
		Email 		: html.EscapeString(SignUpParams.Email),
		CurrConsecLose: 0,
		CurrConsecWin: 0,
		LongConsecLose: 0,
		LongConsecWin: 0,
	}

	insertErr := userCollection.Insert(newUser)

	if (insertErr != nil){
		return nil, insertErr
	}

	return &newUser, nil
}

func (self *UserController)SignIn(params string)(*User, error){
	decoder := json.NewDecoder(strings.NewReader(params))
	var SignInParams UserSignInParams
	decodeErr := decoder.Decode(&SignInParams)

	if (decodeErr != nil){
		return nil, decodeErr
	}
	if (!util.IsValidEmail(SignInParams.Email)){
		return nil, errors.New("MissingEmailField")
	}
	if(len(SignInParams.Password) == 0){
		return nil, errors.New("MissingPasswordField")
	}

	var verifyUser User

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	findErr := userCollection.Find(bson.M{"email" : SignInParams.Email}).One(&verifyUser)

	if (findErr != nil){
		return nil, errors.New("NonexistentUser")
	}

	passwordMatch := util.CheckPasswordHash(SignInParams.Password, verifyUser.PassHash)
	if (!passwordMatch){
		return nil, errors.New("InvalidPassword")
	}

	return &verifyUser, nil
}

func (self *UserController)GetCurrUser(accessToken string)(*User, error){
	uid, err := VerifyAccessToken(accessToken)

	if (err != nil){
		return nil, err
	}
	if (!bson.IsObjectIdHex(uid)){
		return nil, errors.New("InvalidBSONId")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	var returnUser User

	findErr := userCollection.Find(bson.M{ "_id" : bson.ObjectIdHex(uid)}).One(&returnUser)

	if (findErr != nil){
		return nil, findErr
	}

	return &returnUser, nil
}

func (self *UserController)GetUsers(params string)(*[]User, error){
	dec := json.NewDecoder(strings.NewReader(params))
	var userIds UserIdList
	decErr := dec.Decode(&userIds)

	if (decErr != nil){
		return nil, errors.New("CANT DECODE")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	var userIdListBson []bson.ObjectId

	for _, uid := range userIds.IdList{
		if (!bson.IsObjectIdHex(uid)){
			return nil, errors.New("InvalidBSONId")
		}else{
			userIdListBson = append(userIdListBson, bson.ObjectIdHex(uid))
		}
	}

	var userList []User

	findErr := userCollection.Find(bson.M{"_id" : bson.M{"$in" : userIdListBson}}).All(&userList)

	if (findErr != nil){
		return nil, errors.New("CANT FIND USERS FROM BSON")
	}

	return &userList, nil
}

func(self *UserController)DeleteUser(uid bson.ObjectId)(error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))
	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))

	var findUser User
	findUserErr := userCollection.Find(bson.M{"_id":uid}).One(&findUser)

	if (findUserErr != nil){
		return findUserErr
	}

	if (len(findUser.Goals)==0){
		return errors.New("UserHasNoGoals")
	}

	removePostErr := postCollection.Remove(bson.M{"owner":uid})

	if (removePostErr != nil){
		return removePostErr
	}

	removeGoalsErr := goalCollection.Remove(bson.M{"owner":uid})
	if (removeGoalsErr != nil){
		return removeGoalsErr
	}

	removeUserErr := userCollection.Remove(bson.M{"_id":uid})
	if (removeUserErr != nil){
		return removeUserErr
	}

	return nil
}