package fileuploads

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
createMultipartRequest is a helper function to create a new HTTP request
containing a multipart form with a file.
*/
func createMultipartRequest(t *testing.T, fieldName, fileName, fileContent string) *http.Request {
	t.Helper()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part, strings.NewReader(fileContent))
	if err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestReadUploadedFileSuccess(t *testing.T) {
	fieldName := "testfile"
	fileName := "test.txt"
	fileContent := "hello world"
	req := createMultipartRequest(t, fieldName, fileName, fileContent)

	err := ReadUploadedFile(fieldName, req, func(file FileUpload, options *UploadOptions) error {
		if file.info.Filename != fileName {
			t.Errorf("Expected filename '%s', got '%s'", fileName, file.info.Filename)
		}

		content, err := io.ReadAll(file.File)
		if err != nil {
			return fmt.Errorf("could not read file content from callback: %w", err)
		}

		if string(content) != fileContent {
			t.Errorf("Expected file content '%s', got '%s'", fileContent, string(content))
		}

		return nil
	})

	if err != nil {
		t.Fatalf("ReadUploadedFile failed: %v", err)
	}
}

func TestReadUploadedFileNoFile(t *testing.T) {
	req := createMultipartRequest(t, "wrongfield", "test.txt", "content")

	err := ReadUploadedFile("correctfield", req, func(file FileUpload, options *UploadOptions) error {
		t.Fatal("Callback was called when it should not have been")
		return nil
	})

	if err == nil {
		t.Fatal("Expected an error but got nil")
	}

	if !strings.Contains(err.Error(), "error retrieving file info from form") {
		t.Errorf("Expected error to be about retrieving file info, got: %v", err)
	}
}

func TestReadUploadedFileExceedsMaxSize(t *testing.T) {
	fieldName := "testfile"
	fileName := "large.txt"
	fileContent := "this content is definitely larger than 5 bytes"
	req := createMultipartRequest(t, fieldName, fileName, fileContent)

	err := ReadUploadedFile(fieldName, req, func(file FileUpload, options *UploadOptions) error {
		t.Fatal("Callback was called for a file that exceeds max size")
		return nil
	}, WithMaxFileSize(5))

	if err == nil {
		t.Fatal("Expected an error for exceeding max file size, but got nil")
	}
	if !strings.Contains(err.Error(), "exceeds the limit") {
		t.Errorf("Expected error to be about exceeding file size limit, got: %v", err)
	}
}

func TestUploadFileToDirSuccess(t *testing.T) {
	destDir := t.TempDir()
	fieldName := "upload"
	fileName := "data.csv"
	fileContent := "col1,col2"
	template := "file-{{.randomString}}-{{.baseFileName}}"

	req := createMultipartRequest(t, fieldName, fileName, fileContent)
	upload, err := UploadFileToDir(fieldName, req, destDir, template, WithRandomStringSize(5))

	if err != nil {
		t.Fatalf("UploadFileToDir failed: %v", err)
	}

	if upload.SavedFilePath == "" {
		t.Fatal("SavedFilePath was not set")
	}

	/*
	 * Check that the filename matches the template structure
	 */
	baseSavedName := filepath.Base(upload.SavedFilePath)
	if !strings.HasPrefix(baseSavedName, "file-") || !strings.HasSuffix(baseSavedName, "-"+fileName) {
		t.Errorf("Saved filename '%s' does not match template structure", baseSavedName)
	}

	savedContent, err := os.ReadFile(upload.SavedFilePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != fileContent {
		t.Errorf("Saved file content does not match. Expected '%s', got '%s'", fileContent, string(savedContent))
	}
}

func TestUploadFileToDirBadTemplate(t *testing.T) {
	destDir := t.TempDir()
	fieldName := "upload"
	fileName := "data.csv"
	fileContent := "col1,col2"
	badTemplate := "{{.badField}}"

	req := createMultipartRequest(t, fieldName, fileName, fileContent)
	_, err := UploadFileToDir(fieldName, req, destDir, badTemplate)

	if err == nil {
		t.Fatal("Expected a template execution error but got nil")
	}

	if !strings.Contains(err.Error(), "error executing destination file name template") {
		t.Errorf("Expected error about template execution, got: %v", err)
	}
}

func TestUploadFileToDirCreateFails(t *testing.T) {
	// Use a non-existent directory path
	destDir := filepath.Join(t.TempDir(), "non-existent-dir")
	fieldName := "upload"
	fileName := "data.csv"
	fileContent := "col1,col2"
	template := "{{.baseFileName}}"

	req := createMultipartRequest(t, fieldName, fileName, fileContent)
	_, err := UploadFileToDir(fieldName, req, destDir, template)

	if err == nil {
		t.Fatal("Expected a file creation error but got nil")
	}

	if !strings.Contains(err.Error(), "error creating file") {
		t.Errorf("Expected error about creating file, got: %v", err)
	}
}

func TestWithMaxFileSizeOption(t *testing.T) {
	opts := &UploadOptions{}
	size := int64(50 << 20) // 50MB
	opt := WithMaxFileSize(size)
	opt(opts)

	if opts.MaxFileSize != size {
		t.Errorf("Expected MaxFileSize to be %d, but got %d", size, opts.MaxFileSize)
	}
}

func TestWithRandomStringSizeOption(t *testing.T) {
	opts := &UploadOptions{}
	size := 15
	opt := WithRandomStringSize(size)
	opt(opts)

	if opts.RandomStringSize != size {
		t.Errorf("Expected RandomStringSize to be %d, but got %d", size, opts.RandomStringSize)
	}
}
