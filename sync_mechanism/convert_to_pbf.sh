#!/bin/bash
set -e

echo "Converting export.osm to tajikistan-latest.osm.pbf using osmium..."

# Remove old pbf if exists to avoid osmium overwriting error
rm -f tajikistan-latest.osm.pbf

docker run --rm -v $(pwd):/data stefda/osmium-tool osmium cat /data/export.osm -o /data/tajikistan-latest.osm.pbf

echo "Conversion complete!"
