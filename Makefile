.PHONY: feed-map-data

feed-map-data:
	docker run --rm --network franki_map_default -v $$(pwd)/gh-data/tajikistan-latest.osm.pbf:/data.pbf justb4/osm2pgsql:latest osm2pgsql -c -d gis -U postgres -H franken_db -P 5432 /data.pbf


.PHONY: feed-poi-data

feed-poi-data:
	./update_search.sh
