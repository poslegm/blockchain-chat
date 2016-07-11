package server

import (
	"testing"
	"net/http"
	"io/ioutil"
	"net/http/httptest"
	"path"
	"bytes"
)

func TestFileServer(t *testing.T) {
	testDir := "../client"

	homeHandle := createIndexHandler(testDir)
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