package gorm

import (
	"github.com/go-web-framework/gorm"
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

func (g *Gorm) SetSchema() {
}

func (g *Gorm) Insert() {
}

func (g *Gorm) Get() {
}

func (g *Gorm) Update() {
}

func (g *Gorm) Delete() {
}

