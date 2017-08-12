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
)

var MongoSession = mongo.GetMongoSession()
var ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
var UserController = user.NewUserController(MongoSession, ServerLogger)
var GoalController = goal.NewGoalController(MongoSession, ServerLogger)

func main(){
	//p : {"email" : <Email address as string>, "password": <Unhashed password>, "name": <Name>}
	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received New SignUp Request")
		r.ParseForm()
		params := r.Form.Get("p")
		ServerLogger.Debug("params: " + params)

		newUser, err := UserController.SignUp(params)

		if (err != nil){
			ServerLogger.ErrorMsg("Could not Sign Up New User")
			fmt.Fprintf(w, err.Error())
			return
		}

		accessToken, err := user.GetAccessToken(newUser.Id.Hex())

		if (err != nil){
			ServerLogger.ErrorMsg("Could not generate new access Token")
			fmt.Fprintf(w, err.Error())
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
			fmt.Fprintf(w, err.Error())
			return
		}

		accessToken, err := user.GetAccessToken(findUser.Id.Hex())

		if (err != nil){
			ServerLogger.ErrorMsg("Could not generate new access Token")
			fmt.Fprintf(w, err.Error())
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
			fmt.Fprintf(w, err.Error())
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			fmt.Fprintf(w, err.Error())
			return
		}

		goalError := currUserController.AddGoal(params)

		if (goalError != nil){
			ServerLogger.ErrorMsg("Could not add goal to user profile")
			fmt.Fprintf(w, err.Error())
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
			fmt.Fprintf(w, err.Error())
			return
		}

		post, err := GoalController.AddPost(params)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not add post to goal")
			fmt.Fprintf(w, err.Error())
			return
		}

		jsonForm, err := json.Marshal(post)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not marshall JSON post")
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, string(jsonForm))
	})

	//accessToken: <accessToken from cookie>
	//p : {"idList" : <contains array of strings of user ids>}
	http.HandleFunc("/api/addcheckeeof", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Add Goal Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")
		params := r.Form.Get("p")

		addingUser, err := UserController.GetCurrUser(accessToken)
		if (err != nil){
			fmt.Fprintf(w, err.Error())
			return
		}

		currUserController, err := user.NewCurrUserController(addingUser.Id.Hex(),MongoSession, ServerLogger)
		if (err != nil){
			ServerLogger.ErrorMsg("Couldn't get curr user controller")
			fmt.Fprintf(w, err.Error())
			return
		}

		checkersList, addErr := currUserController.AddCheckers(params)
		if (addErr != nil){
			ServerLogger.ErrorMsg("Could not marshall add Checkers to user profile")
			fmt.Fprintf(w, addErr.Error())
			return
		}

		jsonForm, err := json.Marshal(checkersList)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not marshall JSON checkers list")
			fmt.Fprintf(w, err.Error())
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
			fmt.Fprintf(w, err.Error())
			return
		}

		userList, userGetErr := UserController.GetUsers(params)
		if (userGetErr != nil){
			fmt.Fprintf(w, userGetErr.Error())
			return
		}

		jsonForm, jsonErr := util.GetStringJson(userList)
		if (jsonErr != nil){
			fmt.Fprintf(w,jsonErr.Error())
			return
		}
		fmt.Fprintf(w,jsonForm)
	})

	//accessToken: <accessToken from cookie>
	http.HandleFunc("/api/verifyaccesstoken", func(w http.ResponseWriter, r *http.Request){
		ServerLogger.Debug("Received Verify Access Token Request")
		r.ParseForm()
		accessToken := r.Form.Get("accessToken")

		ServerLogger.Debug(accessToken)
		ServerLogger.Debug(r.URL.Hostname())

		uid, err := user.VerifyAccessToken(accessToken)
		if (err != nil){
			ServerLogger.ErrorMsg("Could not verify access token")
			fmt.Fprintf(w, err.Error())
			return
		}

		successMap := map[string]string{"userId":uid}
		jsonForm, err := util.GetStringJson(successMap)
		if (err != nil){
			ServerLogger.ErrorMsg(err.Error())
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w,jsonForm)
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.ListenAndServe(":3000", nil)
}