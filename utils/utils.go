package utils

import (
	"io/ioutil"
	"net/http"
)

type Header struct {
	Key string
	Value string
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func GetRequestBody(url string, method string, headers []Header) []byte {
	req, err := http.NewRequest(method, url, nil)
		CheckError(err)

	if len(headers) != 0 {
		for _, header := range headers {
			req.Header.Set(header.Key, header.Value)
		}
	}	
	res, err := http.DefaultClient.Do(req)
		CheckError(err)
	body, err := ioutil.ReadAll(res.Body)
		CheckError(err)
	return body
}