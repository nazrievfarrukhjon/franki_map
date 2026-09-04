1. Install Docker.
2. Install Docker Compose.
3. Install `make` .
4. git clone repository from git
5. cd franki_map
6. Set Elasticsearch memory requirement: sudo sysctl -w vm.max_map_count=262144
7. Download tajikistan-latest.osm.pbf to franki_map/gh-data/
8. Run: docker-compose up -d
9. Feed map data to PostGIS (for Martin tiles) by running:
  make feed-map-data
10. Feed POI/address data to Elasticsearch (for Pelias search) by running:
  make feed-poi-data
