package main

import (
	"encoding/json"
	"net/http"

	"github.com/eyaegashi/wordTestApp/api"
)

func main() {
	api.CreateTestWordAPI()
	// API実装中のためコメントアウト
	//http.HandleFunc("/testword", testWordAPIHandler)
	//http.ListenAndServe(":8000", nil)
}

func testWordAPIHandler(rw http.ResponseWriter, req *http.Request) {
	testWordApi := api.CreateTestWordAPI()
	data, err := json.Marshal(&testWordApi)
	if err != nil {
		// todo:error処理見直す
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}
