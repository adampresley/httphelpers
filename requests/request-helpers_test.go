package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type BodyTestPayload struct {
	Name string `json:"name" xml:"name"`
	Age  int    `json:"age" xml:"age"`
}

func TestGetAuthorizationBearer(t *testing.T) {
	testCases := []struct {
		name          string
		header        string
		expectedToken string
		expectError   bool
	}{
		{"Valid", "Bearer my-secret-token", "my-secret-token", false},
		{"ValidWithSpaces", "  Bearer   my-secret-token  ", "my-secret-token", false},
		{"ValidLowercase", "bearer my-secret-token", "my-secret-token", false},
		{"NoHeader", "", "", true},
		{"MalformedShort", "Bearer", "", true},
		{"MalformedNoSpace", "Bearertoken", "", true},
		{"NotBearerScheme", "Basic my-secret-token", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			token, err := AuthorizationBearer(req)

			if tc.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect an error, but got: %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token '%s', but got '%s'", tc.expectedToken, token)
				}
			}
		})
	}
}

func TestGetFromRequest(t *testing.T) {
	form := url.Values{}
	form.Add("string", "form_string")
	form.Add("int", "123")
	form.Add("bool", "true")
	form.Add("slice", "form_slice1")
	form.Add("slice", "form_slice2")

	req, _ := http.NewRequest("POST", "/path_value?string=query_string&int=456", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("string", "path_string")
	req.SetPathValue("int_path_only", "789")

	if err := req.ParseForm(); err != nil {
		t.Fatalf("Failed to parse form: %v", err)
	}

	t.Run("StringPrecedence", func(t *testing.T) {
		got := Get[string](req, "string")

		if got != "form_string" {
			t.Errorf("Expected 'form_string' for string precedence, but got '%s'", got)
		}
	})

	t.Run("PathValue", func(t *testing.T) {
		got := Get[int](req, "int_path_only")

		if got != 789 {
			t.Errorf("Expected 789 for path-only int, but got %d", got)
		}
	})

	t.Run("Int", func(t *testing.T) {
		got := Get[int](req, "int")

		if got != 123 {
			t.Errorf("Expected 123 for int, but got %d", got)
		}
	})

	t.Run("Bool", func(t *testing.T) {
		got := Get[bool](req, "bool")

		if !got {
			t.Error("Expected true for bool, but got false")
		}
	})

	t.Run("StringSlice", func(t *testing.T) {
		expected := []string{"form_slice1", "form_slice2"}
		got := Get[[]string](req, "slice")

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected string slice %v, but got %v", expected, got)
		}
	})

	t.Run("ZeroValue", func(t *testing.T) {
		gotString := Get[string](req, "nonexistent")
		if gotString != "" {
			t.Errorf("Expected zero value '' for string, but got '%s'", gotString)
		}

		gotInt := Get[int](req, "nonexistent")
		if gotInt != 0 {
			t.Errorf("Expected zero value 0 for int, but got %d", gotInt)
		}
	})
}

func TestGetStringListFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/?list=a,b,c", nil)
	expected := []string{"a", "b", "c"}
	got := StringListFromRequest(req, "list", ",")

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, but got %v", expected, got)
	}
}

func TestIsHtmx(t *testing.T) {
	t.Run("IsHtmxTrue", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Hx-Request", "true")

		if !IsHtmx(req) {
			t.Error("Expected IsHtmx to be true, but it was false")
		}
	})

	t.Run("IsHtmxFalse", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		if IsHtmx(req) {
			t.Error("Expected IsHtmx to be false, but it was true")
		}
	})
}

func TestReadBody(t *testing.T) {
	payload := BodyTestPayload{Name: "Adam", Age: 30}

	t.Run("JSON", func(t *testing.T) {
		jsonData, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		result, err := Body[BodyTestPayload](req)
		if err != nil {
			t.Fatalf("ReadBody failed for JSON: %v", err)
		}

		if !reflect.DeepEqual(result, payload) {
			t.Errorf("Expected %v, got %v for JSON", payload, result)
		}
	})

	t.Run("XML", func(t *testing.T) {
		xmlData, _ := xml.Marshal(payload)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(xmlData))
		req.Header.Set("Content-Type", "application/xml")

		result, err := Body[BodyTestPayload](req)
		if err != nil {
			t.Fatalf("ReadBody failed for XML: %v", err)
		}

		if !reflect.DeepEqual(result, payload) {
			t.Errorf("Expected %v, got %v for XML", payload, result)
		}
	})

	t.Run("BadJSON", func(t *testing.T) {
		badJSON := `{"name": "Adam", "age": }`
		req := httptest.NewRequest("POST", "/", strings.NewReader(badJSON))
		req.Header.Set("Content-Type", "application/json")

		_, err := Body[BodyTestPayload](req)
		if err == nil {
			t.Fatal("Expected an error for bad JSON, but got nil")
		}
	})

	t.Run("EmptyBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil) // Body is nil
		req.Body = io.NopCloser(strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		_, err := Body[BodyTestPayload](req)
		if err == nil {
			t.Fatal("Expected an error for empty body, but got nil")
		}
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader("text content"))
		req.Header.Set("Content-Type", "text/plain")

		_, err := Body[BodyTestPayload](req)
		if err == nil {
			t.Fatal("Expected an error for unsupported type, but got nil")
		}

		if !strings.Contains(err.Error(), "unsupported content type") {
			t.Errorf("Expected error to contain 'unsupported content type', got: %v", err)
		}
	})
}
