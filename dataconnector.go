package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Temperatures struct {
	Temperature float64
	Weather sql.NullString
	CreatedAt string
}

const (
	hourly = 0
	daily = 1
)

func GetLatestTemperature() (float32, sql.NullString, string) {
	db, err := sql.Open("sqlite3", "/home/pi/temperature_monitor/db/temperature.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := `
		select 
			temperature,
			weather,
			created_at
		from 
			temperature 
		order by created_at desc 
		limit 1
	`

	row := db.QueryRow(q)
	if err != nil {
		log.Fatal(err)
	}

	var temperature float32
	var weather sql.NullString
	var createdAt string
	err = row.Scan(&temperature, &weather, &createdAt)
	return temperature, weather, createdAt
}

func GetTemperatureByTerm(term int) []Temperatures{
	var fromTime time.Time
	switch term {
	case hourly:
		t := time.Now()
		fromTime = t.Add(-time.Hour)
	case daily:
		t := time.Now()
		fromTime = t.Truncate(time.Hour).Add(-time.Duration(t.Hour()) * time.Hour)
		fmt.Println(fromTime)
	default:
		log.Fatal("argument error")
	}

	db, err := sql.Open("sqlite3", "/home/pi/temperature_monitor/db/temperature.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `
		select 
			temperature,
			weather,
			created_at
		from 
			temperature 
		where
			created_at >= ?
		order by created_at asc
	`

	rows, err := db.Query(q, fromTime)
	if err != nil {
		log.Fatal(err)
	}

	tmps := []Temperatures{}
	for rows.Next() {
		var (
			temperature float64
			weather sql.NullString
			createdAt string
		)
		err = rows.Scan(&temperature, &weather, &createdAt)
		t := Temperatures{Temperature:temperature, Weather:weather, CreatedAt:createdAt}
		tmps = append(tmps, t)
	}

	return tmps
}
