/***********************************************************
 *  AirCharts API Golang Implementation (aircharts.api)
 *  Copyright (C) 2018 Daniel A. Hawton daniel@hawton.com
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */


package main

import (
	"database/sql"
	"fmt"
)

type record struct {
	Info *airport `json:"info"`
	Charts map[string][]chart `json:"charts"`
}

type airport struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
	Elevation int `json:"elevation"`
}

type chart struct {
	ID string `json:"id"`
	Chartname string `json:"chartname"`
	Url string `json:"url"`
	Proxy string `json:"proxy"`
}

func getRecord(db *sql.DB, icao string) (*record, error) {
	var err error
	r := record{}
	r.Info, err = getAirport(db, icao)
	if err != nil {
		return nil, err
	}
	r.Charts, err = getCharts(db, icao)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func getAirport(db *sql.DB, icao string) (*airport, error) {
	apt := airport{ID: icao}
	statement := fmt.Sprintf("SELECT name, lat, lon, elevation FROM airports WHERE id='%s'", icao)
	err := db.QueryRow(statement).Scan(&apt.Name, &apt.Lat, &apt.Lon, &apt.Elevation)

	if err != nil {
		fmt.Printf("Got err in getAirport for %s", icao)
		return nil, err
	}

	return &apt, nil
}

func getCharts(db *sql.DB, icao string) (map[string][]chart, error) {
	var err error
	charts := map[string][]chart{}
	charttypes := []string{"General", "SID", "STAR", "Intermediate", "Approach"}
	for _, ct := range charttypes {
		charts[ct], err = getChartsByType(db, icao, ct)
		if err != nil {
			fmt.Printf("Got err in getCharts(%s)", ct)
			return nil, err
		}
	}

	return charts, nil
}

func getChartsByType(db *sql.DB, icao, charttype string) ([]chart, error) {
	statement := fmt.Sprintf("SELECT id, chartname, url FROM charts WHERE charttype='%s' AND (icao='%s' OR iata='%s') ORDER BY chartname", charttype, icao, icao)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Print("Got err in getChartsByType")
		return nil, err
	}

	defer rows.Close()
	charts := make([]chart, 0)
	for rows.Next() {
		var c chart
		if err := rows.Scan(&c.ID, &c.Chartname, &c.Url); err != nil {
			fmt.Print("Got err in getChartsByType row scan")
			return nil, err
		}
		c.Proxy = "https://www.aircharts.org/view/" + c.ID
		charts = append(charts, c)
	}

	return charts, nil
}