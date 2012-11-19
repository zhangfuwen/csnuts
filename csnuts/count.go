package csnuts
import (
	"net/http"
	"appengine"
	"appengine/datastore"
)

type Count struct {
	Value int
}

func setCount(w http.ResponseWriter, r* http.Request)int {
	c:=appengine.NewContext(r)
	k:=datastore.NewKey(c,"aCount","Count",0,nil)
	count:=new(Count)
	if err:=datastore.Get(c,k,count);err!=nil {
		count.Value=1
	}else {
		count.Value++
	}
	if _,err:=datastore.Put(c,k,count);err!=nil {
		return 0// TODO: should probably return error
	}
	return count.Value
}
func DecCount(w http.ResponseWriter, r* http.Request)int {
	c:=appengine.NewContext(r)
	k:=datastore.NewKey(c,"aCount","Count",0,nil)
	count:=new(Count)
	if err:=datastore.Get(c,k,count);err!=nil {
		count.Value=0
	}else {
		count.Value--
	}
	if _,err:=datastore.Put(c,k,count);err!=nil {
		return 0// TODO: should probably return error
	}
	return count.Value
}
func getCount(w http.ResponseWriter, r* http.Request)int {
	c:=appengine.NewContext(r)
	k:=datastore.NewKey(c,"aCount","Count",0,nil)
	count:=new(Count)
	if err:=datastore.Get(c,k,count);err!=nil {
		count.Value=0
	}
	return count.Value
}

