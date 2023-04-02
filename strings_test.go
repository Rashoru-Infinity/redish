package redish

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	redishOptions map[string]string
	redisClientOptions map[string]interface {}
}

func getTestCase() (*TestCase, error) {
	var config TestCase
	var objmap map[string]interface {}
	configFile, err := os.Open("./testcase.json")
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &objmap); err != nil {
		return nil, err
	}
	b, err = json.Marshal(objmap["Redish-Options"].(map[string]interface {}))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &config.redishOptions)
	if err != nil {
		return nil, err
	}
	b, err = json.Marshal(objmap["Redis-Client-Options"].(map[string]interface {}))
	if err != nil {
		return nil, err
	}
	json.Unmarshal(b, &config.redisClientOptions)
	return &config, nil
}

func startServer(ch chan error) {
	var config *TestCase
	config, err := getTestCase()
	if err != nil {
		ch <- err
		return
	}
	port, ok := config.redishOptions["Port"]
	if !ok {
		port = "8080"
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(204), 204)
	})
	http.HandleFunc("/strings/", HandleStrings)
	go func () {
		ch <- http.ListenAndServe(":" + port, nil)
	}()
	req, _ := http.NewRequest("GET", "http://localhost:" +
	port + "/health", nil)
	resp, err := new(http.Client).Do(req)
	if err != nil {
		ch <- err
		return
	}
	if resp.StatusCode != 204 {
		ch <- errors.New(fmt.Sprintf(
			"unexpected HTTP status code %d %s\n",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		))
	}
	ch <- nil
}

func TestInitRedish(t *testing.T) {
	ch := make(chan error)
	go startServer(ch)
	err := <- ch
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
}

func Test1stPostStrings(t *testing.T) {
	var testcase *TestCase
	testcase, err := getTestCase()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if testcase == nil {
		t.Errorf("%s\n", "invalid testcase")
		return
	}
	req, _ := http.NewRequest("POST", "http://localhost:" +
	testcase.redishOptions["Port"] + "/strings/key", strings.NewReader("{\"key\": \"value\"}"))
	b, err := json.Marshal(&testcase.redisClientOptions)	
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	req.Header.Set("Redis-Client-Options", string(b))
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if resp.StatusCode != 201 {
		t.Errorf("unexpected HTTP status code %d %s\n",
			resp.StatusCode, 
			http.StatusText(resp.StatusCode),
		)
		return
	}
}

func Test2ndPostStrings(t *testing.T) {
	var testcase *TestCase
	testcase, err := getTestCase()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if testcase == nil {
		t.Errorf("%s\n", "invalid testcase")
		return
	}
	req, _ := http.NewRequest("POST", "http://localhost:" +
	testcase.redishOptions["Port"] + "/strings/key", strings.NewReader("{\"key\": \"value\"}"))
	b, err := json.Marshal(&testcase.redisClientOptions)	
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	req.Header.Set("Redis-Client-Options", string(b))
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if resp.StatusCode != 405 {
		t.Errorf("unexpected HTTP status code %d %s\n",
			resp.StatusCode, 
			http.StatusText(resp.StatusCode),
		)
		return
	}
}

func TestGetStrings(t *testing.T) {
	var kv map[string]string
	var testcase *TestCase
	testcase, err := getTestCase()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if testcase == nil {
		t.Errorf("%s\n", "invalid testcase")
		return
	}
	req, _ := http.NewRequest("GET", "http://localhost:" +
	testcase.redishOptions["Port"] + "/strings/key", nil)
	b, err := json.Marshal(&testcase.redisClientOptions)	
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	req.Header.Set("Redis-Client-Options", string(b))
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected HTTP status code %d %s\n",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
		return
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if err = json.Unmarshal(b, &kv); err != nil {
		t.Errorf("%s\n", err)
		return
	}
	v, ok := kv["key"]
	if !ok {
		t.Errorf("\"key\" not found\n")
	}
	if v != "value" {
		t.Errorf("got %s, want %s\n", v, "value")
	}
}

func TestPutStrings(t *testing.T) {
	var testcase *TestCase
	testcase, err := getTestCase()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if testcase == nil {
		t.Errorf("%s\n", "invalid testcase")
		return
	}
	req, _ := http.NewRequest("PUT", "http://localhost:" +
	testcase.redishOptions["Port"] + "/strings/key", strings.NewReader("{\"key\": \"value\"}"))
	b, err := json.Marshal(&testcase.redisClientOptions)	
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	req.Header.Set("Redis-Client-Options", string(b))
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if resp.StatusCode != 204 {
		t.Errorf("unexpected HTTP status code %d %s\n",
			resp.StatusCode, 
			http.StatusText(resp.StatusCode),
		)
		return
	}
}

func TestDeleteStrings(t *testing.T) {
	var testcase *TestCase
	testcase, err := getTestCase()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if testcase == nil {
		t.Errorf("%s\n", "invalid testcase")
		return
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:" +
	testcase.redishOptions["Port"] + "/strings/key", nil)
	b, err := json.Marshal(&testcase.redisClientOptions)	
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	req.Header.Set("Redis-Client-Options", string(b))
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	if resp.StatusCode != 204 {
		t.Errorf("unexpected HTTP status code %d %s\n",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
		return
	}
}
