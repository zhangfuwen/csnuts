package csnuts

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"regexp"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
	"time"
	"appengine/user"
	"crypto/md5"
	"io"
)

type User struct {
	Nickname string
	Email    string
	Password   string
	Domains  []string
	Intro    string
}


func init() {
}

func setEmailCookie(w http.ResponseWriter, email string) bool {
	expiration := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "email", Value:email, Expires: expiration}
	http.SetCookie(w, &cookie)
	return true
}
func deleteEmailCookie(w http.ResponseWriter,r *http.Request) {
	/*
	expiration:=time.Now().AddDate(-1,0,0)
	cookie := http.Cookie{Name: "email", Value:"", Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 307)
	*/
	expiration := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{Name: "email", Value: "s", Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 307)
	return
}

func (u * User) isLegal() (b bool, errString string) {
	b=true
	b=isRegularEmail(u.Email)
	if b==false { 
		errString="invalid email."
		return
	}
	b=isRegularNickname(u.Nickname)
	if b==false { 
		errString="invalid nickname,4~20 chars"
		return}
	b=isRegularPassword(u.Password)
	if b==false { 
		errString="invalid password"
		return}
	b=len(u.Domains)<=10
	if b==false { 
		errString="too many domains."
		return}
	for _,domain:=range u.Domains	{
		b=isRegularDomain(domain)
		if b==false { 
			errString="invalid domain"+domain
			return}
	}
	b=isRegularIntro(u.Intro)
	if b==false {
		errString="invalid intro"
	}
	return
}

func isRegularEmail(email string) (b bool) {
	b,_=regexp.MatchString("^[_a-zA-Z0-9][.a-zA-Z0-9_-]{0,20}@[a-zA-Z0-9-]{1,40}(.[a-zA-Z0-9-]{1,20}){1,10}$",email)
	return
}

/* 用户名：数字或字母开头, 4~20个字符*/
func isRegularNickname(nickname string) (b bool) {
//	b,_=regexp.MatchString("^[a-zA-Z0-9_\x7f-\xff][a-zA-Z0-9_\x7f-\xff]{3,19}$", nickname)
//上面那一句想使用汉字做用户名，但无奈不好用，暂时注释掉
	b,_=regexp.MatchString("^[a-zA-Z0-9_][a-zA-Z0-9_]{3,19}$", nickname)
	return
}
func isRegularPassword(password string) (b bool) {
	b,_=regexp.MatchString("^[_a-zA-Z0-9]{8,16}$",password)
	return
}
func isRegularDomain(domain string) (b bool) {
	b=len(domain)<=40
	return
}
func isRegularIntro(intro string) (b bool) {
	b=len(intro)<256
	return
}

func (u * User)saveSession(session *sessions.Session, w http.ResponseWriter,r *http.Request) error {
	session.Values["nickname"] = u.Nickname
	session.Values["email"] = u.Email
	return session.Save(r, w)

}

func registerForm() string{
	return `
		<html>
		    <head>
			<style type="text/css">
			html {
				background:rgb(230,230,230);
			}
			body {
				width:400px;
				margin:20 auto;
				border-radius:5px;
				border:4px solid;
				padding:3px;
				background:white;
			}
			textarea, input {
				height: 32px; 
				border-radius:0px;
				text-decoration: none;
				background: none repeat scroll 0% 0% rgb(241, 241, 241);
				padding: 5px 16px;
				font-size: 110%;
				color: white;
				width:100%;
			}
			</style>

			</head>
			<body>
				<table>
				<form action="/register/" method="post">
				<tr><td>	<label>Nickname(*):</label></td><td><input id="nickname" name="nickname"></input></td></tr>
				<tr><td><label>Email(*):</label></td><td><input id="email" name="email"></input></td></tr>
				<tr><td><label>Password(*):</label></td><td><input id="password" name="password"></input></td></tr>
				<tr><td><label>Domains:</label></td><td><input id="domains" name="domains"></input></td></tr>
					<tr><td><label>Introduction:</label></td><td><textarea id="intro" name="intro"></textarea></td></tr>
					<tr><td><input type="submit" value="submit"  ></input></td></tr>
					</form>
					</table>
			</body>
		</html>
	`
}

func handleMyRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, registerForm())
		return
	}

	//I don't know how to handle FormValue error
	var u User
	c := appengine.NewContext(r)
	u.Nickname= r.FormValue("nickname")
	u.Password = r.FormValue("password")
	u.Email = r.FormValue("email")
	u.Domains = strings.Split(r.FormValue("domains"),",")
	u.Intro = r.FormValue("intro")

	if b,s:=u.isLegal(); b!=true {
		fmt.Fprintf(w,"Invalid User info"+s)
		return
	}
	u.Password=md5Do(u.Password)
	k := datastore.NewKey(c, "User", u.Email, 0, nil)
	_, err := datastore.Put(c, k, &u)
	if err != nil {
		fmt.Fprintf(w, "err put")
		return
	}
	setEmailCookie(w,u.Email)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func loginPage(w http.ResponseWriter, r *http.Request) string{
	c := appengine.NewContext(r)
	url,_:=user.LoginURL(c,"/")
	url="http://localhost:8080/_ah/login?continue=http://localhost:8080/"

	return `<html>
	<head>
	<style type="text/css">
			html {
				background:rgb(230,230,230);
			}
			body {
				width:400px;
				margin:20 auto;
				border-radius:5px;
				border:4px solid;
				padding:3px;
				background:white;
			}
			textarea, input {
				height: 32px; 
				border-radius:0px;
				text-decoration: none;
				background: none repeat scroll 0% 0% rgb(241, 241, 241);
				padding: 5px 16px;
				font-size: 110%;
				color: white;
				width:100%;
			}
	</style>
	</head>
<body>
    <table>
	<form action="/login" method="post">
	<tr><td>
	<label>login:</label></td><td><input id="email" name="email"></input></td></td>
	<tr><td>
	<label>password:</label></td><td><input id="password" name="password"></input></td></tr>
	<tr><td>
	<input type="submit" value="submit" ></input></td></td>
	<tr><td>
	<a href="`+url+`">Login with Google</a></td></tr>
	</form>
</body>
</html>`
}
func handleMyLogin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method != "POST" {
		//FIXME:Put is ok, what is the difference between POST and PUT?
		fmt.Fprintf(w, loginPage(w,r))
		return
	}

	u := new(User)
	email := r.FormValue("email")
	if !isRegularEmail(email) {
		fmt.Fprintf(w,"invalid email:"+email)
		return
	}
	password := r.FormValue("password")
	if !isRegularPassword(password) {
		fmt.Fprintf(w,"invalid password:"+password)
		return
	}
	k := datastore.NewKey(c, "User", email, 0, nil)
	if err := datastore.Get(c, k, u); err != nil {
		fmt.Fprintf(w,"error get user with email:"+email)
		return
	}
	if md5Do(password) != u.Password {
		//TODO: find an appropriate error message to retur
		fmt.Fprintf(w,"wrong password, original:"+u.Password+"try:"+password)
		return
	}
	setEmailCookie(w,email)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func handleMyLogout(w http.ResponseWriter, r *http.Request) {
	deleteEmailCookie(w,r)
//j	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	http.Redirect(w, r, "/", 307)
}
func md5Do(s string) (md5s string) {
	h:=md5.New()
	io.WriteString(h,s)
	fmt.Sprintf(md5s,"%x",h.Sum(nil))
	return
}
