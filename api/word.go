package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/eyaegashi/wordTestApp/config"
	"github.com/eyaegashi/wordTestApp/util"
)

const filterFile = "./filter.json"
const defaultFilter = "Office"

// TestWordResult is the struct to be returned as API result
type TestWordResult struct {
	Result       int
	msg          string
	TestWordInfo TestWordInfo
}

// TestWordInfo is the structure to get test word info
type TestWordInfo struct {
	TestWord        string
	TranslatedWord  string
	ExampleSentence string
}

// to get search filter value
type filters struct {
	Filter []string `json:"filter"`
}

// to get wordlist from API
type candidateTestWords struct {
	Metadata struct {
		SourceLanguage string `json:"sourceLanguage"`
		Provider       string `json:"provider"`
		Limit          int    `json:"limit"`
		Offset         int    `json:"offset"`
		Total          int    `json:"total"`
	} `json:"metadata"`
	Results []struct {
		Word string `json:"word"`
		ID   string `json:"id"`
	} `json:"results"`
}

// to get Japanese translated word
type translatedWord struct {
	Result string `json:"result"`
	Tuc    []struct {
		Phrase struct {
			Text     string `json:"text"`
			Language string `json:"language"`
		} `json:"phrase,omitempty"`
		Meanings []struct {
			Language string `json:"language"`
			Text     string `json:"text"`
		} `json:"meanings,omitempty"`
		MeaningID int64 `json:"meaningId"`
		Authors   []int `json:"authors"`
	} `json:"tuc"`
	Phrase  string `json:"phrase"`
	From    string `json:"from"`
	Dest    string `json:"dest"`
	Authors struct {
		Num1 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"1"`
		Num76 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"76"`
		Num60172 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"60172"`
		Num83042 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"83042"`
		Num83695 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"83695"`
		Num84500 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"84500"`
		Num86065 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"86065"`
		Num86934 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"86934"`
		Num91945 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"91945"`
		Num93369 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"93369"`
		Num97745 struct {
			U   string `json:"U"`
			ID  int    `json:"id"`
			N   string `json:"N"`
			URL string `json:"url"`
		} `json:"97745"`
	} `json:"authors"`
}

// to get example sentence
type exampleSentence struct {
	Metadata struct {
		Provider string `json:"provider"`
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Language        string `json:"language"`
			LexicalCategory string `json:"lexicalCategory"`
			Sentences       []struct {
				Regions   []string `json:"regions"`
				SenseIds  []string `json:"senseIds"`
				Text      string   `json:"text"`
				Domains   []string `json:"domains,omitempty"`
				Registers []string `json:"registers,omitempty"`
			} `json:"sentences"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
}

// CreateTestWordAPI is to get test word info and return result as API
func CreateTestWordAPI() (testWordResult TestWordResult) {
	// get test word
	getWord(&testWordResult)

	if testWordResult.Result != 0 {
		// prepare for goroutine
		//var wg sync.WaitGroup
		//wg.Add(2)

		// get tralnslated word and example sentence by goroutine
		//go getTranslatedWord(&testWordResult, &wg)
		//go getExampleSentence(&testWordResult, &wg)
	}

	return testWordResult
}

// getWord is to get a word from external API
func getWord(testWordResult *TestWordResult) {
	// prepare to get test item
	conf := config.GetConfig()
	filter := getFilter()

	//set wordAPI URL
	// URL例: https://od-api.oxforddictionaries.com:443/api/v1/wordlist/en/domains%3DArt
	url := conf.WordAPI.URL + filter
	fmt.Println(url)
	// send request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("app_id", conf.WordAPI.APIID)
	req.Header.Set("app_key", conf.WordAPI.APIKey)
	req.Header.Set("Accept", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		testWordResult.Result = 0
		testWordResult.msg = "external API request error"
		return
	}

	// get respone body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// convert JSON to wordItem struct
	var word string
	err = json.Unmarshal(body, &word)
	if err != nil {
		// todo: convert error
		testWordResult.Result = 0
		testWordResult.msg = "convert JSON of external API"
	}
	testWordResult.TestWordInfo.TestWord = word
	return
}

func getFilter() string {
	var filterItems filters
	jsonData, err := util.LoadjsonFile(filterFile)
	if err != nil {
		return defaultFilter
	}

	// convert json to struct
	err = json.Unmarshal(jsonData, &filterItems)
	if err != nil {
		return defaultFilter
	}
	index := rand.Intn(len(filterItems.Filter))
	return filterItems.Filter[index]
}

// get the translated word (Japanese)
func getTranslatedWord(testWordResult *TestWordResult, wg *sync.WaitGroup) {
	// todo: get Japanese translated word
	// exmample URL for the external API:
	// https://glosbe.com/gapi_v0_1/translate?from=en&dest=ja&phrase=dog&format=json
	//wg.Done
}

// get the example sentence of the word
func getExampleSentence(testWordResult *TestWordResult, wg *sync.WaitGroup) {
	// todo: get example sentence
	// exmample URL for the external API:
	// https://od-api.oxforddictionaries.com:443/api/v1/entries/en/dog/sentences
	//wg.Done
}
