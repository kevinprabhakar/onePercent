package goal

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"onePercent/util"
	"encoding/json"
	"strings"
	"onePercent/post"
	"onePercent/mongo"
	"errors"
	//"onePercent/user"

	"golang.org/x/net/html"
)

type GoalController struct {
	Session 		*mgo.Session
	Logger 			*util.Logger
}

func NewGoal(uid bson.ObjectId, name string, description string, created int64, updateBy int64)(*Goal){
	return &Goal{
		Id			: bson.NewObjectId(),
		Owner		: uid,
		Posts		: []bson.ObjectId{},
		Name		: name,
		Description	: description,
		Created		: time.Unix(created, 0),
		UpdateBy	: time.Unix(updateBy, 0),
	}
}

func NewGoalController(session *mgo.Session, logger *util.Logger)(*GoalController){
	return &GoalController{session,logger}
}

func (self *GoalController)AddPost(params string)(*post.Post, error){
	//Decoding JSON params
	dec := json.NewDecoder(strings.NewReader(params))
	var postParams post.PostCreateParams
	decErr := dec.Decode(&postParams)

	//Validating JSON params
	if (decErr != nil){
		return nil, decErr
	}
	if (!bson.IsObjectIdHex(postParams.Owner)){
		return nil, errors.New("InvalidBSONId")
	}else{
		userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
		_, findErr := userCollection.Find(bson.M{"_id": bson.ObjectIdHex(postParams.Owner)}).Count()

		if (findErr != nil){
			return nil, findErr
		}
	}
	if (!bson.IsObjectIdHex(postParams.Goal)){
		return nil, errors.New("InvalidBSONId")
	}else{
		goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))
		var findGoal Goal
		findErr := goalCollection.Find(bson.M{"_id":bson.ObjectIdHex(postParams.Goal)}).One(&findGoal)

		if (findErr != nil){
			return nil, findErr
		}
	}
	if (len(postParams.Action) == 0)||(len(postParams.Learning) == 0)||(len(postParams.Feeling) == 0){
		return nil, errors.New("EmptyPostField")
	}

	//Defining Databases and Creating insert post
	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))

	insertPost := post.Post{
		Id		: bson.NewObjectId(),
		Action	: strings.Replace(html.EscapeString(postParams.Action),"&#39;","'",-1),
		Feeling : strings.Replace(html.EscapeString(postParams.Feeling),"&#39;","'",-1),
		Learning: strings.Replace(html.EscapeString(postParams.Learning),"&#39;","'",-1),
		Owner 	: bson.ObjectIdHex(postParams.Owner),
		Goal 	: bson.ObjectIdHex(postParams.Goal),
		Created : time.Unix(postParams.Created, 0),
	}

	//Updating goal collection to include post id
	updateErr := goalCollection.Update(bson.M{"_id" : insertPost.Goal},bson.M{"$push":bson.M{"posts" : insertPost.Id}})
	if (updateErr != nil){
		return nil, updateErr
	}

	//Inserting Post into collection
	insertErr := postCollection.Insert(insertPost)
	if (insertErr != nil){
		return nil, insertErr
	}

	return &insertPost, nil
}

func (self *GoalController)GetAllPosts(goalId string)(*[]post.Post, error){
	if (!bson.IsObjectIdHex(goalId)){
		return nil, errors.New("InvalidBSONId")
	}
	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))

	var postList []post.Post

	findErr := postCollection.Find(bson.M{"goal":bson.ObjectIdHex(goalId)}).Sort("-created").All(&postList)

	if (findErr != nil){
		return nil, findErr
	}

	return &postList, nil
}

