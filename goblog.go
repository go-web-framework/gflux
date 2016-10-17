package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"github.com/go-web-framework/gflux/mux"
	//"./mux"
)

var db *gorm.DB

var templates = template.Must(template.ParseFiles("goblog.html", "page.html"))

// table posts (
//   Post_id: int (autoincrement)
//   Author: varchar(30)
//   Aext: varchar(200)
// )
func main(){

	// open database
	var err error
	db, err = gorm.Open("mysql", "goblog:password@tcp(127.0.0.1:3306)/goblog?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.CreateTable(&Post{})
	// Migrate the schema
 	db.AutoMigrate(&Post{})
	//clear
	db.Delete(Post{})
  
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
	//database call
	/*url := r.URL.Path
	urls := strings.Split(url, "/")
	idURL, err := strconv.Atoi(urls[len(urls)-1])
	if (err != nil){
		idURL = 0
	}*/
	params := mux.GetParams(r)
	idURL, err := strconv.Atoi(params["id"])
	if (err != nil){
		idURL = 0
	}
	var post Post
	db.Where("Post_id = ?", idURL).First(&post)
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
	text := r.FormValue("Text")
	var postList []Post
	db.Find(&postList)
	//store post
	var post = Post{Post_id: len(postList), Author:author, Text:text}
	db.Create(&post)
	
	http.Redirect(w, r, "/home", http.StatusFound)
}


//ajax
type goBlog struct{
	PostList []Post
}

type Post struct{
	gorm.Model
	Post_id 	int 	`gorm:""primary_key"`
	Author 		string `gorm:"type:varchar(20)"`
	Text 			string	`gorm:"type:varchar(200)"`
	
}

type handler404 struct{
}

func (t handler404) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>You've reached a custom 404!</h1>")
	return
}
