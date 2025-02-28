package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MehulSuthar-000/railAPIGin/dbutils"
	"github.com/gin-gonic/gin"

	_ "modernc.org/sqlite"
)

// DB Driver visible to whole program
var DB *sql.DB

// StationResources holds information about locations
type StationResources struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	OpeningTime string `json:"opening_time"`
	ClosingTime string `json:"closing_time"`
}

// GetStation return the station detail
func GetStation(ctx *gin.Context) {
	var station StationResources
	id := ctx.Param("station-id")
	err := DB.QueryRow(
		`SELECT ID, NAME,
			CAST(OPENING_TIME as CHAR),
			CAST(CLOSING_TIME as CHAR) 
		from station 
		where ID=?`,
		id,
	).Scan(
		&station.ID,
		&station.Name,
		&station.OpeningTime,
		&station.ClosingTime,
	)
	if err != nil {
		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			})
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"result": station,
			})

	}
}

// CreateStation handles the POST
func CreateStation(ctx *gin.Context) {
	var station StationResources
	// Parse the body into our resource
	if err := ctx.BindJSON(&station); err == nil {
		// Format Time to go time format
		statement, _ := DB.Prepare(
			`INSERT INTO station 
			(NAME, OPENING_TIME, CLOSING_TIME)
			values (?,?,?)`)
		result, err := statement.Exec(station.Name, station.OpeningTime, station.ClosingTime)
		if err == nil {
			newID, _ := result.LastInsertId()
			station.ID = int(newID)
			ctx.JSON(http.StatusOK, gin.H{
				"result": station,
			})
		} else {
			ctx.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
}

// RemoveStation handles the removing of resource
func RemoveStation(ctx *gin.Context) {
	id := ctx.Param("station-id")
	statement, _ := DB.Prepare("DELETE FROM station WHERE ID=?")
	_, err := statement.Exec(id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.String(http.StatusOK, "")
	}
}

func main() {
	var err error
	DB, err = sql.Open("sqlite", ".\\database\\railapi.db")
	if err != nil {
		log.Println("Driver creation failed!")
	}
	dbutils.Initialize(DB)

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	// Add routes to REST verbs
	router.GET("/v1/stations/:station-id", GetStation)
	router.POST("/v1/stations", CreateStation)
	router.DELETE("/v1/stations/:station-id", RemoveStation)

	// Start serving the application
	router.Run(":8000")
}
