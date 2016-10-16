package mux

import (
	"errors"
	"sort"
	"strings"
	"net/http"
	//"fmt"
)

const (
		MethodAll 	  = "ALL"
        MethodGet     = "GET"
        MethodHead    = "HEAD"
        MethodPost    = "POST"
        MethodPut     = "PUT"
        MethodPatch   = "PATCH" // RFC 5789
        MethodDelete  = "DELETE"
        MethodConnect = "CONNECT"
        MethodOptions = "OPTIONS"
        MethodTrace   = "TRACE"
)

// edge is used to represent an edge node
type edge struct {
	label byte
	node  *node
}

type node struct {
	// Store the value ascect of the KV pair
	val *Route

	// common prefix
	prefix string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges edges

	// True if the prefix ends with forward slash
	// denotes the end of a path, where subpaths can follow
	isAnchor bool

	// the path up to this point - only needed if isAnchor is true
	path string
}

func (n *node) hasValue() bool {
	return n.val != nil
}

func (n *node) addEdge(e edge) {
	n.edges = append(n.edges, e)
	n.edges.Sort()
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
// dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and
// ordered iteration.
type Trie struct {
	root      *node
	size      int
	lastLevel *node
}

// New returns an empty Tree ready-for-use.
func NewTrie() *Trie {
	t := &Trie{
		root: &node{
			path:     "",
			isAnchor: true,
		},
	}
	return t
}

// Len returns the number of elements in the tree.
func (t *Trie) Len() int {
	return t.size
}

// longestPrefix finds the length of the shared prefix
// of two strings.
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

func isWildCardKey(s string) bool {
	return s[0] == '{' && s[len(s)-1] == '}'
}

func cleanWildCardKey(s string) string {
	s = strings.TrimPrefix(s, "{")
	return strings.TrimSuffix(s, "}")
}

func newMethodMap(r *Route) {
	r.handlers = make(map[string]http.Handler)
	r.handlers[r.method] = r.handler
}

func appendMethodMap(r *Route, n *node) bool {


	if _, ok := n.val.handlers[MethodAll]; ok {
		return false //cannot append - ALL methods are already allowed
	}

	if _, ok := n.val.handlers[r.method]; ok {
		return false //cannot append - the method aleady exists
	}

	n.val.handlers[r.method] = r.handler

	return true
}


// Insert is used to add a new entry or update
// an existing entry. Returns true if update successful.
func (t *Trie) insert(r *Route) error {
	toContinue := false
	n := t.root
	parent := t.root
	query := r.path
	var e edge
	var child *node

	// update root controller to act as home controller
	if query == "/" {
		if t.root.val == nil {
			newMethodMap(r)
			t.root.val = r
			return nil
		} else {
			return errors.New("Cannot insert; key already exists")
		}
	}

	//make sure the query doesn't start with forward slash, but does end with one
	if len(query) != 0 && query[len(query)-1] != '/' {
		query = r.path + "/"
	}
	query = strings.TrimPrefix(query, "/")

	for {
		
		if len(query) == 0{
			if appendMethodMap(r, n) {
				return nil
			} else {
				return errors.New("Could not append to existing key")
			}
		}
		
		indexFirstSlash := strings.IndexByte(query, '/')
		if indexFirstSlash < 0 {
			return errors.New("managed to get a query with no forward slash")
		}

		child = parent.getEdge(query[0]) //getEdge returns the child node

		// No edge, create one
		if child == nil {
			if indexFirstSlash+1 == len(query) { //we can just add the whole thing
				newMethodMap(r)
				e = edge{
					label: query[0],
					node: &node{
						val:      r,
						prefix:   query,
						isAnchor: true,
						path:     parent.path + query,
					},
				}
				toContinue = false //the entire query was added to the new node, so we're done here

			} else {
				//no edge, but forward slash in query (not at the end)
				e = edge{
					label: query[0],
					node: &node{ //leave the value nil
						prefix:   query[:indexFirstSlash+1],
						isAnchor: true,
						path:     parent.path + query[:indexFirstSlash+1],
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
				query = query[indexFirstSlash+1:]
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

		n = child
		commonPrefix := longestPrefix(query, n.prefix)
		lastCommonIndex := commonPrefix - 1

		if isWildCardKey(query[:indexFirstSlash]) && commonPrefix != len(n.prefix) {
			return errors.New("attempting to add multiple wildcard keys at same level")
		}

		// if commonPrefix == len(n.prefix), there's no need to split, just add
		// an edge

		if commonPrefix == len(n.prefix) {
			query = query[commonPrefix:]
			parent = n
			continue
		}
		switch {
		case lastCommonIndex < indexFirstSlash && commonPrefix > 0:
			//split n's prefix at commonPrefix, making everything not common a child node. Then start insert again at n. there will no longer be a commonality
			child1_prefix := n.prefix[commonPrefix:]

			//need to create a new edge pointing to a new node with the common prefix
			//add that edge to n
			edge1 := edge{
				label: child1_prefix[0],
				node: &node{
					prefix:   child1_prefix,
					val:      n.val, //the child needs to take the parents Route
					isAnchor: false,
					path:     n.path + child1_prefix,
					edges:    n.edges,
				},
			}
			//update the prefix and val of n
			n.val = nil
			n.prefix = n.prefix[:commonPrefix]
			n.edges = nil
			n.addEdge(edge1)
			t.size++
			parent = n
			query = query[commonPrefix:]
			continue
		case lastCommonIndex == indexFirstSlash && commonPrefix < len(n.prefix):
			//change query to everything after the commonality, and start search from n
			parent = n
			query = query[commonPrefix:]
			continue
		case lastCommonIndex == indexFirstSlash && commonPrefix == len(n.prefix):
			parent = n
			query = query[commonPrefix:]
			continue
			// TODO(?): the node we're attempting to insert already exists,
			// return an error.
		default:
			break
		}
		break
	}
	return nil
}

// assumes query is already formatted so that it does not contain a leading
// forward slash but does end with forward slash
func isValidQuery(s string) bool {
	tokens := strings.Split(s, "/")

	for _, key := range tokens {
		if key != "" && isWildCardKey(key) {
			return false
		}
	}

	return true
}

// Get returns the route, the wilcard values, and
// whether a match was found for the supplied path.
func (t *Trie) Get(s string, method string) (*Route, map[string]string, bool) {
	if s == "/" {
		return t.root.val, nil, true
	}

	// If s does not end with a forward slash, append '/' to s.
	if len(s) != 0 && s[len(s)-1] != '/' {
		s = s + "/"
	}
	s = strings.TrimPrefix(s, "/")
	if !isValidQuery(s) {
		return nil, nil, false
	}

	n, found, remains := t.getLiteral(s, method, t.root)
	if found {
		return n.val, nil, true
	}

	// added now that the path might exist but doesn't
	// have the correct method
	if remains == "" {
		return nil, nil, false
	}

	val, found2, mp := t.getWildCard(remains, method, n)
	
	if found2 {
		return val, mp, true
	}

	return nil, nil, false
}

func (n *node) hasMethodHandler(method string) bool{

	if handl, ok := n.val.handlers[MethodAll]; ok {
		n.val.method = MethodAll
		n.val.handler = handl
		return true
	}

	if handl, ok := n.val.handlers[method]; ok {
		n.val.method = method
		n.val.handler = handl
		return true
	}

	return false
}

func (t *Trie) getLiteral(s string, method string, n *node) (*node, bool, string) {
	var child *node

	search := s
	retNode := n
	remaining := search

	for {
		t.lastLevel = n
		// Check for key exhaution
		if len(search) == 0 {
			if n.hasValue() && n.hasMethodHandler(method) {
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
			if n.isAnchor {
				retNode = n
				remaining = search
			}
		} else {
			break
		}
	}

	return retNode, false, remaining

}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Trie) getWildCard(s string, method string, n *node) (*Route, bool, map[string]string) {
	if n == nil {
		return nil, false, nil
	}

	var rNode *node
	var remains string
	found := false
	search := s
	m := make(map[string]string)

	for len(search) != 0 && n != nil {
		if len(search) == 0 && n.hasValue()  && n.hasMethodHandler(method) {
			return n.val, true, m
		}

		indexFirstSlash := strings.IndexByte(search, '/')
		replacedText := search[:indexFirstSlash]

		//look to see if there's an egde from
		//the current node with a {
		//if so, ensure the pointed to node is a wildcard key (should be)
		//prepend the wildcard key to search[indexFirstSlash:]
		//add the key and replaced text to the map

		n = n.getEdge('{')
		if n == nil {
			return nil, false, nil
		}

		wcFirstSlash := strings.IndexByte(n.prefix, '/')

		if indexFirstSlash < 0 || !isWildCardKey(n.prefix[:wcFirstSlash]) {
			return nil, false, nil
		}

		wc := n.prefix[:wcFirstSlash]
		wcKey := cleanWildCardKey(wc)
		m[wcKey] = replacedText
		search = wc + search[indexFirstSlash:]

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
			rNode, found, remains = t.getLiteral(search, method, n)
			if found {
				m[wcKey] = replacedText
				return rNode.val, found, m
			} else {
				search = remains
				n = rNode
				continue
			}
		} else {
			return nil, false, nil
		}

	}

	if n == nil {
		return nil, false, nil
	}

	if n.hasValue()  && n.hasMethodHandler(method) {
		// n.val.Tokens = append(n.val.Tokens, search)
		return n.val, true, m
	} else {
		return nil, false, nil
	}
}

// walks the tree, returning true if s found
func (t *Trie) Walk(s string) bool {
	n := t.root
	var child *node

	if s == "/" {
		return true
	}

	//if s does not end with a forward slash, append '/' to s
	if len(s) != 0 && s[len(s)-1] != '/' {
		s = s + "/"
	}
	s = strings.TrimPrefix(s, "/")

	search := s

	for {
		t.lastLevel = n
		// Check for key exhaution
		if len(search) == 0 {
			if n.hasValue() {
				return true
			}
			break
		}

		// Look for an edge
		child = n.getEdge(search[0])
		if child == nil {
			return false
		}
		n = child
		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}

	return false
}
