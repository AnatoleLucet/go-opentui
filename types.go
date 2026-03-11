package opentui

/*
#include "opentui.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// MeasureResult represents the result of text measurement
type MeasureResult struct {
	LineCount uint32
	MaxWidth  uint32
}

// CursorState represents the current cursor state
type CursorState struct {
	X        uint32
	Y        uint32
	Visible  bool
	Style    uint8 // 0=block, 1=line, 2=underline
	Blinking bool
	Color    RGBA
}

// LogicalCursor represents a logical cursor position (for EditBuffer)
type LogicalCursor struct {
	Row    uint32
	Col    uint32
	Offset uint32
}

// VisualCursor represents a visual cursor position (for EditorView)
type VisualCursor struct {
	VisualRow  uint32
	VisualCol  uint32
	LogicalRow uint32
	LogicalCol uint32
	Offset     uint32
}

// LineInfo represents line information from TextBufferView/EditorView
type LineInfo struct {
	Starts   []uint32 // Line start offsets
	Widths   []uint32 // Line widths
	Sources  []uint32 // Source line indices
	Wraps    []uint32 // Wrap flags
	MaxWidth uint32   // Maximum line width
}

// lineInfoFromC converts a C LineInfo struct to Go LineInfo
func lineInfoFromC(info *C.LineInfo) LineInfo {
	if info == nil {
		return LineInfo{}
	}

	result := LineInfo{
		MaxWidth: uint32(info.max_width),
	}

	if info.starts_len > 0 && info.starts_ptr != nil {
		result.Starts = cArrayToSlice((*uint32)(unsafe.Pointer(info.starts_ptr)), int(info.starts_len))
	}
	if info.widths_len > 0 && info.widths_ptr != nil {
		result.Widths = cArrayToSlice((*uint32)(unsafe.Pointer(info.widths_ptr)), int(info.widths_len))
	}
	if info.sources_len > 0 && info.sources_ptr != nil {
		result.Sources = cArrayToSlice((*uint32)(unsafe.Pointer(info.sources_ptr)), int(info.sources_len))
	}
	if info.wraps_len > 0 && info.wraps_ptr != nil {
		result.Wraps = cArrayToSlice((*uint32)(unsafe.Pointer(info.wraps_ptr)), int(info.wraps_len))
	}

	return result
}

// Highlight represents a syntax highlight range
type Highlight struct {
	Start    uint32
	End      uint32
	StyleID  uint32
	Priority uint8
	HLRef    uint16
}

// StyledChunk represents a chunk of styled text
type StyledChunk struct {
	Text       string
	Foreground *RGBA
	Background *RGBA
	Attributes uint32
	URL        string // Optional hyperlink URL
}

// Cell represents a single terminal cell with character, colors, and attributes
type Cell struct {
	Char       rune   // Unicode character
	Foreground RGBA   // Foreground color
	Background RGBA   // Background color
	Attributes uint32 // Text attributes (bold, italic, etc.) with optional link ID in upper bits
}

// Text attributes constants
const (
	AttrBold      uint32 = 1 << 0
	AttrDim       uint32 = 1 << 1
	AttrItalic    uint32 = 1 << 2
	AttrUnderline uint32 = 1 << 3
	AttrBlink     uint32 = 1 << 4
	AttrReverse   uint32 = 1 << 5
	AttrStrike    uint32 = 1 << 6
)

// Stats holds renderer statistics
type Stats struct {
	Time              float64
	FPS               uint32
	FrameCallbackTime float64
}

// MemoryStats holds memory usage statistics
type MemoryStats struct {
	HeapUsed     uint32
	HeapTotal    uint32
	ArrayBuffers uint32
}

// BoxOptions holds options for drawing boxes
type BoxOptions struct {
	Sides          BorderSides
	Fill           bool
	Title          string
	TitleAlignment TextAlignment
	BorderChars    BorderChars
}

type BorderChars struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
	TopT        rune
	BottomT     rune
	LeftT       rune
	RightT      rune
	Cross       rune
}

// SuperSampleFormat defines pixel formats for super-sampling
type SuperSampleFormat uint8

const (
	FormatRGBA SuperSampleFormat = iota
	FormatRGB
	FormatBGRA
	FormatBGR
)

// HitTestResult represents the result of a mouse hit test
type HitTestResult struct {
	ID    uint32
	Found bool
}

// Error represents an OpenTUI error
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// newError creates a new OpenTUI error
func newError(msg string) error {
	return &Error{Message: msg}
}

// EncodedChar represents a unicode character with its display width
type EncodedChar struct {
	CharCode uint32
	Width    uint8
}

// finalizer is a helper to set up automatic cleanup for CGO objects
func setFinalizer[T any](obj *T, cleanup func(*T)) {
	if obj != nil {
		runtime.SetFinalizer(obj, func(o *T) { cleanup(o) })
	}
}

// clearFinalizer removes the finalizer from an object
func clearFinalizer[T any](obj *T) {
	if obj != nil {
		runtime.SetFinalizer(obj, nil)
	}
}

// sliceToC converts a Go slice to C array parameters
func sliceToC[T any](slice []T) (*T, C.size_t) {
	if len(slice) == 0 {
		return nil, 0
	}
	return (*T)(unsafe.Pointer(&slice[0])), C.size_t(len(slice))
}

// cArrayToSlice converts a C array to a Go slice (read-only view)
func cArrayToSlice[T any](ptr *T, length int) []T {
	if ptr == nil || length == 0 {
		return nil
	}
	return unsafe.Slice(ptr, length)
}

// runesToC converts a rune slice to uint32 C array
func runesToC(runes []rune) *C.uint32_t {
	if len(runes) == 0 {
		return nil
	}
	// Convert runes to uint32
	uint32s := make([]uint32, len(runes))
	for i, r := range runes {
		uint32s[i] = uint32(r)
	}
	return (*C.uint32_t)(unsafe.Pointer(&uint32s[0]))
}

// UnicodeMethod represents the unicode width calculation method
type UnicodeMethod uint8

const (
	UnicodeMethodWcwidth UnicodeMethod = 0
	UnicodeMethodUnicode UnicodeMethod = 1
)

// Capabilities represents terminal capabilities
type Capabilities struct {
	KittyKeyboard             bool          // Terminal supports Kitty keyboard protocol
	KittyGraphics             bool          // Terminal supports Kitty graphics protocol
	RGB                       bool          // Terminal supports 24-bit true color
	Unicode                   UnicodeMethod // Unicode width calculation method
	SGRPixels                 bool          // Terminal supports SGR pixel mouse reporting
	ColorSchemeUpdates        bool          // Terminal supports color scheme change notifications
	ExplicitWidth             bool          // Terminal supports explicit character width
	ScaledText                bool          // Terminal supports scaled text
	Sixel                     bool          // Terminal supports Sixel graphics
	FocusTracking             bool          // Terminal supports focus tracking events
	Sync                      bool          // Terminal supports synchronized output
	BracketedPaste            bool          // Terminal supports bracketed paste mode
	Hyperlinks                bool          // Terminal supports hyperlinks (OSC 8)
	ExplicitCursorPositioning bool          // Terminal supports explicit cursor positioning
	TermName                  string        // Terminal name
	TermVersion               string        // Terminal version
	TermFromXTVersion         bool          // Terminal name/version from XTVERSION response
}

// Border styles
var (
	BorderCharsSingle = BorderChars{
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
		TopT:        '┬',
		BottomT:     '┴',
		LeftT:       '├',
		RightT:      '┤',
		Cross:       '┼',
	}
	BorderCharsDouble = BorderChars{
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
		TopT:        '╦',
		BottomT:     '╩',
		LeftT:       '╠',
		RightT:      '╣',
		Cross:       '╬',
	}
	BorderCharsRounded = BorderChars{
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
		Horizontal:  '─',
		Vertical:    '│',
		TopT:        '┬',
		BottomT:     '┴',
		LeftT:       '├',
		RightT:      '┤',
		Cross:       '┼',
	}
	BorderCharsHeavy = BorderChars{
		TopLeft:     '┏',
		TopRight:    '┓',
		BottomLeft:  '┗',
		BottomRight: '┛',
		Horizontal:  '━',
		Vertical:    '┃',
		TopT:        '┳',
		BottomT:     '┻',
		LeftT:       '┣',
		RightT:      '┫',
		Cross:       '╋',
	}
)
