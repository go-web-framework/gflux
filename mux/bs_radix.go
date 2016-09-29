package mux

import (
	"sort"
	"strings"
	"net/http"
	//"fmt"
	"errors"
)

type Route struct {
	Path       string
	Middleware []Middleware
	Handler    http.Handler
	Methods    []string // Allowed HTTP methods.
}

// edge is used to represent an edge node
type edge struct {
	label byte
	node  *node
}

type node struct {
	// Store the value ascect of the KV pair
	val *Route

	// prefix is the common prefix we ignore
	prefix string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges edges
}

func (n *node) hasValue() bool {
	return n.val != nil
}

func (n *node) addEdge(e edge) {
	n.edges = append(n.edges, e)
	n.edges.Sort()
}

func (n *node) replaceEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})
	if idx < num && n.edges[idx].label == e.label {
		n.edges[idx].node = e.node
		return
	}
	panic("replacing missing edge")
}

func (n *node) getEdge(label byte) *node {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		return n.edges[idx].node
	}
	return nil
}

func (n *node) delEdge(label byte) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		copy(n.edges[idx:], n.edges[idx+1:])
		n.edges[len(n.edges)-1] = edge{}
		n.edges = n.edges[:len(n.edges)-1]
	}
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e edges) Sort() {
	sort.Sort(e)
}

// Tree implements a radix tree. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and
// ordered iteration,
type Trie struct {
	root *node
	size int
	last_level *node;
}

// New returns an empty Tree
func NewTrie() *Trie {
	t := &Trie{root: &node{}}
	return t
}

// Len is used to return the number of elements in the tree
func (t *Trie) Len() int {
	return t.size
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

// NewRoute return a pointer to a Route instance and call save() on it
func (t *Trie) NewRoute(url string, h http.Handler, mid []Middleware, methods []string) *Route {
	  r := &Route{
		Path: url, 
		Handler: h,
		Middleware: mid,
		Methods: methods,
	}
	t.insert(r)
	return r
}

// NewRoute return a pointer to a Route instance and call save() on it
func (t *Trie) UpdateRouteMethods(path string, method ...string) bool {
	val, found := t.Get(path)
	if !found || val.Path != path {
		return false
	}

	methods := []string{}
	val.Methods = append(methods, method...)
	return true
}

// Insert is used to add a newentry or update
// an existing entry. Returns true if update successful.
func (t *Trie) insert(r *Route) error {
	toContinue := false
	n := t.root
	parent := t.root
	query := r.Path
	var e edge
	var child *node

	//make sure the query doesn't start with forward slash, but does end with one
	if len(query) != 0 && query[len(query) - 1] != '/' {
		query = r.Path + "/"
	}	
	query = strings.TrimPrefix(query, "/")
	
	for {
		// Handle key exhaution
		// what is our desired behavior if the length of query == 0? update that node?
		if len(query) == 0 {			
			return errors.New("not yet sure what to do here")
		}

		child = parent.getEdge(query[0])  //getEdge returns the child node
		indexFirstSlash := strings.IndexByte(query, '/')
		if indexFirstSlash < 0 {
			return errors.New("managed to get a query with no forward slash")
		}
		
		// No edge, create one
		if child == nil {
			if indexFirstSlash + 1 == len(query) { //we can just add the whole thing
				e = edge{
					label: query[0],
					node: &node{
						val: r,
						prefix: query, 
					},
				}
				toContinue = false //the entire query was added to the new node, so we're done here
		
			} else {
				//no edge, but forward slash in query (not at the end)
				e = edge{
					label: query[0],
					node: &node{ //leave the value nil
						prefix: query[:indexFirstSlash + 1] , 
					},
				}
				toContinue = true
			}
			
			//create a new edge pointing to a new node
			//add that edge to the parents edge array
			//increment tree size
			//there was no edge in the parent, so the edge we're creating gets added to the parent
			parent.addEdge(e)
			t.size++

			if toContinue {
				query = query[indexFirstSlash + 1:]
				parent = e.node
				continue
			} else {
				break
			}		
		}
		
		// We found an edge where the label matches query[0]
		// Determine longest prefix of the search key on match
		// look at the child and its prefix
		// the longest node we want to create has a prefix to the forward slash
		// there has to be a common Prefix, so what do we do?
		// if commonPrefix < firstSlash, we have to split at commonPrefix and continue with updated query string
		// if firstSlash < commonPrefix, we have to split at firstSlash and continue with updated query string
		// if firstSlash == commonPrefix AND common prefix < len(query), split at commonPrefix and continue with updated query string
		// if firstSlash == commonPrefix AND common prefix == len(query), can just create a new node with the whole query and be done
		// to split and continue, need to create two new edges/nodes, update the prefix of the current node, and add the edges to the current node (child)
		// nodes should only have a value if they're at the end
		//commonPrefix should never be greater than indexFirstSlash, as all nodes end on a slash if they have one
		//commonPrefix is not the index; last common character is commonPrefix - 1
				
		n = child;
		commonPrefix := longestPrefix(query, n.prefix)
		lastCommonIndex := commonPrefix - 1
		
		switch {
		case lastCommonIndex < indexFirstSlash && commonPrefix > 0: 
			//split n's prefix at commonPrefix, making everything not common a child node. Then start insert again at n. there will no longer be a commonality
			child1_prefix := n.prefix[lastCommonIndex :]
			edge1 := edge{
					label: child1_prefix[0],
					node: &node{
						prefix: child1_prefix,
						val: n.val,  //the child needs to take the parents Route
					},
				}
			//update the prefix and val of n
			n.val = nil
			n.prefix = n.prefix[: commonPrefix ]
			n.addEdge(edge1)
			t.size++
			parent = n
			query = query[commonPrefix  :]
			continue
		case lastCommonIndex == indexFirstSlash && commonPrefix < len(n.prefix):
			//change query to everything after the commonality, and start search from n
			parent = n
			query = query[commonPrefix :]
			continue
		case lastCommonIndex == indexFirstSlash  && commonPrefix == len(n.prefix):
			parent = n
			query = query[commonPrefix :]
			continue
		}	
	}
	return errors.New("Unable to insert")
}




// Delete is used to delete a key, returning any errors
func (t *Trie) Delete(s string) error {
	var parent *node
	var label byte
	n := t.root
	search := s
	for {
		// Check for key exhaution
		if len(search) == 0 {
			if !n.hasValue() {
				break
			}
			goto DELETE
		}

		// Look for an edge
		parent = n
		label = search[0]
		n = n.getEdge(label)
		if n == nil {
			break
		}

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return errors.New("Could not delete. Key not found")

DELETE:
	// Delete the value
	n.val = nil
	t.size--

	// Check if we should delete this node from the parent
	if parent != nil && len(n.edges) == 0 {
		parent.delEdge(label)
	}

	// Check if we should merge this node
	if n != t.root && len(n.edges) == 1 {
		n.mergeChild()
	}

	// Check if we should merge the parent's other child
	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.hasValue() {
		parent.mergeChild()
	}

	return nil
}

func (n *node) mergeChild() {
	e := n.edges[0]
	child := e.node
	n.prefix = n.prefix + child.prefix
	n.val = child.val
	n.edges = child.edges
}

func (t *Trie) Get(s string) (*Route, bool) {
	//if s does not end with a forward slash, append '/' to s
	if len(s) != 0 && s[len(s) - 1] != '/' {
		s = s + "/"
	}
	s = strings.TrimPrefix(s, "/") 
	
	n, found, remains := t.getLiteral(s)
	if found {
		return n.val, true
	}

	val,found2 := t.getWildCard(remains, n)
	if found2 {
		return val, true
	} 

	return nil, false
}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Trie) getLiteral(s string) (*node, bool, string) {
	n := t.root
	var child *node
	search := s

	for {
		t.last_level = n;
		// Check for key exhaution
		if len(search) == 0 {
			if n.hasValue() {
				return n, true, ""
			}
			break
		}

		// Look for an edge
		child = n.getEdge(search[0])
		if child == nil {
			return n, false, search
		}
		n = child
		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
				search = search[len(n.prefix):]
		} else {
			break
		}
	}

	return n, false, search
}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Trie) getWildCard(s string, n *node) (*Route, bool) {
	if n == nil {
		return nil, false
	}

	search := strings.TrimPrefix(s, n.prefix)
	if len(search) == 0 {
		if n.hasValue() {
			return n.val, true
		} else {
			return nil, false
		}
	}

	index := strings.LastIndexByte(search, '/')
	if index == -1 {
		search += "/"
		index = len(search) - 1;
	}

	search =  "*" + search[index:]
	if len(search) == 0  && n.hasValue() {
			return n.val, true
	}

	// Look for an edge
	n = n.getEdge(search[0])
	if n == nil {
			return nil, false
	}
	
	// Consume the search prefix
	if strings.HasPrefix(search, n.prefix) {
		search = search[len(n.prefix):]
		return t.getWildCard(search, n)
	} else {
			return nil, false
	}
}

