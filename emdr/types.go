package emdr

import (
	"reflect"
	"strings"
	"log"
	"os"
	"strconv"
	"encoding/json"
)

type EmdrMessage struct {
	ResultType string "resultType"
}

type EmdrOrderRowSet struct {
	TypeId      int64 "typeID"
	RegionID    int64 "regionID"
	GeneratedAt string "generatedAt"
	Rows        [][]interface{} "rows" // Format not document friendly.
	DataRows    []EmdrOrderDocument
}

type EmdrOrderMessage struct {
	EmdrMessage
	UploadKeys  []map[string]string "uploadKeys"
	Generator   map[string]string "generator"
	CurrentTime string "currentTime"
	Version     string "version"
	RowSets     []EmdrOrderRowSet "rowsets"
	Columns     []string "columns"
}


type EmdrOrderDocument struct {
	RegionId               int64
	RegionName             string
	TypeId                 int64
	TypeName               string
	Price                  float64
	VolRemaining           int64
	Range                  int64
	OrderId                int64
	VolEntered             int64
	MinVolume              int64
	Bid                    bool
	IssueDate              string
	Duration               int64
	StationID              int64
	StationName            string
	SolarSystemID          int64
	SolarSystemName        string
}

func (doc *EmdrOrderDocument) String() string {
	if false {
		return strconv.FormatFloat(doc.Price, 'f', 4, 64)
	}

	result, _ := json.Marshal(doc)
	return string(result)
}

func (order *EmdrOrderMessage) mapRows() {
	// Rewrite the data sets to a key:value format.
	// Loop the sets.
	for setIndex, setValue := range order.RowSets {
		order.RowSets[setIndex].DataRows = make([]EmdrOrderDocument, len(order.RowSets[setIndex].Rows))

		// Loop the rows.
		for rowIndex, rowValue := range setValue.Rows {
			// Create a new mapping.
			mapping := EmdrOrderDocument{}
			for index, name := range order.Columns {
				setOrderFieldValue(name, &mapping, rowValue[index])
			}

			mapping.RegionId = setValue.RegionID
			mapping.TypeId = setValue.TypeId

			region, _ := ReadRegion(strconv.FormatInt(mapping.RegionId, 10))
			mapping.RegionName = region.Name

			solarsystem, _ := ReadSolarSystem(strconv.FormatInt(mapping.SolarSystemID, 10))
			mapping.SolarSystemName = solarsystem.Name

			order.RowSets[setIndex].DataRows[rowIndex] = mapping
		}

		// Remove the old data
		order.RowSets[setIndex].Rows = nil
	}
}

func setOrderFieldValue(name string, item *EmdrOrderDocument, value interface{}) {
	s := reflect.ValueOf(item).Elem()
	field := s.FieldByNameFunc(func(fieldName string) bool {
		return strings.EqualFold(fieldName, name)
	})

	if field.IsValid() {
		switch v := value.(type) {
		case string:
			field.SetString(v)
		case int64:
			field.SetInt(v)
		case float64:
			switch field.Interface().(type) {
			case int64:
				field.SetInt(int64(v))
			default:
				field.SetFloat(v)
			}
		case bool:
			field.SetBool(v)
		default:
			log.Panicf("%T", v)
		}
	} else {
		log.Println(name)
		log.Printf("%T", field)
		log.Printf("%T", value)
		os.Exit(1)
	}
}

type EmdrHistoryMessage struct {
	EmdrMessage
	UploadKeys  []map[string]string "uploadKeys"
	Generator   map[string]string "generator"
	CurrentTime string "currentTime"
	Version     string "version"
	RowSets     []EmdrHistoryRowSet "rowsets"
	Columns     []string "columns"
}

type EmdrHistoryRowSet struct {
	GeneratedAt string "generatedAt"
	RegionId    int64 "regionID"
	TypeId      int64 "typeID"
	Rows        [][]interface{} "rows"
	DataRows    []EmdrHistoryDocument
}

type EmdrHistoryDocument struct {
	RegionId       int64
	RegionName     string
	TypeId         int64
	TypeName       string
	Date           string
	Orders         int64
	Quantity       int64
	Low            float64
	High           float64
	Average        float64
}

func (order *EmdrHistoryMessage) mapRows() {
	// Rewrite the data sets to a key:value format.
	// Loop the sets.
	for setIndex, setValue := range order.RowSets {
		order.RowSets[setIndex].DataRows = make([]EmdrHistoryDocument, len(order.RowSets[setIndex].Rows))

		// Loop the rows.
		for rowIndex, rowValue := range setValue.Rows {
			// Create a new mapping.
			mapping := EmdrHistoryDocument{}
			for index, name := range order.Columns {
				setHistoryFieldValue(name, &mapping, rowValue[index])
			}
			order.RowSets[setIndex].DataRows[rowIndex] = mapping
		}

		// Remove the old data
		order.RowSets[setIndex].Rows = nil
	}
}

func setHistoryFieldValue(name string, item *EmdrHistoryDocument, value interface{}) {
	s := reflect.ValueOf(item).Elem()
	field := s.FieldByNameFunc(func(fieldName string) bool {
		return strings.EqualFold(fieldName, name)
	})

	if field.IsValid() {
		switch v := value.(type) {
		case string:
			field.SetString(v)
		case int64:
			field.SetInt(v)
		case float64:
			switch field.Interface().(type) {
			case int64:
				field.SetInt(int64(v))
			default:
				field.SetFloat(v)
			}
		case bool:
			field.SetBool(v)
		default:
			log.Panicf("%T", v)
		}
	} else {
		log.Println(name)
		log.Printf("%T", field)
		log.Printf("%T", value)
		os.Exit(1)
	}
}
