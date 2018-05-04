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
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strings"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}
func (a *App) Run(addr string) {
	fmt.Println("Starting ListenAndServe...")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/v2/Airport/{id:[A-Za-z0-9,]+}", a.getAirports).Methods("GET")
}

func (a *App) getAirports(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := vars["id"]

	data = strings.Replace(data, " ", "", -1)

	records := make(map[string]record)

	var ids []string
	if strings.Contains(data, ",") {
		ids = strings.Split(data, ",")
	} else {
		ids = append(ids, data)
	}

	for _, id := range ids {
		record, err := getRecord(a.DB, id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "Airport Not Found")
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		records[id] = *record
	}

	respondWithJSON(w, http.StatusOK, records)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}