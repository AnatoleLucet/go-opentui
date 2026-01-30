package opentui

/*
#include "opentui.h"
#include <stdlib.h>
*/
import "C"

// SyntaxStyle wraps the SyntaxStyle from the C library.
// It manages named syntax highlighting styles.
type SyntaxStyle struct {
	ptr *C.SyntaxStyle
}

// NewSyntaxStyle creates a new syntax style registry.
func NewSyntaxStyle() *SyntaxStyle {
	ptr := C.createSyntaxStyle()
	if ptr == nil {
		return nil
	}

	ss := &SyntaxStyle{ptr: ptr}
	setFinalizer(ss, func(ss *SyntaxStyle) { ss.Close() })
	return ss
}

// Close releases the syntax style's resources.
func (ss *SyntaxStyle) Close() error {
	if ss.ptr != nil {
		clearFinalizer(ss)
		C.destroySyntaxStyle(ss.ptr)
		ss.ptr = nil
	}
	return nil
}

// Valid checks if the syntax style is still valid.
func (ss *SyntaxStyle) Valid() bool {
	return ss.ptr != nil
}

// Register registers a new syntax style and returns its ID.
// fg and bg can be nil to use default colors.
func (ss *SyntaxStyle) Register(name string, fg, bg *RGBA, attributes uint32) uint32 {
	if ss.ptr == nil || len(name) == 0 {
		return 0
	}

	namePtr, nameLen := stringToC(name)

	var fgPtr, bgPtr *C.float
	if fg != nil {
		fgPtr = fg.toCFloat()
	}
	if bg != nil {
		bgPtr = bg.toCFloat()
	}

	return uint32(C.syntaxStyleRegister(ss.ptr, namePtr, nameLen, fgPtr, bgPtr, C.uint32_t(attributes)))
}

// ResolveByName returns the style ID for a given name.
// Returns 0 if the style is not found (0 is also a valid style ID, so check GetStyleCount).
func (ss *SyntaxStyle) ResolveByName(name string) uint32 {
	if ss.ptr == nil || len(name) == 0 {
		return 0
	}

	namePtr, nameLen := stringToC(name)
	return uint32(C.syntaxStyleResolveByName(ss.ptr, namePtr, nameLen))
}

// GetStyleCount returns the number of registered styles.
func (ss *SyntaxStyle) GetStyleCount() int {
	if ss.ptr == nil {
		return 0
	}
	return int(C.syntaxStyleGetStyleCount(ss.ptr))
}
