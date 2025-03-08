package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

var allowedFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func ValidateFile(fileHeader *multipart.FileHeader) error {
	ext := filepath.Ext(fileHeader.Filename)
	if !allowedFormats[ext] {
		return errors.New("unsupported file format, only JPG and PNG are allowed")
	}
	return nil
}
