package user

import (
	"gopkg.in/mgo.v2"
	"onePercent/util"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"onePercent/mongo"
	"encoding/json"
	"strings"
	//"onePercent/goal"
	"onePercent/goal"
)

type CurrUserController struct {
	Uid 		bson.ObjectId
	Session		*mgo.Session
	Logger 		*util.Logger
}

func NewCurrUserController(uid string, session *mgo.Session, logger *util.Logger) (*CurrUserController, error){
	if (!bson.IsObjectIdHex(uid)){
		return nil, errors.New("Invalid BSON Id")
	}
	if (session == nil){
		return nil, errors.New("Mongo Session Nil")
	}
	if (logger == nil){
		return nil, errors.New("Logger is Nil")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(session))

	var verifyUser User
	findErr:= userCollection.Find(bson.M{"_id" : bson.ObjectIdHex(uid)}).One(&verifyUser)

	if (findErr != nil)&&(findErr != mgo.ErrNotFound){
		return nil, findErr
	}else if (findErr == mgo.ErrNotFound){
		return nil, errors.New("Invalid User Id")
	}

	return &CurrUserController{bson.ObjectIdHex(uid), session, logger}, nil
}

func (self *CurrUserController)AddGoal(params string) (error){
	decoder := json.NewDecoder(strings.NewReader(params))
	var goalParams GoalCreateParams
	decodeErr := decoder.Decode(&goalParams)

	if (decodeErr != nil){
		return decodeErr
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))

	var currUser User
	findErr := userCollection.Find(bson.M{"_id":self.Uid}).One(&currUser)

	if (findErr != nil){
		return findErr
	}

	insertGoal := goal.NewGoal(currUser.Id,goalParams.Name,goalParams.Description,goalParams.Created,goalParams.UpdateBy)

	insertGoalError := goalCollection.Insert(insertGoal)

	if (insertGoalError != nil){
		return insertGoalError
	}

	updateErr := userCollection.Update(bson.M{"_id":self.Uid},bson.M{"$push":bson.M{"goals" : insertGoal.Id}})

	if (updateErr != nil){
		return updateErr
	}

	return nil
}

func (self *CurrUserController)AddCheckers(params string)(*[]User, error){
	dec := json.NewDecoder(strings.NewReader(params))
	var checkersList UserIdList
	decodeErr := dec.Decode(&checkersList)

	if (decodeErr != nil){
		return nil, decodeErr
	}

	for _, id := range checkersList.IdList {
		if (!bson.IsObjectIdHex(id)){
			return nil, errors.New("User Id given isn't valid BSON id")
		}
	}

	userIdList := make([]bson.ObjectId,0)

	for _, id := range checkersList.IdList {
		userIdList = append(userIdList, bson.ObjectIdHex(id))
	}
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(mongo.GetMongoSession()))

	UserList := make([]User, 0)

	findUsersErr := userCollection.Find(bson.M{"_id" : bson.M{"$in" : userIdList}}).All(&UserList)

	if (findUsersErr != nil){
		return nil, findUsersErr
	}

	for _, User := range UserList {
		if (User.Id == self.Uid){
			return nil, errors.New("Error: Cannot add self as checker")
		}
	}

	matchQuery := bson.M{"_id" : self.Uid}
	changeQuery := bson.M{"$addToSet" : bson.M{"checkeeOf" : bson.M{"$each" : userIdList}}}

	updateErr := userCollection.Update(matchQuery,changeQuery)

	if (updateErr != nil){
		return nil, updateErr
	}
	return &UserList, nil
}