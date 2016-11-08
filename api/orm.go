// package api allows for easy creation of REST APIs
package api

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

type Orm struct {
	DB *sql.DB
}

func InitDB(driver string, filepath string) *Orm {
	db, err := sql.Open(driver, filepath)
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
	o.DB.Close()
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
	table := "CREATE TABLE IF NOT EXISTS " + name + "(Id VARCHAR(200) NOT NULL PRIMARY KEY, "
	for i := 0; i < t.NumField()-1; i++ {
		table = table + t.Field(i).Name + " VARCHAR(200), "
	}
	table = table + t.Field(t.NumField()-1).Name + " VARCHAR(200));"

	_, err := o.DB.Exec(table)
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
	err := o.DB.QueryRow("SELECT * FROM " + tableName + " WHERE Id=?", id).Scan(fields...)
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

func (o *Orm) FindAll(t reflect.Type, tableName string) []interface{} {
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
	rows, err := o.DB.Query("SELECT * FROM " + tableName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Create output
	var ret []interface{}
	for rows.Next() {
		// Scan rows into fields
		if err := rows.Scan(fields...); err != nil {
			panic(err)
		}
		
		// Scan fields into retVal interface
		retVal := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			v1 := retVal.Field(i)
			v2 := reflect.ValueOf(fields[i+1])

			v1.SetString(v2.Elem().Interface().(string))
		}
		
		// Append retVal interface to list of returned interface objects
		ret = append(ret, retVal.Interface())
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	
	return ret
}

func (o *Orm) FindAndDelete(t reflect.Type, tableName string, id string) interface{} {
	ret := o.Find(t, tableName, id)
	
	// Don't try deleting if doesnt exist
	if ret == nil {
		return nil
	}
	
	// Create statement
	stmt, err  := o.DB.Prepare("DELETE FROM " + tableName + " WHERE Id=?")
	if err != nil {
		panic(err)
	}
	
	// Execute deletion statment
	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	
	return ret
}
