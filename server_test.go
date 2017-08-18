package main

import (
	"testing"
	"onePercent/goal"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"onePercent/user"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"onePercent/mongo"
	"onePercent/post"
	"onePercent/util"
	//"onePercent/notif"
	"onePercent/notif"
)

var TestUserController = user.NewUserController(MongoSession, ServerLogger)
var TestNotifController = notif.NewNotifController(MongoSession,ServerLogger)

func TestSendEmail(t *testing.T){
	TestNotifController.CheckAndSend()
}

func TestGetStringJson(t *testing.T){
	sampleGoalCreateParams := goal.NewGoal(bson.NewObjectId(),"Name","description",0,0)

	jsonForm, err := util.GetStringJson(sampleGoalCreateParams)

	if (err != nil){
		t.Fail()
		return
	}
	fmt.Println(string(jsonForm))
}

func TestSignUp(t *testing.T){
	sampleUserSignUpParams := user.UserSignUpParams{
		Email: "wubble@gmail.com",
		Password: "bottle",
		Name: "Bruce Wayne",
	}

	jsonForm, err := util.GetStringJson(sampleUserSignUpParams)

	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	params := url.Values{}
	params.Set("p",jsonForm)

	resp, err := http.PostForm("http://localhost:3000/api/signup", params)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(string(respBody))
}

func TestSignIn(t *testing.T){
	sampleUserSignInParams := user.UserSignInParams{
		Email : "kevin.surya@gmail.com",
		Password: "bottle",
	}

	jsonForm, err := util.GetStringJson(sampleUserSignInParams)

	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	params := url.Values{}
	params.Set("p",jsonForm)

	resp, err := http.PostForm("http://localhost:3000/api/signin", params)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(string(respBody))
}

func TestAddGoal(t *testing.T){
	sampleUserSignInParams := user.UserSignInParams{
		Email : "kevin.surya@gmail.com",
		Password: "bottle",
	}

	jsonForm, err := util.GetStringJson(sampleUserSignInParams)

	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	params := url.Values{}
	params.Set("p",jsonForm)

	resp, err := http.PostForm("http://localhost:3000/api/signin", params)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	accessToken := string(respBody)
	fmt.Println(accessToken)

	addingUser, err := TestUserController.GetCurrUser(accessToken)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	sampleGoalCreateParams := user.GoalCreateParams{
		Owner: addingUser.Id.Hex(),
		Name: "Sample Name",
		Description: "Sample Description",
		Created: time.Now().Unix(),
		UpdateBy: time.Now().Unix(),
	}

	goalJsonForm, err := util.GetStringJson(sampleGoalCreateParams)

	if (err != nil){
		t.Fail()
		return
	}

	urlParams2 := url.Values{}
	urlParams2.Set("accessToken", accessToken)
	urlParams2.Set("p", goalJsonForm)

	resp2, err := http.PostForm("http://localhost:3000/api/addgoal", urlParams2)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody2, err := ioutil.ReadAll(resp2.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(string(respBody2))

}

func TestVerifyAccess(t *testing.T){
	a, b := user.VerifyAccessToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDI2NDU5MzUsInVpZCI6IjU5OGM5MTg4NzkxZTRiNTBlZWYxMmUzZSJ9.T6uSQFKwJHRJXF0joRwUDPnFL3Zny93ky-U_vd-IQ8Q")
	if (b != nil){
		fmt.Println(b.Error())
	}
	fmt.Println(a)
}

func TestAddPost(t *testing.T){
	sampleUserSignInParams := user.UserSignInParams{
		Email : "kevin.surya@gmail.com",
		Password: "bottle",
	}

	jsonForm, err := util.GetStringJson(sampleUserSignInParams)

	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	params := url.Values{}
	params.Set("p",jsonForm)

	resp, err := http.PostForm("http://localhost:3000/api/signin", params)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	accessToken := string(respBody)

	uid, _ := user.VerifyAccessToken(accessToken)

	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(MongoSession))

	var findGoal goal.Goal

	findErr := goalCollection.Find(bson.M{"owner" : bson.ObjectIdHex(uid)}).One(&findGoal)

	if (findErr != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	samplePostCreateParams := post.PostCreateParams{
		Action: " Sample Action",
		Feeling: "Sample Feeling",
		Learning: "Sample Learning",
		Owner: uid,
		Goal: findGoal.Id.Hex(),
		Created: time.Now().Unix(),
	}

	jsonPostParams, err := util.GetStringJson(samplePostCreateParams)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	urlValues2 := url.Values{}
	urlValues2.Set("accessToken", accessToken)
	urlValues2.Set("p", jsonPostParams)

	resp2, err := http.PostForm("http://localhost:3000/api/addpost", urlValues2)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody2, err := ioutil.ReadAll(resp2.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(string(respBody2))
}

func TestAddCheckee(t *testing.T){
	sampleUserSignInParams := user.UserSignInParams{
		Email : "kevin.surya@gmail.com",
		Password: "bottle",
	}

	jsonForm, err := util.GetStringJson(sampleUserSignInParams)

	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	params := url.Values{}
	params.Set("p",jsonForm)

	resp, err := http.PostForm("http://localhost:3000/api/signin", params)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	accessToken := string(respBody)

	userCollection := mongo.GetUserCollection(mongo.GetDataBase(MongoSession))

	var CheckerUser user.User

	userCollection.Find(bson.M{"email":"wubble@gmail.com"}).One(&CheckerUser)

	idList := []string{CheckerUser.Id.Hex()}

	CheckersList := user.UserIdList{idList}

	jsonForm2, _ := util.GetStringJson(CheckersList)

	urlValues2 := url.Values{}
	urlValues2.Set("accessToken", accessToken)
	urlValues2.Set("p", jsonForm2)

	resp2, err := http.PostForm("http://localhost:3000/api/addcheckeeof", urlValues2)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	respBody2, err := ioutil.ReadAll(resp2.Body)
	if (err != nil){
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	fmt.Println(string(respBody2))
}

func TestAccessTokenMessage(t *testing.T){
	from := "kevin.surya@gmail.com"
	to := "prabhakk@usc.edu"
	token, err := user.GetMessageAccessToken(from, to)
	if (err != nil){
		fmt.Println(err.Error())
		return
	}

	fmt.Println(token)

	a,b,c := user.VerifyMessageAccessToken(token)
	if (c != nil){
		fmt.Println(c.Error())
		return
	}

	fmt.Println(a + " " + b)
}