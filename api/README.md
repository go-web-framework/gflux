gflux/api
===

# Usage
See apimain.go currently in gflux/ (will later move to examples)

# Drawbacks

## User Drawbacks
* Only supports strings right now.
* Api IDs are abstracted away from the user.
* Name of variables will be the same name in the database as well as used for requests. In go this could be a minor drawback as the variable names will have to be capitalized.
* Only support mysql and sqlite3 right now (Not REALLY a drawback).

## Code Drawbacks
* Using our own orm and database management
* Only works for strings at the moment
* Case-sensitive for some things.
* Order of struct variables MIGHT matter to the code

# TODOs
* POST handlers
* A lot of things
* Basically only GET is done
* Change imports to github.com/go-web-framework
* If the name gflux changes, this might affect some panic statements
* Move apimain.go to examples
* Should panic() and Println() be log.Fatal() and log.Printf()?