package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	coingecko "github.com/piotrksiazek/fomo-sapiens/coingecko"
	models "github.com/piotrksiazek/fomo-sapiens/models"
	reddit "github.com/piotrksiazek/fomo-sapiens/reddit"
	"github.com/piotrksiazek/fomo-sapiens/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var templates *template.Template
var db *gorm.DB
func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	templates = template.Must(template.ParseGlob("templates/*.html"))

	db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&models.Day{})

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/getDays/{count}", getDays).Methods("GET")

	

	// http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func getDays(w http.ResponseWriter, r *http.Request) {
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
	// args := utils.RequestArgs{"after" : "2d", "before": "1d", "limit" : "100"}
	// nids := reddit.GetPostIds("Bitcoin", args)
	// comments := nids.GetCommentsManyPosts("Bitcoin")
	// sentiment := reddit.GetSentiment(comments)

	// price := coingecko.GetCurrentPrice("bitcoin", "usd")
	// predictedPrice := 34000
	
	// day := models.Day{CreationDay: time.Now(), RealPrice: price, PredictedPrice: predictedPrice, Sentiment: sentiment}
	// db.Create(&day)
	// populateWithHistoricalData()
	// fmt.Println(coingecko.GetHistoricalPrice("bitcoin", "usd", coingecko.Date{Day:"12", Month: "12", Year: "2013"}))
	// fmt.Println(price, sentiment)
}

func populateWithHistoricalData() {

	for i:= 0 ; i<100; i++ {
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
		predictedPrice := 34000
		
		dayModel := models.Day{CreationDay: date, RealPrice: price, PredictedPrice: predictedPrice, Sentiment: sentiment}
		db.Create(&dayModel)
	}
}