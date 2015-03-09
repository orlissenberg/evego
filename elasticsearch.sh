#!/bin/sh

#curl -XDELETE 'http://localhost:9200/eve/history'
#curl -XPUT 'http://localhost:9200/eve/_mapping/order' -d @./elasticsearch/history_map.json

curl -XDELETE 'http://localhost:9200/eve/order'
echo
curl -XPUT 'http://localhost:9200/eve/_mapping/order' -d @./elasticsearch/order_map.json
echo
