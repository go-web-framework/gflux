package gorm

import (
	"github.com/go-web-framework/gorm"
//	_ "github.com/jinzhu/gorm/dialects/mysql"
//	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
//	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// Gorm is the gorm wrapper
type Gorm struct {
	db gorm
}

func New(db_driver string, db_file string) *Gorm, err {
	g, err := gorm.Open(db_driver, db_file)
	return &Gorm{
		db: &g
	}, err
}

func (g *Gorm) Close() {
	g.db.Close()
}

func (g *Gorm) SetSchema(typePointer interface{}) {
}

func (g *Gorm) Insert() {
}

func (g *Gorm) Get() {
}

func (g *Gorm) Update() {
}

func (g *Gorm) Delete() {
}

