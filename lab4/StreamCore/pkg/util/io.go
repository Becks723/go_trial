package util

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func mime(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	mime := http.DetectContentType(buf)
	return mime, nil
}
func IsValidImage(fileHeader *multipart.FileHeader) bool {
	s, err := mime(fileHeader)
	if err != nil {
		return false
	}
	return strings.HasPrefix(s, "image/")
}

func IsValidVideo(fileHeader *multipart.FileHeader) bool {
	s, err := mime(fileHeader)
	if err != nil {
		return false
	}
	return strings.HasPrefix(s, "video/")
}

func SaveFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	err = os.MkdirAll(filepath.Dir(dst), 0750)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

const (
	Mb = 1 << 20
)

func ToMb(bytes int64) float64 {
	if bytes < 0 {
		return 0
	}
	return float64(bytes) / float64(Mb)
}

func ToByte(mb float64) int64 {
	if mb < 0 {
		return 0
	}
	return int64(mb * Mb)
}
