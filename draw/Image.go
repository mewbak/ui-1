// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"github.com/richardwilkes/errs"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/geom"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"sync"
	"unsafe"
	// #cgo pkg-config: pangocairo
	// #include <cairo.h>
	"C"
)

type imgRef struct {
	img   *Image
	count int
}

type fsKey struct {
	fs   http.FileSystem
	path string
}

// Image represents a set of pixels that can be drawn to a graphics.Context.
type Image struct {
	id         int
	disabledID int
	width      int
	height     int
	img        unsafe.Pointer
	key        interface{}
}

// ImageData is the raw information that makes up an Image.
type ImageData struct {
	Width  int
	Height int
	Pixels []color.Color
}

var (
	imageRegistryLock sync.Mutex
	nextImageID       = 1
	imageRegistry     = make(map[interface{}]*imgRef)
)

func loadFromStream(key interface{}, stream io.ReadCloser) (ref *imgRef, err error) {
	defer stream.Close()
	var simg image.Image
	if simg, _, err = image.Decode(stream); err != nil {
		return nil, errs.Wrap(err)
	}
	bounds := simg.Bounds()
	cimg := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32, C.int(bounds.Dx()), C.int(bounds.Dy()))
	stride := int(C.cairo_image_surface_get_stride(cimg)) / 4
	pixels := (*[1 << 30]color.Color)(unsafe.Pointer(C.cairo_image_surface_get_data(cimg)))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := simg.At(x, y).RGBA()
			pixels[(y-bounds.Min.Y)*stride+(x-bounds.Min.X)] = color.RGBA(int(r>>8), int(g>>8), int(b>>8), float64(a>>8)/255)
		}
	}
	C.cairo_surface_mark_dirty(cimg)
	ref = &imgRef{img: &Image{id: nextImageID, width: bounds.Dx(), height: bounds.Dy(), img: unsafe.Pointer(cimg), key: key}}
	nextImageID++
	return ref, nil
}

// AcquireImageFromFile attempts to load an image from the file system.
func AcquireImageFromFile(fs http.FileSystem, path string) (img *Image, err error) {
	imageRegistryLock.Lock()
	defer imageRegistryLock.Unlock()
	var ref *imgRef
	var ok bool
	key := fsKey{fs: fs, path: path}
	if ref, ok = imageRegistry[key]; !ok {
		var file http.File
		if file, err = fs.Open(path); err != nil {
			return nil, errs.Wrap(err)
		}
		if ref, err = loadFromStream(key, file); err != nil {
			return nil, err
		}
		imageRegistry[key] = ref
	}
	ref.count++
	return ref.img, nil
}

// AcquireImageFromURL attempts to load an image from a URL.
func AcquireImageFromURL(url string) (img *Image, err error) {
	imageRegistryLock.Lock()
	defer imageRegistryLock.Unlock()
	var ref *imgRef
	var ok bool
	if ref, ok = imageRegistry[url]; !ok {
		var resp *http.Response
		if resp, err = http.Get(url); err != nil {
			return nil, errs.Wrap(err)
		}
		if ref, err = loadFromStream(url, resp.Body); err != nil {
			return nil, err
		}
		imageRegistry[url] = ref
	}
	ref.count++
	return ref.img, nil
}

// AcquireImageFromID attempts to find an already loaded image by its ID and return it. Returns nil
// if it cannot be found.
func AcquireImageFromID(id int) *Image {
	imageRegistryLock.Lock()
	defer imageRegistryLock.Unlock()
	if r, ok := imageRegistry[id]; ok {
		r.count++
		return r.img
	}
	return nil
}

// AcquireImageFromData creates a new image from the specified data.
func AcquireImageFromData(data *ImageData) (img *Image, err error) {
	cimg := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32, C.int(data.Width), C.int(data.Height))
	stride := int(C.cairo_image_surface_get_stride(cimg)) / 4
	pixels := (*[1 << 30]color.Color)(unsafe.Pointer(C.cairo_image_surface_get_data(cimg)))
	for y := 0; y < data.Height; y++ {
		for x := 0; x < data.Width; x++ {
			pixels[y*stride+x] = data.Pixels[y*data.Width+x].Premultiply()
		}
	}
	C.cairo_surface_mark_dirty(cimg)
	imageRegistryLock.Lock()
	defer imageRegistryLock.Unlock()
	ref := &imgRef{img: &Image{id: nextImageID, width: data.Width, height: data.Height, img: unsafe.Pointer(cimg), key: nextImageID}, count: 1}
	imageRegistry[nextImageID] = ref
	nextImageID++
	return ref.img, nil
}

// AcquireImageArea creates a new image from an area within this image.
func (img *Image) AcquireImageArea(bounds geom.Rect) (image *Image, err error) {
	return AcquireImageFromData(img.DataFromArea(bounds))
}

// AcquireDisabled returns an image based on this image which is desaturated and ghosted to
// represent a disabled state.
func (img *Image) AcquireDisabled() (image *Image, e error) {
	image = AcquireImageFromID(img.disabledID)
	if image != nil {
		return image, nil
	}
	data := img.Data()
	for i := range data.Pixels {
		p := data.Pixels[i]
		v := int((p.Luminance() * 255) + 0.5)
		data.Pixels[i] = color.RGBA(v, v, v, p.AlphaIntensity()*0.4)
	}
	if image, e = AcquireImageFromData(data); e == nil {
		img.disabledID = image.id
	}
	return image, e
}

// ID returns the underlying ID of the image.
func (img *Image) ID() int {
	return img.id
}

// Size returns the size of the image.
func (img *Image) Size() geom.Size {
	return geom.Size{Width: float64(img.width), Height: float64(img.height)}
}

// Data extracts the raw image data.
func (img *Image) Data() *ImageData {
	data := &ImageData{Width: img.width, Height: img.height, Pixels: make([]color.Color, img.width*img.height)}
	stride := int(C.cairo_image_surface_get_stride(img.img)) / 4
	pixels := (*[1 << 30]color.Color)(unsafe.Pointer(C.cairo_image_surface_get_data(img.img)))
	for y := 0; y < img.height; y++ {
		for x := 0; x < img.width; x++ {
			data.Pixels[y*img.width+x] = pixels[y*stride+x].Unpremultiply()
		}
	}
	return data
}

// DataFromArea extracts the raw image data from an area within an image.
func (img *Image) DataFromArea(bounds geom.Rect) *ImageData {
	width := int(bounds.Width)
	height := int(bounds.Height)
	data := &ImageData{Width: width, Height: height, Pixels: make([]color.Color, width*height)}
	stride := int(C.cairo_image_surface_get_stride(img.img)) / 4
	pixels := (*[1 << 30]color.Color)(unsafe.Pointer(C.cairo_image_surface_get_data(img.img)))
	baseX := int(bounds.X)
	baseY := int(bounds.Y)
	outsidePixel := color.RGBA(0, 0, 0, 0)
	for y := 0; y < height; y++ {
		yy := y + baseY
		for x := 0; x < width; x++ {
			var pixel color.Color
			xx := x + baseX
			if xx < 0 || xx >= img.width || yy < 0 || yy >= img.height {
				pixel = outsidePixel
			} else {
				pixel = pixels[yy*stride+xx].Unpremultiply()
			}
			data.Pixels[y*width+x] = pixel
		}
	}
	return data
}

// PlatformPtr returns a pointer to the underlying platform-specific data.
func (img *Image) PlatformPtr() unsafe.Pointer {
	return img.img
}

// Release releases the image. If no other client is using the image, then the underlying OS
// resources for the image will be disposed of.
func (img *Image) Release() {
	imageRegistryLock.Lock()
	defer imageRegistryLock.Unlock()
	if ref, ok := imageRegistry[img.key]; ok {
		ref.count--
		if ref.count > 0 {
			return
		}
		delete(imageRegistry, img.key)
	}
	if img.img != nil {
		C.cairo_surface_destroy(img.img)
	}
}
