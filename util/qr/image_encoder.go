package qr

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type encoder struct {
	encoder encoderFn
	writer  io.Writer
}

func (e encoder) Encode(img image.Image) error {
	return e.encoder.Encode(e.writer, img)
}

func PngEncoder(w io.Writer) ImageEncoder {
	enc := png.Encoder{CompressionLevel: png.BestCompression}

	return &encoder{
		writer: w,
		encoder: func(w io.Writer, img image.Image) error {
			return enc.Encode(w, img)
		},
	}
}

func JpegEncoder(w io.Writer) ImageEncoder {
	return &encoder{
		writer: w,
		encoder: func(w io.Writer, img image.Image) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: 80})
		},
	}
}

type encoderFn func(io.Writer, image.Image) error

func (ef encoderFn) Encode(w io.Writer, img image.Image) error {
	return ef(w, img)
}
