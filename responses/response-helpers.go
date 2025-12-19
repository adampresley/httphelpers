package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Html writes content to the response writer with a text/html header.
*/
func Html(w http.ResponseWriter, status int, value any) {
	write(w, "text/html", status, value)
}

/*
HtmlOK is a convenience wrapper to send a 200 with an arbitrary HTML body
*/
func HtmlOK(w http.ResponseWriter, value any) {
	Html(w, http.StatusOK, value)
}

/*
IsSuccessRange returns true if the status code falls within the 200-299 range.
*/
func IsSuccessRange(status int) bool {
	return status >= 200 && status < 300
}

/*
Json writes JSON content to the response writer. If there is an error
marshalling the value, it writes a 500 status code with a generic error message.
*/
func Json(w http.ResponseWriter, status int, value any) {
	var (
		err error
		b   []byte
	)

	w.Header().Set("Content-Type", "application/json")

	if b, err = json.Marshal(value); err != nil {
		b, _ = json.Marshal(struct {
			Message    string `json:"message"`
			Suggestion string `json:"suggestion"`
		}{
			Message:    "Error marshaling value for writing",
			Suggestion: "See error log for more information",
		})

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%s", string(b))
		return
	}

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%s", string(b))
}

/*
JsonOK is a convenience wrapper to send a 200 with an
arbitrary structure.
*/
func JsonOK(w http.ResponseWriter, value any) {
	Json(w, http.StatusOK, value)
}

/*
JsonBadRequest is a convenience wrapper to send a 400 with an
arbitrary structure.
*/
func JsonBadRequest(w http.ResponseWriter, value any) {
	Json(w, http.StatusBadRequest, value)
}

/*
JsonInternalServerError is a convenience wrapper to send a 500 with an
arbitrary structure.
*/
func JsonInternalServerError(w http.ResponseWriter, value any) {
	Json(w, http.StatusInternalServerError, value)
}

/*
JsonErrorMessage is a convenience wrapper to send a JSON body with
the specified status code, and a body that looks like this:

	{"message": "<message> goes here"}
*/
func JsonErrorMessage(w http.ResponseWriter, status int, message string) {
	result := make(map[string]string)
	result["message"] = message

	Json(w, status, result)
}

/*
JsonUnauthorized is a convenience wrapper to send a 401 with an
arbitrary value.
*/
func JsonUnauthorized(w http.ResponseWriter, value any) {
	Json(w, http.StatusUnauthorized, value)
}

/*
Text writes content to the response writer with a text/plain header.
*/
func Text(w http.ResponseWriter, status int, value any) {
	write(w, "text/plain", status, value)
}

/*
TextOK is a convenience wrapper to send a 200 with an
arbitrary text body.
*/
func TextOK(w http.ResponseWriter, value any) {
	Text(w, http.StatusOK, value)
}

/*
TextBadRequest is a convenience wrapper to send a 400 with an
arbitrary text body.
*/
func TextBadRequest(w http.ResponseWriter, value any) {
	Text(w, http.StatusBadRequest, value)
}

/*
TextBadRequest is a convenience wrapper to send a 500 with an
arbitrary text body.
*/
func TextInternalServerError(w http.ResponseWriter, value any) {
	Text(w, http.StatusInternalServerError, value)
}

/*
TextBadRequest is a convenience wrapper to send a 401 with an
arbitrary text body.
*/
func TextUnauthorized(w http.ResponseWriter, value any) {
	Text(w, http.StatusUnauthorized, value)
}

func write(w http.ResponseWriter, contentType string, status int, value any) {
	w.Header().Set("Content-Type", contentType)

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%v", value)
}
