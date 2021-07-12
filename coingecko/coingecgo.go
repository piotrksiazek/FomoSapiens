package coingecko

import (
	"encoding/json"
	"fmt"

	utils "github.com/piotrksiazek/fomo-sapiens/utils"
)

type Date struct {
	Day string
	Month string
	Year string
}

func (d Date) String() string {
	return fmt.Sprint(d.Day, "-", d.Month, "-" ,d.Year)
}


var baseUrl string = `https://api.coingecko.com/api/v3/`

func GetCurrentPrice(crypto string, currency string) int {
	url := baseUrl + "simple/price?ids=" + crypto + "&" + "vs_currencies=" + currency

	body := utils.GetRequestBody(url, "GET", nil)

	var c map[string]interface{}

	json.Unmarshal(body, &c)

	if price, ok := c[crypto].(map[string]interface{})[currency].(float64); ok {
		return int(price)
	}
	return -1 //no asset can be worth -1 dollars, signifies that error occured
}

func GetHistoricalPrice(crypto string, currency string, date Date) int {
	url := baseUrl + "coins/" + crypto + "/history?" + "date=" + date.String()
	body := utils.GetRequestBody(url, "GET", nil)

	var c map[string]interface{}

	json.Unmarshal(body, &c)

	if price, ok := c["market_data"].(map[string]interface{})["current_price"].(map[string]interface{})[currency].(float64); ok {
		return int(price)
	}
	return -1
}