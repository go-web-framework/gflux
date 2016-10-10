package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"./mux" // TODO: change this to github.com/go-web-framework/gflux/mux to allow 'go install' and 'go get'
)

var db *sql.DB

var templates = template.Must(template.ParseFiles("goblog.html", "page.html"))

//postid: int
//author: varchar(30)
//text: varchar(200)
func main(){
	db, _ = sql.Open("mysql", "goblog:password@tcp(127.0.0.1:3306)/goblog")
	testMux := mux.New()
	homeHandler := homeHandler{}
	pageHandler := pageHandler{}
	testHandler3 := handler404{}
	testMux.Handle("/home", nil, homeHandler)
	testMux.Handle("/page/*", nil, pageHandler)
	testMux.Handle("/page/new/*", nil, pageHandler)
	testMux.AllowMethod("/page/new/*", "Post")
	testMux.SetNotFound(testHandler3)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", testMux)
	
}

type homeHandler struct{
}

func (t homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
		//make database call
		rows, err := db.Query("SELECT * from posts")
		postList := make([]Post, 10)
		for rows.Next(){
		var id int
		var author string
		var text string
		err = rows.Scan(&id, &author, &text)
		postList = append(postList, Post{ID: id, Author: author, Text: text})
	}
	 err = templates.ExecuteTemplate(w, "goblog.html", &goBlog{PostList: postList})
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
	rows, err := db.Query("SELECT * FROM posts WHERE post_id = ?", idURL)
	if (err != nil){
		fmt.Fprintf(w, "<h1>No such post</h1>")
		return
	}
	//hack, will only be one
	var id int
	var author string
	var text string
	for rows.Next(){
		err = rows.Scan(&id, &author, &text)
	}
	post := Post{Author: author, Text: text, ID: id}
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
	//latest in database more secure	
	//idText := r.FormValue("id")
	//postNum, err := strconv.Atoi(idText)
	text := r.FormValue("text")
	//store post
	stmt, err := db.Prepare("INSERT posts SET author=?,text=?")
	_, err = stmt.Exec(author, text)
	if (err != nil){
		panic(err);
	}
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
