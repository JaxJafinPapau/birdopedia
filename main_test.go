package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	s "strings"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	// creates an HTTP handler from our main package
	hf := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. This line executes the handler.
	hf.ServeHTTP(recorder, req)

	// Create error message for status code to check assertions

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Hello, World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned wrong body: got %v want %v", actual, expected)
	}
}

// The recorder is like a mini-browser, it accepts the result of the request
func TestRouter(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/hello")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected := "Hello, World!"

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

// sad path
func TestRouterForBadMethodOnExistingRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}
func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/fuzzywuzzywuzabear", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Status should be 404, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	received := s.Replace(respString, "\n", "", -1)
	expected := "404 page not found"

	if received != expected {
		t.Errorf("Response should be %s, got %s", expected, received)
	}
}
func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}

func TestGetBirdsHandler(t *testing.T) {
	mockStore := InitMockStore()

	mockStore.On("GetBirds").Return([]*Bird{
		{"sparrow", "a small harmless bird"},
	}, nil).Once()
	req, err := http.NewRequest("GET", "/birds", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getBirdsHandler)

	hf.ServeHTTP(recorder, req)

	status := recorder.Code

	if status != http.StatusOK {
		t.Errorf("Status should be 200, got %d", status)
	}

	expected := Bird{"sparrow", "a small harmless bird"}
	birdList := []Bird{}
	err = json.NewDecoder(recorder.Body).Decode(&birdList)

	if err != nil {
		t.Fatal(err)
	}

	actual := birdList[0]

	if actual != expected {
		t.Errorf("unexpected body, got %v want %v", actual, expected)
	}

	mockStore.AssertExpectations(t)
}

func TestCreateBirdHandler(t *testing.T) {
	mockStore := InitMockStore()

	mockStore.On("CreateBird", &Bird{"eagle", "a sharp eyed bird of prey"}).Return(nil)

	form := newCreateBirdForm()
	req, err := http.NewRequest("POST", "/birds", bytes.NewBufferString(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	// body := strings.NewReader("species=dodo&description=a dumb extinct bird")
	// req, err := http.NewRequest("POST", "/bird", body)
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createBirdHandler)

	hf.ServeHTTP(recorder, req)

	status := recorder.Code

	if status != http.StatusFound {
		t.Errorf("Status should be %d, got %d", http.StatusOK, status)
	}
	mockStore.AssertExpectations(t)
}

func newCreateBirdForm() *url.Values {
	form := url.Values{}
	form.Set("species", "eagle")
	form.Set("description", "a sharp eyed bird of prey")
	return &form
}
