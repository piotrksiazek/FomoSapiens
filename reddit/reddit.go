package reddit

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cdipaolo/sentiment"
	utils "github.com/piotrksiazek/fomo-sapiens/utils"
)

var baseUrl string = "https://www.reddit.com/r/"

var pushiftBaseUrl string = "https://api.pushshift.io/reddit/search/"


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
		// time.Sleep(time.Second * 2) //avoid being blocked by reddit for too frequent requests
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

type PushiftPost struct {
	Score int `json:"score"`
	Body string `json:"body"`
	Author string `json:"author"`
}

type PushiftPosts struct {
	Data []struct {
		PushiftPost
	}`json:"data"`
}

func GetTopCommentFromDay(d int) PushiftPost { //d says from how many days ago should the post be searched
	after := strconv.Itoa(d) //parse d int to d in string format
	before := strconv.Itoa(d-1)
	url:= "https://api.pushshift.io/reddit/search/comment/?subreddit=Bitcoin&subreddit=CryptoCurrency&after=" + after + "d&before=" + before +"d&size=500"
	body := utils.GetRequestBody(url, "GET", []utils.Header{})
	posts := PushiftPosts{}
	json.Unmarshal([]byte(body), &posts)

	var max int
	var maxIndex int

	//find comment with the highest score among others
	for index, content := range posts.Data {
		fmt.Println(content.Body)
		tmp := content.PushiftPost.Score
		if tmp > max{
			max = tmp
			maxIndex = index
		}
	}
	return posts.Data[maxIndex].PushiftPost
}