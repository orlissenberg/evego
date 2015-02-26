package emdr

type EmdrMessage struct {
	ResultType string "resultType"
}

type EmdrOrderRowSet struct {
	TypeId      int64 "typeID"
	RegionID    int64 "regionID"
	GeneratedAt string "generatedAt"
	Rows        [][]interface{} "rows" // Format not document friendly.
	DataRows    []map[string]interface{}
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

func (order *EmdrOrderMessage) mapRows() {
	// Rewrite the data sets to a key:value format.
	// Loop the sets.
	for setIndex, setValue := range order.RowSets {
		order.RowSets[setIndex].DataRows = make([]map[string]interface{}, len(order.RowSets[setIndex].Rows))

		// Loop the rows.
		for rowIndex, rowValue := range setValue.Rows {
			// Create a new mapping.
			mapping := make(map[string]interface{})
			for index, name := range order.Columns {
				mapping[name] = rowValue[index]
			}
			order.RowSets[setIndex].DataRows[rowIndex] = mapping
		}

		// Remove the old data
		order.RowSets[setIndex].Rows = nil
	}
}

type EmdrHistoryMessage struct {
	EmdrMessage
	UploadKeys  []map[string]string "uploadKeys"
	Generator   map[string]string "generator"
	CurrentTime string "currentTime"
	Version     string "version"
	RowSets     []EmdrOrderRowSet "rowsets"
	Columns     []string "columns"
}

type EmdrHistoryRowSet struct {
	GeneratedAt string "generatedAt"
	RegionID    int64 "regionID"
	TypeId      int64 "typeID"
	Rows        [][]interface{} "rows"
	DataRows    []map[string]interface{}
}

func (history *EmdrHistoryMessage) mapRows() {
	// Rewrite the data sets to a key:value format.
	// Loop the sets.
	for setIndex, setValue := range history.RowSets {
		history.RowSets[setIndex].DataRows = make([]map[string]interface{}, len(history.RowSets[setIndex].Rows))

		// Loop the rows.
		for rowIndex, rowValue := range setValue.Rows {
			// Create a new mapping.
			mapping := make(map[string]interface{})
			for index, name := range history.Columns {
				mapping[name] = rowValue[index]
			}
			history.RowSets[setIndex].DataRows[rowIndex] = mapping
		}

		// Remove the old data
		history.RowSets[setIndex].Rows = nil
	}
}
