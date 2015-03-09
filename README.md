#EveGo

Go applications to import Eve Online data into Elasticsearch.

- InvTypes (SqlServer)
- Regions (SqlServer)
- Stations (Sqlite)
- Solar systems (Sqlite)
- Orders (EMDR)

##Only build.

    ./build.sh

##Only run.

    ./run.sh -verbose=vvv -es_host=elacticserver.here:9200 sqlite

##Build and run.

    ./build.sh -verbose=vvv -es_host=elacticserver.here:9200 emdr

##Build & run SqlServer data import.

    ./wineve.cmd

