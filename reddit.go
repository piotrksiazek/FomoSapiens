package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var baseUrl string = "https://www.reddit.com/r/"


type Post struct {
	Data struct {
		Children []struct {
			Data struct {
				Id string `json:"id"`
				Name string `json:"name"`
			} `json"data"`
		} `json"children`
	} `json"data"`
}

type NameAndId struct {
	Name string
	Id string
}


func getPostIds(subreddit string) []NameAndId {
	// url = url + subreddit + "/" + "new.json"
	url := baseUrl + subreddit + "/" + "new.json"

	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	req.Header.Set("User-Agent", "My_unique_user_agent")

	res, err := http.DefaultClient.Do(req)
	checkError(err)

	body, err := ioutil.ReadAll(res.Body)
	checkError(err)

	result := Post{}
	json.Unmarshal([]byte(body), &result)

	var ids []NameAndId

	for _, post := range result.Data.Children {
		nid := NameAndId{}
		nid.Id = post.Data.Id
		nid.Name = post.Data.Name
		ids = append(ids, nid)
	}

	return ids
}

func getComments(nid NameAndId, subreddit string) []string {
	var url string = baseUrl + subreddit + "/comments/" + nid.Id + "/" + nid.Name + ".json"

	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	req.Header.Set("User-Agent", "Hello_me")

	res, err := http.DefaultClient.Do(req)
	checkError(err)

	body, err := ioutil.ReadAll(res.Body)
	checkError(err)
	
	r := regexp.MustCompile(`"body"\s*:\s*"([^"]+)`)
	matches := r.FindAllStringSubmatch(string(body), -1)

	var result []string
	for _, v := range matches {
		result = append(result, v[1])
		fmt.Println(v[1])
	}

	return result
}

// "body"\s*:\s*"([^"]+)