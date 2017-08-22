package main

import(
	"onePercent/util"
	"io/ioutil"
	"os"
	"net/http"
	"fmt"
	"onePercent/mongo"
	"onePercent/user"
	"onePercent/goal"
	"encoding/json"
	"onePercent/notif"
)

var port = os.Getenv("PORT")

var MongoSession = mongo.GetMongoSession(true)
var ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
var UserController = user.NewUserController(MongoSession, ServerLogger)
var GoalController = goal.NewGoalController(MongoSession, ServerLogger)
var NotifController = notif.NewNotifController(MongoSession, ServerLogger)


func main(){
	go NotifController.CheckAndSend()

	//p : {"email" : <Email address as string>, "password": <Unhashed password>, "name": <Name>}
	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received New SignUp Request")
		r.ParseForm()
		params := r.Form.Get("p")
		ServerLogger.Debug("params: " + params)

		newUser, err := UserController.SignUp(params)

		if (err != nil){
			ServerLogger.ErrorMsg("Could not Sign Up New User")
			util.CustomError(w, err.Error(),400)
			return
		}

		accessToken, err := user.GetAccessToken(newUser.Id.Hex())

		if (err != nil){
			ServerLogger.ErrorMsg("Could not generate new access Token")
			util.CustomError(w, err.Error(),400)

			return
		}

		fmt.Fprintf(w, accessToken)
	})

	//p : {"email" : <Email address as string>, "password": <Unhashed password>}
	http.HandleFunc("/api/signin",func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received New SignIn Request")

		r.ParseForm()
		params := r.Form.Get("p")
		ServerLogger.Debug("params: " + params)

		findUser, err := UserController.SignIn(params)

		if (err != nil){
			ServerLogger.ErrorMsg("Could not sign in User")
			util.CustomError(w, err.Error(),400)
			return
		}

		accessToken, err := user.GetAccessToken(findUser.Id.Hex())

		if (err != nil){
			ServerLogger.ErrorMsg("Could not generate new access Token")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, accessToken)
	})

	//accessToken: <accessToken from cookie>
	//p : {"owner": <currUser UID as string>, "name":<Name of Goal as string>, "description":<Description of goal as string>,
	//     "created": <current time in unix seconds int64>, "updateBy":<time to update by in seconds>}
	http.HandleFunc("/api/addgoal", func (w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Add Goal Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		goalError := currUserController.AddGoal(params)

		if (goalError != nil){
			ServerLogger.ErrorMsg("Could not add goal to user profile")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, util.GetNoDataSuccessResponse())
	})

	//accessToken: <accessToken from cookie>
	//p : {"action":<action as string>, "feeling":<feeling as string>, "learning":<learning as string>, "owner":<owner uid as string>,
	//     "goal":<goal id as string>, "created":<int64 unix time of creation>}
	http.HandleFunc("/api/addpost", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Add Post Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")
		ServerLogger.Debug("params: " + params)

		_, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("Invalid Access Token")
			util.CustomError(w, err.Error(),400)
			return
		}

		post, err := GoalController.AddPost(params)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not add post to goal")
			util.CustomError(w, err.Error(),400)
			return
		}

		jsonForm, err := json.Marshal(post)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not marshall JSON post")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, string(jsonForm))
	})

	//accessToken: <accessToken from cookie>
	//p : {"idList" : <contains array of strings of user ids>}
	http.HandleFunc("/api/addcheckeeof", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Add Checkee Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")

		ServerLogger.Debug(params)

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		checkersList, addErr := currUserController.AddCheckers(params)
		if (addErr != nil){
			ServerLogger.ErrorMsg("Could not marshall add Checkers to user profile")
			util.CustomError(w, addErr.Error(),400)
			return
		}
		jsonForm, err := json.Marshal(checkersList)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not marshall JSON checkers list")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, string(jsonForm))
	})

	//accessToken: <accessToken from cookie>
	//p : {"idList" : <contains array of strings of user ids>}
	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")

		_, err := user.VerifyAccessToken(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		userList, userGetErr := UserController.GetUsers(params)
		if (userGetErr != nil){
			util.CustomError(w, userGetErr.Error(),400)
			return
		}

		jsonForm, jsonErr := util.GetStringJson(userList)
		if (jsonErr != nil){
			util.CustomError(w, jsonErr.Error(),400)
			return
		}
		fmt.Fprintf(w,jsonForm)
	})

	//accessToken: <accessToken from cookie>
	http.HandleFunc("/api/goal", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Goal Info Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get user")
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		goalList, err := currUserController.GetGoals()

		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user goals")
			util.CustomError(w, err.Error(),400)
			return
		}

		jsonForm, jsonErr := util.GetStringJson(goalList)
		if (jsonErr != nil){
			ServerLogger.ErrorMsg("Couldn't get server messages")
			util.CustomError(w, jsonErr.Error(),400)
			return
		}

		ServerLogger.Debug(jsonForm)
		fmt.Fprintf(w,jsonForm)
	})

	http.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Dashboard Info Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		dashboard, err := currUserController.GetDashboard()

		jsonForm, err := util.GetStringJson(dashboard)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't marshall JSON")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, jsonForm)
	})

	//accessToken: <accessToken from cookie>
	http.HandleFunc("/api/verifyaccesstoken", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Verify Access Token Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")

		uid, err := user.VerifyAccessToken(accessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("InvalidAccessToken")
			util.CustomError(w, err.Error(),400)
			return
		}

		successMap := map[string]string{"userId":uid}
		jsonForm, err := util.GetStringJson(successMap)
		if (err != nil){
			ServerLogger.ErrorMsg(err.Error())
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w,jsonForm)
	})


	http.HandleFunc("/api/removecheckeeof", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Remove Checkee Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		emailRemove := r.Form.Get("checkee")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		removeErr := currUserController.RemoveCheckeeOf(emailRemove)

		if (removeErr != nil){
			ServerLogger.ErrorMsg("Couldn't get remove checkee of ")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, util.GetNoDataSuccessResponse())
	})

	http.HandleFunc("/api/changeemail", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Change Email Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		newEmail := r.Form.Get("email")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		changeEmailErr := currUserController.ChangeEmail(newEmail)
		if (changeEmailErr != nil){
			ServerLogger.ErrorMsg("Couldn't change user email")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w,util.GetNoDataSuccessResponse())
	})

	http.HandleFunc("/api/changepassword", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Change Password Request")
		r.ParseForm()

		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")

		ServerLogger.Debug(params)

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			util.CustomError(w, err.Error(),400)
			return
		}

		changePassErr := currUserController.ChangePassword(params)

		if (changePassErr != nil){
			ServerLogger.ErrorMsg("Couldn't change User Password")
			util.CustomError(w, changePassErr.Error(),400)
			return
		}

		fmt.Fprintf(w, util.GetNoDataSuccessResponse())
	})

	http.HandleFunc("/api/verifymessageaccesstoken", func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		messageAccessToken := r.Form.Get("messageAccessToken")

		from, to, messageError := user.VerifyMessageAccessToken(messageAccessToken)
		if (messageError != nil){
			ServerLogger.ErrorMsg("Couldn't verify message accessToken")
			util.CustomError(w, messageError.Error(),400)
			return
		}

		returnEmails := user.MessageEmails{from,to}

		jsonForm, jsonErr := util.GetStringJson(returnEmails)
		if (jsonErr != nil){
			ServerLogger.ErrorMsg("Couldn't marshall message email addresses")
			util.CustomError(w, jsonErr.Error(),400)
			return
		}

		fmt.Fprintf(w, jsonForm)
	})

	http.HandleFunc("/api/sendUserEmail", func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		messageAccessToken := r.Form.Get("messageAccessToken")
		fromEmail := r.Form.Get("fromEmail")
		toEmail := r.Form.Get("toEmail")
		messageSubject := r.Form.Get("messageSubject")
		messageText := r.Form.Get("messageText")

		_,_, err := user.VerifyMessageAccessToken(messageAccessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get email addresses from access token")
			util.CustomError(w, err.Error(),400)
			return
		}

		sendErr := NotifController.SendUserEmail(fromEmail,toEmail,messageSubject,messageText)
		if (sendErr != nil){
			ServerLogger.ErrorMsg("Couldn't get email addresses from access token")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, util.GetNoDataSuccessResponse())

	})

	http.HandleFunc("/api/deleteaccount", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Delete User Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get User from Access Token")
			util.CustomError(w, err.Error(),400)
			return
		}

		deleteUserErr:= UserController.DeleteUser(addingUser.Id)
		if (deleteUserErr != nil){
			ServerLogger.ErrorMsg("Couldn't Delete User")
			util.CustomError(w, deleteUserErr.Error(),400)
			return
		}

		fmt.Fprintf(w,util.GetNoDataSuccessResponse())

	})

	http.HandleFunc("/api/getwinstreak", func(w http.ResponseWriter, r *http.Request) {
		ServerLogger.Debug("Received Win Streak Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		goalId := r.Form.Get("goalId")

		_, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			util.CustomError(w, err.Error(),400)
			return
		}

		allPosts, allPostsErr := GoalController.GetAllPosts(goalId)
		if (allPostsErr != nil){
			ServerLogger.ErrorMsg("Couldn't Get Posts for Goal Id")
			util.CustomError(w, allPostsErr.Error(),400)
			return
		}

		longW, currW, err := NotifController.GetStreakStats(allPosts)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't Get Streak Stats")
			util.CustomError(w, err.Error(),400)
			return
		}

		returnColl := goal.StreakReturn{longW,currW}
		jsonFrom, jsonErr := util.GetStringJson(returnColl)

		if (jsonErr != nil){
			ServerLogger.ErrorMsg("Couldn't Marshall JSON")
			util.CustomError(w, err.Error(),400)
			return
		}

		fmt.Fprintf(w, jsonFrom)
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.ListenAndServe(":"+port, nil)
}