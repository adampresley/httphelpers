package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type jsonTestPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestHtml(t *testing.T) {
	w := httptest.NewRecorder()
	status := http.StatusAccepted
	body := "<p>Hello World</p>"

	Html(w, status, body)

	if w.Code != status {
		t.Errorf("Expected status code %d, got %d", status, w.Code)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/html" {
		t.Errorf("Expected Content-Type 'text/html', got '%s'", contentType)
	}

	if w.Body.String() != body {
		t.Errorf("Expected body '%s', got '%s'", body, w.Body.String())
	}
}

func TestHtmlOK(t *testing.T) {
	w := httptest.NewRecorder()
	body := "<h1>OK</h1>"

	HtmlOK(w, body)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/html" {
		t.Errorf("Expected Content-Type 'text/html', got '%s'", contentType)
	}
}

func TestIsSuccessRange(t *testing.T) {
	testCases := []struct {
		status   int
		expected bool
	}{
		{199, false},
		{200, true},
		{299, true},
		{300, false},
		{404, false},
		{500, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Status%d", tc.status), func(t *testing.T) {
			if got := IsSuccessRange(tc.status); got != tc.expected {
				t.Errorf("IsSuccessRange(%d) = %v; want %v", tc.status, got, tc.expected)
			}
		})
	}
}

func TestJson(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload := jsonTestPayload{Name: "Test", Age: 10}
		expectedBody, _ := json.Marshal(payload)

		Json(w, http.StatusOK, payload)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		if w.Body.String() != string(expectedBody) {
			t.Errorf("Expected body '%s', got '%s'", string(expectedBody), w.Body.String())
		}
	})

	t.Run("SuccessWithStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload := jsonTestPayload{Name: "Created", Age: 1}

		Json(w, http.StatusCreated, payload)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("MarshallingError", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Channels cannot be marshalled to JSON
		payload := make(chan int)

		Json(w, http.StatusOK, payload)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d on marshal error, got %d", http.StatusInternalServerError, w.Code)
		}

		if !strings.Contains(w.Body.String(), "Error marshaling value") {
			t.Errorf("Expected error message in body, got: %s", w.Body.String())
		}
	})
}

func TestJsonConvenience(t *testing.T) {
	payload := jsonTestPayload{Name: "Test", Age: 99}

	testCases := []struct {
		name           string
		handlerFunc    func(http.ResponseWriter, any)
		expectedStatus int
	}{
		{"JsonOK", JsonOK, http.StatusOK},
		{"JsonBadRequest", JsonBadRequest, http.StatusBadRequest},
		{"JsonInternalServerError", JsonInternalServerError, http.StatusInternalServerError},
		{"JsonUnauthorized", JsonUnauthorized, http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.handlerFunc(w, payload)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), payload.Name) {
				t.Errorf("Expected body to contain payload, but it did not. Body: %s", w.Body.String())
			}
		})
	}
}

func TestJsonErrorMessage(t *testing.T) {
	var body map[string]string
	w := httptest.NewRecorder()
	status := http.StatusNotFound
	message := "The requested resource was not found"

	JsonErrorMessage(w, status, message)

	if w.Code != status {
		t.Errorf("Expected status code %d, got %d", status, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	expected := map[string]string{"message": message}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Expected body %v, got %v", expected, body)
	}
}

func TestText(t *testing.T) {
	w := httptest.NewRecorder()
	status := http.StatusAccepted
	body := "Hello World"

	Text(w, status, body)

	if w.Code != status {
		t.Errorf("Expected status code %d, got %d", status, w.Code)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", contentType)
	}

	if w.Body.String() != body {
		t.Errorf("Expected body '%s', got '%s'", body, w.Body.String())
	}
}

func TestTextConvenience(t *testing.T) {
	body := "some text"

	testCases := []struct {
		name           string
		handlerFunc    func(http.ResponseWriter, any)
		expectedStatus int
	}{
		{"TextOK", TextOK, http.StatusOK},
		{"TextBadRequest", TextBadRequest, http.StatusBadRequest},
		{"TextInternalServerError", TextInternalServerError, http.StatusInternalServerError},
		{"TextUnauthorized", TextUnauthorized, http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.handlerFunc(w, body)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != body {
				t.Errorf("Expected body '%s', got '%s'", body, w.Body.String())
			}
		})
	}
}
