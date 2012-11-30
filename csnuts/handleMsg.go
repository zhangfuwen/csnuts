package csnuts

import (
	"net/http"
	"text/template"
	"strconv"
    "io"
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func handleMsg(w http.ResponseWriter, r *http.Request) {
	var pageData articlePage
    pageData.SiteBase=Site
    pageData.CurrPageBase=Site+"/msg/"
	if r.Method != "GET" || r.URL.Path != "/" {
//		serve404(w)
//		return
	}

	c := appengine.NewContext(r)
    if !user.IsAdmin(c) {
        pageData.IsAdmin=false
    }else {
        pageData.IsAdmin=true
    }

	u:=user.Current(c)
    pageData.U=u
	if u==nil {
		url,_:=user.LoginURL(c,"/")
		pageData.Loginbar="<a href=\""+url+"\">Login with google</a>"
	} else {
		url,_:=user.LogoutURL(c,"/")
		pageData.Loginbar="Welcome,"+u.String()+"(<a href=\""+url+"\">Logout</a>)"
	}

	id,err:=strconv.ParseInt(r.FormValue("id"),0,64)
    if err!=nil {
        c.Errorf("%v",err)
        return
    }
/*    m:=new(Message)
	k:=datastore.NewKey(c,"Message","",id,nil)
	if err=datastore.Get(c,k,m);err!=nil {
        c.Errorf("%v",err)
    //    serve404(w)
        return
    }
    */
    pageData.Msg=getMessage(r,id)
    /*
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, m)
    */

    // display ID
    pageData.Msg.ID=id;
    // Display Comments
    for _,cmtid:=range pageData.Msg.Comments {
        pageData.Cmts=append(pageData.Cmts,getComment(r,cmtid))
    }
    // tagcloud
    tags:=new([]*Tag)
	q:= datastore.NewQuery("aTag").Order("-Count").Limit(100)
	ks,err:=q.GetAll(c,tags)
	if err != nil {
		serveError(c, w, err)
		return
	}
    for i,k:=range ks {
        tagcount:=new(TagCount)
        tagcount.TagName=k.StringID()
        tagcount.Count=(*tags)[i].Count
        pageData.TagCloud=append(pageData.TagCloud, tagcount)
    }
    //end tagcloud
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msgPage,err:= template.ParseFiles(templatePath+"msg.html")
	if(err!=nil) {
		c.Errorf("%v", err)
		return
	}
    pageData.Art=Msg2Art(pageData.Msg)
	if err := msgPage.Execute(w, pageData); err != nil {
		c.Errorf("%v", err)
	}
}

func getMessage(r *http.Request,id int64) *Message {
	c := appengine.NewContext(r)
	k:=datastore.NewKey(c,"aMessage","",id,nil)
    m:=new(Message)
	if err:=datastore.Get(c,k,m);err!=nil {
        c.Errorf("%v",err)
        return nil
    }
    m.ID=k.IntID()
    return m
}

func putMessage(r *http.Request,id int64,m *Message) bool {
	c := appengine.NewContext(r)
	k:=datastore.NewKey(c,"aMessage","",id,nil)
	if _,err:=datastore.Put(c,k,m);err!=nil {
        c.Errorf("%v",err)
        return false
    }
    return true
}

func insertMessage(r *http.Request,m *Message) *datastore.Key {
	c := appengine.NewContext(r)
	if k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "aMessage", nil), m); err != nil {
		return nil
	}else {
        return k
    }
    return nil
}

func handleMsgGood(w http.ResponseWriter, r *http.Request) {
    Msg:=new(Message)
	c := appengine.NewContext(r)

	id,err:=strconv.ParseInt(r.FormValue("id"),0,64)
    if err!=nil {
        c.Errorf("%v",err)
        return
    }

    Msg=getMessage(r,id)
    if Msg==nil {
	    w.WriteHeader(http.StatusNotFound)
    	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    	io.WriteString(w, "NotFound")
        return
    }
    Msg.Good++
    putMessage(r,id,Msg)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "OK")
    return
}

func handleMsgDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    if !user.IsAdmin(c) {
        //Redirect
        w.WriteHeader(http.StatusBadRequest)
	    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	    io.WriteString(w, "Sorry, only site admin can do this.")
        return
    }
	id,err:=strconv.ParseInt(r.FormValue("id"),0,64)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
	    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	    io.WriteString(w, "BadRequest")
        return
    }
	k:=datastore.NewKey(c,"aMessage","",id,nil)
    err=datastore.Delete(c,k)
    if err!=nil {
	    w.WriteHeader(http.StatusNotFound)
	    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	    io.WriteString(w, "NotFound")
        return
    }
    DecCount(w,r)
    w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "OK")
    return
}

func handleMsgQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
//	c.Errorf(r.URL.Path)
	next,err:=strconv.Atoi(r.FormValue("next"))
	q := datastore.NewQuery("aMessage").Order("-Date").Offset(next).Limit(10)
	var msgs []*Message
	ks, err:= q.GetAll(c, &msgs)
	if err != nil {
		serveError(c, w, err)
		return
	}
	if len(msgs)<=0 {
		serve404(w)
		return
	}
/*	for _,m :=range msgs {
		m.Content=strings.Replace(m.Content,"\n","<br>",-1)
	}*/
    for i,_:=range msgs {
        msgs[i].ID=ks[i].IntID()
        msgs[i].Content=[]byte(SubstrByByte(string(msgs[i].Content),lenSummery))
    }
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page,err:= template.ParseFiles(templatePath+"query.html",templatePath+"articles.html")
	if(err!=nil) {
		c.Errorf("%v", err)
		return
	}
	if err := page.Execute(w, Msgs2Arts(msgs)); err != nil {
		c.Errorf("%v", err)
	}
}
