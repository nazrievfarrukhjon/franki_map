#!/bin/bash
# Script for operators to update the Search Engine after editing custom_places.csv

echo "Starting Pelias Search Engine update..."


echo "Creating Pelias Elasticsearch Schema..."
docker run --rm   --network franki_map_default   -v /Users/farruqnazri/Desktop/franki_map/pelias.json:/code/pelias.json   pelias/schema npm run create_index

echo "Importing CSV data..."
docker run --rm \
  --network franki_map_default \
  -v $(pwd)/pelias.json:/code/pelias.json \
  -v $(pwd)/gh-data:/data \
  pelias/csv-importer npm start

echo ""
echo "✅ Update complete! The new POIs and streets are now searchable."
