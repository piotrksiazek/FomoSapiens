package main

import (
	"encoding/json"
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

type NamesAndIds []NameAndId


func getPostIds(subreddit string) NamesAndIds {
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

	var ids NamesAndIds

	for _, post := range result.Data.Children {
		nid := NameAndId{}
		nid.Id = post.Data.Id
		nid.Name = post.Data.Name
		ids = append(ids, nid)
	}

	return ids
}

func getCommentsSinglePost(nid NameAndId, subreddit string, c chan string) []string {
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
		c <- v[1]
	}

	return result
}

// func (nids NamesAndIds) getCommentsManyPosts(c chan string, subrettit string) {
// 	for _, nid := range nids {
// 		go getCommentsSinglePost(nid, "Bitcoin", c)
// 	}
// 	// time.Sleep(4 * time.Second)
// 	// close(c)
// }

func (nids NamesAndIds) getCommentsManyPosts(subrettit string) []string {
	var result []string
	c := make(chan string)

	for _, nid := range nids {
		go func(ch chan string, n NameAndId) {
				getCommentsSinglePost(n, "Bitcoin", ch)
				// fmt.Println(<-ch)
			}(c, nid)
		}
	for i:=0; i<len(nids); i++ {
		var comment string = <-c
		result = append(result, comment)
	}
	return result
}