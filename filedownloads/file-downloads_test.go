package filedownloads

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDownloadCsv(t *testing.T) {
	w := httptest.NewRecorder()
	filename := "test.csv"
	content := []byte("col1,col2\nval1,val2")

	err := DownloadCsv(w, filename, content)
	if err != nil {
		t.Fatalf("DownloadCsv failed: %v", err)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/csv" {
		t.Errorf("Expected Content-Type 'text/csv', got '%s'", contentType)
	}

	expectedDisposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	if disposition := w.Header().Get("Content-Disposition"); disposition != expectedDisposition {
		t.Errorf("Expected Content-Disposition '%s', got '%s'", expectedDisposition, disposition)
	}

	expectedLength := fmt.Sprintf("%d", len(content))
	if length := w.Header().Get("Content-Length"); length != expectedLength {
		t.Errorf("Expected Content-Length '%s', got '%s'", expectedLength, length)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Errorf("Response body does not match content. Got '%s', expected '%s'", w.Body.String(), string(content))
	}
}

func TestDownloadCsvFile(t *testing.T) {
	content := []byte("csv,content")

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Reset the file pointer to the beginning
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	w := httptest.NewRecorder()
	filename := "download.csv"

	err = DownloadCsvFile(w, filename, tmpFile)
	if err != nil {
		t.Fatalf("DownloadCsvFile failed: %v", err)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "text/csv" {
		t.Errorf("Expected Content-Type 'text/csv', got '%s'", contentType)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Errorf("Response body does not match content. Got '%s', expected '%s'", w.Body.String(), string(content))
	}
}

func TestDownloadFile(t *testing.T) {
	content := []byte("file content")

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	w := httptest.NewRecorder()
	filename := "download.txt"
	contentType := "text/plain"

	err = DownloadFile(w, filename, contentType, tmpFile)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	if gotContentType := w.Header().Get("Content-Type"); gotContentType != contentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", contentType, gotContentType)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Errorf("Response body does not match content. Got '%s', expected '%s'", w.Body.String(), string(content))
	}
}

func TestDownloadFileStatFails(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Close the file so Stat will fail
	tmpFile.Close()

	w := httptest.NewRecorder()
	err = DownloadFile(w, "any.txt", "text/plain", tmpFile)

	if err == nil {
		t.Fatal("Expected an error from DownloadFile, but got nil")
	}

	if !strings.Contains(err.Error(), "error getting file info") {
		t.Errorf("Expected error to contain 'error getting file info', but got: %v", err)
	}
}

func TestStreamBytes(t *testing.T) {
	w := httptest.NewRecorder()
	filename := "test.bin"
	contentType := "application/octet-stream"
	content := []byte{1, 2, 3, 4, 5}

	err := StreamBytes(w, filename, contentType, content)
	if err != nil {
		t.Fatalf("StreamBytes failed: %v", err)
	}

	if gotContentType := w.Header().Get("Content-Type"); gotContentType != contentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", contentType, gotContentType)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Error("Response body does not match content")
	}
}

func TestStreamBytesWriteFails(t *testing.T) {
	w := &failingWriter{}
	content := []byte{1, 2, 3}

	err := StreamBytes(w, "test.bin", "application/octet-stream", content)
	if err == nil {
		t.Fatal("Expected an error from StreamBytes, but got nil")
	}

	if !strings.Contains(err.Error(), "simulated write error") {
		t.Errorf("Expected error to contain 'simulated write error', got: %v", err)
	}
}

func TestStreamContent(t *testing.T) {
	w := httptest.NewRecorder()
	filename := "test.txt"
	contentType := "text/plain"
	content := "hello world"
	contentReader := strings.NewReader(content)
	size := int64(len(content))

	err := StreamContent(w, filename, contentType, contentReader, size)
	if err != nil {
		t.Fatalf("StreamContent failed: %v", err)
	}

	if gotContentType := w.Header().Get("Content-Type"); gotContentType != contentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", contentType, gotContentType)
	}

	if w.Body.String() != content {
		t.Errorf("Response body does not match content. Got '%s', expected '%s'", w.Body.String(), content)
	}
}

func TestStreamContentReadFails(t *testing.T) {
	w := httptest.NewRecorder()
	contentReader := &failingReader{}

	err := StreamContent(w, "test.txt", "text/plain", contentReader, 100)
	if err == nil {
		t.Fatal("Expected an error from StreamContent, but got nil")
	}

	if !strings.Contains(err.Error(), "simulated read error") {
		t.Errorf("Expected error to contain 'simulated read error', got: %v", err)
	}
}

// failingReader is a helper to simulate errors when reading content.
type failingReader struct{}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

/*
failingWriter is a helper to simulate errors when writing to a response.
*/
type failingWriter struct {
	httptest.ResponseRecorder
}

func (fw *failingWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("simulated write error")
}

func (fw *failingWriter) ReadFrom(r io.Reader) (int64, error) {
	return 0, fmt.Errorf("simulated ReadFrom error")
}
