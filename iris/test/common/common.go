package common

import (
	"io/ioutil"
	"bytes"
	"net/http"
	"github.com/chengwenxi/blockchain/iris/test/types"
	"encoding/json"
	"log"
)

func DoGet(url string) []byte {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	var result types.Result
	json.Unmarshal(body, &result)
	if result.Error != "" {
		log.Println(result.Error)
		return nil
	}
	return body
}

func DoPost(url string, data []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	var result types.Result
	json.Unmarshal(body, &result)
	if result.Error != "" {
		log.Println(result.Error)
		return nil
	}
	return body
}
