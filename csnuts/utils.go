package csnuts

import (
	"time"
    "strings"
    "regexp"
)

type Message struct {
	ID      int64
	Title	string
	Author  string
	Content []byte
	Tags []string
    Menus []string
	Date    time.Time
	Views   int
	Comments  []int64
	Good    int
	Bad     int
}
type Article struct {
	ID      int64
	Title	string
	Author  string
	Content string
	Tags []string
    Menus []string
	Date    time.Time
	Views   int
	Comments  []int64
	Good    int
	Bad     int
}
func Msg2Art(m *Message) *Article {
	return &Article{
		ID:      m.ID,
		Title:   m.Title,
		Author:  m.Author,
		Content: string(m.Content),
		Tags:   m.Tags,
        Menus:  m.Menus,
        Date:   m.Date,
        Views: m.Views,
        Comments: m.Comments,
        Good: m.Good,
        Bad: m.Bad,
    }
}

func Msgs2Arts(ms []*Message) []*Article {
    as:=make([]*Article, len(ms))
    for i,m:=range ms {
       as[i]=Msg2Art(m)
    }
    return as
}


func processMsgContent(m *Message) {
    lines:=strings.Split(string(m.Content),"\n")
    reH3,_:=regexp.Compile("^\\s*?={3}([\\s\\S]*?)={3}\\s*?$")
    reH2,_:=regexp.Compile("^\\s*?==([\\s\\S]*?)==\\s*?$")
    reImg,_:=regexp.Compile("^\\s*?@@([\\s\\S]*?)@@\\s*?$")
    reLink,_:=regexp.Compile("^\\s*?@([\\s\\S]*?)@\\s*?$")
    //reCode,_:=regexp.Compile("^[code]$")
    //reEndCode,_:=regexp.Compile("^[/code]$")
    menuItems:=make([]string,0,10)
//    var pre bool
//    var codeStart int
//    pre=false
    for index,_:=range lines {
        //code
        /* decided to make it simple
        if reCode.MatchString(line[index]) {
            pre=true
            codeStart=index
            continue
        }
        if reEndCode.MatchString(line[index]) {
            if pre==false {
                continue
            }
            pre=false
            line[codeStart]="<pre>"
            line[index]="</pre>"
            continue
        }
        */
        // H3 
        item:=reH3.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"=","",-1))
            lines[index]="<h3 id=\""+item+"\"><em>"+item+"</em></h3>"
            menuItems=append(menuItems,"<a href=\"#"+item+"\">&nbsp;&nbsp;|--"+item+"</a>")
            continue
        }

        // H2
        item=reH2.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"=","",-1))
            menuItems=append(menuItems,"<a href=\"#"+item+"\">"+item+"</a>")
            lines[index]="<h2 id=\""+item+"\"><em>"+item+"</em></h2>"
            continue
        }
        //Image
        item=reImg.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"@","",-1))
            lines[index]="<div class=\"m\"><img src=\""+item+"\" /></div><br/>"
            continue
        }
        //Link
        item=reLink.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"@","",-1))
            lines[index]="<a href=\""+item+"\" target=\"_blank\">"+item+"</a><br/>"
            continue
        }
        // P
        lines[index]=lines[index]+"<br/>"
    }
    m.Content=[]byte(strings.Join(lines,""))
    m.Menus=menuItems
    //BBcode
	m.Content=[]byte(strings.Replace(string(m.Content),"\n","<br/>",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[code]","<pre class=\"prettyprint\">",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[/code]","</pre>",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[h2]","<h2>",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[/h2]","</h2>",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[swf]","<div class=\"m\"><embed width=\"610\" height=\"498\" type=\"application/x-shockwave-flash\" allowfullscreen=\"true\" wmode=\"transparent\" src=\"",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[flv]","<div class=\"m\"><embed width=\"610\" height=\"498\" type=\"application/x-shockwave-flash\" allowfullscreen=\"true\" wmode=\"transparent\" src=\"",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[/swf]","\"></div><br/>",-1))
	m.Content=[]byte(strings.Replace(string(m.Content),"[/flv]","\"></div><br/>",-1))
    //BBcode end 
}

func processCmtContent(m *Comment) {
    lines:=strings.Split(string(m.Content),"\n")
    reH3,_:=regexp.Compile("^\\s*?={3}([\\s\\S]*?)={3}\\s*?$")
    reH2,_:=regexp.Compile("^\\s*?==([\\s\\S]*?)==\\s*?$")
    menuItems:=make([]string,0,10)
    for index,_:=range lines {
        // H3 
        item:=reH3.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"=","",-1))
            lines[index]="<h3 id=\""+item+"\">"+item+"</h3>"
            menuItems=append(menuItems,"<a href=\"#"+item+"\">|--"+item+"</a>")
            continue
        }

        // H2
        item=reH2.FindString(lines[index])
        if item!="" {
            item=strings.TrimSpace(strings.Replace(item,"=","",-1))
            menuItems=append(menuItems,"<a href=\"#"+item+"\">"+item+"</a>")
            lines[index]="<h2 id=\""+item+"\">"+item+"</h2>"
            continue
        }
        // P
        lines[index]="<p>"+lines[index]+"</p>"
    }
    m.Content=strings.Join(lines,"")
    m.Menus=menuItems
    //BBcode
	m.Content=strings.Replace(m.Content,"\n","<br>",-1)
	m.Content=strings.Replace(m.Content,"[h2]","<h2>",-1)
	m.Content=strings.Replace(m.Content,"[/h2]","</h2>",-1)
    //BBcode end 
}

func DeScript(s  string) string {
	re,_:=regexp.Compile("\\<[\\S\\s]+?\\>")
	s=re.ReplaceAllStringFunc(s,strings.ToLower)
	//de css
	re,_=regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	s=re.ReplaceAllString(s,"")
	//de script
	re,_=regexp.Compile("\\<script[\\S\\s]+?\\>")
	s=re.ReplaceAllString(s,"")
	//de html
	re,_=regexp.Compile("\\<[\\S\\s]+?\\>")
	s=re.ReplaceAllString(s,"\n")
	//de \n\n->\n
	re,_=regexp.Compile("s{2,}")
	s=re.ReplaceAllString(s,"\n")
	return s
}

func SubstrByByte(str string, length int) string {
    if len(str)<length {
        return str
    }
    bs := []byte(str)[:length]
    bl := 0
    for i:=len(bs)-1; i>=0; i-- {
        switch {
                case bs[i] >= 0 && bs[i] <= 127:
                    return string(bs[:i+1])
                case bs[i] >= 128 && bs[i] <= 191:
                        bl++;
                case bs[i] >= 192 && bs[i] <= 253:
                        cl := 0
                    switch {
                         case bs[i] & 252 == 252:
                                cl = 6
                        case bs[i] & 248 == 248:
                             cl = 5
                        case bs[i] & 240 == 240:
                        cl = 4
                        case bs[i] & 224 == 224:
                            cl = 3
                         default:
                                cl = 2
                    }
                    if bl+1 == cl {
                    return string(bs[:i+cl])
                    }
                    return string(bs[:i])
            }
        }
    return ""
}
