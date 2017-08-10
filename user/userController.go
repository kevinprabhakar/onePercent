package user

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"strings"
	"onePercent/mongo"
	"onePercent/util"
	"errors"
)

type UserController struct {
	session 		*mgo.Session
}

type UserSignUpParams struct {
	Name 			string 			`json:"name"`
	Password 		string 			`json:"password"`
	Email 			string 			`json:"email"`
}

type UserSignInParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
}
func NewUserController(session *mgo.Session)(*UserController, error) {
	return &UserController{session}
}

func (self *UserController)SignUp(params string)(User, error){
	decoder := json.NewDecoder(strings.NewReader(params))
	var SignUpParams UserSignUpParams
	decodeErr := decoder.Decode(&SignUpParams)

	if (decodeErr != nil){
		return User{}, decodeErr
	}

	if (!util.IsValidEmail(SignUpParams.Email)){
		return User{}, errors.New("InvalidEmailAddress")
	}

	if(len(SignUpParams.Password) < 6){
		return User{}, errors.New("Password is too short")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.session))

	var findUser User

	err := userCollection.Find(bson.M{ "email": SignUpParams.Email}).One(findUser)

	if (err != mgo.ErrNotFound){
		return "", errors.New("User already exists")
	}

	passHash, err := util.HashPassword(SignUpParams.Password)
	if (err != nil){
		return User{}, err
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
		return User{}, insertErr
	}

	return newUser, nil
}

func (self *UserController)SignIn(params string)(User, error){
	decoder := json.NewDecoder(strings.NewReader(params))
	var SignInParams UserSignInParams
	decodeErr := decoder.Decode(&SignInParams)

	if (decodeErr != nil){
		return User{}, decodeErr
	}
	if (!util.IsValidEmail(SignInParams.Email)){
		return User{}, errors.New("Email Field Missing")
	}
	if(len(SignInParams.Password) == 0){
		return User{}, errors.New("Password Field Missing")
	}

	var verifyUser User

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.session))

	findErr := userCollection.Find(bson.M{"email" : SignInParams.Email}).One(&verifyUser)

	if (findErr != nil){
		return User{}, errors.New("User doesn't exist")
	}

	passwordMatch := util.CheckPasswordHash(SignInParams.Password, verifyUser.PassHash)
	if (!passwordMatch){
		return User{}, errors.New("Invalid Password")
	}

	return verifyUser, nil
}

func (self *UserController)GetUser(accessToken string)(User, error){
	uid, err := VerifyAccessToken(accessToken)

	if (err != nil){
		return User{}, err
	}
	if (!bson.IsObjectIdHex(uid)){
		return User{}, errors.New("Invalid BSON Id")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.session))

	var returnUser User

	findErr := userCollection.Find(bson.M{ "_id" : bson.ObjectIdHex(uid)}).One(&returnUser)

	if (findErr != nil){
		return User{}, findErr
	}
	return returnUser, nil

}