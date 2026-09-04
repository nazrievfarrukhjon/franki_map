package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

type Tag struct {
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

type Node struct {
	ID        int64     `xml:"id,attr"`
	Version   int       `xml:"version,attr"`
	Lat       float64   `xml:"lat,attr"`
	Lon       float64   `xml:"lon,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`
	Changeset int64     `xml:"changeset,attr"`
	UID       int       `xml:"uid,attr"`
	User      string    `xml:"user,attr"`
	Tags      []Tag     `xml:"tag"`
}

type Nd struct {
	Ref int64 `xml:"ref,attr"`
}

type Way struct {
	ID        int64     `xml:"id,attr"`
	Version   int       `xml:"version,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`
	Changeset int64     `xml:"changeset,attr"`
	UID       int       `xml:"uid,attr"`
	User      string    `xml:"user,attr"`
	Nds       []Nd      `xml:"nd"`
	Tags      []Tag     `xml:"tag"`
}

type Osm struct {
	XMLName   xml.Name `xml:"osm"`
	Version   string   `xml:"version,attr"`
	Generator string   `xml:"generator,attr"`
	Nodes     []Node   `xml:"node"`
	Ways      []Way    `xml:"way"`
}

func parseTags(tagsJSON sql.NullString) []Tag {
	var tags []Tag
	if !tagsJSON.Valid || tagsJSON.String == "" {
		return tags
	}

	var tagMap map[string]string
	if err := json.Unmarshal([]byte(tagsJSON.String), &tagMap); err != nil {
		log.Println("Error parsing tags:", err)
		return tags
	}

	for k, v := range tagMap {
		tags = append(tags, Tag{K: k, V: v})
	}
	return tags
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5455")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", dbUser, dbPass, dbName, dbHost, dbPort)
	log.Printf("Connecting to APIDB at %s:%s...", dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	osm := Osm{
		Version:   "0.6",
		Generator: "go-osm-extract",
	}

	// Extract Nodes
	rows, err := db.Query(`SELECT id, version, user_id, tstamp, changeset_id, hstore_to_json(tags), ST_Y(geom), ST_X(geom) FROM nodes`)
	if err != nil {
		log.Println("No nodes found or error:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var n Node
			var tagsJSON sql.NullString
			if err := rows.Scan(&n.ID, &n.Version, &n.UID, &n.Timestamp, &n.Changeset, &tagsJSON, &n.Lat, &n.Lon); err != nil {
				log.Fatal(err)
			}
			n.User = "operator"
			n.Tags = parseTags(tagsJSON)
			osm.Nodes = append(osm.Nodes, n)
		}
	}

	// Extract Ways
	wayRows, err := db.Query(`SELECT id, version, user_id, tstamp, changeset_id, hstore_to_json(tags), nodes FROM ways`)
	if err != nil {
		log.Println("No ways found or error:", err)
	} else {
		defer wayRows.Close()
		for wayRows.Next() {
			var w Way
			var tagsJSON sql.NullString
			var ndsArray []int64
			if err := wayRows.Scan(&w.ID, &w.Version, &w.UID, &w.Timestamp, &w.Changeset, &tagsJSON, pq.Array(&ndsArray)); err != nil {
				log.Fatal(err)
			}
			w.User = "operator"
			w.Tags = parseTags(tagsJSON)
			for _, ndID := range ndsArray {
				w.Nds = append(w.Nds, Nd{Ref: ndID})
			}
			osm.Ways = append(osm.Ways, w)
		}
	}

	outFile, err := os.Create("export.osm")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	outFile.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "  ")
	if err := encoder.Encode(osm); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully exported APIDB edits to export.osm")
}
