package datadump

import (
	"strconv"
	elastic "github.com/mattbaird/elastigo/lib"
)

type EveTypeIdList struct {
	list map[string]EveTypeId
}

func (item *EveTypeIdList) Keys() (keys []string) {
	keys = make([]string, len(item.list))
	for key, _ := range item.list {
		keys = append(keys, key)
	}

	return
}

type EveTypeId struct {
	TypeID         int64
	GraphicID      int32 "graphicID"
	Radius         float32 "radius"
	SoundID        int32 "soundID"
	IconId         int32 "iconID"
	SofFactionName string "sofFactionName"
	FactionID      int32 "factionID"
	Masteries map[string][]int32 "masteries"
	Traits map[string]map[string]EveTrait "traits"
}

type EveTrait struct {
	Bonus     int32 "bonus"
	BonusText string "bonusText"
	UnitID    int32 "unitID"
}

type EveTypeIdWriter interface {
	Write(EveTypeId) (err error)
}

type ElasticEveTypeIdWriter struct {}

func (writer *ElasticEveTypeIdWriter) Write(t EveTypeId) (err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	_, err = c.Index("eve", "type_id", strconv.FormatInt(t.TypeID, 10), nil, t)

	return
}
