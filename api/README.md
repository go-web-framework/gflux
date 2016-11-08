gflux/api
===

# Usage
See apimain.go currently in gflux/ (will later move to examples)

```go
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
```

## Resources

```go
resPosts := a.NewResource("posts", &Post{})
```
Resources allow GET, POST and DELETE requests by default. Creating a new resource creates a new table in the specified database and IDs are controlled by gflux/api. In the above example, Post is a user-defined type. A pointer to an empty Post is passed to the function as a schema-type for the function.

The user can override (or add new) handlers for methods. There are two sets of handlers which are explained in context to the above example of ```resPosts```:

ItemHandlers:
```go
resPosts.Handlers["GET"] = func(
    DBObject interface{},
    w http.responseWriter,
    accepts []string){
    // OVERRIDE FUNCTION HERE
})
```
ItemHandlers refers to the handlers at ```/posts/{id}``` for requests. Here, ```DBObject``` is the item with the specified id if found. Otherwise it is ```nil```.

CollectionHandlers:
```go
resPosts.Handlers["GET"] = func(
    DBObjects []interface{},
    w http.responseWriter,
    accepts []string){
    // OVERRIDE FUNCTION HERE
})
```
CollectionHandlers refers to handlers at ```/posts``` for requests. Here, ```DBObjects``` is the complete list of objects in the database if any exist.

## Default Handlers

As mentioned before, only GET, POST and DELETE are supported. The user must add Handlers to handle other request methods or allow for accepts other than ```application/json```.

TODO: Table showing the responses gflux/api provides by default for CollectionHandlers and ItemHandlers.

## Database Support

gflux/api supports sqlite3 databases as well as mysql databases.

sqlite3:
```go
a := api.New("sqlite3", "test.db")
defer a.Close()
```

mysql:
```go
a := api.New("mysql", "user:password@/test")
defer a.Close()
```

We are currently using the drivers at [github.com/mattn/go-sqlite3]([github.com/mattn/go-sqlite3]) for sqlite3 and [github.com/go-sql-driver/mysql](github.com/go-sql-driver/mysql) for mysql.

# Drawbacks

## User Drawbacks
* The user cannot already have a table with the same name as a resource in the database provided to gflux/api.
* Only supports strings right now.
* Api IDs are abstracted away from the user.
* Name of variables will be the same name in the database as well as used for requests. In go this could be a minor drawback as the variable names will have to be capitalized.
* Only support mysql and sqlite3 right now (Not REALLY a drawback).
* Only supports application/json right now.
* Only supports POST, GET and DELETE by default.

## Code Drawbacks
* Using our own orm and database management
* Only works for strings at the moment
* Case-sensitive for some things
* Order of struct variables MIGHT matter to the code
* Current implementation still requires all methods to be processed in ServeHTTP() in resource.go (The user would not be able to handle DELETE from scratch in their override handlers since gflux/api takes care of the database queries before calling the override handlers).

# TODOs
* POST handlers
* Allow()
* Change imports to github.com/go-web-framework
* If the name gflux changes, this might affect some panic statements
* Move apimain.go to examples
* Should panic() and Println() be log.Fatal() and log.Printf()?
