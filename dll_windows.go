//go:build windows
// +build windows

package opentui

// Static linking is now used for Windows.
// The opentui.lib file is linked directly into the binary at compile time.
// No DLL extraction or runtime loading is required.
