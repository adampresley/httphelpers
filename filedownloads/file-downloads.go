package filedownloads

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
DownloadCsv downloads a byte slice to the provided HTTP response writer as a CSV file.
*/
func DownloadCsv(w http.ResponseWriter, filename string, content []byte) error {
	return StreamBytes(w, filename, "text/csv", content)
}

/*
DownloadCsvFile downloads an os.File to the provided HTTP response writer as a CSV file.
*/
func DownloadCsvFile(w http.ResponseWriter, filename string, file *os.File) error {
	return DownloadFile(w, filename, "text/csv", file)
}

/*
DownloadFile downloads an os.File to the provided HTTP response writer.
*/
func DownloadFile(w http.ResponseWriter, filename, contentType string, file *os.File) error {
	var (
		err      error
		fileInfo os.FileInfo
	)

	if fileInfo, err = file.Stat(); err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	return StreamContent(w, filename, contentType, file, fileInfo.Size())
}

/*
StreamBytes writes a byte slice to the provided HTTP response writer.
*/
func StreamBytes(w http.ResponseWriter, filename, contentType string, content []byte) error {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))

	_, err := w.Write(content)
	return err
}

/*
StreamContent writes an io.Reader to the provided HTTP response writer.
*/
func StreamContent(w http.ResponseWriter, filename, contentType string, content io.Reader, size int64) error {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))

	_, err := io.Copy(w, content)
	return err
}
