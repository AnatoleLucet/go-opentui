package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// EditorView wraps the EditorView from the C library.
// It provides a view for editing text with viewport, cursor, and selection management.
type EditorView struct {
	ptr        *C.EditorView
	editBuffer *EditBuffer // keep reference to prevent GC
}

// NewEditorView creates a new editor view for the given edit buffer.
func NewEditorView(eb *EditBuffer, viewportWidth, viewportHeight uint32) *EditorView {
	if eb == nil || eb.ptr == nil {
		return nil
	}

	ptr := C.createEditorView(eb.ptr, C.uint32_t(viewportWidth), C.uint32_t(viewportHeight))
	if ptr == nil {
		return nil
	}

	view := &EditorView{ptr: ptr, editBuffer: eb}
	setFinalizer(view, func(v *EditorView) { v.Close() })
	return view
}

// Close releases the editor view's resources.
func (v *EditorView) Close() error {
	if v.ptr != nil {
		clearFinalizer(v)
		C.destroyEditorView(v.ptr)
		v.ptr = nil
		v.editBuffer = nil
	}
	return nil
}

// Valid checks if the editor view is still valid.
func (v *EditorView) Valid() bool {
	return v.ptr != nil
}

// EditBuffer returns the underlying EditBuffer.
func (v *EditorView) EditBuffer() *EditBuffer {
	return v.editBuffer
}

// GetTextBufferView returns the underlying TextBufferView.
func (v *EditorView) GetTextBufferView() *TextBufferView {
	if v.ptr == nil {
		return nil
	}
	tbvPtr := C.editorViewGetTextBufferView(v.ptr)
	if tbvPtr == nil {
		return nil
	}
	// Note: This is managed by the EditorView, don't set a finalizer
	return &TextBufferView{ptr: tbvPtr}
}

// SetViewportSize sets the viewport dimensions.
func (v *EditorView) SetViewportSize(width, height uint32) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetViewportSize(v.ptr, C.uint32_t(width), C.uint32_t(height))
}

// SetViewport sets the viewport position and size.
func (v *EditorView) SetViewport(x, y, width, height uint32, moveCursor bool) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetViewport(v.ptr, C.uint32_t(x), C.uint32_t(y), C.uint32_t(width), C.uint32_t(height), C.bool(moveCursor))
}

// GetViewport returns the current viewport position and size.
func (v *EditorView) GetViewport() (x, y, width, height uint32, ok bool) {
	if v.ptr == nil {
		return 0, 0, 0, 0, false
	}
	var outX, outY, outWidth, outHeight C.uint32_t
	success := bool(C.editorViewGetViewport(v.ptr, &outX, &outY, &outWidth, &outHeight))
	if !success {
		return 0, 0, 0, 0, false
	}
	return uint32(outX), uint32(outY), uint32(outWidth), uint32(outHeight), true
}

// ClearViewport clears the viewport settings.
func (v *EditorView) ClearViewport() {
	if v.ptr == nil {
		return
	}
	C.editorViewClearViewport(v.ptr)
}

// SetScrollMargin sets the scroll margin.
func (v *EditorView) SetScrollMargin(margin float32) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetScrollMargin(v.ptr, C.float(margin))
}

// SetWrapMode sets the text wrapping mode.
func (v *EditorView) SetWrapMode(mode uint8) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetWrapMode(v.ptr, C.uint8_t(mode))
}

// GetVirtualLineCount returns the number of visible virtual lines.
func (v *EditorView) GetVirtualLineCount() uint32 {
	if v.ptr == nil {
		return 0
	}
	return uint32(C.editorViewGetVirtualLineCount(v.ptr))
}

// GetTotalVirtualLineCount returns the total number of virtual lines.
func (v *EditorView) GetTotalVirtualLineCount() uint32 {
	if v.ptr == nil {
		return 0
	}
	return uint32(C.editorViewGetTotalVirtualLineCount(v.ptr))
}

// SetSelection sets a selection range.
func (v *EditorView) SetSelection(start, end uint32, bgColor, fgColor *RGBA) {
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
	C.editorViewSetSelection(v.ptr, C.uint32_t(start), C.uint32_t(end), bg, fg)
}

// ResetSelection clears the current selection.
func (v *EditorView) ResetSelection() {
	if v.ptr == nil {
		return
	}
	C.editorViewResetSelection(v.ptr)
}

// GetSelection returns the current selection as a packed uint64.
func (v *EditorView) GetSelection() uint64 {
	if v.ptr == nil {
		return 0
	}
	return uint64(C.editorViewGetSelection(v.ptr))
}

// SetLocalSelection sets a local selection by coordinates.
func (v *EditorView) SetLocalSelection(anchorX, anchorY, focusX, focusY int32, bgColor, fgColor *RGBA, updateCursor, followCursor bool) bool {
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
	return bool(C.editorViewSetLocalSelection(v.ptr, C.int32_t(anchorX), C.int32_t(anchorY), C.int32_t(focusX), C.int32_t(focusY), bg, fg, C.bool(updateCursor), C.bool(followCursor)))
}

// UpdateSelection updates the selection end position.
func (v *EditorView) UpdateSelection(end uint32, bgColor, fgColor *RGBA) {
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
	C.editorViewUpdateSelection(v.ptr, C.uint32_t(end), bg, fg)
}

// UpdateLocalSelection updates the local selection by coordinates.
func (v *EditorView) UpdateLocalSelection(anchorX, anchorY, focusX, focusY int32, bgColor, fgColor *RGBA, updateCursor, followCursor bool) bool {
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
	return bool(C.editorViewUpdateLocalSelection(v.ptr, C.int32_t(anchorX), C.int32_t(anchorY), C.int32_t(focusX), C.int32_t(focusY), bg, fg, C.bool(updateCursor), C.bool(followCursor)))
}

// ResetLocalSelection clears the local selection.
func (v *EditorView) ResetLocalSelection() {
	if v.ptr == nil {
		return
	}
	C.editorViewResetLocalSelection(v.ptr)
}

// GetSelectedText returns the currently selected text.
func (v *EditorView) GetSelectedText(maxLen int) string {
	if v.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editorViewGetSelectedTextBytes(v.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// DeleteSelectedText deletes the currently selected text.
func (v *EditorView) DeleteSelectedText() {
	if v.ptr == nil {
		return
	}
	C.editorViewDeleteSelectedText(v.ptr)
}

// GetCursor returns the logical cursor position (row, col).
func (v *EditorView) GetCursor() (row, col uint32) {
	if v.ptr == nil {
		return 0, 0
	}
	var outRow, outCol C.uint32_t
	C.editorViewGetCursor(v.ptr, &outRow, &outCol)
	return uint32(outRow), uint32(outCol)
}

// GetVisualCursor returns the visual cursor position.
func (v *EditorView) GetVisualCursor() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetVisualCursor(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// SetCursorByOffset sets the cursor position by offset.
func (v *EditorView) SetCursorByOffset(offset uint32) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetCursorByOffset(v.ptr, C.uint32_t(offset))
}

// MoveUpVisual moves the cursor up visually.
func (v *EditorView) MoveUpVisual() {
	if v.ptr == nil {
		return
	}
	C.editorViewMoveUpVisual(v.ptr)
}

// MoveDownVisual moves the cursor down visually.
func (v *EditorView) MoveDownVisual() {
	if v.ptr == nil {
		return
	}
	C.editorViewMoveDownVisual(v.ptr)
}

// GetNextWordBoundary returns the position of the next word boundary.
func (v *EditorView) GetNextWordBoundary() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetNextWordBoundary(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// GetPrevWordBoundary returns the position of the previous word boundary.
func (v *EditorView) GetPrevWordBoundary() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetPrevWordBoundary(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// GetEOL returns the position of the end of the current line.
func (v *EditorView) GetEOL() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetEOL(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// GetVisualSOL returns the position of the start of the visual line.
func (v *EditorView) GetVisualSOL() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetVisualSOL(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// GetVisualEOL returns the position of the end of the visual line.
func (v *EditorView) GetVisualEOL() VisualCursor {
	if v.ptr == nil {
		return VisualCursor{}
	}
	var cursor C.VisualCursor
	C.editorViewGetVisualEOL(v.ptr, &cursor)
	return VisualCursor{
		VisualRow:  uint32(cursor.visual_row),
		VisualCol:  uint32(cursor.visual_col),
		LogicalRow: uint32(cursor.logical_row),
		LogicalCol: uint32(cursor.logical_col),
		Offset:     uint32(cursor.offset),
	}
}

// GetText returns the entire text content.
func (v *EditorView) GetText(maxLen int) string {
	if v.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editorViewGetText(v.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// GetLineInfo returns line information for the view.
func (v *EditorView) GetLineInfo() LineInfo {
	if v.ptr == nil {
		return LineInfo{}
	}
	var info C.LineInfo
	C.editorViewGetLineInfoDirect(v.ptr, &info)
	return lineInfoFromC(&info)
}

// GetLogicalLineInfo returns logical line information for the view.
func (v *EditorView) GetLogicalLineInfo() LineInfo {
	if v.ptr == nil {
		return LineInfo{}
	}
	var info C.LineInfo
	C.editorViewGetLogicalLineInfoDirect(v.ptr, &info)
	return lineInfoFromC(&info)
}

// SetTabIndicator sets the tab indicator character.
func (v *EditorView) SetTabIndicator(indicator rune) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetTabIndicator(v.ptr, C.uint32_t(indicator))
}

// SetTabIndicatorColor sets the tab indicator color.
func (v *EditorView) SetTabIndicatorColor(color RGBA) {
	if v.ptr == nil {
		return
	}
	C.editorViewSetTabIndicatorColor(v.ptr, color.toCFloat())
}

// SetPlaceholderStyledText sets placeholder text with styling.
func (v *EditorView) SetPlaceholderStyledText(chunks []StyledChunk) {
	if v.ptr == nil || len(chunks) == 0 {
		return
	}

	// Filter out empty chunks
	var validChunks []StyledChunk
	for _, chunk := range chunks {
		if len(chunk.Text) > 0 {
			validChunks = append(validChunks, chunk)
		}
	}
	if len(validChunks) == 0 {
		return
	}

	chunkCount := C.size_t(len(validChunks))
	cChunks := (*C.StyledChunk)(C.malloc(C.size_t(unsafe.Sizeof(C.StyledChunk{})) * chunkCount))
	if cChunks == nil {
		return
	}
	defer C.free(unsafe.Pointer(cChunks))

	slice := unsafe.Slice(cChunks, len(validChunks))
	fgColors := make([]*C.float, 0, len(validChunks))
	bgColors := make([]*C.float, 0, len(validChunks))
	urlPtrs := make([]unsafe.Pointer, 0, len(validChunks))
	defer func() {
		for _, p := range fgColors {
			C.free(unsafe.Pointer(p))
		}
		for _, p := range bgColors {
			C.free(unsafe.Pointer(p))
		}
		for _, p := range urlPtrs {
			C.free(p)
		}
	}()

	for i, chunk := range validChunks {
		textBytes := []byte(chunk.Text)

		slice[i].text_ptr = (*C.uint8_t)(unsafe.Pointer(&textBytes[0]))
		slice[i].text_len = C.size_t(len(textBytes))
		slice[i].attributes = C.uint32_t(chunk.Attributes)
		slice[i].fg = nil
		slice[i].bg = nil
		slice[i].link_ptr = nil
		slice[i].link_len = 0

		// Handle URL/link if provided
		if chunk.URL != "" {
			urlBytes := []byte(chunk.URL)
			cUrl := C.malloc(C.size_t(len(urlBytes)))
			if cUrl != nil {
				copy(unsafe.Slice((*byte)(cUrl), len(urlBytes)), urlBytes)
				slice[i].link_ptr = (*C.uint8_t)(cUrl)
				slice[i].link_len = C.size_t(len(urlBytes))
				// Store to free later
				urlPtrs = append(urlPtrs, cUrl)
			}
		}

		if chunk.Foreground != nil {
			fg := (*C.float)(C.malloc(4 * C.size_t(unsafe.Sizeof(C.float(0)))))
			if fg != nil {
				fgSlice := unsafe.Slice(fg, 4)
				fgSlice[0] = C.float(chunk.Foreground.R)
				fgSlice[1] = C.float(chunk.Foreground.G)
				fgSlice[2] = C.float(chunk.Foreground.B)
				fgSlice[3] = C.float(chunk.Foreground.A)
				slice[i].fg = fg
				fgColors = append(fgColors, fg)
			}
		}
		if chunk.Background != nil {
			bg := (*C.float)(C.malloc(4 * C.size_t(unsafe.Sizeof(C.float(0)))))
			if bg != nil {
				bgSlice := unsafe.Slice(bg, 4)
				bgSlice[0] = C.float(chunk.Background.R)
				bgSlice[1] = C.float(chunk.Background.G)
				bgSlice[2] = C.float(chunk.Background.B)
				bgSlice[3] = C.float(chunk.Background.A)
				slice[i].bg = bg
				bgColors = append(bgColors, bg)
			}
		}
	}

	C.editorViewSetPlaceholderStyledText(v.ptr, cChunks, chunkCount)
}

// DrawEditorView draws an EditorView to the buffer.
func (b *Buffer) DrawEditorView(view *EditorView, x, y int32) error {
	if b.ptr == nil {
		return newError("buffer is closed")
	}
	if view == nil || view.ptr == nil {
		return newError("view is nil or closed")
	}

	C.bufferDrawEditorView(b.ptr, view.ptr, C.int32_t(x), C.int32_t(y))
	return nil
}
