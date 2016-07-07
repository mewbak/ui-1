// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package image

import (
	"github.com/richardwilkes/go-ui/color"
	"github.com/richardwilkes/go-ui/geom"
	"unsafe"
)

// #cgo darwin LDFLAGS: -framework Cocoa
// #include <stdlib.h>
// #include <CoreFoundation/CoreFoundation.h>
// #include <CoreGraphics/CoreGraphics.h>
// #include <ImageIO/ImageIO.h>
import "C"

func newImageFromBytes(buffer []byte) *Image {
	data := C.CFDataCreate(nil, (*C.UInt8)(&buffer[0]), C.CFIndex(len(buffer)))
	defer C.CFRelease(data)
	if imgSrc := C.CGImageSourceCreateWithData(data, nil); imgSrc != nil {
		defer C.CFRelease(imgSrc)
		if image := C.CGImageSourceCreateImageAtIndex(imgSrc, 0, nil); image != nil {
			return createImage(image)
		}
	}
	return nil
}

func newImageFromURL(url string) *Image {
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	cURLRef := C.CFURLCreateWithBytes(nil, (*C.UInt8)(unsafe.Pointer(cURL)), C.CFIndex(len(url)), C.kCFStringEncodingUTF8, nil)
	defer C.CFRelease(cURLRef)
	if imgSrc := C.CGImageSourceCreateWithURL(cURLRef, nil); imgSrc != nil {
		defer C.CFRelease(imgSrc)
		if image := C.CGImageSourceCreateImageAtIndex(imgSrc, 0, nil); image != nil {
			return createImage(image)
		}
	}
	return nil
}

func newImageFromData(data *Data) *Image {
	colorspace := C.CGColorSpaceCreateWithName(C.kCGColorSpaceGenericRGB)
	defer C.CGColorSpaceRelease(colorspace)

	length := len(data.Pixels)
	pixels := make([]color.Color, length)
	copy(pixels, data.Pixels)

	// Perform alpha pre-multiplication, since macOS requires it
	for i := 0; i < length; i++ {
		p := pixels[i]
		a := p >> 24
		r := ((((p >> 16) & 0xFF) * a) >> 8) << 16
		g := ((((p >> 8) & 0xFF) * a) >> 8) << 8
		b := ((p & 0xFF) * a) >> 8
		a <<= 24
		pixels[i] = a | r | g | b
	}

	dataProvider := C.CGDataProviderCreateWithData(nil, unsafe.Pointer(&pixels[0]), C.size_t(length*4), nil)
	defer C.CGDataProviderRelease(dataProvider)
	if image := C.CGImageCreate(C.size_t(data.Width), C.size_t(data.Height), 8, 32, C.size_t(data.Width*4), colorspace, C.kCGBitmapByteOrder32Host|C.kCGImageAlphaPremultipliedFirst, dataProvider, nil, false, C.kCGRenderingIntentDefault); image != nil {
		return createImage(image)
	}
	return nil
}

func newImageFromImage(other *Image, bounds geom.Rect) *Image {
	if image := C.CGImageCreateWithImageInRect(other.img, C.CGRectMake(C.CGFloat(bounds.X), C.CGFloat(bounds.Y), C.CGFloat(bounds.Width), C.CGFloat(bounds.Height))); image != nil {
		return createImage(image)
	}
	return nil
}

func createImage(img C.CGImageRef) *Image {
	return &Image{size: geom.Size{Width: float32(C.CGImageGetWidth(img)), Height: float32(C.CGImageGetHeight(img))}, img: unsafe.Pointer(img)}
}

func (img *Image) dispose() {
	C.CGImageRelease(img.img)
	img.img = nil
}

// Data extracts the raw image data.
func (img *Image) Data() *Data {
	size := img.Size()
	width := C.size_t(size.Width)
	height := C.size_t(size.Height)
	pixels := make([]color.Color, width*height)
	colorspace := C.CGColorSpaceCreateWithName(C.kCGColorSpaceGenericRGB)
	defer C.CGColorSpaceRelease(colorspace)
	context := C.CGBitmapContextCreate(unsafe.Pointer(&pixels[0]), width, height, 8, width*4, colorspace, C.uint32_t(C.kCGBitmapByteOrder32Host|C.kCGImageAlphaPremultipliedFirst))
	if context == nil {
		return nil
	}
	defer C.CGContextRelease(context)
	C.CGContextDrawImage(context, C.CGRectMake(0, 0, C.CGFloat(width), C.CGFloat(height)), img.img)

	// Remove alpha pre-multiplication, since macOS always returns the data that way, but we need it unmodified.
	// Note that this can be slightly lossy. Unfortunately, the options to request a non-pre-multiplied alpha
	// image error out. :-(
	length := len(pixels)
	for i := 0; i < length; i++ {
		p := pixels[i]
		a := p >> 24
		if a == 0 {
			pixels[i] = a << 24
		} else {
			r := ((((p >> 16) & 0xFF) << 8) / a) << 16
			g := ((((p >> 8) & 0xFF) << 8) / a) << 8
			b := ((p & 0xFF) << 8) / a
			a <<= 24
			pixels[i] = a | r | g | b
		}
	}

	return &Data{Width: int(width), Height: int(height), Pixels: pixels}
}
