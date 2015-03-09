package no_orm

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"database/sql"
	"log"
	"encoding/json"
	"os"
)

func Transfer() {
	// connection string
	connString := `server=localhost;user id=eve;password=eve;database=eve_data`

	// create connection
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query(`
	SELECT
		typeID,
      	groupID,
      	typeName,
      	description,
      	mass,
      	volume,
      	capacity,
      	portionSize,
      	raceID,
      	basePrice,
      	published,
      	marketGroupID,
      	chanceOfDuplicating
	FROM eve_data.dbo.invTypes
	`)


	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer rows.Close()

	// fmt.Printf("%T", rows)

	invTypeList := make([]EveInvType, 0)
	for rows.Next() {
		var inv EveInvType;
		rows.Scan(
			&inv.TypeID,
			&inv.GroupID,
			&inv.TypeName,
			&inv.Description,
			&inv.Mass,
			&inv.Volume,
			&inv.Capacity,
			&inv.PortionSize,
			&inv.RaceID,
			&inv.BasePrice,
			&inv.Published,
			&inv.MarketGroupID,
			&inv.ChanceOfDuplicating,
		)
		invTypeList = append(invTypeList, inv)
	}

	file, _ := os.Create("Invtypes.json")
	defer file.Close()

	data, _ := json.MarshalIndent(invTypeList, "", "    ");
	file.WriteString(string(data))

	fmt.Printf("Exported %d records.", len(invTypeList))
}

type EveInvType struct {
	TypeID int32
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
