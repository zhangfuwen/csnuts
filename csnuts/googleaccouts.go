package csnuts

import (
	"net/http"
	"strings"
	"text/template"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type postPage struct {
	QueryBase string
	SiteBase  string
	Loginbar  string
	TagCloud  []*TagCount
	U         *user.User
	UserName  string
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	var pageData postPage
	userName := ""
	if u != nil { //a google user
		userName = u.String()
		pageData.U = u
	} else { //not a google user
		//is it a local user?
		cookie, err := r.Cookie("email")
		if err == nil {
			userName = cookie.Value
			pageData.UserName = userName
		} else { //no logged in yet

			badRequest(w, "Only login user can post messages.")
			return
		}
	}
	if r.Method != "POST" {
		postPage, err := template.ParseFiles(templatePath + "post.html")
		if err != nil {
			c.Errorf("%v", err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//模板文件login/signup.html中没有需要套用的数据
		if err := postPage.Execute(w, pageData); err != nil {
			c.Errorf("%v", err)
		}
		return
	}
	if err := r.ParseForm(); err != nil {
		serveError(c, w, err)
		return
	}
	tagsString := r.FormValue("tags")
	m := &Message{
		ID:      0,
		Title:   template.HTMLEscapeString(r.FormValue("title")),
		Author:  template.HTMLEscapeString(r.FormValue("author")),
		Content: []byte(template.HTMLEscapeString(r.FormValue("content"))),
		Tags:    strings.Split(template.HTMLEscapeString(tagsString), ","),
		Date:    time.Now(),
		Views:   0,
		Good:    0,
		Bad:     0,
	}
	if badTitle(m.Title) || badAuthor(m.Author) || badContent(string(m.Content)) || badTag(tagsString) {
		badRequest(w, "Input too long")
		return
	}

	processMsgContent(m)
	//TODO: build References and Referedby list
	if u := user.Current(c); u != nil {
		m.Author = userName
		//TODO: hook this message under user's msglist
	}
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "aMessage", nil), m)
	if err != nil {
		serveError(c, w, err)
		return
	}
	putMsgTags(r, k.IntID(), m.Tags)
	setCount(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
