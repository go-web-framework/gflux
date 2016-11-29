package main

import (
   "github.com/go-web-framework/gflux/api"
)

type Post struct {
   Title string
   Author string
}

func main() {
   a := api.New("sqlite3", "test.db")
   defer a.Close()

   a.NewResource("posts", &Post{})

   a.Serve()
}
