package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//struct to hold refs of router and database
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type Building struct {
	ID            uint    `json:"id"`
	DoittID       string  `json:"doittId" gorm:"unique;not null"`
	Bin           string  `json:"bin" gorm:"not null"`
	ConstructYear int     `json:"constructYear" gorm:"index:building_year"`
	HeightRoof    float64 `json:"heightRoof"`
	Area          float64 `json:"area"`
}

// create database connection and set up routing
func (svc *App) Initialize(user, password, dbname, host, port, sslmode string) {
	dbConfig := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbname, host, port, sslmode)

	var err error
	svc.DB, err = gorm.Open("postgres", dbConfig)

	if err != nil {
		panic("failed to connect database")
	}

	svc.DB.LogMode(true)

	svc.DB.AutoMigrate(&Building{})

	svc.Router = mux.NewRouter()
	svc.initializeRoutes()
}

// run application
func (svc *App) Run(addr string) {
	handler := cors.Default().Handler(svc.Router)

	log.Fatal(http.ListenAndServe(addr, handler))

	defer svc.DB.Close()
}

// initialize routes into router that call methods on requests
func (svc *App) initializeRoutes() {
	// all endpoints support pagination by params ?from=[id]&limit=[1-100]

	// get buildings of a specific year
	svc.Router.HandleFunc("/v1/buildings/years/{year:[0-9]+}", svc.GetBuildingsOfYear).Methods("GET")

	// get all buildings
	svc.Router.HandleFunc("/v1/buildings", svc.GetBuildings).Methods("GET")

	// get a building by id
	svc.Router.HandleFunc("/v1/buildings/{id:[0-9]+}", svc.GetBuilding).Methods("GET")

	// get statistics results for buildings of a specific year
	svc.Router.HandleFunc("/v1/buildings/stats/years/{year}", svc.GetSummaryByYear).Methods("GET")

	// get statistics results for buildings group by years
	svc.Router.HandleFunc("/v1/buildings/stats/years", svc.GetSummaryByYears).Methods("GET")
}

// send a payload of JSON content
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a JSON error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func handleQueryError(err error, w http.ResponseWriter) {
	if gorm.IsRecordNotFoundError(err) {
		respondWithError(w, http.StatusNotFound, "Record not found")
	} else {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

// params year validator
func getParamYear(paramsYear string) (int, error) {
	year, err := strconv.Atoi(paramsYear)
	if err != nil || (err == nil && year <= 0) {
		err = errors.New("year must be positive")
	}
	return year, err
}

// Query Scopes
func OfConstructYear(constructYear int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&Building{ConstructYear: constructYear})
	}
}

func Page(fromId, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id >= ?", fromId).Limit(limit)
	}
}

// limit max 100 records
// params ?from=0&limit=100
func (svc *App) PagedDB(r *http.Request) *gorm.DB {
	params := r.URL.Query()
	fromId, err := strconv.Atoi(params.Get("from"))
	if err != nil || fromId < 0 {
		fromId = 0
	}

	const MaxPageSize = 100

	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil || limit > MaxPageSize || limit <= 0 {
		limit = MaxPageSize
	}

	return svc.DB.Scopes(Page(fromId, limit))
}

// handle get building request
func (svc *App) GetBuildings(w http.ResponseWriter, r *http.Request) {
	var buildings []Building
	if err := svc.PagedDB(r).Find(&buildings).Error; err != nil {
		handleQueryError(err, w)
	} else {
		respondWithJSON(w, http.StatusOK, buildings)
	}

}

func (svc *App) GetBuilding(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var building Building
	if err := svc.DB.First(&building, params["id"]).Error; err != nil {
		handleQueryError(err, w)
	} else {
		respondWithJSON(w, http.StatusOK, building)
	}
}

func (svc *App) GetBuildingsOfYear(w http.ResponseWriter, r *http.Request) {
	var year = -1
	var err error
	params := mux.Vars(r)
	// validate param["year"]
	if year, err = getParamYear(params["year"]); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var buildings []Building
	if err := svc.PagedDB(r).Scopes(OfConstructYear(year)).Find(&buildings).Error; err != nil {
		handleQueryError(err, w)
	} else {
		respondWithJSON(w, http.StatusOK, buildings)
	}

}

type Stats struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Count int     `json:"count"`
}

type SummaryOfYear struct {
	Year   int   `json:"year"`
	Height Stats `json:"height"`
	Area   Stats `json:"area"`
}

func buildingSummaryByYears(db *gorm.DB) ([]SummaryOfYear, error) {
	rows, err := db.Table("buildings").
		Select("construct_year, MIN(height_roof), MAX(height_roof), AVG(height_roof), MIN(area), MAX(area), AVG(area), count(area)").
		Group("construct_year").Order("construct_year desc").Rows()

	var summaries []SummaryOfYear
	if err == nil {
		for rows.Next() {
			var (
				year = 0
				minH = 0.0
				maxH = 0.0
				avgH = 0.0

				minA = 0.0
				maxA = 0.0
				avgA = 0.0

				count = 0
			)
			rows.Scan(&year, &minH, &maxH, &avgH, &minA, &maxA, &avgA, &count)
			summary := SummaryOfYear{
				Year:   year,
				Height: Stats{Min: minH, Max: maxH, Avg: avgH, Count: count},
				Area:   Stats{Min: minA, Max: maxA, Avg: avgA, Count: count},
			}
			summaries = append(summaries, summary)
		}
	}
	return summaries, err
}

func (svc *App) GetSummaryByYears(w http.ResponseWriter, r *http.Request) {
	summaries, err := buildingSummaryByYears(svc.DB)
	if err != nil {
		handleQueryError(err, w)
	} else {
		respondWithJSON(w, http.StatusOK, summaries)
	}
}

func (svc *App) GetSummaryByYear(w http.ResponseWriter, r *http.Request) {
	var year = -1
	var err error

	// TODO: use request validator
	params := mux.Vars(r)
	if year, err = getParamYear(params["year"]); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid year")
		return
	}

	summaries, err := buildingSummaryByYears(svc.DB.Scopes(OfConstructYear(year)))
	if err != nil {
		handleQueryError(err, w)
	} else if len(summaries) == 0 {
		respondWithError(w, http.StatusNotFound, "Not found summary for year "+params["year"])
	} else {
		respondWithJSON(w, http.StatusOK, summaries[0])
	}

}
