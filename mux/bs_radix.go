package mux

import (
	"sort"
	"strings"
	"net/http"
)

type Route struct {
	Path       string
	Middleware []Middleware
	Handler    http.Handler
	Methods    []string // Allowed HTTP methods.
}

// leafNode is used to represent a value
type leafNode struct {
	key string
	val *Route
}

// edge is used to represent an edge node
type edge struct {
	label byte
	node  *node
}

type node struct {
	// leaf is used to store possible leaf
	leaf *leafNode

	// prefix is the common prefix we ignore
	prefix string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges edges
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
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

	//var tester byte= "{:Id}/"
	/*
	if idx < num && n.edges[idx].label == "{:Id}/" {
		return n.edges[idx].node
	}
	*/
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
// an existing entry. Returns if updated.
func (t *Trie) insert(r *Route) bool {
	var parent *node
	n := t.root
	search := r.Path

	//if search query does not end with a forward slash, append '/' to s
	if len(search) != 0 && search[len(search) - 1] != '/' {
		search = r.Path + "/"
	}
	
	for {
		// Handle key exhaution
		if len(search) == 0 {
			if n.isLeaf() {
				n.leaf.val = r
				return true
			}

			n.leaf = &leafNode{
				key: r.Path,
				val: r,
			}
			t.size++
			return false
		}

		// Look for the edge
		parent = n
		
		n = n.getEdge(search[0])

		// No edge, create one
		if n == nil {
			e := edge{
				label: search[0],
				node: &node{
					leaf: &leafNode{
						key: r.Path,
						val: r,
					},
					prefix: search,
					
				},
			}
			parent.addEdge(e)
			t.size++
			return false
		}

		// Determine longest prefix of the search key on match
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		t.size++
		child := &node{
			prefix: search[:commonPrefix],

		}
		parent.replaceEdge(edge{
			label: search[0],
			node:  child,
		})

		// Restore the existing node
		child.addEdge(edge{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]
		
		// Create a new leaf node
		leaf := &leafNode{
			key: r.Path,
			val: r,
		}

		// If the new key is a subset, add to to this node
		search = search[commonPrefix:]
		if len(search) == 0 {
			child.leaf = leaf
			return false
		}

		// Create a new edge for the node
		child.addEdge(edge{
			label: search[0],
			node: &node{
				leaf:   leaf,
				prefix: search,
				
			},
			
		})
		return false
	}
}


// Delete is used to delete a key, returning the previous
// value and if it was deleted
func (t *Trie) Delete(s string) (bool) {
	var parent *node
	var label byte
	n := t.root
	search := s
	for {
		// Check for key exhaution
		if len(search) == 0 {
			if !n.isLeaf() {
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
	return false

DELETE:
	// Delete the leaf
	n.leaf = nil
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
	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
		parent.mergeChild()
	}

	return true
}

func (n *node) mergeChild() {
	e := n.edges[0]
	child := e.node
	n.prefix = n.prefix + child.prefix
	n.leaf = child.leaf
	n.edges = child.edges
}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Trie) Get(s string) (*Route, bool) {
	n := t.root
	search := s
	//if s does not end with a forward slash, append '/' to s
	if len(s) != 0 && search[len(s) - 1] != '/' {
		search = s + "/"
	}

	for {
		// Check for key exhaution
		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.val, true
			}
			break
		}

		// Look for an edge
		n = n.getEdge(search[0])
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
	//val, found := t.GetWildCard(s)
	//if found {
	//	return val, true
	//} else {
		return nil, false
	//}
	
}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Trie) GetWildCard(s string) (*Route, bool) {
	n := t.root
	search := strings.TrimSuffix(s, "/")
	index := strings.LastIndexByte(search, '/')
	search = search[:index]
	search += "/{:Id}/"

	for {
		// Check for key exhaution
		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.val, true
			}
			break
		}

		// Look for an edge
		n = n.getEdge(search[0])
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

	return nil, false
}




