package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TextBuffer wraps the TextBuffer from the C library.
// It stores styled text content that can be rendered via a TextBufferView.
type TextBuffer struct {
	ptr *C.TextBuffer
	// Keep references to byte slices to prevent GC from collecting them
	// The C library stores pointers to this data
	dataRefs [][]byte
}

// NewTextBuffer creates a new text buffer.
// widthMethod: 0 = wcwidth, 1 = Unicode standard width
func NewTextBuffer(widthMethod uint8) *TextBuffer {
	ptr := C.createTextBuffer(C.uint8_t(widthMethod))
	if ptr == nil {
		return nil
	}

	tb := &TextBuffer{ptr: ptr, dataRefs: make([][]byte, 0)}
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
	// Clear our references since the C buffer is cleared
	tb.dataRefs = tb.dataRefs[:0]
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

	// Convert string to byte slice and keep a reference to prevent GC
	// The C library stores a pointer to this data
	bytes := []byte(text)
	tb.dataRefs = append(tb.dataRefs, bytes)

	C.textBufferAppend(tb.ptr, (*C.uint8_t)(unsafe.Pointer(&bytes[0])), C.size_t(len(bytes)))
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

// SetDefaultAttributes sets the default text attributes.
func (tb *TextBuffer) SetDefaultAttributes(attributes uint32) {
	if tb.ptr == nil {
		return
	}
	C.textBufferSetDefaultAttributes(tb.ptr, (*C.uint32_t)(unsafe.Pointer(&attributes)))
}

// ResetDefaults clears all default styling.
func (tb *TextBuffer) ResetDefaults() {
	if tb.ptr == nil {
		return
	}
	C.textBufferResetDefaults(tb.ptr)
}

// GetPlainText returns the plain text content of the buffer.
func (tb *TextBuffer) GetPlainText(maxLen int) string {
	if tb.ptr == nil || maxLen <= 0 {
		return ""
	}

	buffer := make([]byte, maxLen)
	actualLen := C.textBufferGetPlainText(tb.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}

	return string(buffer[:actualLen])
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

// RegisterMemBuffer registers a memory buffer and returns its ID.
// Returns 0xFFFF on failure.
func (tb *TextBuffer) RegisterMemBuffer(data []byte, owned bool) uint16 {
	if tb.ptr == nil || len(data) == 0 {
		return 0xFFFF
	}
	return uint16(C.textBufferRegisterMemBuffer(tb.ptr, (*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), C.bool(owned)))
}

// ReplaceMemBuffer replaces the content of a registered memory buffer.
func (tb *TextBuffer) ReplaceMemBuffer(id uint8, data []byte, owned bool) bool {
	if tb.ptr == nil || len(data) == 0 {
		return false
	}
	return bool(C.textBufferReplaceMemBuffer(tb.ptr, C.uint8_t(id), (*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), C.bool(owned)))
}

// ClearMemRegistry clears all registered memory buffers.
func (tb *TextBuffer) ClearMemRegistry() {
	if tb.ptr == nil {
		return
	}
	C.textBufferClearMemRegistry(tb.ptr)
}

// SetTextFromMem sets the text buffer content from a registered memory buffer.
func (tb *TextBuffer) SetTextFromMem(id uint8) {
	if tb.ptr == nil {
		return
	}
	C.textBufferSetTextFromMem(tb.ptr, C.uint8_t(id))
}

// AppendFromMemId appends text from a registered memory buffer.
func (tb *TextBuffer) AppendFromMemId(id uint8) {
	if tb.ptr == nil {
		return
	}
	C.textBufferAppendFromMemId(tb.ptr, C.uint8_t(id))
}

// LoadFile loads text content from a file.
func (tb *TextBuffer) LoadFile(path string) bool {
	if tb.ptr == nil || len(path) == 0 {
		return false
	}
	pathPtr, pathLen := stringToC(path)
	return bool(C.textBufferLoadFile(tb.ptr, pathPtr, pathLen))
}

// GetTextRange returns text in the specified offset range.
func (tb *TextBuffer) GetTextRange(startOffset, endOffset uint32, maxLen int) string {
	if tb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.textBufferGetTextRange(tb.ptr, C.uint32_t(startOffset), C.uint32_t(endOffset), (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// GetTextRangeByCoords returns text in the specified coordinate range.
func (tb *TextBuffer) GetTextRangeByCoords(startRow, startCol, endRow, endCol uint32, maxLen int) string {
	if tb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.textBufferGetTextRangeByCoords(tb.ptr, C.uint32_t(startRow), C.uint32_t(startCol), C.uint32_t(endRow), C.uint32_t(endCol), (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// SetSyntaxStyle sets the syntax style for highlighting.
func (tb *TextBuffer) SetSyntaxStyle(style *SyntaxStyle) {
	if tb.ptr == nil {
		return
	}
	var stylePtr *C.SyntaxStyle
	if style != nil {
		stylePtr = style.ptr
	}
	C.textBufferSetSyntaxStyle(tb.ptr, stylePtr)
}

// AddHighlight adds a highlight to a specific line.
func (tb *TextBuffer) AddHighlight(lineIdx uint32, highlight Highlight) {
	if tb.ptr == nil {
		return
	}
	chl := C.Highlight{
		start:    C.uint32_t(highlight.Start),
		end:      C.uint32_t(highlight.End),
		style_id: C.uint32_t(highlight.StyleID),
		priority: C.uint8_t(highlight.Priority),
		hl_ref:   C.uint16_t(highlight.HLRef),
	}
	C.textBufferAddHighlight(tb.ptr, C.uint32_t(lineIdx), &chl)
}

// AddHighlightByCharRange adds a highlight by character range.
func (tb *TextBuffer) AddHighlightByCharRange(highlight Highlight) {
	if tb.ptr == nil {
		return
	}
	chl := C.Highlight{
		start:    C.uint32_t(highlight.Start),
		end:      C.uint32_t(highlight.End),
		style_id: C.uint32_t(highlight.StyleID),
		priority: C.uint8_t(highlight.Priority),
		hl_ref:   C.uint16_t(highlight.HLRef),
	}
	C.textBufferAddHighlightByCharRange(tb.ptr, &chl)
}

// RemoveHighlightsByRef removes all highlights with the given reference.
func (tb *TextBuffer) RemoveHighlightsByRef(hlRef uint16) {
	if tb.ptr == nil {
		return
	}
	C.textBufferRemoveHighlightsByRef(tb.ptr, C.uint16_t(hlRef))
}

// ClearLineHighlights clears all highlights on a specific line.
func (tb *TextBuffer) ClearLineHighlights(lineIdx uint32) {
	if tb.ptr == nil {
		return
	}
	C.textBufferClearLineHighlights(tb.ptr, C.uint32_t(lineIdx))
}

// ClearAllHighlights clears all highlights.
func (tb *TextBuffer) ClearAllHighlights() {
	if tb.ptr == nil {
		return
	}
	C.textBufferClearAllHighlights(tb.ptr)
}

// GetHighlightCount returns the total number of highlights.
func (tb *TextBuffer) GetHighlightCount() uint32 {
	if tb.ptr == nil {
		return 0
	}
	return uint32(C.textBufferGetHighlightCount(tb.ptr))
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
// Returns the actual width and height (line count) needed, and success status.
func (v *TextBufferView) MeasureForDimensions(width, height uint32) (outWidth, outHeight uint32, success bool) {
	if v.ptr == nil {
		return 0, 0, false
	}
	var result C.MeasureResult
	cResult := C.textBufferViewMeasureForDimensions(v.ptr, C.uint32_t(width), C.uint32_t(height), &result)
	success = bool(cResult)
	return uint32(result.max_width), uint32(result.line_count), success
}

// Valid checks if the view is still valid.
func (v *TextBufferView) Valid() bool {
	return v.ptr != nil
}

// Buffer returns the underlying TextBuffer.
func (v *TextBufferView) Buffer() *TextBuffer {
	return v.buffer
}

// SetSelection sets a selection range with optional colors.
func (v *TextBufferView) SetSelection(start, end uint32, bgColor, fgColor *RGBA) {
	if v.ptr == nil {
		return
	}
	var bg, fg *C.float
	if bgColor != nil {
		bg = bgColor.toCFloat()
	}
	if fgColor != nil {
		fg = fgColor.toCFloat()
	}
	C.textBufferViewSetSelection(v.ptr, C.uint32_t(start), C.uint32_t(end), bg, fg)
}

// ResetSelection clears the current selection.
func (v *TextBufferView) ResetSelection() {
	if v.ptr == nil {
		return
	}
	C.textBufferViewResetSelection(v.ptr)
}

// GetSelectionInfo returns packed selection information.
func (v *TextBufferView) GetSelectionInfo() uint64 {
	if v.ptr == nil {
		return 0
	}
	return uint64(C.textBufferViewGetSelectionInfo(v.ptr))
}

// SetLocalSelection sets a local selection by coordinates.
func (v *TextBufferView) SetLocalSelection(anchorX, anchorY, focusX, focusY int32, bgColor, fgColor *RGBA) bool {
	if v.ptr == nil {
		return false
	}
	var bg, fg *C.float
	if bgColor != nil {
		bg = bgColor.toCFloat()
	}
	if fgColor != nil {
		fg = fgColor.toCFloat()
	}
	return bool(C.textBufferViewSetLocalSelection(v.ptr, C.int32_t(anchorX), C.int32_t(anchorY), C.int32_t(focusX), C.int32_t(focusY), bg, fg))
}

// UpdateSelection updates the selection end position.
func (v *TextBufferView) UpdateSelection(end uint32, bgColor, fgColor *RGBA) {
	if v.ptr == nil {
		return
	}
	var bg, fg *C.float
	if bgColor != nil {
		bg = bgColor.toCFloat()
	}
	if fgColor != nil {
		fg = fgColor.toCFloat()
	}
	C.textBufferViewUpdateSelection(v.ptr, C.uint32_t(end), bg, fg)
}

// UpdateLocalSelection updates the local selection by coordinates.
func (v *TextBufferView) UpdateLocalSelection(anchorX, anchorY, focusX, focusY int32, bgColor, fgColor *RGBA) bool {
	if v.ptr == nil {
		return false
	}
	var bg, fg *C.float
	if bgColor != nil {
		bg = bgColor.toCFloat()
	}
	if fgColor != nil {
		fg = fgColor.toCFloat()
	}
	return bool(C.textBufferViewUpdateLocalSelection(v.ptr, C.int32_t(anchorX), C.int32_t(anchorY), C.int32_t(focusX), C.int32_t(focusY), bg, fg))
}

// ResetLocalSelection clears the local selection.
func (v *TextBufferView) ResetLocalSelection() {
	if v.ptr == nil {
		return
	}
	C.textBufferViewResetLocalSelection(v.ptr)
}

// GetSelectedText returns the currently selected text.
func (v *TextBufferView) GetSelectedText(maxLen int) string {
	if v.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.textBufferViewGetSelectedText(v.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// GetPlainText returns the plain text content of the view.
func (v *TextBufferView) GetPlainText(maxLen int) string {
	if v.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.textBufferViewGetPlainText(v.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// SetTabIndicator sets the tab indicator character.
func (v *TextBufferView) SetTabIndicator(indicator rune) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetTabIndicator(v.ptr, C.uint32_t(indicator))
}

// SetTabIndicatorColor sets the tab indicator color.
func (v *TextBufferView) SetTabIndicatorColor(color RGBA) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetTabIndicatorColor(v.ptr, color.toCFloat())
}

// SetTruncate enables or disables text truncation.
func (v *TextBufferView) SetTruncate(truncate bool) {
	if v.ptr == nil {
		return
	}
	C.textBufferViewSetTruncate(v.ptr, C.bool(truncate))
}

// GetLineInfo returns line information for the view.
func (v *TextBufferView) GetLineInfo() LineInfo {
	if v.ptr == nil {
		return LineInfo{}
	}
	var info C.LineInfo
	C.textBufferViewGetLineInfoDirect(v.ptr, &info)
	return lineInfoFromC(&info)
}

// GetLogicalLineInfo returns logical line information for the view.
func (v *TextBufferView) GetLogicalLineInfo() LineInfo {
	if v.ptr == nil {
		return LineInfo{}
	}
	var info C.LineInfo
	C.textBufferViewGetLogicalLineInfoDirect(v.ptr, &info)
	return lineInfoFromC(&info)
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
