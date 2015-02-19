package emdr

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
