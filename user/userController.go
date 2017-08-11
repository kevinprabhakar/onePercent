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
		return nil, errors.New("Password is too short")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	var findUser User

	err := userCollection.Find(bson.M{ "email": SignUpParams.Email}).One(&findUser)

	if (err != mgo.ErrNotFound){
		return nil, errors.New("User with this email Id already exists")
	}

	passHash, err := util.HashPassword(SignUpParams.Password)
	if (err != nil){
		return nil, err
	}

	newUser := User{
		Id			: bson.NewObjectId(),
		Name		: SignUpParams.Name,
		PassHash	: passHash,
		Goals 		: []bson.ObjectId{},
		CheckerOf	: []bson.ObjectId{},
		CheckeeOf   : []bson.ObjectId{},
		Email 		: SignUpParams.Email,
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
		return nil, errors.New("Email Field Missing")
	}
	if(len(SignInParams.Password) == 0){
		return nil, errors.New("Password Field Missing")
	}

	var verifyUser User

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	findErr := userCollection.Find(bson.M{"email" : SignInParams.Email}).One(&verifyUser)

	if (findErr != nil){
		return nil, errors.New("User doesn't exist")
	}

	passwordMatch := util.CheckPasswordHash(SignInParams.Password, verifyUser.PassHash)
	if (!passwordMatch){
		return nil, errors.New("Invalid Password")
	}

	return &verifyUser, nil
}

func (self *UserController)GetCurrUser(accessToken string)(*User, error){
	uid, err := VerifyAccessToken(accessToken)

	if (err != nil){
		return nil, err
	}
	if (!bson.IsObjectIdHex(uid)){
		return nil, errors.New("Invalid BSON Id")
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
		return nil, decErr
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	var userIdListBson []bson.ObjectId

	for _, uid := range userIds.IdList{
		if (!bson.IsObjectIdHex(uid)){
			return nil, errors.New("Invalid BSON Id")
		}else{
			userIdListBson = append(userIdListBson, bson.ObjectIdHex(uid))
		}
	}

	var userList []User

	findErr := userCollection.Find(bson.M{"_id" : bson.M{"$in" : userIdListBson}}).All(&userList)

	if (findErr != nil){
		return nil, findErr
	}

	return &userList, nil
}