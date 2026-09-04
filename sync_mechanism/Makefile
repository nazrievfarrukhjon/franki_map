include .env
export

.PHONY: sync-all extract-data convert-pbf update-triad

sync-all: extract-data convert-pbf update-triad
	@echo "✅ APIDB Synchronization completely finished!"

extract-data:
	@echo "1. Extracting data from APIDB..."
	go run extract.go

convert-pbf:
	@echo "2. Converting OSM XML to highly compressed PBF..."
	chmod +x convert_to_pbf.sh
	./convert_to_pbf.sh

update-triad:
	@echo "3. Updating Map Triad infrastructure..."
	# Copy to GraphHopper and PostGIS data directory
	cp tajikistan-latest.osm.pbf ../frankenmap/gh-data/
	
	# Update PostGIS database (Martin Map Tiles)
	@echo "Updating Map Tiles (Martin)..."
	cd ../frankenmap && make feed-map-data
	
	# Update GraphHopper
	@echo "Updating Routing Engine (GraphHopper)..."
	# We MUST delete the routing cache so it rebuilds from the new map file
	sudo rm -rf ../frankenmap/gh-data/default-gh
	docker restart franken_graphhopper
	
	# Update Pelias
	@echo "Updating Search Engine (Pelias)..."
	cd ../frankenmap && make feed-poi-data
