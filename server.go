package toggl

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

// TestJSONFileResponse test response
type TestJSONFileResponse struct {
	StatusCode int
	JSONPath   string
}

// TestCreateResponseHandler test handler
func TestCreateResponseHandler(res TestJSONFileResponse, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(pwd + res.JSONPath)
		src, err := ioutil.ReadFile(pwd + res.JSONPath)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(res.StatusCode)
		w.Write(src)
		return
	}
}
