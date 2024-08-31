package qr

import (
	"image"
	"image/color"
)

type Settable interface {
	Set(x, y int, c color.Color)
}

func Square(i image.Image) image.Image {
	size := i.Bounds().Max.X
	if i.Bounds().Max.Y > size {
		size = i.Bounds().Max.Y
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	halfX := (size - i.Bounds().Max.X) / 2
	halfY := (size - i.Bounds().Max.Y) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if y > halfY && y < i.Bounds().Max.Y+halfY && x > halfX && x < i.Bounds().Max.X+halfX {
				img.Set(x, y, i.At(x-halfX, y-halfY))
			}
		}
	}
	return img
}

func Background(i image.Image, bgc color.Color) image.Image {
	g, ok := i.(Settable)
	if !ok {
		return i
	}

	alphaOffset := uint32(2000)

	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			col := i.At(x, y)
			_, _, _, alpha := col.RGBA()
			if alpha < alphaOffset {
				col = bgc
			}
			g.Set(x, y, col)
		}
	}
	return i
}

func AddMargin(i image.Image, size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, i.Bounds().Max.X+size*2, i.Bounds().Max.Y+size*2))
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			if y > size && y < i.Bounds().Max.Y+size && x > size && x < i.Bounds().Max.X+size {
				img.Set(x, y, i.At(x-size, y-size))
			}
		}
	}
	return img
}

func Resize(src image.Image, size int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			dst.Set(x, y, src.At(x*src.Bounds().Max.X/size, y*src.Bounds().Max.Y/size))
		}
	}
	return dst
}
