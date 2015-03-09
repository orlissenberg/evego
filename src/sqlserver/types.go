package sqlserver

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
