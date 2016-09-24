package mux
/*
import (
	"errors"
	"net/http"
	"golang.org/x/net/context"
	"strings"
)


var Stop = errors.New("exit request handling")
var urlVar = "^r"

// Middleware represents an HTTP middlware function.
type Middleware func(w http.ResponseWriter, r *http.Request) error

type tree interface {
	Get(p string) (*Entry, bool)
	Set(p string, m *Entry)
}

type treeImpl map[string]*Entry

func (t treeImpl) Get(p string) (e *Entry, ok bool) {
	e, ok = t[p]
	return
}

func (t treeImpl) Set(p string, m *Entry) {
	t[p] = m
}

type Entry struct {
	Path       string
	Middleware []Middleware
	Handler    http.Handler
	Methods    []string // Allowed HTTP methods.
}

func (e *Entry) Allow(methods ...string) *Entry {
	e.Methods = append(e.Methods, methods...)
	return e
}

// Mux is a serve mux.
type Mux struct {
	tree Tree
}

func NewMux() *Mux {
	return &Mux{
		tree: New(),
	}
}

//The original path with vars is saved in the struct.
//The path vars are replaced with a common "^r" string and are put into the tree
//This method seems like a hack; maybe find a better one
func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Entry {
	//clean url
	e := &Entry{
		Path:       path,
		Middleware: mw,
		Handler:    h,
		Methods:    nil,
	}
	//m.tree.Set(path, e)
	pathArray := strings.Split(path, "/")
	insertPath := ""
	if (len(pathArray) == 0){
		insertPath = "/"
	}
	for _, pathVal := range pathArray{
		strings.Join(insertPath, "/")
		if (pathVal[0] == ":"){		
			strings.Join(insertPath, urlVar)
		}
		else{
			strings.Join(insertpath, pathVal)
		}
	}
	m.tree.Insert(insertPath, e)
	return e
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	getPath, varValues := serveHandleURL(r.URL.Path)
	e, _ := m.tree.Get(getPath);
	//add to request
	ctx := r.Context()
	ctx = context.WithValue(ctx, 0, RequestVars{vars: varValues})
	r = r.WithContext(ctx)
	//iterator best way?
	for _, middleware := range e.Middleware{
		err := middleware(w, r)
		if (err != nil){ //err == Stop
  		w.WriteHeader(403)
			return
		}
	}
	//method check?
	e.Handler.ServeHTTP(w, r)
}

func serveHandleURL(path string) (getPath string, varValues map[string]string){
	pathArray := strings.Split(path, "/")
	getPath := ""
	for index, pathVal := range pathArray{
		strings.Join(getPath, "/")
		tempGetPath := getPath
		strings.Join(tempGetPath, pathVal)
		_, gotEntry := m.tree.Get(tempGetPath)
		if (gotEntry == nil){
			tempGetPathVar := getPath
			strings.Join(tempGetPathVar, urlVar)
			tempEntry, gotEntry := m.tree.Get(tempGetPathVar)
			if (gotEntry == nil){
				//404
				return "/404", _
			} else{
				//variable in url
				varArray := strings.Split(tempEntry.path, "/")
				tempVarURL := varArray[len(varArray - 1)]
				//take off ':'
				tempVarURL = tempVarURL[1:len(tempVarURL)]
				varValues[tempVarURL] = pathVal
				strings.Join(getPath, urlVar)
			}
		} else{
			strings.Join(getPath, pathVal)
		}
	}
	return
}

type RequestVars struct{
	vars map[string]string
}

*/
