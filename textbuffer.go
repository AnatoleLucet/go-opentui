package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"

// TextBuffer wraps the TextBuffer from the C library.
// It stores styled text content that can be rendered via a TextBufferView.
type TextBuffer struct {
	ptr *C.TextBuffer
}

// NewTextBuffer creates a new text buffer.
// widthMethod: 0 = wcwidth, 1 = Unicode standard width
func NewTextBuffer(length uint32, widthMethod uint8) *TextBuffer {
	ptr := C.createTextBuffer(C.uint32_t(length), C.uint8_t(widthMethod))
	if ptr == nil {
		return nil
	}

	tb := &TextBuffer{ptr: ptr}
	setFinalizer(tb, func(tb *TextBuffer) { tb.Close() })
	return tb
}

// Close releases the text buffer's resources.
func (tb *TextBuffer) Close() error {
	if tb.ptr != nil {
		clearFinalizer(tb)
		C.destroyTextBuffer(tb.ptr)
		tb.ptr = nil
	}
	return nil
}

// Reset clears the text buffer content.
func (tb *TextBuffer) Reset() {
	if tb.ptr == nil {
		return
	}
	C.textBufferReset(tb.ptr)
}

// Clear clears the text buffer.
func (tb *TextBuffer) Clear() {
	if tb.ptr == nil {
		return
	}
	C.textBufferClear(tb.ptr)
}

// Append adds text to the buffer.
func (tb *TextBuffer) Append(text string) {
	if tb.ptr == nil || len(text) == 0 {
		return
	}

	textPtr, textLen := stringToC(text)
	if textPtr == nil {
		return
	}

	C.textBufferAppend(tb.ptr, textPtr, C.uint32_t(textLen))
}

// GetLength returns the total display width of the text.
func (tb *TextBuffer) GetLength() uint32 {
	if tb.ptr == nil {
		return 0
	}
	return uint32(C.textBufferGetLength(tb.ptr))
}

// GetByteSize returns the byte size of the text.
func (tb *TextBuffer) GetByteSize() uint32 {
	if tb.ptr == nil {
		return 0
	}
	return uint32(C.textBufferGetByteSize(tb.ptr))
}

// GetLineCount returns the number of lines.
func (tb *TextBuffer) GetLineCount() uint32 {
	if tb.ptr == nil {
		return 0
	}
	return uint32(C.textBufferGetLineCount(tb.ptr))
}

// SetDefaultFg sets the default foreground color.
func (tb *TextBuffer) SetDefaultFg(fg RGBA) {
	if tb.ptr == nil {
		return
	}
	C.textBufferSetDefaultFg(tb.ptr, fg.toCFloat())
}

// SetDefaultBg sets the default background color.
func (tb *TextBuffer) SetDefaultBg(bg RGBA) {
	if tb.ptr == nil {
		return
	}
	C.textBufferSetDefaultBg(tb.ptr, bg.toCFloat())
}

// ResetDefaults clears all default styling.
func (tb *TextBuffer) ResetDefaults() {
	if tb.ptr == nil {
		return
	}
	C.textBufferResetDefaults(tb.ptr)
}

// GetTabWidth returns the tab width.
func (tb *TextBuffer) GetTabWidth() uint8 {
	if tb.ptr == nil {
		return 0
	}
	return uint8(C.textBufferGetTabWidth(tb.ptr))
}

// SetTabWidth sets the tab width.
func (tb *TextBuffer) SetTabWidth(width uint8) {
	if tb.ptr == nil {
		return
	}
	C.textBufferSetTabWidth(tb.ptr, C.uint8_t(width))
}

// Valid checks if the text buffer is still valid.
func (tb *TextBuffer) Valid() bool {
	return tb.ptr != nil
}

// TextBufferView wraps the TextBufferView from the C library.
// It handles text wrapping, viewport, and rendering.
type TextBufferView struct {
	ptr    *C.TextBufferView
	buffer *TextBuffer // keep reference to prevent GC
}

// WrapMode constants
const (
	WrapModeNone = 0 // No wrapping
	WrapModeChar = 1 // Wrap at character boundaries
	WrapModeWord = 2 // Wrap at word boundaries
)

// NewTextBufferView creates a view for the given text buffer.
func NewTextBufferView(tb *TextBuffer) *TextBufferView {
	if tb == nil || tb.ptr == nil {
		return nil
	}

	ptr := C.createTextBufferView(tb.ptr)
	if ptr == nil {
		return nil
	}

	view := &TextBufferView{ptr: ptr, buffer: tb}
	setFinalizer(view, func(v *TextBufferView) { v.Close() })
	return view
}

// Close releases the view's resources.
func (v *TextBufferView) Close() error {
	if v.ptr != nil {
		clearFinalizer(v)
		C.destroyTextBufferView(v.ptr)
		v.ptr = nil
		v.buffer = nil
	}
	return nil
}

// SetWrapWidth sets the width for text wrapping.
func (v *TextBufferView) SetWrapWidth(width uint32) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetWrapWidth(v.ptr, C.uint32_t(width))
}

// SetWrapMode sets the wrapping mode (WrapModeNone, WrapModeChar, WrapModeWord).
func (v *TextBufferView) SetWrapMode(mode uint8) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetWrapMode(v.ptr, C.uint8_t(mode))
}

// SetViewportSize sets the viewport dimensions for clipping.
func (v *TextBufferView) SetViewportSize(width, height uint32) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetViewportSize(v.ptr, C.uint32_t(width), C.uint32_t(height))
}

// SetViewport sets the viewport position and size.
func (v *TextBufferView) SetViewport(x, y, width, height uint32) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetViewport(v.ptr, C.uint32_t(x), C.uint32_t(y), C.uint32_t(width), C.uint32_t(height))
}

// GetVirtualLineCount returns the number of virtual lines (after wrapping).
func (v *TextBufferView) GetVirtualLineCount() uint32 {
	if v.ptr == nil {
		return 0
	}
	return uint32(C.textBufferViewGetVirtualLineCount(v.ptr))
}

// MeasureForDimensions measures text for the given dimensions.
// Returns the actual width and height (line count) needed.
func (v *TextBufferView) MeasureForDimensions(width, height uint32) (outWidth, outHeight uint32) {
	if v.ptr == nil {
		return 0, 0
	}
	var lineCount, maxWidth C.uint32_t
	C.textBufferViewMeasureForDimensions(v.ptr, C.uint32_t(width), C.uint32_t(height), &lineCount, &maxWidth)
	return uint32(maxWidth), uint32(lineCount)
}

// Valid checks if the view is still valid.
func (v *TextBufferView) Valid() bool {
	return v.ptr != nil
}

// Buffer returns the underlying TextBuffer.
func (v *TextBufferView) Buffer() *TextBuffer {
	return v.buffer
}

// DrawTextBufferView draws a TextBufferView to the buffer.
func (b *Buffer) DrawTextBufferView(view *TextBufferView, x, y int32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if view == nil || view.ptr == nil {
		return newError("view is nil or closed")
	}

	C.bufferDrawTextBufferView(b.ptr, view.ptr, C.int32_t(x), C.int32_t(y))
	return nil
}
