#!/bin/bash
# Script for operators to update the Search Engine after editing custom_places.csv

echo "Starting Pelias Search Engine update..."

docker run --rm \
  --network frankenmap_default \
  -v $(pwd)/pelias.json:/code/pelias.json \
  -v $(pwd)/gh-data:/data \
  pelias/csv-importer npm start

echo ""
echo "✅ Update complete! The new POIs and streets are now searchable."
