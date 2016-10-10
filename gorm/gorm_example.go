package main

import (
	"gflux/gorm"
)

type Product struct {
	Code string
	Price uint
}

func main() {
	// create sqlite3 database named test.db
	db, err := gorm.New("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	
	// set database schema
	db.SetSchema(&Product{})
	
	// insert item into database
	db.Insert(Product{
		Code: "L1212"
		Price: 1000
	})
	
	// get item from database
	product := db.Get("Code", "L1212") // returns a Product*
	
	// update the obtained item
	db.Update(product, "Price", 2000) // argument is a Product*
	
	// delete an item
	db.Delete(product) // argument is a Product*
}
	
	
