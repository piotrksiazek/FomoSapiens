package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/onatm/clockwerk"
	coingecko "github.com/piotrksiazek/fomo-sapiens/coingecko"
	models "github.com/piotrksiazek/fomo-sapiens/models"
	reddit "github.com/piotrksiazek/fomo-sapiens/reddit"
	"github.com/piotrksiazek/fomo-sapiens/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var templates *template.Template
var db *gorm.DB

type DailyJob struct {}

func main() {
	r := mux.NewRouter()

	//Daily fetch and ananyle data from reddit and coingecko apis
	var dailyJob DailyJob
	cw := clockwerk.New()
	cw.Every(24 * time.Hour).Do(dailyJob)
	cw.Start()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	templates = template.Must(template.ParseGlob("templates/*.html"))

	db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	db.AutoMigrate(&models.Day{})

	r.HandleFunc("/", indexHandler).Methods("GET")

	r.HandleFunc("/getDays/{count}", getDaysHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

//
func (d DailyJob) Run() {
	args := utils.RequestArgs{"after" : "2d", "before": "1d", "limit" : "100"}
	nids := reddit.GetPostIds("Bitcoin", args)
	comments := nids.GetCommentsManyPosts("Bitcoin")
	sentiment := reddit.GetSentiment(comments)

	price := coingecko.GetCurrentPrice("bitcoin", "usd")
	predictedPrice := 34000
	
	day := models.Day{CreationDay: time.Now(), RealPrice: price, PredictedPrice: predictedPrice, Sentiment: sentiment}
	db.Create(&day)
}

func getDaysHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	numberOfDays, _ := strconv.Atoi(params["count"])
	condition := time.Now().AddDate(0,0, -numberOfDays)

	var days []models.Day
	db.Where("creation_day > ?", condition).Find(&days)

	json.NewEncoder(w).Encode(days)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func populateWithHistoricalData(howManyDays int) {

	for i:= 0 ; i<howManyDays; i++ {
		after := strconv.Itoa(2 + i) + "d"
		before := strconv.Itoa(1 + i) + "d"
		args := utils.RequestArgs{"after" : after, "before": before, "limit" : "100"}
		nids := reddit.GetPostIds("Bitcoin", args)
		
		comments := nids.GetCommentsManyPosts("Bitcoin")
		
		sentiment := reddit.GetSentiment(comments)

		date := time.Now().AddDate(0,0,-i)
		day := utils.AddLeadingZeroIfSingleDigit(strconv.Itoa(date.Day()))
		month := utils.AddLeadingZeroIfSingleDigit(strconv.Itoa(int(date.Month())))
		year := utils.AddLeadingZeroIfSingleDigit(strconv.Itoa(date.Year()))

		price := coingecko.GetHistoricalPrice("bitcoin", "usd", coingecko.Date{Day: day, Month: month, Year: year})
		predictedPrice := 0
		
		dayModel := models.Day{CreationDay: date, RealPrice: price, PredictedPrice: predictedPrice, Sentiment: sentiment}
		db.Create(&dayModel)
		fmt.Println("Day: ", i)
	}
}