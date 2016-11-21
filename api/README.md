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
resPosts.SetItemHandler("GET", func(
    w http.ResponseWriter,
    r *http.Request,
    DBObject interface{}){
    // OVERRIDE FUNCTION HERE
}))
```
ItemHandlers refers to the handlers at ```/posts/{id}``` for requests. Here, ```DBObject``` is the item with the specified id if found. Otherwise it is ```nil```.

CollectionHandlers:
```go
resPosts.SetCollectionHandler("GET", func(
    w http.ResponseWriter,
    r *http.Request,
    DBObjects []interface{}){
    // OVERRIDE FUNCTION HERE
}))
```
CollectionHandlers refers to handlers at ```/posts``` for requests. Here, ```DBObjects``` is the complete list of objects in the database if any exist.

## Default Handlers

As mentioned before, only GET, POST and DELETE are supported. The user must add Handlers to handle other request methods or allow for accepts other than ```application/json```.

The following shows the default behavior for CollectionHandlers and ItemHandlers:

### ItemHandlers

* GET: Returns the object as json with response 200 StatusOK. If requested ID does not exist, returns error object as json with response 404 StatusNotFound.
* POST: Program returns 404 StatusNotFound as POST ItemHandler is undefined.
* DELETE: Returns the deleted object as json with response 200 StatusOK. If requested ID does not exist, returns error object as json with response 404 StatusNotFound.

### CollectionHandlers

* GET: Returns an array of all objects in the database as json with response 200 StatusOK. If database is empty, returns error object as json with response 404 StatusNotFound.
* POST: Returns the posted object as json with 201 StatusCreated. If request was unprocessable, returns error object as json with response 422 StatusUnprocessableEntity.
* DELETE: Program returns 404 StatusNotFound as DELETE CollectionHandler is undefined.

## Allowing Methods

Methods can be allowed and disallowed. For example, if the user did not want their API resource `resPosts` to allow DELETE methods, the following line could be called:

```go
resPosts.Disallow("DELETE")
```
Now whenever a DELETE request is made on ```/posts``` or ```/posts/{id}```, 404 StatusNotFound will be returned.

GET, POST and DELETE are allowed by default. If the user wanted to allow PATCH and PUT requests, for example, the following line would help:

```go
resPosts.Allow("PATCH", "PUT")
```

Of course, the user would also have to implement the handlers for these requests using SetItemHandler and/or SetCollectionHandler.

Both ```Allow()``` and ```Disallow()``` accept variadic arguments.


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

# API Calls

## http

```
http GET localhost:8080/{name}
```
```
http GET localhost:8080/{name}/{id}
```
```
http DELETE localhost:8080/{name}/{id}
```
```
http POST http://localhost:8080/{name} Field1=value1 Field2=value2
```

## curl

```
curl -X GET localhost:8080/{name}
```
```
curl -X GET localhost:8080/{name}/{id}
```
```
curl -X DELETE localhost:8080/{name}/{id}
```
```
curl -H "Content-Type: application/json" -d '{"Field1":"value1","Field2":"value2"}' http://localhost:8080/{name}
```


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
* Move apimain.go to examples
* Create a logger interface so that panic() and Println() can be log.Fatal() and log.Printf()
