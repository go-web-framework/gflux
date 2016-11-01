package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
//	"strconv"
//	"github.com/go-web-framework/gflux/mux"
	"./mux"
	"github.com/google/uuid"
)

var db *gorm.DB

var templates = template.Must(template.ParseFiles("goblog.html", "page.html"))

// table posts (
//   Post_id: varchar(30)
//   Author: varchar(30)
//   Aext: varchar(200)
// )
func main(){

	// open database
	var err error
	db, err = gorm.Open("sqlite3", "goblog.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	
	// Migrate the schema
 	db.AutoMigrate(&Post{})
 	
	testMux := mux.New()
	homeHandler := homeHandler{}
	pageHandler := pageHandler{}
	newHandler := newHandler{}
	testMux.GET("/home", nil, homeHandler)
	testMux.GET("/page/{id}", nil, pageHandler)
	testMux.POST("/page/new", nil, newHandler)
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
	params := mux.GetParams(r)
	idURL := params["id"]
	var post Post
	db.Where("Post_id = ?", idURL).First(&post)
	err := templates.ExecuteTemplate(w, "page.html", &post)
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
	text := r.FormValue("Text")
	var postList []Post
	db.Find(&postList)
	//store post
	var post = Post{Author:author, Text:text, Post_id:uuid.New().String()}
	db.Create(&post)
	
	http.Redirect(w, r, "/home", http.StatusFound)
}


//ajax
type goBlog struct{
	PostList []Post
}

type Post struct{
	gorm.Model
	Post_id 	string 	`gorm:""primary_key"`
	Author 		string `gorm:"type:varchar(20)"`
	Text 			string	`gorm:"type:varchar(200)"`
}

type handler404 struct{
}

func (t handler404) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>You've reached a custom 404!</h1>")
	return
}
