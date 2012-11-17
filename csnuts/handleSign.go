package csnuts

import (
	"net/http"
	"text/template"
	"time"
    "strings"
	"appengine"
	"appengine/datastore"
	"appengine/user"

)

func handleSign(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		serveError(c, w, err)
		return
	}
	m := &Message{
		ID:      0,
		Title:   template.HTMLEscapeString(r.FormValue("title")),
		Author:  template.HTMLEscapeString(r.FormValue("author")),
		Content: []byte(template.HTMLEscapeString(r.FormValue("content"))),
		Tags: strings.Split(template.HTMLEscapeString(r.FormValue("tags")),","),
		Date:    time.Now(),
		Views:   0,
		Good:    0,
		Bad:     0,
	}
    if badTitle(m.Title) || badAuthor(m.Author) || badContent(string(m.Content)) || badTag(r.FormValue("tags")) {
        badRequest(w)
    }

    processMsgContent(m)
	//TODO: build References and Referedby list
	if u := user.Current(c); u != nil {
		m.Author = u.String()
	//TODO: hook this message under user's msglist
	}
    k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Message", nil), m)
	if err != nil {
		serveError(c, w, err)
		return
	}
    putMsgTags(r,k.IntID(),m.Tags)
	setCount(w,r)
	http.Redirect(w, r, "/", http.StatusFound)
}

