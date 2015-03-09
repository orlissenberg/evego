package emdr

import (
	"os"
	"log"
	"encoding/json"
	redis "github.com/garyburd/redigo/redis"
)

type RedisEmdrWriter struct {}

func (writer *RedisEmdrWriter) Write(message []byte) (err error) {
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

func (writer *RedisEmdrWriter) WriteOrder(message []byte) (err error) {
	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)

	// order.mapRows()

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	js, err := json.Marshal(order)
	if err == nil {
		if _, err := c.Do("LPUSH", redis.Args{}.Add("orderlist").Add(js)...); err != nil {
			log.Fatalln(err)
		}

		log.Println("Order")
	} else {
		log.Println(err)
	}

	return
}

func (writer *RedisEmdrWriter) WriteHistory(message []byte) (err error) {
	log.Fatalln("NOT_IMPLEMENTED")
	return
}

func (writer *RedisEmdrWriter) DeleteAll() (err error){
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	defer c.Close()

	_, err = c.Do("DEL", redis.Args{}.Add("orderlist"));
	return
}

