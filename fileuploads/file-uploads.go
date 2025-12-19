package fileuploads

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type FileUpload struct {
	File multipart.File
	info *multipart.FileHeader

	FileName      string
	Size          int64
	Header        textproto.MIMEHeader
	Ext           string
	SavedFile     string
	SavedFilePath string
}

type UploadOptions struct {
	MaxFileSize      int64
	RandomStringSize int
}

type UploadOption func(o *UploadOptions)

func ReadUploadedFile(fieldName string, r *http.Request, callback func(file FileUpload, options *UploadOptions) error, options ...UploadOption) error {
	var (
		err          error
		uploadedFile multipart.File
		info         *multipart.FileHeader
	)

	opts := &UploadOptions{
		MaxFileSize:      10 << 20,
		RandomStringSize: 10,
	}

	for _, opt := range options {
		opt(opts)
	}

	if err = r.ParseMultipartForm(opts.MaxFileSize); err != nil {
		return fmt.Errorf("error parsing form data: %w", err)
	}

	if uploadedFile, info, err = r.FormFile(fieldName); err != nil {
		return fmt.Errorf("error retrieving file info from form: %w", err)
	}

	defer uploadedFile.Close()

	if info.Size > opts.MaxFileSize {
		return fmt.Errorf("file size of %d bytes exceeds the limit of %d bytes", info.Size, opts.MaxFileSize)
	}

	callbackData := FileUpload{
		File: uploadedFile,
		info: info,
	}

	return callback(callbackData, opts)
}

func UploadFileToDir(fieldName string, r *http.Request, destDir string, destFileNameTemplate string, options ...UploadOption) (FileUpload, error) {
	result := &FileUpload{}

	err := ReadUploadedFile(fieldName, r, func(file FileUpload, options *UploadOptions) error {
		var (
			err      error
			tt       *template.Template
			tempFile *os.File
		)

		result = &file

		/*
		 * Craft a destination file name using the provided template
		 */
		destFileName := strings.Builder{}

		templateData := map[string]any{
			"fileName":     file.info.Filename,
			"baseFileName": filepath.Base(file.info.Filename),
			"ext":          filepath.Ext(file.info.Filename),
			"size":         file.info.Size,
			"randomString": randomString(options.RandomStringSize),
		}

		if tt, err = template.New("fileupload").Option("missingkey=error").Parse(destFileNameTemplate); err != nil {
			return fmt.Errorf("error parsing destination file name template: %w", err)
		}

		if err = tt.Execute(&destFileName, templateData); err != nil {
			return fmt.Errorf("error executing destination file name template: %w", err)
		}

		/*
		 * Move the uploaded file to the destination directory
		 */
		finalPath := filepath.Join(destDir, destFileName.String())

		// Security check to ensure the file path is within the destination directory
		if !strings.HasPrefix(finalPath, filepath.Clean(destDir)) {
			return fmt.Errorf("invalid destination file path attempted: %s", destFileName.String())
		}

		if tempFile, err = os.Create(finalPath); err != nil {
			return fmt.Errorf("error creating file '%s': %w", destFileName.String(), err)
		}

		defer tempFile.Close()

		if _, err = io.Copy(tempFile, file.File); err != nil {
			tempFile.Close()
			return fmt.Errorf("error copying uploaded file '%s' to temporary file '%s': %w", file.info.Filename, destFileName.String(), err)
		}

		if absolutePath, err := filepath.Abs(tempFile.Name()); err == nil {
			result.SavedFilePath = absolutePath
		}

		result.Ext = filepath.Ext(tempFile.Name())
		result.SavedFile = tempFile.Name()

		return nil
	}, options...)

	return *result, err
}

func WithMaxFileSize(size int64) UploadOption {
	return func(o *UploadOptions) {
		o.MaxFileSize = size
	}
}

func WithRandomStringSize(size int) UploadOption {
	return func(o *UploadOptions) {
		o.RandomStringSize = size
	}
}

func randomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for index := range b {
		b[index] = characters[seed.Intn(len(characters))]
	}

	return string(b)
}
