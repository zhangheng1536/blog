package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/labstack/echo"
	"net/http"
	"html/template"
	"github.com/zhangheng1536/blog/util"
)

type Kkey struct {
	Name  string
	Value string
}

var db = new(leveldb.DB)

func main() {
	e := echo.New()
	t := &util.Template{
		Templates: template.Must(template.ParseGlob("view/*.html")),
	}

	e.Renderer = t
	e.GET("/edit", editPage)
	e.GET("/show", showPage)
	e.POST("/save", savePage)
	e.Logger.Fatal(e.Start(":1212"))
}
func editPage(con echo.Context) error {
	return con.Render(http.StatusOK, "edit", "zhheng")

}
func showPage(con echo.Context) error {
	result := make([]*Kkey, 0)
	db, err := leveldb.OpenFile("blog/test", new(opt.Options))
	if nil != err {
		fmt.Errorf("leveldb open file err:%v", err.Error())
	}
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	if iter.First(){
		kkey := new(Kkey)
		kkey.Name = string(iter.Key())
		kkey.Value = string(iter.Value())
		result = append(result, kkey)
	}
	for iter.Next() {
		kkey := new(Kkey)
		kkey.Name = string(iter.Key())
		kkey.Value = string(iter.Value())
		result = append(result, kkey)
	}

	return con.Render(http.StatusOK, "show", result)

}

func savePage(con echo.Context) error {
	b := new(Blog)
	if err := con.Bind(b); err != nil {
		return nil;
	}
	db, err := leveldb.OpenFile("blog/test", new(opt.Options))
	if nil != err {
		fmt.Errorf("leveldb open file err:%v", err.Error())
	}
	defer db.Close()
	db.Put([]byte(b.Title),[]byte( b.Text),nil)
	result := make([]*Kkey, 0)

	kkey := new(Kkey)
	kkey.Name = b.Title
	kkey.Value = b.Text
	result = append(result, kkey)
	return con.Render(http.StatusOK, "show", result)

}

// User
type Blog struct {
	Title string `json:"title" form:"title" query:"title"`
	Text  string `json:"text" form:"text" query:"text"`
}
