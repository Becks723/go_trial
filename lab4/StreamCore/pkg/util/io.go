package util

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
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

func IsValidImage(raw []byte) bool {
	t := http.DetectContentType(raw[:512])
	if !strings.HasPrefix(t, "image/") {
		return false
	}
	return true
}

func IsValidVideo(raw []byte) bool {
	t := http.DetectContentType(raw[:512])
	if !strings.HasPrefix(t, "video/") {
		return false
	}
	return true
}

func SaveFile(raw []byte, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0750)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, raw, 0644)
}

const (
	Mb = 1 << 20
)

func ToMb(bytes int) float64 {
	if bytes < 0 {
		return 0
	}
	return float64(bytes) / float64(Mb)
}

func ToByte(mb float64) int {
	if mb < 0 {
		return 0
	}
	return int(mb * Mb)
}

func ReadRequiredFormFile(c *app.RequestContext, key string) ([]byte, error) {
	fh, err := c.FormFile(key)
	if err != nil {
		return nil, err
	}
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadOptionalFormFile may return (nil, nil)
func ReadOptionalFormFile(c *app.RequestContext, key string) ([]byte, error) {
	fh, err := c.FormFile(key)
	if errors.Is(err, protocol.ErrMissingFile) { // client did not provide a file
		return nil, nil
	}
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
