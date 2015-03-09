#EveGo

Go applications to import Eve Online data into Elasticsearch.

- InvTypes (SqlServer)
- Regions (SqlServer)
- Stations (Sqlite)
- Solar systems (Sqlite)
- Orders (EMDR)

**Build**

    ./build.sh

**Run**

    ./run.sh -verbose=vvv -es_host=elacticserver.here:9200 sqlite

**Build & run.**

    ./build.sh -verbose=vvv -es_host=elacticserver.here:9200 emdr

**Build & run SqlServer data import (Windows).**

    ./wineve.cmd

#### References

[Eve Market Data Relay (EMDR)](https://eve-market-data-relay.readthedocs.org/en/latest/)

