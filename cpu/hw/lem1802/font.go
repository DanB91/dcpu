// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package lem1802

import (
	"bytes"
	"errors"
	"github.com/jteeuwen/dcpu/cpu"
	"image"
	_ "image/png"
	"io"
)

// DefaultFont defines the default character set for the LEM1802 monitor.
// 128 characters at 2 words/character. 
func DefaultFont() []cpu.Word {
	b := bytes.NewBuffer(lem1802font_png)
	font, _ := LoadFont(b)
	return font
}

// LoadFont loads a PNG image from the given stream and interprets it as a font.
//
// The supplied image should be 128x32 pixels where each
// glyph occupies 4x8 pixels. This gives 128 glyphs.
// 
// The loader expects a black background with
// glyphs drawn in white. To be more precise, the glyphs need a red channel
// with a value > 0.
//
// Each word splits into two rows of eight bits.
// Giving a 4x8 grid for each character.
//
// For example, the character 'F' is encoded as follows:
//
//     word 1 = 0xfe12 = 1111111000010010
//     word 2 = 0x0200 = 0000001000000000
//
// When split into octets:
//
//     word 1 = 11111110 (0xfe)
//              00010010 (0x12)
//     word 2 = 00000010 (0x02)
//              00000000 (0x00)
func LoadFont(r io.Reader) (font []cpu.Word, err error) {
	var img image.Image

	img, _, err = image.Decode(r)
	if err != nil {
		return
	}

	bounds := img.Bounds()
	size := bounds.Size()

	if size.X != 128 || size.Y != 32 {
		return nil, errors.New("Invalid image size. Expected 128x32")
	}

	var index int
	var char [4]cpu.Word

	font = make([]cpu.Word, 256)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 8 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 4 {
			char[0] = translate(img, x+0, y)
			char[1] = translate(img, x+1, y)
			char[2] = translate(img, x+2, y)
			char[3] = translate(img, x+3, y)

			font[index+0] = char[0]<<8 | char[1]
			font[index+1] = char[2]<<8 | char[3]
			index += 2
		}
	}

	return
}

// translate retrieves the red component for each pixel in
// the given column. It then constructs an 8 bit value from
// the 8 red components.
func translate(img image.Image, x, y int) cpu.Word {
	var c [8]uint32
	c[0], _, _, _ = img.At(x, y+0).RGBA()
	c[1], _, _, _ = img.At(x, y+1).RGBA()
	c[2], _, _, _ = img.At(x, y+2).RGBA()
	c[3], _, _, _ = img.At(x, y+3).RGBA()
	c[4], _, _, _ = img.At(x, y+4).RGBA()
	c[5], _, _, _ = img.At(x, y+5).RGBA()
	c[6], _, _, _ = img.At(x, y+6).RGBA()
	c[7], _, _, _ = img.At(x, y+7).RGBA()

	return cpu.Word(((c[0]&1)<<0 | (c[1]&1)<<1 |
		(c[2]&1)<<2 | (c[3]&1)<<3 |
		(c[4]&1)<<4 | (c[5]&1)<<5 |
		(c[6]&1)<<6 | (c[7]&1)<<7) & 0xff)
}
