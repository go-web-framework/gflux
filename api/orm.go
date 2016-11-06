package api

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
)

type Orm struct {
	db *sql.DB
}

func InitDB(filepath string) *Orm {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		fmt.Println("error opening database : ")
		panic(err)
	}
	if db == nil {
		panic(filepath + " not found")
	}
	return &Orm{db}
}

// Close database
func (o *Orm) Close() {
	o.db.Close()
}

// Create a table if it does not exist
// Table has the same schema as input struct plus an Id field
// Panics if input type is not a struct with at least one field
func (o *Orm) CreateTable(name string, schemaType interface{}) *Orm {
	t := reflect.TypeOf(schemaType)

	// If ptr, dereference schema type
	if t.Kind().String() == "ptr" {
		t = t.Elem()
	}

	// Panics
	if t.Kind().String() != "struct" {
		panic("Table creation : schema type must be a struct")
	}
	if t.NumField() < 1 {
		panic("Table creation : schema type empty")
	}

	// Create table
	table := "CREATE TABLE IF NOT EXISTS " + name + "(Id TEXT NOT NULL PRIMARY KEY, "
	for i := 0; i < t.NumField()-1; i++ {
		table = table + t.Field(i).Name + " TEXT, "
	}
	table = table + t.Field(t.NumField()-1).Name + " TEXT);"

	_, err := o.db.Exec(table)
	if err != nil {
		fmt.Println("Table creation error : ")
		panic(err)
	}

	return o
}
