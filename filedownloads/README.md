# File Downloads

This package contains helpers for downloading files and content to an HTTP client.

## DownloadCsv

**DownloadCsv** downloads a byte slice to the provided HTTP response writer as a CSV file.

```go
// w is a http.ResponseWriter
// filename is "report.csv"
// content is a []byte slice
err := filedownloads.DownloadCsv(w, filename, content)
```

## DownloadCsvFile

**DownloadCsvFile** downloads a file from the file system to the provided HTTP response writer as a CSV file.

```go
// w is a http.ResponseWriter
// filename is "report.csv"
// file is an *os.File
err := filedownloads.DownloadCsvFile(w, filename, file)
```

## DownloadFile

**DownloadFile** downloads a file from the file system to the provided HTTP response writer. You must provide the content type.

```go
// w is a http.ResponseWriter
// filename is "image.png"
// contentType is "image/png"
// file is an *os.File
err := filedownloads.DownloadFile(w, filename, contentType, file)
```

## StreamBytes

**StreamBytes** writes a byte slice to the provided HTTP response writer. This is a generic function that sets the appropriate headers for downloading.

```go
// w is a http.ResponseWriter
// filename is "data.bin"
// contentType is "application/octet-stream"
// content is a []byte slice
err := filedownloads.StreamBytes(w, filename, contentType, content)
```

## StreamContent

**StreamContent** writes an `io.Reader` to the provided HTTP response writer. This is a generic function that sets the appropriate headers for downloading.

```go
// w is a http.ResponseWriter
// filename is "archive.zip"
// contentType is "application/zip"
// content is an io.Reader
// size is an int64
err := filedownloads.StreamContent(w, filename, contentType, content, size)
```
