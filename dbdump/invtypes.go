package dbdump

import (
	"log"
	"strconv"
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	gorm "github.com/jinzhu/gorm"
	elastic "github.com/mattbaird/elastigo/lib"
)

func ReadSqlServer() {
	db, _ := gorm.Open("mssql", `server=127.0.0.1;user id=eve;password=eve;database=eve_data`)
	db.LogMode(true)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()

	// Disable table name's pluralization
	db.SingularTable(true)

	// Hiding dev code 〜(￣△￣〜)
	if false {
		firstType := EveInvType{}
		query := db.Find(&firstType, 19)
		firstType.ElasticWrite()

		firstStation := EveStation{}
		query = db.Find(&firstStation, 60000004)
		log.Println(firstStation.String())
		log.Println(query.Error)
	}

	// Import inv types
	var types []EveInvType
	db.Find(&types)
	for _, t := range types {
		t.ElasticWrite()
	}

	// Import stations
	var stations []EveStation
	db.Find(&stations)
	for _, s := range stations {
		s.ElasticWrite()
	}
}

func (e EveStation) TableName() string {
	return "eve_data.dbo.staStations"
}

type EveStation struct {
	Id                         int64 `gorm:"column:stationID;primary_key:yes"`
	Security                   int8 `gorm:"column:security"`
	DockingCostPerVolume       float32 `gorm:"column:dockingCostPerVolume"`
	MaxShipVolumeDockable      float32 `gorm:"column:maxShipVolumeDockable"`
	OfficeRentalCost           int32 `gorm:"column:officeRentalCost"`
	OperationId                int64 `gorm:"column:operationID"`
	StationTypeId              int64 `gorm:"column:stationTypeID"`
	CorporationId              int64 `gorm:"column:corporationID"`
	SolarSystemId              int64 `gorm:"column:solarSystemID"`
	ConstellationId            int64 `gorm:"column:constellationID"`
	RegionId                   int64 `gorm:"column:regionID"`
	Name                       string `gorm:"column:stationName"`
	X                          float32 `gorm:"column:x"`
	Y                          float32 `gorm:"column:y"`
	Z                          float32 `gorm:"column:z"`
	ReprocessingEfficiency     float32 `gorm:"column:reprocessingEfficiency"`
	ReprocessingStationsTake   float32 `gorm:"column:reprocessingStationsTake"`
	ReprocessingHangarFlag     int8 `gorm:"column:reprocessingHangarFlag"`
}

func (station *EveStation) String() string {
	if false {
		data, _ := json.MarshalIndent(station, "", "    ");
		return string(data)
	}

	return station.Name
}

func (station *EveStation) ElasticWrite() {
	c := elastic.NewConn()
	c.Hosts = []string{"192.168.33.48:9200"}

	_, err := c.Index("eve", "station", strconv.FormatInt(int64(station.Id), 10), nil, station)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(station.String())
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

	_, err := c.Index("eve", "invtype", strconv.FormatInt(int64(invType.Id), 10), nil, invType)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(invType.String())
}
