package csnuts

import (
	"net/http"
	"text/template"
	"time"
    "strconv"
	"appengine"
	"appengine/datastore"
	"appengine/user"

)

type Comment struct {
    ID int64
    Author string
    Content string
    Menus []string
    Date time.Time
    Comments []int64
    Good int
    Bad int
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		serveError(c, w, err)
		return
	}
	m:= &Comment{
		ID:      0,
		Author:  template.HTMLEscapeString(r.FormValue("newauthor")),
		Content: template.HTMLEscapeString(r.FormValue("content")),
		Date:    time.Now(),
		Good:    0,
		Bad:     0,
	}
    if badAuthor(m.Author)|| badComment(m.Content) {
        badRequest(w,"Input too long!")
        return
    }
    processCmtContent(m)

	//TODO: build References and Referedby list
	if u := user.Current(c); u != nil {
		m.Author = u.String()
	//TODO: hook this message under user's msglist
	}
	k,err:= datastore.Put(c, datastore.NewIncompleteKey(c, "aComment", nil), m)
	if err != nil {
		serveError(c, w, err)
		return
	}

	msgid,err:=strconv.ParseInt(r.FormValue("cmtid"),0,64)
//
    msg:=getMessage(r,msgid)
    msg.Comments=append(msg.Comments,k.IntID())
    putMessage(r,msgid,msg)
	http.Redirect(w, r, "/msg/?id="+strconv.FormatInt(msgid,10), http.StatusFound)
}

func getComment(r *http.Request,id int64) *Comment {
	c := appengine.NewContext(r)
	k:=datastore.NewKey(c,"aComment","",id,nil)
    cmt:=new(Comment)
	if err:=datastore.Get(c,k,cmt);err!=nil {
        c.Errorf("%v",err)
        return nil
    }
    return  cmt
}

func putComment(r *http.Request,id int64,cmt *Comment) bool {
	c := appengine.NewContext(r)
	k:=datastore.NewKey(c,"aComment","",id,nil)
	if _,err:=datastore.Put(c,k,cmt);err!=nil {
        c.Errorf("%v",err)
        return false
    }
    return true
}

func insertComment(r *http.Request,cmt *Comment) *datastore.Key {
	c := appengine.NewContext(r)
	if k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "aComment", nil), cmt); err != nil {
		return nil
	}else {
        return k
    }
    return nil
}

