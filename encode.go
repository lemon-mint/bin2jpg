package bin2jpg

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/andybalholm/brotli"
)

func ImageEncode(data []byte) image.Image {

	var compressed bytes.Buffer
	compressed.Write([]byte{0, 0, 0, 0})
	w := brotli.NewWriterLevel(&compressed, brotli.BestCompression)
	w.Write(data)
	w.Close()
	data = compressed.Bytes()
	binary.LittleEndian.PutUint32(data[:4], uint32(len(data)-4))

	var width = int(math.RoundToEven(math.Sqrt(float64(len(data)*8*5))/5.0)) * 5
	if width < 10 {
		width = 30
	}

	img := image.NewGray(image.Rect(0, 0, width, (len(data)*8/width+1)*5+3))
	var x, y int
	bound := img.Bounds()
	bound_x := bound.Max.X

	for i, b := range data {
		bit := byte(b)
		for j := 0; j < 8; j++ {
			for k := 0; k < 5; k++ {
				if x%bound_x == 0 && i != 0 {
					y += 1
					x = 0
				}
				if (bit & (1 << uint(j))) != 0 {
					img.SetGray(x, y, color.Gray{255})
				} else {
					img.SetGray(x, y, color.Gray{0})
				}
				x++
			}
		}
	}

	return img
}

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r < 65536/2 && g < 65536/2 && b < 65536/2
}

func ImageDecode(img image.Image) ([]byte, error) {
	var data []byte
	bound := img.Bounds()
	bound_x := bound.Max.X
	var x, y int
	y = -1

	readByte := func() byte {
		var b byte

		for i := 0; i < 8; i++ {
			var black, white uint8
			for j := 0; j < 5; j++ {
				if x%bound_x == 0 {
					y += 1
					x = 0
				}
				//println("x:", x, "y:", y)
				if isBlack(img.At(x, y)) {
					black++
				} else {
					white++
				}
				x++
			}
			if black < white {
				b |= 1 << uint(i)
			}
		}
		return b
	}

	var length [4]byte = [4]byte{readByte(), readByte(), readByte(), readByte()}
	l := binary.LittleEndian.Uint32(length[:])

	for i := 0; i < int(l); i++ {
		data = append(data, readByte())
	}

	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}
