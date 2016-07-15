package server

import (
	"testing"
	"net/http"
	"io/ioutil"
	"net/http/httptest"
	"path"
	"bytes"
	"fmt"
	"time"
)

func TestFileServer(t *testing.T) {
	testDir := "../client"
	testPort := "8081"

	testHomePage(t, testDir)
	testAssets(t, testDir, testPort)
}

func testHomePage(t *testing.T, testDir string) {
	homeHandle := createPageHandler(testDir)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	actual, _ := ioutil.ReadAll(w.Body)
	index, _ := ioutil.ReadFile(path.Join(testDir, "index.html"))

	if bytes.Compare(actual, index) != 0 {
		t.Errorf("Home page didn't return index file")
	}
}

func testAssets(t *testing.T, testDir string, testPort string) {
	go Run(testDir, testPort)

	time.Sleep(8000)

	response, err := http.Get("http://127.0.0.1:" + testPort + "/js/index.js")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, _ := ioutil.ReadAll(response.Body)
	index, _ := ioutil.ReadFile(path.Join(testDir, "js", "index.js"))

	if bytes.Compare(actual, index) != 0 {
		t.Errorf("Assets didn't return requested file")

		fmt.Println("ACTUAL")
		fmt.Println(string(actual))
		fmt.Println("============================================")
		fmt.Println("REQUESTED")
		fmt.Println(string(index))
	}
}