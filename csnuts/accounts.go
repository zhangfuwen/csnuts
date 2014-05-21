package csnuts

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"text/template"
	"time"

	"github.com/gorilla/sessions"

	"appengine"
	"appengine/datastore"
)

type User struct {
	Nickname string
	Email    string
	Password string
	Domains  []string
	Intro    string
}

func init() {
}

func setEmailCookie(w http.ResponseWriter, email string) bool {
	expiration := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "email", Value: email, Expires: expiration}
	http.SetCookie(w, &cookie)
	return true
}
func deleteEmailCookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{Name: "email", Value: "s", Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 307)
	return
}

func (u *User) isLegal() (b bool, errString string) {
	b = true
	b = isRegularEmail(u.Email)
	if b == false {
		errString = "invalid email."
		return
	}
	b = isRegularPassword(u.Password)
	if b == false {
		errString = "invalid password"
		return
	}
	return
}

func isRegularEmail(email string) (b bool) {
	b, _ = regexp.MatchString("^[_a-zA-Z0-9][.a-zA-Z0-9_-]{0,20}@[a-zA-Z0-9-]{1,40}(.[a-zA-Z0-9-]{1,20}){1,10}$", email)
	return
}

/* 用户名：数字或字母开头, 4~20个字符*/
func isRegularNickname(nickname string) (b bool) {
	//	b,_=regexp.MatchString("^[a-zA-Z0-9_\x7f-\xff][a-zA-Z0-9_\x7f-\xff]{3,19}$", nickname)
	//上面那一句想使用汉字做用户名，但无奈不好用，暂时注释掉
	b, _ = regexp.MatchString("^[a-zA-Z0-9_][a-zA-Z0-9_]{3,19}$", nickname)
	return
}
func isRegularPassword(password string) (b bool) {
	b, _ = regexp.MatchString("^[_a-zA-Z0-9]{8,16}$", password)
	return
}
func isRegularDomain(domain string) (b bool) {
	b = len(domain) <= 40
	return
}
func isRegularIntro(intro string) (b bool) {
	b = len(intro) < 256
	return
}

func (u *User) saveSession(session *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	session.Values["nickname"] = u.Nickname
	session.Values["email"] = u.Email
	return session.Save(r, w)

}

func handleMyRegister(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	type Thispage struct {
		Errormsg string
	}
	var thispage Thispage
	signupPage, err := template.ParseFiles(templatePath + "login/signup.html")
	if err != nil {
		c.Errorf("%v", err)
		return
	}
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//模板文件login/signup.html中没有需要套用的数据
		if err := signupPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}

	//I don't know how to handle FormValue error
	var u User
	u.Password = r.FormValue("password")
	u.Email = r.FormValue("email")

	if b, _ := u.isLegal(); b != true {
		thispage.Errormsg = "Invalid user information."
		if err := signupPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	u.Password = md5Do(u.Password)
	k := datastore.NewKey(c, "User", u.Email, 0, nil)
	_, err = datastore.Put(c, k, &u)
	if err != nil {
		thispage.Errormsg = "Internal Error: Can't put data into database."
		if err := signupPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	setEmailCookie(w, u.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleMyLogin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	type Thispage struct {
		Errormsg string
	}
	var thispage Thispage
	loginPage, err := template.ParseFiles(templatePath + "login/login.html")
	if err != nil {
		c.Errorf("%v", err)
		return
	}
	if r.Method != "POST" {
		//FIXME:Put is ok, what is the difference between POST and PUT?
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//模板文件login/login.html中没有需要套用的数据
		if err := loginPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}

	u := new(User)
	email := r.FormValue("email")
	if !isRegularEmail(email) {
		thispage.Errormsg = "Invalid email address."
		if err := loginPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	password := r.FormValue("password")
	if !isRegularPassword(password) {
		thispage.Errormsg = "Invalid password."
		if err := loginPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	k := datastore.NewKey(c, "User", email, 0, nil)
	if err := datastore.Get(c, k, u); err != nil {
		thispage.Errormsg = "Account does not exist."
		if err := loginPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	if md5Do(password) != u.Password {
		//TODO: find an appropriate error message to retur
		thispage.Errormsg = "Wrong password."
		if err := loginPage.Execute(w, thispage); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	setEmailCookie(w, email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleMyLogout(w http.ResponseWriter, r *http.Request) {
	deleteEmailCookie(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func md5Do(s string) (md5s string) {
	h := md5.New()
	io.WriteString(h, s)
	fmt.Sprintf(md5s, "%x", h.Sum(nil))
	return
}
