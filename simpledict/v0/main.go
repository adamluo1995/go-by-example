package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type QueryResponse struct {
	Word         string
	TransType    string
	Explanations []string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "%v\n", len(os.Args))
		os.Exit(1)
	}
	rsp := QueryResponse{Word: os.Args[1], TransType: os.Args[2]}
	err := query(&rsp)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", rsp)
	os.Exit(0)
}

func query(rsp *QueryResponse) (err error) {
	client := &http.Client{}
	// var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)
	req_str := DictRequest{TransType: rsp.TransType, Source: rsp.Word}
	buf, err := json.Marshal(req_str)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Set("authority", "api.interpreter.caiyunai.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,ca;q=0.6,zh-TW;q=0.5")
	req.Header.Set("app-name", "xy")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("os-type", "web")
	req.Header.Set("os-version", "")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return err
	}
	if resp.StatusCode != 200 {
		log.Fatal("Bad Code: ", resp.StatusCode, "body", string(bodyText))
		return err
	}

	var dictResponse DictResponse

	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
		return err
	}

	rsp.Explanations = dictResponse.Dictionary.Explanations
	return nil
}
