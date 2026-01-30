package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// Buffer wraps the OptimizedBuffer from the C library.
// It represents a 2D array of terminal cells for efficient rendering.
type Buffer struct {
	ptr     *C.OptimizedBuffer
	managed bool // true if buffer is managed by renderer
}

// WidthMethod constants for Unicode width calculation
const (
	WidthMethodWCWidth = 0 // Use wcwidth for width calculation
	WidthMethodUnicode = 1 // Use Unicode standard width calculation
)

// NewBuffer creates a new buffer with the specified dimensions.
// If respectAlpha is true, the buffer will handle alpha blending.
// The widthMethod parameter controls how text width is calculated (use WidthMethodUnicode for full Unicode support).
func NewBuffer(width, height uint32, respectAlpha bool, widthMethod uint8) *Buffer {
	return NewBufferWithID(width, height, respectAlpha, widthMethod, "")
}

// NewBufferWithID creates a new buffer with the specified dimensions and ID.
// If respectAlpha is true, the buffer will handle alpha blending.
// The widthMethod parameter controls how text width is calculated (use WidthMethodUnicode for full Unicode support).
// The id parameter is used for debugging/identification purposes.
func NewBufferWithID(width, height uint32, respectAlpha bool, widthMethod uint8, id string) *Buffer {
	if width == 0 || height == 0 {
		return nil
	}

	var idPtr *C.uint8_t
	var idLen C.size_t
	if len(id) > 0 {
		idPtr = (*C.uint8_t)(unsafe.Pointer(&[]byte(id)[0]))
		idLen = C.size_t(len(id))
	}

	ptr := C.createOptimizedBuffer(C.uint32_t(width), C.uint32_t(height), C.bool(respectAlpha), C.uint8_t(widthMethod), idPtr, idLen)
	if ptr == nil {
		return nil
	}

	b := &Buffer{ptr: ptr, managed: false}
	setFinalizer(b, func(b *Buffer) { b.Close() })
	return b
}

// Close releases the buffer's resources.
// After calling Close, the buffer should not be used.
// Note: Buffers obtained from a renderer are managed automatically and don't need to be closed.
func (b *Buffer) Close() error {
	if b.ptr != nil && !b.managed {
		clearFinalizer(b)
		C.destroyOptimizedBuffer(b.ptr)
		b.ptr = nil
	}
	return nil
}

// Width returns the buffer width in cells.
func (b *Buffer) Width() (uint32, error) {
	if b.ptr == nil {
		return 0, newError("buffer is closed")
	}
	return uint32(C.getBufferWidth(b.ptr)), nil
}

// Height returns the buffer height in cells.
func (b *Buffer) Height() (uint32, error) {
	if b.ptr == nil {
		return 0, newError("buffer is closed")
	}
	return uint32(C.getBufferHeight(b.ptr)), nil
}

// Size returns the buffer dimensions.
func (b *Buffer) Size() (uint32, uint32, error) {
	if b.ptr == nil {
		return 0, 0, newError("buffer is closed")
	}
	w := uint32(C.getBufferWidth(b.ptr))
	h := uint32(C.getBufferHeight(b.ptr))
	return w, h, nil
}

// Clear fills the entire buffer with the specified background color.
func (b *Buffer) Clear(bg RGBA) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferClear(b.ptr, bg.toCFloat())
	return nil
}

// GetRespectAlpha returns whether the buffer respects alpha values.
func (b *Buffer) GetRespectAlpha() (bool, error) {
	if b.ptr == nil {
		return false, newError("buffer is closed")
	}
	return bool(C.bufferGetRespectAlpha(b.ptr)), nil
}

// SetRespectAlpha sets whether the buffer should respect alpha values.
func (b *Buffer) SetRespectAlpha(respectAlpha bool) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferSetRespectAlpha(b.ptr, C.bool(respectAlpha))
	return nil
}

// DrawText draws text at the specified position with the given colors and attributes.
func (b *Buffer) DrawText(text string, x, y uint32, fg RGBA, bg *RGBA, attributes uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}

	textPtr, textLen := stringToC(text)
	if textPtr == nil {
		return nil // Empty string, nothing to draw
	}

	var bgPtr *C.float
	if bg != nil {
		bgPtr = bg.toCFloat()
	}

	C.bufferDrawText(b.ptr, textPtr, textLen, C.uint32_t(x), C.uint32_t(y), fg.toCFloat(), bgPtr, C.uint32_t(attributes))
	return nil
}

// SetCellWithAlphaBlending sets a single cell with alpha blending support.
func (b *Buffer) SetCellWithAlphaBlending(x, y uint32, char rune, fg, bg RGBA, attributes uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferSetCellWithAlphaBlending(b.ptr, C.uint32_t(x), C.uint32_t(y), C.uint32_t(char), fg.toCFloat(), bg.toCFloat(), C.uint32_t(attributes))
	return nil
}

// FillRect fills a rectangular area with the specified background color.
func (b *Buffer) FillRect(x, y, width, height uint32, bg RGBA) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferFillRect(b.ptr, C.uint32_t(x), C.uint32_t(y), C.uint32_t(width), C.uint32_t(height), bg.toCFloat())
	return nil
}

// DrawPackedBuffer draws packed buffer data at the specified position.
func (b *Buffer) DrawPackedBuffer(data []byte, posX, posY, terminalWidthCells, terminalHeightCells uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if len(data) == 0 {
		return nil
	}

	dataPtr, dataLen := sliceToC(data)
	C.bufferDrawPackedBuffer(b.ptr, (*C.uint8_t)(unsafe.Pointer(dataPtr)), dataLen,
		C.uint32_t(posX), C.uint32_t(posY), C.uint32_t(terminalWidthCells), C.uint32_t(terminalHeightCells))
	return nil
}

// DrawSuperSampleBuffer draws super-sampled pixel data for high-resolution graphics.
func (b *Buffer) DrawSuperSampleBuffer(x, y uint32, pixelData []byte, format SuperSampleFormat, alignedBytesPerRow uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if len(pixelData) == 0 {
		return nil
	}

	dataPtr, dataLen := sliceToC(pixelData)
	C.bufferDrawSuperSampleBuffer(b.ptr, C.uint32_t(x), C.uint32_t(y),
		(*C.uint8_t)(unsafe.Pointer(dataPtr)), dataLen, C.uint8_t(format), C.uint32_t(alignedBytesPerRow))
	return nil
}

// DrawBox draws a box with optional borders and title.
func (b *Buffer) DrawBox(x, y int32, width, height uint32, options BoxOptions, borderColor, backgroundColor RGBA) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}

	// Convert border characters to C array
	borderChars := runesToC([]rune{
		options.BorderChars.TopLeft,
		options.BorderChars.TopRight,
		options.BorderChars.BottomLeft,
		options.BorderChars.BottomRight,
		options.BorderChars.Horizontal,
		options.BorderChars.Vertical,
		options.BorderChars.TopT,
		options.BorderChars.BottomT,
		options.BorderChars.LeftT,
		options.BorderChars.RightT,
		options.BorderChars.Cross,
	})

	// Pack options
	packed := packBorderOptions(options.Sides, options.Fill, uint8(options.TitleAlignment))

	// Handle title
	var titlePtr *C.uint8_t
	var titleLen C.uint32_t
	if options.Title != "" {
		ptr, len := stringToC(options.Title)
		titlePtr = ptr
		titleLen = C.uint32_t(len)
	}

	C.bufferDrawBox(b.ptr, C.int32_t(x), C.int32_t(y), C.uint32_t(width), C.uint32_t(height),
		borderChars, packed, borderColor.toCFloat(), backgroundColor.toCFloat(), titlePtr, titleLen)
	return nil
}

// Resize changes the buffer dimensions.
// This may invalidate any existing content.
func (b *Buffer) Resize(width, height uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if width == 0 || height == 0 {
		return newError("invalid dimensions")
	}
	C.bufferResize(b.ptr, C.uint32_t(width), C.uint32_t(height))
	return nil
}

// DrawFrameBuffer draws another buffer onto this buffer at the specified position.
func (b *Buffer) DrawFrameBuffer(destX, destY int32, frameBuffer *Buffer, sourceX, sourceY, sourceWidth, sourceHeight uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if frameBuffer == nil || frameBuffer.ptr == nil {
		return newError("frame buffer is nil or closed")
	}

	C.drawFrameBuffer(b.ptr, C.int32_t(destX), C.int32_t(destY), frameBuffer.ptr,
		C.uint32_t(sourceX), C.uint32_t(sourceY), C.uint32_t(sourceWidth), C.uint32_t(sourceHeight))
	return nil
}

// GetDirectAccess returns direct access to the buffer's internal arrays.
// This is an advanced feature for performance-critical operations.
// The returned slices are valid until the buffer is resized or closed.
func (b *Buffer) GetDirectAccess() (*DirectAccess, error) {
	if b.ptr == nil {
		return nil, newError("buffer is closed")
	}

	width, height, err := b.Size()
	if err != nil {
		return nil, err
	}

	size := int(width * height)

	charPtr := C.bufferGetCharPtr(b.ptr)
	fgPtr := C.bufferGetFgPtr(b.ptr)
	bgPtr := C.bufferGetBgPtr(b.ptr)
	attrPtr := C.bufferGetAttributesPtr(b.ptr)

	return &DirectAccess{
		Chars:      cArrayToSlice((*uint32)(charPtr), size),
		Foreground: cArrayToSlice((*RGBA)(unsafe.Pointer(fgPtr)), size),
		Background: cArrayToSlice((*RGBA)(unsafe.Pointer(bgPtr)), size),
		Attributes: cArrayToSlice((*uint32)(unsafe.Pointer(attrPtr)), size),
		Width:      width,
		Height:     height,
	}, nil
}

// DirectAccess provides direct access to buffer internal arrays for performance-critical operations.
// Warning: This is an advanced feature. Modifying these slices directly bypasses normal safety checks.
type DirectAccess struct {
	Chars      []uint32 // Character codes (Unicode code points)
	Foreground []RGBA   // Foreground colors
	Background []RGBA   // Background colors
	Attributes []uint32 // Text attributes (including link IDs in upper bits)
	Width      uint32   // Buffer width
	Height     uint32   // Buffer height
}

// GetCell returns the cell at the specified coordinates using direct access.
func (da *DirectAccess) GetCell(x, y uint32) (*Cell, error) {
	if x >= da.Width || y >= da.Height {
		return nil, newError("coordinates out of bounds")
	}

	index := y*da.Width + x
	return &Cell{
		Char:       rune(da.Chars[index]),
		Foreground: da.Foreground[index],
		Background: da.Background[index],
		Attributes: da.Attributes[index],
	}, nil
}

// SetCell sets the cell at the specified coordinates using direct access.
func (da *DirectAccess) SetCell(x, y uint32, cell Cell) error {
	if x >= da.Width || y >= da.Height {
		return newError("coordinates out of bounds")
	}

	index := y*da.Width + x
	da.Chars[index] = uint32(cell.Char)
	da.Foreground[index] = cell.Foreground
	da.Background[index] = cell.Background
	da.Attributes[index] = cell.Attributes
	return nil
}

// Valid checks if the buffer is still valid (not closed).
func (b *Buffer) Valid() bool {
	return b.ptr != nil
}

// PushScissorRect pushes a scissor rect onto the stack.
// All subsequent draw operations will be clipped to this rect.
func (b *Buffer) PushScissorRect(x, y int32, width, height uint32) {
	if b.ptr == nil {
		return
	}
	C.bufferPushScissorRect(b.ptr, C.int32_t(x), C.int32_t(y), C.uint32_t(width), C.uint32_t(height))
}

// PopScissorRect pops the top scissor rect from the stack.
func (b *Buffer) PopScissorRect() {
	if b.ptr == nil {
		return
	}
	C.bufferPopScissorRect(b.ptr)
}

// ClearScissorRects clears all scissor rects from the stack.
func (b *Buffer) ClearScissorRects() {
	if b.ptr == nil {
		return
	}
	C.bufferClearScissorRects(b.ptr)
}

// PushOpacity pushes an opacity value onto the opacity stack.
func (b *Buffer) PushOpacity(opacity float32) {
	if b.ptr == nil {
		return
	}
	C.bufferPushOpacity(b.ptr, C.float(opacity))
}

// PopOpacity pops the top opacity value from the stack.
func (b *Buffer) PopOpacity() {
	if b.ptr == nil {
		return
	}
	C.bufferPopOpacity(b.ptr)
}

// GetCurrentOpacity returns the current opacity value.
func (b *Buffer) GetCurrentOpacity() float32 {
	if b.ptr == nil {
		return 1.0
	}
	return float32(C.bufferGetCurrentOpacity(b.ptr))
}

// ClearOpacity clears the opacity stack.
func (b *Buffer) ClearOpacity() {
	if b.ptr == nil {
		return
	}
	C.bufferClearOpacity(b.ptr)
}

// SetCell sets a single cell without alpha blending.
func (b *Buffer) SetCell(x, y uint32, char rune, fg, bg RGBA, attributes uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferSetCell(b.ptr, C.uint32_t(x), C.uint32_t(y), C.uint32_t(char), fg.toCFloat(), bg.toCFloat(), C.uint32_t(attributes))
	return nil
}

// DrawChar draws a character at the specified position.
func (b *Buffer) DrawChar(x, y int32, char rune, fg, bg RGBA, attributes uint32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	C.bufferDrawChar(b.ptr, C.int32_t(x), C.int32_t(y), C.uint32_t(char), fg.toCFloat(), bg.toCFloat(), C.uint32_t(attributes))
	return nil
}

// GetID returns the buffer's ID.
func (b *Buffer) GetID(maxLen int) string {
	if b.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.bufferGetId(b.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// GetRealCharSize returns the real character size of the buffer.
func (b *Buffer) GetRealCharSize() uint32 {
	if b.ptr == nil {
		return 0
	}
	return uint32(C.bufferGetRealCharSize(b.ptr))
}

// WriteResolvedChars writes resolved characters to the output buffer.
func (b *Buffer) WriteResolvedChars(output []byte, addLineBreaks bool) uint32 {
	if b.ptr == nil || len(output) == 0 {
		return 0
	}
	return uint32(C.bufferWriteResolvedChars(b.ptr, (*C.uint8_t)(unsafe.Pointer(&output[0])), C.size_t(len(output)), C.bool(addLineBreaks)))
}

// DrawGrayscaleBuffer draws grayscale intensity data.
func (b *Buffer) DrawGrayscaleBuffer(posX, posY int32, intensities []float32, srcWidth, srcHeight uint32, fg, bg *RGBA) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if len(intensities) == 0 {
		return nil
	}
	var fgPtr, bgPtr *C.float
	if fg != nil {
		fgPtr = fg.toCFloat()
	}
	if bg != nil {
		bgPtr = bg.toCFloat()
	}
	C.bufferDrawGrayscaleBuffer(b.ptr, C.int32_t(posX), C.int32_t(posY), (*C.float)(unsafe.Pointer(&intensities[0])), C.uint32_t(srcWidth), C.uint32_t(srcHeight), fgPtr, bgPtr)
	return nil
}

// DrawGrayscaleBufferSupersampled draws supersampled grayscale intensity data.
func (b *Buffer) DrawGrayscaleBufferSupersampled(posX, posY int32, intensities []float32, srcWidth, srcHeight uint32, fg, bg *RGBA) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if len(intensities) == 0 {
		return nil
	}
	var fgPtr, bgPtr *C.float
	if fg != nil {
		fgPtr = fg.toCFloat()
	}
	if bg != nil {
		bgPtr = bg.toCFloat()
	}
	C.bufferDrawGrayscaleBufferSupersampled(b.ptr, C.int32_t(posX), C.int32_t(posY), (*C.float)(unsafe.Pointer(&intensities[0])), C.uint32_t(srcWidth), C.uint32_t(srcHeight), fgPtr, bgPtr)
	return nil
}
