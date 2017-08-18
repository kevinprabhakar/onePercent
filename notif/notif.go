package notif

import (
	"gopkg.in/gomail.v2"
	"onePercent/mongo"

	"time"
	"onePercent/user"
	"gopkg.in/mgo.v2"
	"onePercent/util"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"onePercent/goal"
	"onePercent/post"
	"github.com/pkg/errors"
	"bytes"
	"html/template"

)

type NotifController struct{
	Session 		*mgo.Session
	Logger 			*util.Logger
}

func NewNotifController(session *mgo.Session, logger *util.Logger)(*NotifController){
	return &NotifController{session, logger}
}

type EmailTemplateData struct{
	UserName 		string
	PartnerName 	string
	MessageURL 		string
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (self *NotifController)SendEmail(email string, name string, userName string, userEmail string) (error){
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "one.percent.emailer@gmail.com", "One Percent Project")
	m.SetAddressHeader("To", email, name)
	m.SetHeader("Subject", "Hello!")

	token, err := user.GetMessageAccessToken(email, userEmail)
	if (err != nil){
		return err
	}

	urlLink := "http://localhost:3000/message.html?messageAccessToken="+token

	newTemplate := EmailTemplateData{UserName:userName,PartnerName:name,MessageURL:urlLink}

	sendingMessage, err := ParseTemplate("notif/emailTemplate.html", newTemplate)
	if (err != nil){
		return err
	}

	m.SetBody("text/html", sendingMessage)

	d := gomail.NewPlainDialer("smtp.gmail.com", 465, "one.percent.emailer@gmail.com", "Heb1Pet!")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (self *NotifController)GetAndCheckUsers(midnight time.Time)(error){
	var userList []user.User
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))
	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))
	findErr := userCollection.Find(nil).All(&userList)

	if (findErr != nil){
		return findErr
	}

	var emailErr error

	if (len(userList)==0){
		return errors.New("No Users In Database")
	}

	for _, user := range userList{
		if (len(user.Goals)!=0) {
			var UserGoal goal.Goal

			findGoalErr := goalCollection.Find(bson.M{"_id":user.Goals[0]}).One(&UserGoal)
			if (findGoalErr != nil) {
				return findGoalErr
			}
			if (len(UserGoal.Posts)!=0){
				lastPost := UserGoal.Posts[len(UserGoal.Posts) - 1]

				var userLastPost post.Post
				findPostErr := postCollection.Find(bson.M{"_id":lastPost}).One(&userLastPost)

				if (findPostErr != nil) {
					return findPostErr
				}

				lastDay := midnight.AddDate(0, 0, -1)

				if (!userLastPost.Created.After(lastDay)) {
					for _, checker := range user.CheckeeOf {
						emailErr = self.SendEmail(checker.Email, checker.Name, user.Name, user.Email)
						self.Logger.Debug("Email sent to " + checker.Email + " for " + user.Name)
					}
				}
			}
		}

	}
	if (emailErr != nil){
		return emailErr
	}
	return nil

}

func (self *NotifController)CheckAndSend(){
	for {
		hour := time.Now().Hour()
		self.Logger.Debug("Doing Notify Partners Time Check")
		if (hour == 0){
			self.Logger.Debug("Notify Partners Time Check Validated")
			midnight := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),0,0,0,0,time.Now().Location())
			checkErr := self.GetAndCheckUsers(midnight)
			if (checkErr != nil){
				fmt.Println("Checking And Sending Failed")
				self.Logger.Debug(checkErr.Error())
			}
		}
		time.Sleep(1 * time.Hour)
	}
}

func (self *NotifController)SendUserEmail(fromEmail string, toEmail string, messageSubject string, messageText string) (error){
	userCollection := mongo.GetUserCollection(mongo.GetDataBase(self.Session))
	var findUser user.User
	findErr := userCollection.Find(bson.M{"email":toEmail}).One(&findUser)

	if (findErr != nil){
		return findErr
	}

	var name string
	for _, partner := range findUser.CheckeeOf{
		if (partner.Email == fromEmail){
			name = partner.Name
		}
	}

	if (name == ""){
		return errors.New("No user exists with this email address")
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "one.percent.emailer@gmail.com", "One Percent Project")
	m.SetAddressHeader("To", toEmail, name)
	m.SetHeader("Subject", messageSubject)


	m.SetBody("text/plain", messageText)

	d := gomail.NewPlainDialer("smtp.gmail.com", 465, "one.percent.emailer@gmail.com", "Heb1Pet!")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}