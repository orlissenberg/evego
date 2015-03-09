#!/bin/sh

HOST="localhost"

#curl -XDELETE "http://$HOST:9200/eve/history"
#echo
#curl -XPUT "http://localhost:9200/eve/_mapping/order" -d @./elasticsearch/history_map.json
#echo

#curl -XDELETE "http://$HOST:9200/eve/invtype"
#echo
#curl -XPUT "http://$HOST:9200/eve/_mapping/invtype" -d @./elasticsearch/invtype_map.json
#echo

curl -XDELETE "http://$HOST:9200/eve/order"
echo
curl -XPUT "http://$HOST:9200/eve/_mapping/order" -d @./elasticsearch/order_map.json
echo
