package reddit

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/cdipaolo/sentiment"
	utils "github.com/piotrksiazek/fomo-sapiens/utils"
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


func GetPostIds(subreddit string, args utils.RequestArgs) NamesAndIds {
	url := baseUrl + subreddit + "/" + "new" + ".json"
	url = utils.AddRequestArgs(url, args)
	fmt.Println(url)

	header := utils.Header{Key:"User-Agent", Value:"My_unique_user_agent"}
	headers := []utils.Header{header}
	body := utils.GetRequestBody(url, "GET", headers)

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

	header := utils.Header{Key:"User-Agent", Value:"My_unique_user_agent"}
	headers := []utils.Header{header}
	body := utils.GetRequestBody(url, "GET", headers)
	
	r := regexp.MustCompile(`"body"\s*:\s*"([^"]+)`)
	matches := r.FindAllStringSubmatch(string(body), -1)

	var result []string
	for _, v := range matches {
		c <- v[1]
	}

	return result
}

func (nids NamesAndIds) GetCommentsManyPosts(subrettit string) []string {
	var result []string
	c := make(chan string)

	for _, nid := range nids {
		time.Sleep(time.Second * 2) //avoid being blocked by reddit for too frequent requests
		go func(ch chan string, n NameAndId) {
				getCommentsSinglePost(n, "Bitcoin", ch)
			}(c, nid)
		}
	for i:=0; i<len(nids); i++ {
		var comment string = <-c
		result = append(result, comment)
	}
	return result
}

func GetSentiment(comments []string) int { //returns percentage of positive-sentiment comments
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	var analysis *sentiment.Analysis
	sentimentAccum := 0
	for _, comment := range comments{
		analysis = model.SentimentAnalysis(comment, sentiment.English)
		if analysis.Score == 1{
			sentimentAccum++
		}
	}
	totalNumberOfComments := len(comments)
	return int((float64(sentimentAccum)/float64(totalNumberOfComments))*100)
}