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
	"onePercent/post"
	"golang.org/x/net/html"
)

type CurrUserController struct {
	Uid 		bson.ObjectId
	Session		*mgo.Session
	Logger 		*util.Logger
}

func NewCurrUserController(uid string, session *mgo.Session, logger *util.Logger) (*CurrUserController, error){
	if (!bson.IsObjectIdHex(uid)){
		return nil, errors.New("InvalidBSONId")
	}
	if (session == nil){
		return nil, errors.New("InvalidMongoSession")
	}
	if (logger == nil){
		return nil, errors.New("NilLogger")
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(session))

	var verifyUser User
	findErr:= userCollection.Find(bson.M{"_id" : bson.ObjectIdHex(uid)}).One(&verifyUser)

	if (findErr != nil)&&(findErr != mgo.ErrNotFound){
		return nil, findErr
	}else if (findErr == mgo.ErrNotFound){
		return nil, errors.New("NonexistentUser")
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

	insertGoal := goal.NewGoal(currUser.Id,
		strings.Replace(html.EscapeString(goalParams.Name),"&#39;","'",-1),
		strings.Replace(html.EscapeString(goalParams.Description),"&#39;","'",-1),
		goalParams.Created,
		goalParams.UpdateBy)

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

func (self *CurrUserController)AddCheckers(params string)(*[]Checker, error){
	dec := json.NewDecoder(strings.NewReader(params))
	var checkersList CheckerEmailName
	decodeErr := dec.Decode(&checkersList)

	if (decodeErr != nil){
		return nil, decodeErr
	}

	for index, _ := range checkersList.CheckerList {
		if (!util.IsValidEmail(checkersList.CheckerList[index].Email)){
			return nil, errors.New("InvalidEmailId")
		}
		if (len(checkersList.CheckerList[index].Name)==0){
			return nil,errors.New("InvalidName")
		}
		checkersList.CheckerList[index].Email = strings.Replace(html.EscapeString(checkersList.CheckerList[index].Email),"&#39;","'",-1)
		checkersList.CheckerList[index].Name = strings.Replace(html.EscapeString(checkersList.CheckerList[index].Name),"&#39;","'",-1)

	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))


	matchQuery := bson.M{"_id" : self.Uid}
	changeQuery := bson.M{"$addToSet" : bson.M{"checkeeOf" : bson.M{"$each" : checkersList.CheckerList}}}

	updateErr := userCollection.Update(matchQuery,changeQuery)

	if (updateErr != nil){
		return nil, updateErr
	}
	return &checkersList.CheckerList, nil
}

func (self *CurrUserController)RemoveCheckeeOf(email string)(error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	err := userCollection.Update(bson.M{"_id":self.Uid},bson.M{"$pull":bson.M{"checkeeOf":bson.M{"email":email}}})

	if (err != nil){
		self.Logger.Debug(err.Error())
		return err
	}
	return nil
}

func (self *CurrUserController)ChangeEmail(email string)(error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))

	findQuery := bson.M{"email":email}

	userCount, countErr := userCollection.Find(findQuery).Count()

	if (!util.IsValidEmail(email)){
		return errors.New("InvalidEmail")
	}

	if (countErr != nil){
		self.Logger.Debug(countErr.Error())
		return countErr
	}
	if (userCount != 0){
		errNew := errors.New("UserAlreadyExists")
		self.Logger.Debug(errNew.Error())
		return countErr
	}

	err := userCollection.Update(bson.M{"_id":self.Uid},bson.M{"$set":bson.M{"email" : html.EscapeString(email)}})

	if (err != nil){
		self.Logger.Debug(err.Error())
		return err
	}
	return nil
}

func (self *CurrUserController)ChangePassword(params string)(error){
	dec := json.NewDecoder(strings.NewReader(params))
	var changeParams ChangePasswordParams
	decErr := dec.Decode(&changeParams)

	if (decErr != nil){
		return decErr
	}

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	var findUser User
	findErr := userCollection.Find(bson.M{"_id":self.Uid}).One(&findUser)

	if (findErr != nil){
		return findErr
	}

	if (!util.CheckPasswordHash(changeParams.OldPassword, findUser.PassHash)){
		return errors.New("InvalidOldPassword")
	}

	newPassHash, passHashErr := util.HashPassword(changeParams.NewPassword)

	if (passHashErr != nil){
		return passHashErr
	}

	updateErr := userCollection.Update(bson.M{"_id": self.Uid},bson.M{"$set":bson.M{"passHash":newPassHash}})

	if (updateErr != nil){
		return updateErr
	}

	return nil
}

func(self *CurrUserController)GetDashboard()(*[]GoalPosts, error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))
	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))

	var findUser User

	findErr := userCollection.Find(bson.M{"_id":self.Uid}).One(&findUser)

	if (findErr != nil){
		return nil, findErr
	}

	var Goals []goal.Goal

	findGoalsErr := goalCollection.Find(bson.M{"_id":bson.M{"$in" : findUser.Goals}}).All(&Goals)

	if (findGoalsErr != nil){
		return nil, findGoalsErr
	}

	returnCollection := make([]GoalPosts,0)

	for _, goal := range Goals{
		var Posts []post.Post
		postFindErr := postCollection.Find(bson.M{"_id":bson.M{"$in":goal.Posts}}).All(&Posts)

		if (postFindErr != nil){
			return nil, postFindErr
		}

		returnCollection = append(returnCollection,GoalPosts{goal,Posts})
	}

	for index := range returnCollection{
		returnCollection[index].Goal.Description = strings.Replace(returnCollection[index].Goal.Description,"&#39;","'",-1)
		returnCollection[index].Goal.Name = strings.Replace(returnCollection[index].Goal.Name,"&#39;","'",-1)
	}
	return &returnCollection, nil
}

func(self *CurrUserController)GetGoals()(*[]goal.Goal, error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))

	var findUser User

	findErr := userCollection.Find(bson.M{"_id":self.Uid}).One(&findUser)

	if (findErr != nil){
		return nil, findErr
	}

	var Goals []goal.Goal

	findGoalsErr := goalCollection.Find(bson.M{"_id":bson.M{"$in" : findUser.Goals}}).All(&Goals)

	if (findGoalsErr != nil){
		return nil, findGoalsErr
	}

	return &Goals, nil
}
