package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"github.com/go-web-framework/gflux/mux"
)

var db *gorm.DB

var templates = template.Must(template.ParseFiles("goblog.html", "page.html"))

//postid: int
//author: varchar(30)
//text: varchar(200)
func main(){
	db, _ = gorm.Open("mysql", "goblog:password@tcp(127.0.0.1:3306)/goblog")
	//defer db.Close()
	testMux := mux.New()
	homeHandler := homeHandler{}
	pageHandler := pageHandler{}
	testHandler3 := handler404{}
	testMux.Handle("/home", nil, homeHandler)
	testMux.Handle("/page/*", nil, pageHandler)
	testMux.Handle("/page/new/*", nil, pageHandler).Allow("Post")
	testMux.SetNotFound(testHandler3)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", testMux)
	
}

type homeHandler struct{
}

func (t homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//make database call
	var postList []Post
	db.Find(&postList)
	err := templates.ExecuteTemplate(w, "goblog.html", &goBlog{PostList: postList})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	return
}

type pageHandler struct{
}

func (t pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//database call
	url := r.URL.Path
	urls := strings.Split(url, "/")
	idURL, err := strconv.Atoi(urls[len(urls)-1])
	if (err != nil){
		idURL = 0
	}

	var post Post
	db.Where("ID = ?", idURL).First(&post)
	err = templates.ExecuteTemplate(w, "page.html", &post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	return
}

type newHandler struct{
}

func (t newHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	author := r.FormValue("Author")
	if (author == ""){
		author = "anon"
	}
	text := r.FormValue("text")
	
	//store post
	var post = Post{Author:author, Text:text}
	db.Create(&post)
	
	http.Redirect(w, r, "/home", http.StatusFound)
}


//ajax
type goBlog struct{
	PostList []Post
}

type Post struct{
	Author string
	Text string
	ID int
}

type handler404 struct{
}

func (t handler404) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>You've reached a custom 404!</h1>")
	return
}
