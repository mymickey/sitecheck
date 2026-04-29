package main

import (
	"bytes"
	"image"
	"image/png"

	"golang.org/x/image/draw"
)

func RenderTrayIcon(logo []byte) ([]byte, error) {
	const size = 32
	canvas := image.NewRGBA(image.Rect(0, 0, size, size))

	if len(logo) > 0 {
		source, err := png.Decode(bytes.NewReader(logo))
		if err != nil {
			return nil, err
		}
		draw.ApproxBiLinear.Scale(canvas, image.Rect(2, 2, size-2, size-2), source, source.Bounds(), draw.Over, nil)
	}

	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}


