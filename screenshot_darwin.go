package screenshot

// CREDITS : https://github.com/vova616/screenshot/blob/master/screenshot_darwin.go
// :p, my repo actually just for linux (based on my need)

import (
	// #cgo LDFLAGS: -framework CoreGraphics
	// #cgo LDFLAGS: -framework CoreFoundation
	// #include <CoreGraphics/CoreGraphics.h>
	// #include <CoreFoundation/CoreFoundation.h>
	"C"
	"image"
	"reflect"
	"unsafe"
)

func ScreenRect() (image.Rectangle, error) {
	displayID := C.CGMainDisplayID()
	width := int(C.CGDisplayPixelsWide(displayID))
	height := int(C.CGDisplayPixelsHigh(displayID))
	return image.Rect(0, 0, width, height), nil
}

func CaptureScreen() (*image.RGBA, error) {
	rect, err := ScreenRect()
	if err != nil {
		return nil, err
	}
	return CaptureRect(rect)
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	displayID := C.CGMainDisplayID()
	width := int(C.CGDisplayPixelsWide(displayID))
	rawData := C.CGDataProviderCopyData(C.CGImageGetDataProvider(C.CGDisplayCreateImage(displayID)))

	length := int(C.CFDataGetLength(rawData))
	ptr := unsafe.Pointer(C.CFDataGetBytePtr(rawData))

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = length
	hdrp.Cap = length

	imageBytes := make([]byte, length)

	for i := 0; i < length; i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	C.CFRelease(rawData)

	img := &image.RGBA{Pix: imageBytes, Stride: 4 * width, Rect: rect}
	return img, nil
}