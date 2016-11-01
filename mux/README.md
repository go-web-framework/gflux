**MUX** 
<p>
The mux is responsible for handling HTTP requests on the server. Users can specify a list of functions to run upon receiving an incoming HTTP request for a particular URI path. For example, a user can specify that handler functions A, B, and C should be executed upon receiving a HTTP request for the “/profile” URI path, dependant on the request method.
</p><p>
URI specific information is stored in a Route struct as follows:

```
type Route struct {
    path       string
    middleware []Middleware
    handler    http.Handler
    methods    []string
}
```

<br>

The mux uses a radix tree for handling each Route object.  Two operations are supported: insert a Route, and retrieve a previously inserted Route. The tree will also disallow inserting multiple Route objects that have the same URI. To improve retrieval speeds and save memory, the tree consolidates its data into an efficient structure after each insert operation.  
</p><p>
The Mux allows 'wildcard' enties with open/close brackets. Ex /{id}

**Example**
```go
package main

import(
  "net/http"
  "github.com/go-web-framework/gflux/mux"
)

func main () {
  m := mux.New()
  

  // mux.Get, Post, etc
  // arguments: path, middleware, hanlder
  // middleware can be nil
	m.GET("/", nil, http.HandlerFunc(HomeHandler))
	m.GET("/page/{id}", nil, http.HandlerFunc(PageHandler))

  m.NotFound = http.HandlerFunc(NotFoundHandler)

  http.ListenAndServe(":8080", m)
}

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
  s := "Hello"
  val := []byte(s)
 rw.Write(val)
}

func PageHandler(rw http.ResponseWriter, req *http.Request) {
  s := "Page"
  val := []byte(s)
 rw.Write(val)
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
  s := "Nope, Sorry"
  val := []byte(s)
 rw.Write(val)
}

```
