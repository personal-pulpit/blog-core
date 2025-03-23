package file

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"

	"github.com/chai2010/webp"
)

const (
	MaxFileSize = 1024 * 1024 * 3 // 3MB
)

var (
	ErrFileTooLarge = errors.New("file size is too large")
)

func ConvertToWebP(imageData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	var webpBuf bytes.Buffer
	err = webp.Encode(&webpBuf, img, &webp.Options{Quality: 80}) // Adjust quality as needed
	if err != nil {
		return nil, err
	}

	return webpBuf.Bytes(), nil
}

func CheckFile(file *multipart.FileHeader) ([]byte, error) {
	if file.Size > MaxFileSize {
		return nil, errors.New("file too large")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	webpData, err := ConvertToWebP(data)
	if err != nil {
		return nil, err
	}

	return webpData, nil
}

func GenerateFileName(ID uint) string {

	return fmt.Sprintf("user%d.webp", ID)
}
