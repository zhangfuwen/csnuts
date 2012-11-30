package csnuts

import (
	"net/http"
    "text/template"
    "strconv"
//    "io"
//	"fmt"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Tag struct {
    Count int
    IDs []int64

}
func putMsgTags(r *http.Request,id int64,tags []string) bool {
	c := appengine.NewContext(r)
    for _,tagstring:=range tags {
        if tagstring!="" {
	        k:=datastore.NewKey(c,"aTag",tagstring,0,nil)
            tag:=new(Tag)
	        datastore.Get(c,k,tag)// whatever it returns, don't care
            tag.IDs=append(tag.IDs,id)
            tag.Count=len(tag.IDs)
	        if _,err:=datastore.Put(c,k,tag);err!=nil {
                c.Errorf("%v",err)
                return false
            }
         }
    }
    return true
}

func getIDsByTag(r *http.Request,tagstring string) []int64 {
	c := appengine.NewContext(r)
    k:=datastore.NewKey(c,"aTag",tagstring,0,nil)
    tag:=new(Tag)
    if err:=datastore.Get(c,k,tag);err!=nil {
        c.Errorf("%v",err)
        return nil
    }
    return tag.IDs
}

func handleTaggedMsgs(w http.ResponseWriter, r * http.Request) {
	var pageData listPage
    pageData.SiteBase=Site
	pageData.NumMsgs.Value=getCount(w,r)
	if r.Method != "GET" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	u:=user.Current(c)
    pageData.U=u
	if u==nil {
		url,_:=user.LoginURL(c,"/")
		pageData.Loginbar="<a href=\""+url+"\">Login with google</a>"
	} else {
		url,_:=user.LogoutURL(c,"/")
		pageData.Loginbar="Welcome,"+u.String()+"(<a href=\""+url+"\">Logout</a>)"
	}
    tagstring:=r.FormValue("tag")
    pageData.Tag=tagstring
    pageData.QueryBase=Site+"/tagquery/?tag="+tagstring+"&"
    pageData.Msgs=handleTagQueryReturn(w,r,tagstring,0)
    // tagcloud
    tags:=new([]*Tag)
	q := datastore.NewQuery("aTag").Order("-Count").Limit(100)
	ks, err := q.GetAll(c, tags)
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

    for i,_:=range pageData.Msgs {
        pageData.Msgs[i].Content=[]byte(SubstrByByte(string(pageData.Msgs[i].Content),lenSummery))
    }
    pageData.Arts=Msgs2Arts(pageData.Msgs)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tagPage,err:= template.ParseFiles(templatePath+"msglist.html",templatePath+"articles.html",templatePath+"header.html",templatePath+"footer.html")
	if(err!=nil) {
		c.Errorf("%v", err)
		return
	}
	if err := tagPage.Execute(w, pageData); err != nil {
		c.Errorf("%v", err)
	}
}

func GetTaggedMsgs(w http.ResponseWriter, r * http.Request) {
	var pageData listPage
	pageData.NumMsgs.Value=getCount(w,r)
	if r.Method != "GET" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
    tagstring:=r.FormValue("tag")
    ids:=getIDsByTag(r,tagstring)
    for _,id:=range ids {
        msg:=getMessage(r,id)
        if msg!=nil {
            pageData.Msgs=append(pageData.Msgs,getMessage(r,id))
        }
	}

    for i,_:=range pageData.Msgs {
        pageData.Msgs[i].Content=[]byte(SubstrByByte(string(pageData.Msgs[i].Content),lenSummery))
    }
    pageData.Arts=Msgs2Arts(pageData.Msgs)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tagPage,err:= template.ParseFiles(templatePath+"msglist.html",templatePath+"articles.html",templatePath+"headertag.html",templatePath+"footer.html")
	if(err!=nil) {
		c.Errorf("%v", err)
		return
	}
	if err := tagPage.Execute(w, pageData); err != nil {
		c.Errorf("%v", err)
	}
}
func handleTagQueryReturn(w http.ResponseWriter, r *http.Request,tag string,next int) []*Message {
    tagstring:=r.FormValue("tag")
    ids:=getIDsByTag(r,tagstring)
//	c.Errorf(r.URL.Path)
	var msgs []*Message
    for cnt,id:=range ids[next:] {
        if cnt>=10 {
            break
        }
        msg:=getMessage(r,id)
        if msg!=nil {
            msgs=append(msgs,getMessage(r,id))
        }
	}
    return msgs
}

func handleTagQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
    tagstring:=r.FormValue("tag")
    ids:=getIDsByTag(r,tagstring)
//	c.Errorf(r.URL.Path)
	next,err:=strconv.Atoi(r.FormValue("next"))
	var msgs []*Message
    for cnt,id:=range ids[next:] {
        if cnt>=10 {
            break
        }
        msg:=getMessage(r,id)
        if msg!=nil {
            msgs=append(msgs,msg)
        }
	}
	if len(msgs)<=0 {
		serve404(w)
		return
	}
/*	for _,m :=range msgs {
		m.Content=strings.Replace(m.Content,"\n","<br>",-1)
	}*/
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
