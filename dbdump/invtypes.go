package main

import (
	"fmt"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	db, _ := gorm.Open("mssql", `server=localhost;user id=eve;password=eve;database=eve_data`)
	db.LogMode(true)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()

	// Disable table name's pluralization
	db.SingularTable(true)

	firstType := EveInvType{}
	query := db.Find(&firstType, 620)

	fmt.Println(query.Error)
	fmt.Println(firstType)
}

func (e EveInvType) TableName() string {
	return "eve_data.dbo.invTypes"
}

type EveInvType struct {
	TypeID int32 `gorm:"column:typeID;primary_key:yes"`
	GroupID int32
	TypeName string
	Description string
	Mass float32
	Volume float32
	Capacity float32
	PortionSize int32
	RaceID int32
	BasePrice float32
	Published bool
	MarketGroupID int32
	ChanceOfDuplicating float32
}
