package dbdump

import (
	"log"
	"strconv"
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	gorm "github.com/jinzhu/gorm"
	elastic "github.com/mattbaird/elastigo/lib"
)

func ReadInvTypes() {
	db, _ := gorm.Open("mssql", `server=127.0.0.1;user id=eve;password=eve;database=eve_data`)
	db.LogMode(true)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()

	// Disable table name's pluralization
	db.SingularTable(true)

	var users []EveInvType
	db.Find(&users)
	for _, user := range users {
		user.ElasticWrite()
	}

	if false {
		firstType := EveInvType{}
		query := db.Find(&firstType, 620)
		firstType.ElasticWrite()

		log.Println(query.Error)
		log.Println(firstType)
	}
}

func (e EveInvType) TableName() string {
	return "eve_data.dbo.invTypes"
}

type EveInvType struct {
	Id                  int32 `gorm:"column:typeID;primary_key:yes"`
	GroupId             int64 `gorm:"column:groupID"`
	Name                string `gorm:"column:typeName"`
	Description         string `gorm:"column:description"`
	Mass                float32
	Volume              float32
	Capacity            float32
	PortionSize         int32 `gorm:"column:portionSize"`
	RaceId              int32 `gorm:"column:raceID"`
	BasePrice           float32 `gorm:"column:basePrice"`
	Published           bool
	MarketGroupId       int32 `gorm:"column:marketGroupID"`
	ChanceOfDuplicating float32 `gorm:"column:chanceOfDuplicating"`
}

func (invType *EveInvType) String() string {
	if false {
		data, _ := json.MarshalIndent(invType, "", "    ");
		return string(data)
	}

	return invType.Name
}

func (invType *EveInvType) ElasticWrite() {
	c := elastic.NewConn()
	c.Hosts = []string{"192.168.33.48:9200"}

	_, err := c.Index("eve", "invtypes", strconv.FormatInt(int64(invType.Id), 10), nil, invType)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(invType.String())
}
