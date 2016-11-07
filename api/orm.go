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
// TODO: Right now it only works with strings
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

func (o *Orm) Find(t reflect.Type, tableName string, id string) interface{} {
	// Panics
	if t.Kind().String() != "struct" {
		panic("Table find : schema type must be a struct")
	}
	if t.NumField() < 1 {
		panic("Table find : schema type empty")
	}

	// Initialize query structure
	var fields []interface{}

	var idd string
	fields = append(fields, &idd)
	for i := 1; i < t.NumField()+1; i++ {
		var a string
		fields = append(fields, &a)
	}

	// Query
	err := o.db.QueryRow("SELECT * FROM "+tableName+" WHERE Id=?", id).Scan(fields...)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	// Create output
	retVal := reflect.New(t).Elem()
	for i := 0; i < t.NumField(); i++ {
		v1 := retVal.Field(i)
		v2 := reflect.ValueOf(fields[i+1])

		v1.SetString(v2.Elem().Interface().(string))
	}

	ret := retVal.Interface()
	return ret
}
