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
	"os"
	"github.com/joho/godotenv"
	"log"
	"fmt"
)

func main() {
	fmt.Println(`
aircharts.api Copyright (C) 2018 Daniel A. Hawton daniel@hawton.com
This program comes with ABSOLUTELY NO WARRANTY; This is free software,
and you are welcome to redistribute it under certain conditions; See
LICENSE file for more information.`)
	fmt.Println("")

    err := godotenv.Load()
    if err != nil {
    	log.Fatal("Error loading .env file")
	}
	a := App{}
	a.Initialize(os.Getenv("DB_USERNAME"),os.Getenv("DB_PASSWORD"),os.Getenv("DB_DATABASE"))
	if os.Getenv("PORT") != "" {
		a.Run(":" + os.Getenv("PORT"))
	} else {
		a.Run(":8080")
	}
}