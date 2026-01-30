package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// EditBuffer wraps the EditBuffer from the C library.
// It provides a rich text editing buffer with undo/redo support.
type EditBuffer struct {
	ptr *C.EditBuffer
}

// NewEditBuffer creates a new edit buffer.
// widthMethod: 0 = wcwidth, 1 = Unicode standard width
func NewEditBuffer(widthMethod uint8) *EditBuffer {
	ptr := C.createEditBuffer(C.uint8_t(widthMethod))
	if ptr == nil {
		return nil
	}

	eb := &EditBuffer{ptr: ptr}
	setFinalizer(eb, func(eb *EditBuffer) { eb.Close() })
	return eb
}

// Close releases the edit buffer's resources.
func (eb *EditBuffer) Close() error {
	if eb.ptr != nil {
		clearFinalizer(eb)
		C.destroyEditBuffer(eb.ptr)
		eb.ptr = nil
	}
	return nil
}

// Valid checks if the edit buffer is still valid.
func (eb *EditBuffer) Valid() bool {
	return eb.ptr != nil
}

// GetTextBuffer returns the underlying TextBuffer.
func (eb *EditBuffer) GetTextBuffer() *TextBuffer {
	if eb.ptr == nil {
		return nil
	}
	tbPtr := C.editBufferGetTextBuffer(eb.ptr)
	if tbPtr == nil {
		return nil
	}
	// Note: This TextBuffer is managed by the EditBuffer, don't set a finalizer
	return &TextBuffer{ptr: tbPtr}
}

// GetID returns the edit buffer's ID.
func (eb *EditBuffer) GetID() uint16 {
	if eb.ptr == nil {
		return 0
	}
	return uint16(C.editBufferGetId(eb.ptr))
}

// DebugLogRope outputs debug information about the rope structure.
func (eb *EditBuffer) DebugLogRope() {
	if eb.ptr == nil {
		return
	}
	C.editBufferDebugLogRope(eb.ptr)
}

// SetText sets the entire text content.
func (eb *EditBuffer) SetText(text string) {
	if eb.ptr == nil {
		return
	}
	if len(text) == 0 {
		C.editBufferSetText(eb.ptr, nil, 0)
		return
	}
	textPtr, textLen := stringToC(text)
	C.editBufferSetText(eb.ptr, textPtr, textLen)
}

// SetTextFromMem sets text from a registered memory buffer.
func (eb *EditBuffer) SetTextFromMem(memId uint8) {
	if eb.ptr == nil {
		return
	}
	C.editBufferSetTextFromMem(eb.ptr, C.uint8_t(memId))
}

// ReplaceText replaces the current selection or inserts at cursor.
func (eb *EditBuffer) ReplaceText(text string) {
	if eb.ptr == nil {
		return
	}
	if len(text) == 0 {
		C.editBufferReplaceText(eb.ptr, nil, 0)
		return
	}
	textPtr, textLen := stringToC(text)
	C.editBufferReplaceText(eb.ptr, textPtr, textLen)
}

// ReplaceTextFromMem replaces text from a registered memory buffer.
func (eb *EditBuffer) ReplaceTextFromMem(memId uint8) {
	if eb.ptr == nil {
		return
	}
	C.editBufferReplaceTextFromMem(eb.ptr, C.uint8_t(memId))
}

// GetText returns the entire text content.
func (eb *EditBuffer) GetText(maxLen int) string {
	if eb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editBufferGetText(eb.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// InsertChar inserts a character at the cursor position.
func (eb *EditBuffer) InsertChar(char string) {
	if eb.ptr == nil || len(char) == 0 {
		return
	}
	charPtr, charLen := stringToC(char)
	C.editBufferInsertChar(eb.ptr, charPtr, charLen)
}

// InsertText inserts text at the cursor position.
func (eb *EditBuffer) InsertText(text string) {
	if eb.ptr == nil || len(text) == 0 {
		return
	}
	textPtr, textLen := stringToC(text)
	C.editBufferInsertText(eb.ptr, textPtr, textLen)
}

// DeleteChar deletes the character after the cursor.
func (eb *EditBuffer) DeleteChar() {
	if eb.ptr == nil {
		return
	}
	C.editBufferDeleteChar(eb.ptr)
}

// DeleteCharBackward deletes the character before the cursor.
func (eb *EditBuffer) DeleteCharBackward() {
	if eb.ptr == nil {
		return
	}
	C.editBufferDeleteCharBackward(eb.ptr)
}

// DeleteRange deletes text in the specified range.
func (eb *EditBuffer) DeleteRange(startRow, startCol, endRow, endCol uint32) {
	if eb.ptr == nil {
		return
	}
	C.editBufferDeleteRange(eb.ptr, C.uint32_t(startRow), C.uint32_t(startCol), C.uint32_t(endRow), C.uint32_t(endCol))
}

// NewLine inserts a new line at the cursor position.
func (eb *EditBuffer) NewLine() {
	if eb.ptr == nil {
		return
	}
	C.editBufferNewLine(eb.ptr)
}

// DeleteLine deletes the current line.
func (eb *EditBuffer) DeleteLine() {
	if eb.ptr == nil {
		return
	}
	C.editBufferDeleteLine(eb.ptr)
}

// Clear clears all content.
func (eb *EditBuffer) Clear() {
	if eb.ptr == nil {
		return
	}
	C.editBufferClear(eb.ptr)
}

// MoveCursorLeft moves the cursor left.
func (eb *EditBuffer) MoveCursorLeft() {
	if eb.ptr == nil {
		return
	}
	C.editBufferMoveCursorLeft(eb.ptr)
}

// MoveCursorRight moves the cursor right.
func (eb *EditBuffer) MoveCursorRight() {
	if eb.ptr == nil {
		return
	}
	C.editBufferMoveCursorRight(eb.ptr)
}

// MoveCursorUp moves the cursor up.
func (eb *EditBuffer) MoveCursorUp() {
	if eb.ptr == nil {
		return
	}
	C.editBufferMoveCursorUp(eb.ptr)
}

// MoveCursorDown moves the cursor down.
func (eb *EditBuffer) MoveCursorDown() {
	if eb.ptr == nil {
		return
	}
	C.editBufferMoveCursorDown(eb.ptr)
}

// GetCursor returns the current cursor position (row, col).
func (eb *EditBuffer) GetCursor() (row, col uint32) {
	if eb.ptr == nil {
		return 0, 0
	}
	var outRow, outCol C.uint32_t
	C.editBufferGetCursor(eb.ptr, &outRow, &outCol)
	return uint32(outRow), uint32(outCol)
}

// SetCursor sets the cursor position.
func (eb *EditBuffer) SetCursor(row, col uint32) {
	if eb.ptr == nil {
		return
	}
	C.editBufferSetCursor(eb.ptr, C.uint32_t(row), C.uint32_t(col))
}

// SetCursorToLineCol sets the cursor to a specific line and column.
func (eb *EditBuffer) SetCursorToLineCol(row, col uint32) {
	if eb.ptr == nil {
		return
	}
	C.editBufferSetCursorToLineCol(eb.ptr, C.uint32_t(row), C.uint32_t(col))
}

// SetCursorByOffset sets the cursor by character offset.
func (eb *EditBuffer) SetCursorByOffset(offset uint32) {
	if eb.ptr == nil {
		return
	}
	C.editBufferSetCursorByOffset(eb.ptr, C.uint32_t(offset))
}

// GotoLine moves the cursor to the beginning of a line.
func (eb *EditBuffer) GotoLine(line uint32) {
	if eb.ptr == nil {
		return
	}
	C.editBufferGotoLine(eb.ptr, C.uint32_t(line))
}

// GetCursorPosition returns detailed cursor position information.
func (eb *EditBuffer) GetCursorPosition() LogicalCursor {
	if eb.ptr == nil {
		return LogicalCursor{}
	}
	var cursor C.LogicalCursor
	C.editBufferGetCursorPosition(eb.ptr, &cursor)
	return LogicalCursor{
		Row:    uint32(cursor.row),
		Col:    uint32(cursor.col),
		Offset: uint32(cursor.offset),
	}
}

// GetNextWordBoundary returns the position of the next word boundary.
func (eb *EditBuffer) GetNextWordBoundary() LogicalCursor {
	if eb.ptr == nil {
		return LogicalCursor{}
	}
	var cursor C.LogicalCursor
	C.editBufferGetNextWordBoundary(eb.ptr, &cursor)
	return LogicalCursor{
		Row:    uint32(cursor.row),
		Col:    uint32(cursor.col),
		Offset: uint32(cursor.offset),
	}
}

// GetPrevWordBoundary returns the position of the previous word boundary.
func (eb *EditBuffer) GetPrevWordBoundary() LogicalCursor {
	if eb.ptr == nil {
		return LogicalCursor{}
	}
	var cursor C.LogicalCursor
	C.editBufferGetPrevWordBoundary(eb.ptr, &cursor)
	return LogicalCursor{
		Row:    uint32(cursor.row),
		Col:    uint32(cursor.col),
		Offset: uint32(cursor.offset),
	}
}

// GetEOL returns the position of the end of the current line.
func (eb *EditBuffer) GetEOL() LogicalCursor {
	if eb.ptr == nil {
		return LogicalCursor{}
	}
	var cursor C.LogicalCursor
	C.editBufferGetEOL(eb.ptr, &cursor)
	return LogicalCursor{
		Row:    uint32(cursor.row),
		Col:    uint32(cursor.col),
		Offset: uint32(cursor.offset),
	}
}

// OffsetToPosition converts an offset to a row/col position.
func (eb *EditBuffer) OffsetToPosition(offset uint32) (LogicalCursor, bool) {
	if eb.ptr == nil {
		return LogicalCursor{}, false
	}
	var cursor C.LogicalCursor
	success := bool(C.editBufferOffsetToPosition(eb.ptr, C.uint32_t(offset), &cursor))
	if !success {
		return LogicalCursor{}, false
	}
	return LogicalCursor{
		Row:    uint32(cursor.row),
		Col:    uint32(cursor.col),
		Offset: uint32(cursor.offset),
	}, true
}

// PositionToOffset converts a row/col position to an offset.
func (eb *EditBuffer) PositionToOffset(row, col uint32) uint32 {
	if eb.ptr == nil {
		return 0
	}
	return uint32(C.editBufferPositionToOffset(eb.ptr, C.uint32_t(row), C.uint32_t(col)))
}

// GetLineStartOffset returns the offset of the start of a line.
func (eb *EditBuffer) GetLineStartOffset(row uint32) uint32 {
	if eb.ptr == nil {
		return 0
	}
	return uint32(C.editBufferGetLineStartOffset(eb.ptr, C.uint32_t(row)))
}

// GetTextRange returns text in the specified offset range.
func (eb *EditBuffer) GetTextRange(startOffset, endOffset uint32, maxLen int) string {
	if eb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editBufferGetTextRange(eb.ptr, C.uint32_t(startOffset), C.uint32_t(endOffset), (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// GetTextRangeByCoords returns text in the specified coordinate range.
func (eb *EditBuffer) GetTextRangeByCoords(startRow, startCol, endRow, endCol uint32, maxLen int) string {
	if eb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editBufferGetTextRangeByCoords(eb.ptr, C.uint32_t(startRow), C.uint32_t(startCol), C.uint32_t(endRow), C.uint32_t(endCol), (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// Undo undoes the last edit operation. Returns the undone text if successful.
func (eb *EditBuffer) Undo(maxLen int) string {
	if eb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editBufferUndo(eb.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// Redo redoes the last undone operation. Returns the redone text if successful.
func (eb *EditBuffer) Redo(maxLen int) string {
	if eb.ptr == nil || maxLen <= 0 {
		return ""
	}
	buffer := make([]byte, maxLen)
	actualLen := C.editBufferRedo(eb.ptr, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.size_t(maxLen))
	if actualLen == 0 {
		return ""
	}
	return string(buffer[:actualLen])
}

// CanUndo returns whether an undo operation is available.
func (eb *EditBuffer) CanUndo() bool {
	if eb.ptr == nil {
		return false
	}
	return bool(C.editBufferCanUndo(eb.ptr))
}

// CanRedo returns whether a redo operation is available.
func (eb *EditBuffer) CanRedo() bool {
	if eb.ptr == nil {
		return false
	}
	return bool(C.editBufferCanRedo(eb.ptr))
}

// ClearHistory clears the undo/redo history.
func (eb *EditBuffer) ClearHistory() {
	if eb.ptr == nil {
		return
	}
	C.editBufferClearHistory(eb.ptr)
}
