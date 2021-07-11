package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Header struct {
	Key string
	Value string
}

type RequestArgs map[string]string


func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func AddRequestArgs(baseUrl string, args RequestArgs) string {
	if args != nil {
		q := url.Values{}
		baseUrl += "?"
		for key, value := range args {
			q.Add(key, value)
		}
		return baseUrl + q.Encode()
	}
	return baseUrl
}

func GetRequestBody(url string, method string, headers []Header) []byte {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}

	if len(headers) != 0 {
		for _, header := range headers {
			req.Header.Add(header.Key, header.Value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	
	return body
}

func AddLeadingZeroIfSingleDigit(number string) string {
	if len(number) == 1 {
		return "0" + number
	}
	return number
}