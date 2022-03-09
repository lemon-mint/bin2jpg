package bin2jpg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/andybalholm/brotli"
	"github.com/lemon-mint/vbox"
)

func ImageEncode(data []byte, key []byte) image.Image {
	var compressBuffer bytes.Buffer
	w := brotli.NewWriterLevel(&compressBuffer, brotli.BestCompression)
	w.Write(data)
	w.Flush()
	w.Close()

	compressed := compressBuffer.Bytes()

	if len(key) > 0 {
		box := vbox.NewBlackBox(key)
		compressed = box.Seal(compressed)
	}

	var length [4]byte
	binary.LittleEndian.PutUint32(length[:], uint32(len(compressed)))

	data = nil
	data = append(data, length[:]...)
	data = append(data, compressed...)

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

var ErrAEADOpenError = errors.New("bin2jpg: aead open error")

func ImageDecode(img image.Image, key []byte) ([]byte, error) {
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

	if len(key) > 0 {
		box := vbox.NewBlackBox(key)
		var ok bool
		data, ok = box.OpenOverWrite(data)
		if !ok {
			return nil, ErrAEADOpenError
		}
	}

	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}
