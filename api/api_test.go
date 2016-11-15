package api

import (
	"testing"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"
)

type Post struct {
	Title  string
	Author string
}

type Ret struct{
	Id string
	Value interface{}
}

func ExampleAPI_NewResource() {
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	a.Serve()
}

func TestAPI_NewResource_NonPointer(t *testing.T) {
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", Post{}) // Post{} instead of &Post{}
}

func TestGET(t *testing.T){
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	newPost := Post{Title: "title", Author: "author"}
	//insert manually
	retPost := a.DB.Insert(newPost, "posts")
	ret := retPost[0].(struct {
		Id    string
		Value interface{}
	})
	fmt.Printf("%x\n", ret.Id)


	r, _ := http.NewRequest("GET", "/posts/" + ret.Id, nil)
	w := httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)

	var post Post
	if err := json.Unmarshal([]byte(w.Body.String()), &post); err != nil{
		t.Error("json decode failed")
	}
	if(post.Title != "title" || post.Author != "author"){
		t.Error("Expected 'title' and 'author'")
	}

	clearTestDB()
}

func TestGETCollection(t *testing.T){
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	newPost := Post{Title: "title", Author: "author"}
	newPost2 := Post{Title: "tiTle2", Author: "Author2"}
	newPost3 := Post{Title: "titLe3", Author: "auThor3"}
	//insert manually
	a.DB.Insert(newPost, "posts")
	a.DB.Insert(newPost2, "posts")
	a.DB.Insert(newPost3, "posts")


	r, _ := http.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)
	
	var posts []struct{
		Id string
		Value interface{}
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &posts); err != nil{
		t.Error("json decode failed")
	}
	post0 := posts[0].Value.(map[string]interface{})
	post1 := posts[1].Value.(map[string]interface{})
	post2 := posts[2].Value.(map[string]interface{})
	
	if(post0["Title"] != "title" || post0["Author"] != "author"){
		t.Error("Expected 'title' and 'author'")
	}
	if(post1["Title"] != "tiTle2" || post1["Author"] != "Author2"){
		t.Error("Expected 'tiTle2' and 'Author2'")
	}
	if(post2["Title"] != "titLe3" || post2["Author"] != "auThor3"){
		t.Error("Expected 'titLe3' and 'auThor3'")
	}
	clearTestDB()
}

func TestDELETE(t *testing.T){
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	newPost := Post{Title: "title", Author: "author"}
	//insert manually
	retPost := a.DB.Insert(newPost, "posts")
	ret := retPost[0].(struct {
		Id    string
		Value interface{}
	})

	newPost2 := Post{Title: "tiTle2", Author: "Author2"}
	newPost3 := Post{Title: "titLe3", Author: "auThor3"}

	a.DB.Insert(newPost2, "posts")
	a.DB.Insert(newPost3, "posts")


	r, _ := http.NewRequest("DELETE", "/posts/" + ret.Id, nil)
	w := httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)
	

	r, _ = http.NewRequest("GET", "/posts", nil)
	w = httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)
	
	var posts []struct{
		Id string
		Value interface{}
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &posts); err != nil{
		t.Error("json decode failed")
	}
	post1 := posts[0].Value.(map[string]interface{})
	post2 := posts[1].Value.(map[string]interface{})
	
	if(len(posts) != 2){
		t.Error("length should be 2")
	}
	if(post1["Title"] != "tiTle2" || post1["Author"] != "Author2"){
		t.Error("Expected 'tiTle2' and 'Author2'")
	}
	if(post2["Title"] != "titLe3" || post2["Author"] != "auThor3"){
		t.Error("Expected 'titLe3' and 'auThor3'")
	}
	clearTestDB()
}

func TestPOST(t *testing.T){
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	newPost := Post{Title: "title", Author: "author"}
	//insert manually
	a.DB.Insert(newPost, "posts")


	newPost2 := Post{Title: "tiTle2", Author: "Author2"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(newPost2)

	r, _ := http.NewRequest("POST", "/posts", b)
	w := httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)

	r, _ = http.NewRequest("GET", "/posts", nil)
	w = httptest.NewRecorder()
	a.mux.ServeHTTP(w, r)

	var posts []struct{
		Id string
		Value interface{}
	}
	if err := json.Unmarshal([]byte(w.Body.String()), &posts); err != nil{
		t.Error("json decode failed")
	}
	post0 := posts[0].Value.(map[string]interface{})
	post1 := posts[1].Value.(map[string]interface{})
	
	if(post0["Title"] != "title" || post0["Author"] != "author"){
		t.Error("Expected 'title' and 'author'")
	}
	if(post1["Title"] != "tiTle2" || post1["Author"] != "Author2"){
		t.Error("Expected 'tiTle2' and 'Author2'")
	}

	clearTestDB()

}

func clearTestDB(){
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil{
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM POSTS")
	if err != nil{
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil{
		panic(err)
	}
}
