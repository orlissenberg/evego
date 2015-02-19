package emdr

import (
	"os"
	"fmt"
	"encoding/json"
	elastic "github.com/mattbaird/elastigo/lib"
	redis "github.com/garyburd/redigo/redis"
)

type ElasticEmdrWriter struct {}

func (writer *ElasticEmdrWriter) Write(message []byte) (err error) {
	var v EmdrMessage
	json.Unmarshal(message, &v)

	switch v.ResultType {
	case "orders":
		err = writer.WriteOrder(message)
	case "history":
		err = writer.WriteHistory(message)
	}

	return
}

func (writer *ElasticEmdrWriter) WriteOrder(message []byte) (err error) {
	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)

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

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	js, err := json.Marshal(order)
	if err == nil {
		if _, err := c.Do("LPUSH", redis.Args{}.Add("orderlist").Add(js)...); err != nil {
			panic(err)
		}

		fmt.Println("Order: " + UnixTimeStampString())
	} else {
		fmt.Println(err)
		// Dump data to file system on error for inspection.
		// ioutil.WriteFile("error_order_"+strconv.FormatInt(time.Now().Unix(), 10)+".json", message, 0644)
	}

	return
}

func (writer *ElasticEmdrWriter) WriteHistory(message []byte) (err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	fmt.Println("History: " + UnixTimeStampString())
	_, err = c.Index("eve", "history", "", nil, string(message))

	return
}

func (writer *ElasticEmdrWriter) DeleteAll() (err error){
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	defer c.Close()

	_, err = c.Do("DEL", redis.Args{}.Add("orderlist"));
	return
}

