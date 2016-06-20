package s3gof3r

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetterRetryRequest(t *testing.T) {
	expectedTries := 3
	actualTries := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualTries++

		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if actualTries < expectedTries {
			http.Error(w, "InternalError", http.StatusInternalServerError)
		} else {
			http.Error(w, "ok", http.StatusOK)
		}
	}))

	r, _, err := b.GetReader("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	g, ok := r.(*getter)
	if !ok {
		t.Fatal("getter type cast failed")
	}
	_, err = g.retryRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err.Error())
	}
	if actualTries != expectedTries {
		t.Fatalf("Expected %d tries, got: %d", expectedTries, actualTries)
	}
}
