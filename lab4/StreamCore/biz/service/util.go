package service

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func parseTIme(timestamp string) (t time.Time, err error) {
	unix, err := strconv.ParseUint(timestamp, 10, 64)
	if err != nil {
		return
	}
	t = time.UnixMilli(int64(unix))
	return
}

func isValidImage(fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return false
	}
	mime := http.DetectContentType(buf)
	if !strings.HasPrefix(mime, "image/") {
		return false
	}
	return true
}

func isValidVideo(fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 512)
	file.Read(buf)
	mimeType := http.DetectContentType(buf)
	if !strings.HasPrefix(mimeType, "video/") {
		return false
	}
	return true
}

func saveFile(fileHeader *multipart.FileHeader, dst string) error {
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
