.PHONY: feed-map-data

feed-map-data:
	mkdir -p gh-data
	wget -nc -O gh-data/tajikistan-latest.osm.pbf https://download.geofabrik.de/asia/tajikistan-latest.osm.pbf
	docker run --rm -e PGPASSWORD=postgres --platform linux/amd64 --network franki_map_default -v $$(pwd)/gh-data/tajikistan-latest.osm.pbf:/data.pbf justb4/osm2pgsql:latest osm2pgsql -c -d gis -U postgres -H franken_db -P 5432 /data.pbf


.PHONY: feed-poi-data

feed-poi-data:
	-docker exec franken_elasticsearch bin/elasticsearch-plugin install --batch analysis-icu
	-docker restart franken_elasticsearch
	sleep 15
	./update_search.sh
