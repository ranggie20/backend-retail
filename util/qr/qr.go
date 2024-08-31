package qr

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

type ImageEncoder interface {
	Encode(image.Image) error
}

// GenerateQR
//
//		qrSize is resulting image width & height
//	 	data is the information that will be generated as QRCode
func GenerateQR(qrSize int, data string, encoder ImageEncoder) error {
	qc, err := qrcode.New(data, qrcode.Highest)
	if err != nil {
		return err
	}

	img := qc.Image(qrSize)

	err = encoder.Encode(img)
	return err
}

// GenerateQRWithLogo
func GenerateQRWithLogo(qrSize int, data string, logo image.Image, encoder ImageEncoder) error {
	qc, err := qrcode.New(data, qrcode.Highest)
	if err != nil {
		return err
	}

	img := qc.Image(qrSize)
	maxLogoSize := uint(qrSize) / 4
	logo = Background(AddMargin(Resize(Square(logo), 900), 62), color.White) // treatment for png with transparent background
	logo = resize.Resize(maxLogoSize, 0, logo, resize.Bilinear)
	overlay(img, logo)

	err = encoder.Encode(img)
	return err
}

func overlay(dst, logo image.Image) {
	dstImg, ok := dst.(Settable)
	if !ok {
		return
	}

	offsetX := (dst.Bounds().Max.X - logo.Bounds().Max.X) / 2
	offsetY := (dst.Bounds().Max.Y - logo.Bounds().Max.Y) / 2
	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			dstImg.Set(x+offsetX, y+offsetY, logo.At(x, y))
		}
	}
}
