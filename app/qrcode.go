package app

import (
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	"image/draw"
	"log"
	"os"
)

// DecodeQRCodeFromImage decode image QR
func DecodeQRCodeFromImage(imagePath string) (string, error) {
	// open file image
	file, err := os.Open(imagePath)
	if err != nil {
		log.Println("Lỗi khi mở file:", err) // Log lỗi khi mở file
		return "", fmt.Errorf("lỗi khi mở file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Lỗi khi đóng file:", err) // Log lỗi khi đóng file
		}
	}(file)

	// decode image
	decodedImg, _, err := image.Decode(file)
	if err != nil {
		log.Println("Lỗi khi giải mã hình ảnh:", err) // Log lỗi khi giải mã hình ảnh
		return "", fmt.Errorf("lỗi khi giải mã hình ảnh: %v", err)
	}

	// convert *image.RGBA
	rgbaImg := image.NewRGBA(decodedImg.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), decodedImg, image.Point{}, draw.Src)

	// ready BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(rgbaImg)
	if err != nil {
		log.Println("Lỗi khi tạo BinaryBitmap:", err) // Log lỗi khi tạo BinaryBitmap
		return "", fmt.Errorf("lỗi khi tạo BinaryBitmap: %v", err)
	}

	// handler qrcode
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		log.Println("Lỗi khi giải mã QR code:", err) // Log lỗi khi giải mã QR code
		return "", fmt.Errorf("lỗi khi giải mã QR code: %v", err)
	}

	return result.String(), nil
}
