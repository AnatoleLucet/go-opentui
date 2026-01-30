#ifndef OPENTUI_H
#define OPENTUI_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

// Opaque type definitions
typedef struct CliRenderer CliRenderer;
typedef struct OptimizedBuffer OptimizedBuffer;
typedef struct TextBuffer TextBuffer;
typedef struct TextBufferView TextBufferView;
typedef struct EditBuffer EditBuffer;
typedef struct EditorView EditorView;
typedef struct SyntaxStyle SyntaxStyle;

// RGBA color type - array of 4 floats [r, g, b, a]
typedef float RGBA[4];

// ============================================================================
// Common Structures
// ============================================================================

// Unicode width method enum
typedef enum {
    UnicodeMethodWcwidth = 0,
    UnicodeMethodUnicode = 1,
} UnicodeMethod;

// Terminal capabilities structure
typedef struct {
    uint8_t kitty_keyboard;
    uint8_t kitty_graphics;
    uint8_t rgb;
    uint8_t unicode;               // UnicodeMethod enum
    uint8_t sgr_pixels;
    uint8_t color_scheme_updates;
    uint8_t explicit_width;
    uint8_t scaled_text;
    uint8_t sixel;
    uint8_t focus_tracking;
    uint8_t sync;
    uint8_t bracketed_paste;
    uint8_t hyperlinks;
    uint8_t explicit_cursor_positioning;
    const char* term_name;
    uint64_t term_name_len;
    const char* term_version;
    uint64_t term_version_len;
    uint8_t term_from_xtversion;
} Capabilities;

// Cursor state
typedef struct {
    uint32_t x;
    uint32_t y;
    bool visible;
    uint8_t style;
    bool blinking;
    float r;
    float g;
    float b;
    float a;
} CursorState;

// Measurement result
typedef struct {
    uint32_t line_count;
    uint32_t max_width;
} MeasureResult;

// Logical cursor position (for EditBuffer)
typedef struct {
    uint32_t row;
    uint32_t col;
    uint32_t offset;
} LogicalCursor;

// Visual cursor position (for EditorView)
typedef struct {
    uint32_t visual_row;
    uint32_t visual_col;
    uint32_t logical_row;
    uint32_t logical_col;
    uint32_t offset;
} VisualCursor;

// Line information structure
typedef struct {
    const uint32_t* starts_ptr;
    uint32_t starts_len;
    const uint32_t* widths_ptr;
    uint32_t widths_len;
    const uint32_t* sources_ptr;
    uint32_t sources_len;
    const uint32_t* wraps_ptr;
    uint32_t wraps_len;
    uint32_t max_width;
} LineInfo;

// Highlight structure
typedef struct {
    uint32_t start;
    uint32_t end;
    uint32_t style_id;
    uint8_t priority;
    uint16_t hl_ref;
} Highlight;

// Styled text chunk for textBufferSetStyledText
typedef struct {
    const uint8_t* text_ptr;
    size_t text_len;
    const float* fg;      // nullable
    const float* bg;      // nullable
    uint32_t attributes;
    uint32_t link_id;     // 0 if no link
} StyledChunk;

// Encoded unicode character
typedef struct {
    uint32_t char_code;
    uint8_t width;
} EncodedChar;

// ============================================================================
// Callback Types
// ============================================================================

typedef void (*LogCallback)(uint8_t level, const uint8_t* msgPtr, size_t msgLen);
typedef void (*EventCallback)(const uint8_t* namePtr, size_t nameLen, const uint8_t* dataPtr, size_t dataLen);

// ============================================================================
// Global/Utility Functions
// ============================================================================

void setLogCallback(LogCallback callback);
void setEventCallback(EventCallback callback);
size_t getArenaAllocatedBytes(void);

// ============================================================================
// Renderer Management Functions
// ============================================================================

CliRenderer* createRenderer(uint32_t width, uint32_t height, bool testing);
void destroyRenderer(CliRenderer* renderer);
void setUseThread(CliRenderer* renderer, bool useThread);
void setBackgroundColor(CliRenderer* renderer, const float* color);
void setRenderOffset(CliRenderer* renderer, uint32_t offset);
void updateStats(CliRenderer* renderer, double time, uint32_t fps, double frameCallbackTime);
void updateMemoryStats(CliRenderer* renderer, uint32_t heapUsed, uint32_t heapTotal, uint32_t arrayBuffers);
OptimizedBuffer* getNextBuffer(CliRenderer* renderer);
OptimizedBuffer* getCurrentBuffer(CliRenderer* renderer);
void render(CliRenderer* renderer, bool force);
void resizeRenderer(CliRenderer* renderer, uint32_t width, uint32_t height);
void enableMouse(CliRenderer* renderer, bool enableMovement);
void disableMouse(CliRenderer* renderer);
void queryPixelResolution(CliRenderer* renderer);
void setTerminalTitle(CliRenderer* renderer, const uint8_t* title, size_t titleLen);
void suspendRenderer(CliRenderer* renderer);
void resumeRenderer(CliRenderer* renderer);
void writeOut(CliRenderer* renderer, const uint8_t* data, size_t dataLen);

// ============================================================================
// Buffer Management Functions
// ============================================================================

OptimizedBuffer* createOptimizedBuffer(uint32_t width, uint32_t height, bool respectAlpha, uint8_t widthMethod, const uint8_t* id, size_t idLen);
void destroyOptimizedBuffer(OptimizedBuffer* buffer);
void destroyFrameBuffer(OptimizedBuffer* frameBuffer);
uint32_t getBufferWidth(OptimizedBuffer* buffer);
uint32_t getBufferHeight(OptimizedBuffer* buffer);

// ============================================================================
// Buffer Drawing Functions
// ============================================================================

void bufferClear(OptimizedBuffer* buffer, const float* bg);
uint32_t* bufferGetCharPtr(OptimizedBuffer* buffer);
float* bufferGetFgPtr(OptimizedBuffer* buffer);
float* bufferGetBgPtr(OptimizedBuffer* buffer);
uint32_t* bufferGetAttributesPtr(OptimizedBuffer* buffer);
size_t bufferGetId(OptimizedBuffer* buffer, uint8_t* outPtr, size_t maxLen);
uint32_t bufferGetRealCharSize(OptimizedBuffer* buffer);
uint32_t bufferWriteResolvedChars(OptimizedBuffer* buffer, uint8_t* outputPtr, size_t outputLen, bool addLineBreaks);
bool bufferGetRespectAlpha(OptimizedBuffer* buffer);
void bufferSetRespectAlpha(OptimizedBuffer* buffer, bool respectAlpha);
void bufferDrawText(OptimizedBuffer* buffer, const uint8_t* text, size_t textLen, uint32_t x, uint32_t y, const float* fg, const float* bg, uint32_t attributes);
void bufferSetCellWithAlphaBlending(OptimizedBuffer* buffer, uint32_t x, uint32_t y, uint32_t char_code, const float* fg, const float* bg, uint32_t attributes);
void bufferSetCell(OptimizedBuffer* buffer, uint32_t x, uint32_t y, uint32_t char_code, const float* fg, const float* bg, uint32_t attributes);
void bufferDrawChar(OptimizedBuffer* buffer, int32_t x, int32_t y, uint32_t char_code, const float* fg, const float* bg, uint32_t attributes);
void bufferFillRect(OptimizedBuffer* buffer, uint32_t x, uint32_t y, uint32_t width, uint32_t height, const float* bg);
void bufferDrawPackedBuffer(OptimizedBuffer* buffer, const uint8_t* data, size_t dataLen, uint32_t posX, uint32_t posY, uint32_t terminalWidthCells, uint32_t terminalHeightCells);
void bufferDrawSuperSampleBuffer(OptimizedBuffer* buffer, uint32_t x, uint32_t y, const uint8_t* pixelData, size_t len, uint8_t format, uint32_t alignedBytesPerRow);
void bufferDrawGrayscaleBuffer(OptimizedBuffer* buffer, int32_t posX, int32_t posY, const float* intensities, uint32_t srcWidth, uint32_t srcHeight, const float* fg, const float* bg);
void bufferDrawGrayscaleBufferSupersampled(OptimizedBuffer* buffer, int32_t posX, int32_t posY, const float* intensities, uint32_t srcWidth, uint32_t srcHeight, const float* fg, const float* bg);
void bufferDrawBox(OptimizedBuffer* buffer, int32_t x, int32_t y, uint32_t width, uint32_t height, const uint32_t* borderChars, uint32_t packedOptions, const float* borderColor, const float* backgroundColor, const uint8_t* title, uint32_t titleLen);
void bufferResize(OptimizedBuffer* buffer, uint32_t width, uint32_t height);
void drawFrameBuffer(OptimizedBuffer* target, int32_t destX, int32_t destY, OptimizedBuffer* frameBuffer, uint32_t sourceX, uint32_t sourceY, uint32_t sourceWidth, uint32_t sourceHeight);

// Scissor rect functions
void bufferPushScissorRect(OptimizedBuffer* buffer, int32_t x, int32_t y, uint32_t width, uint32_t height);
void bufferPopScissorRect(OptimizedBuffer* buffer);
void bufferClearScissorRects(OptimizedBuffer* buffer);

// Opacity functions
void bufferPushOpacity(OptimizedBuffer* buffer, float opacity);
void bufferPopOpacity(OptimizedBuffer* buffer);
float bufferGetCurrentOpacity(OptimizedBuffer* buffer);
void bufferClearOpacity(OptimizedBuffer* buffer);

// ============================================================================
// Cursor Functions
// ============================================================================

void setCursorPosition(CliRenderer* renderer, int32_t x, int32_t y, bool visible);
void setCursorStyle(CliRenderer* renderer, const uint8_t* style, size_t styleLen, bool blinking);
void setCursorColor(CliRenderer* renderer, const float* color);
void getCursorState(CliRenderer* renderer, CursorState* outState);

// ============================================================================
// Terminal Capability Functions
// ============================================================================

void getTerminalCapabilities(CliRenderer* renderer, Capabilities* caps);
void processCapabilityResponse(CliRenderer* renderer, const uint8_t* response, size_t responseLen);
void setHyperlinksCapability(CliRenderer* renderer, bool enabled);

// ============================================================================
// Hyperlink Functions
// ============================================================================

uint32_t linkAlloc(const uint8_t* url, size_t urlLen);
size_t linkGetUrl(uint32_t id, uint8_t* outPtr, size_t maxLen);
uint32_t attributesWithLink(uint32_t baseAttributes, uint32_t linkId);
uint32_t attributesGetLinkId(uint32_t attributes);
void clearGlobalLinkPool(void);

// ============================================================================
// Debug and Utility Functions
// ============================================================================

void setDebugOverlay(CliRenderer* renderer, bool enabled, uint8_t corner);
void clearTerminal(CliRenderer* renderer);
void dumpHitGrid(CliRenderer* renderer);
void dumpBuffers(CliRenderer* renderer, int64_t timestamp);
void dumpStdoutBuffer(CliRenderer* renderer, int64_t timestamp);

// ============================================================================
// Hit Grid Functions
// ============================================================================

void addToHitGrid(CliRenderer* renderer, int32_t x, int32_t y, uint32_t width, uint32_t height, uint32_t id);
void clearCurrentHitGrid(CliRenderer* renderer);
void hitGridPushScissorRect(CliRenderer* renderer, int32_t x, int32_t y, uint32_t width, uint32_t height);
void hitGridPopScissorRect(CliRenderer* renderer);
void hitGridClearScissorRects(CliRenderer* renderer);
void addToCurrentHitGridClipped(CliRenderer* renderer, int32_t x, int32_t y, uint32_t width, uint32_t height, uint32_t id);
uint32_t checkHit(CliRenderer* renderer, uint32_t x, uint32_t y);
bool getHitGridDirty(CliRenderer* renderer);

// ============================================================================
// Keyboard and Terminal Setup Functions
// ============================================================================

void enableKittyKeyboard(CliRenderer* renderer, uint8_t flags);
void disableKittyKeyboard(CliRenderer* renderer);
void setKittyKeyboardFlags(CliRenderer* renderer, uint8_t flags);
uint8_t getKittyKeyboardFlags(CliRenderer* renderer);
void setupTerminal(CliRenderer* renderer, bool useAlternateScreen);

// ============================================================================
// TextBuffer Functions
// ============================================================================

TextBuffer* createTextBuffer(uint8_t widthMethod);
void destroyTextBuffer(TextBuffer* textBuffer);
void textBufferReset(TextBuffer* textBuffer);
void textBufferClear(TextBuffer* textBuffer);
uint32_t textBufferGetLength(TextBuffer* textBuffer);
uint32_t textBufferGetByteSize(TextBuffer* textBuffer);
uint32_t textBufferGetLineCount(TextBuffer* textBuffer);
void textBufferAppend(TextBuffer* textBuffer, const uint8_t* text, size_t textLen);
void textBufferSetDefaultFg(TextBuffer* textBuffer, const float* fg);
void textBufferSetDefaultBg(TextBuffer* textBuffer, const float* bg);
void textBufferSetDefaultAttributes(TextBuffer* textBuffer, const uint32_t* attr);
void textBufferResetDefaults(TextBuffer* textBuffer);
uint8_t textBufferGetTabWidth(TextBuffer* textBuffer);
void textBufferSetTabWidth(TextBuffer* textBuffer, uint8_t width);
size_t textBufferGetPlainText(TextBuffer* textBuffer, uint8_t* outPtr, size_t maxLen);
uint16_t textBufferRegisterMemBuffer(TextBuffer* textBuffer, const uint8_t* data, size_t dataLen, bool owned);
bool textBufferReplaceMemBuffer(TextBuffer* textBuffer, uint8_t id, const uint8_t* data, size_t dataLen, bool owned);
void textBufferClearMemRegistry(TextBuffer* textBuffer);
void textBufferSetTextFromMem(TextBuffer* textBuffer, uint8_t id);
void textBufferAppendFromMemId(TextBuffer* textBuffer, uint8_t id);
bool textBufferLoadFile(TextBuffer* textBuffer, const uint8_t* path, size_t pathLen);

// Text range retrieval
size_t textBufferGetTextRange(TextBuffer* textBuffer, uint32_t startOffset, uint32_t endOffset, uint8_t* outPtr, size_t maxLen);
size_t textBufferGetTextRangeByCoords(TextBuffer* textBuffer, uint32_t startRow, uint32_t startCol, uint32_t endRow, uint32_t endCol, uint8_t* outPtr, size_t maxLen);

// Styled text
void textBufferSetStyledText(TextBuffer* textBuffer, const StyledChunk* chunks, size_t chunkCount);

// Syntax highlighting
void textBufferSetSyntaxStyle(TextBuffer* textBuffer, SyntaxStyle* style);
void textBufferAddHighlight(TextBuffer* textBuffer, uint32_t lineIdx, const Highlight* highlight);
void textBufferAddHighlightByCharRange(TextBuffer* textBuffer, const Highlight* highlight);
void textBufferRemoveHighlightsByRef(TextBuffer* textBuffer, uint16_t hlRef);
void textBufferClearLineHighlights(TextBuffer* textBuffer, uint32_t lineIdx);
void textBufferClearAllHighlights(TextBuffer* textBuffer);
uint32_t textBufferGetHighlightCount(TextBuffer* textBuffer);
size_t textBufferGetLineHighlightsPtr(TextBuffer* textBuffer, uint32_t lineIdx, Highlight** outPtr);
void textBufferFreeLineHighlights(const Highlight* ptr, size_t count);

// ============================================================================
// TextBufferView Functions
// ============================================================================

TextBufferView* createTextBufferView(TextBuffer* textBuffer);
void destroyTextBufferView(TextBufferView* view);
void textBufferViewSetWrapWidth(TextBufferView* view, uint32_t width);
void textBufferViewSetWrapMode(TextBufferView* view, uint8_t mode);
void textBufferViewSetViewportSize(TextBufferView* view, uint32_t width, uint32_t height);
void textBufferViewSetViewport(TextBufferView* view, uint32_t x, uint32_t y, uint32_t width, uint32_t height);
uint32_t textBufferViewGetVirtualLineCount(TextBufferView* view);
bool textBufferViewMeasureForDimensions(TextBufferView* view, uint32_t width, uint32_t height, MeasureResult* outResult);

// Selection functions
void textBufferViewSetSelection(TextBufferView* view, uint32_t start, uint32_t end, const float* bgColor, const float* fgColor);
void textBufferViewResetSelection(TextBufferView* view);
uint64_t textBufferViewGetSelectionInfo(TextBufferView* view);
bool textBufferViewSetLocalSelection(TextBufferView* view, int32_t anchorX, int32_t anchorY, int32_t focusX, int32_t focusY, const float* bgColor, const float* fgColor);
void textBufferViewUpdateSelection(TextBufferView* view, uint32_t end, const float* bgColor, const float* fgColor);
bool textBufferViewUpdateLocalSelection(TextBufferView* view, int32_t anchorX, int32_t anchorY, int32_t focusX, int32_t focusY, const float* bgColor, const float* fgColor);
void textBufferViewResetLocalSelection(TextBufferView* view);

// Text retrieval
size_t textBufferViewGetSelectedText(TextBufferView* view, uint8_t* outPtr, size_t maxLen);
size_t textBufferViewGetPlainText(TextBufferView* view, uint8_t* outPtr, size_t maxLen);

// Line information
void textBufferViewGetLineInfoDirect(TextBufferView* view, LineInfo* outInfo);
void textBufferViewGetLogicalLineInfoDirect(TextBufferView* view, LineInfo* outInfo);

// Tab indicator
void textBufferViewSetTabIndicator(TextBufferView* view, uint32_t indicator);
void textBufferViewSetTabIndicatorColor(TextBufferView* view, const float* color);

// Truncation
void textBufferViewSetTruncate(TextBufferView* view, bool truncate);

// Drawing
void bufferDrawTextBufferView(OptimizedBuffer* buffer, TextBufferView* view, int32_t x, int32_t y);

// ============================================================================
// EditBuffer Functions
// ============================================================================

EditBuffer* createEditBuffer(uint8_t widthMethod);
void destroyEditBuffer(EditBuffer* editBuffer);
TextBuffer* editBufferGetTextBuffer(EditBuffer* editBuffer);
uint16_t editBufferGetId(EditBuffer* editBuffer);
void editBufferDebugLogRope(EditBuffer* editBuffer);

// Text manipulation
void editBufferSetText(EditBuffer* editBuffer, const uint8_t* text, size_t textLen);
void editBufferSetTextFromMem(EditBuffer* editBuffer, uint8_t memId);
void editBufferReplaceText(EditBuffer* editBuffer, const uint8_t* text, size_t textLen);
void editBufferReplaceTextFromMem(EditBuffer* editBuffer, uint8_t memId);
size_t editBufferGetText(EditBuffer* editBuffer, uint8_t* outPtr, size_t maxLen);
void editBufferInsertChar(EditBuffer* editBuffer, const uint8_t* charPtr, size_t charLen);
void editBufferInsertText(EditBuffer* editBuffer, const uint8_t* text, size_t textLen);
void editBufferDeleteChar(EditBuffer* editBuffer);
void editBufferDeleteCharBackward(EditBuffer* editBuffer);
void editBufferDeleteRange(EditBuffer* editBuffer, uint32_t startRow, uint32_t startCol, uint32_t endRow, uint32_t endCol);
void editBufferNewLine(EditBuffer* editBuffer);
void editBufferDeleteLine(EditBuffer* editBuffer);
void editBufferClear(EditBuffer* editBuffer);

// Cursor control
void editBufferMoveCursorLeft(EditBuffer* editBuffer);
void editBufferMoveCursorRight(EditBuffer* editBuffer);
void editBufferMoveCursorUp(EditBuffer* editBuffer);
void editBufferMoveCursorDown(EditBuffer* editBuffer);
void editBufferGetCursor(EditBuffer* editBuffer, uint32_t* outRow, uint32_t* outCol);
void editBufferSetCursor(EditBuffer* editBuffer, uint32_t row, uint32_t col);
void editBufferSetCursorToLineCol(EditBuffer* editBuffer, uint32_t row, uint32_t col);
void editBufferSetCursorByOffset(EditBuffer* editBuffer, uint32_t offset);
void editBufferGotoLine(EditBuffer* editBuffer, uint32_t line);
void editBufferGetCursorPosition(EditBuffer* editBuffer, LogicalCursor* outCursor);

// Word/line boundaries
void editBufferGetNextWordBoundary(EditBuffer* editBuffer, LogicalCursor* outCursor);
void editBufferGetPrevWordBoundary(EditBuffer* editBuffer, LogicalCursor* outCursor);
void editBufferGetEOL(EditBuffer* editBuffer, LogicalCursor* outCursor);
bool editBufferOffsetToPosition(EditBuffer* editBuffer, uint32_t offset, LogicalCursor* outCursor);
uint32_t editBufferPositionToOffset(EditBuffer* editBuffer, uint32_t row, uint32_t col);
uint32_t editBufferGetLineStartOffset(EditBuffer* editBuffer, uint32_t row);

// Text range retrieval
size_t editBufferGetTextRange(EditBuffer* editBuffer, uint32_t startOffset, uint32_t endOffset, uint8_t* outPtr, size_t maxLen);
size_t editBufferGetTextRangeByCoords(EditBuffer* editBuffer, uint32_t startRow, uint32_t startCol, uint32_t endRow, uint32_t endCol, uint8_t* outPtr, size_t maxLen);

// Undo/Redo
size_t editBufferUndo(EditBuffer* editBuffer, uint8_t* outPtr, size_t maxLen);
size_t editBufferRedo(EditBuffer* editBuffer, uint8_t* outPtr, size_t maxLen);
bool editBufferCanUndo(EditBuffer* editBuffer);
bool editBufferCanRedo(EditBuffer* editBuffer);
void editBufferClearHistory(EditBuffer* editBuffer);

// ============================================================================
// EditorView Functions
// ============================================================================

EditorView* createEditorView(EditBuffer* editBuffer, uint32_t viewportWidth, uint32_t viewportHeight);
void destroyEditorView(EditorView* view);
TextBufferView* editorViewGetTextBufferView(EditorView* view);

// Viewport
void editorViewSetViewportSize(EditorView* view, uint32_t width, uint32_t height);
void editorViewSetViewport(EditorView* view, uint32_t x, uint32_t y, uint32_t width, uint32_t height, bool moveCursor);
bool editorViewGetViewport(EditorView* view, uint32_t* outX, uint32_t* outY, uint32_t* outWidth, uint32_t* outHeight);
void editorViewClearViewport(EditorView* view);
void editorViewSetScrollMargin(EditorView* view, float margin);
void editorViewSetWrapMode(EditorView* view, uint8_t mode);
uint32_t editorViewGetVirtualLineCount(EditorView* view);
uint32_t editorViewGetTotalVirtualLineCount(EditorView* view);

// Selection
void editorViewSetSelection(EditorView* view, uint32_t start, uint32_t end, const float* bgColor, const float* fgColor);
void editorViewResetSelection(EditorView* view);
uint64_t editorViewGetSelection(EditorView* view);
bool editorViewSetLocalSelection(EditorView* view, int32_t anchorX, int32_t anchorY, int32_t focusX, int32_t focusY, const float* bgColor, const float* fgColor, bool updateCursor, bool followCursor);
void editorViewUpdateSelection(EditorView* view, uint32_t end, const float* bgColor, const float* fgColor);
bool editorViewUpdateLocalSelection(EditorView* view, int32_t anchorX, int32_t anchorY, int32_t focusX, int32_t focusY, const float* bgColor, const float* fgColor, bool updateCursor, bool followCursor);
void editorViewResetLocalSelection(EditorView* view);
size_t editorViewGetSelectedTextBytes(EditorView* view, uint8_t* outPtr, size_t maxLen);
void editorViewDeleteSelectedText(EditorView* view);

// Cursor
void editorViewGetCursor(EditorView* view, uint32_t* outRow, uint32_t* outCol);
void editorViewGetVisualCursor(EditorView* view, VisualCursor* outCursor);
void editorViewSetCursorByOffset(EditorView* view, uint32_t offset);
void editorViewMoveUpVisual(EditorView* view);
void editorViewMoveDownVisual(EditorView* view);

// Word/line boundaries
void editorViewGetNextWordBoundary(EditorView* view, VisualCursor* outCursor);
void editorViewGetPrevWordBoundary(EditorView* view, VisualCursor* outCursor);
void editorViewGetEOL(EditorView* view, VisualCursor* outCursor);
void editorViewGetVisualSOL(EditorView* view, VisualCursor* outCursor);
void editorViewGetVisualEOL(EditorView* view, VisualCursor* outCursor);

// Text retrieval
size_t editorViewGetText(EditorView* view, uint8_t* outPtr, size_t maxLen);

// Line information
void editorViewGetLineInfoDirect(EditorView* view, LineInfo* outInfo);
void editorViewGetLogicalLineInfoDirect(EditorView* view, LineInfo* outInfo);

// Placeholder and styling
void editorViewSetPlaceholderStyledText(EditorView* view, const StyledChunk* chunks, size_t chunkCount);
void editorViewSetTabIndicator(EditorView* view, uint32_t indicator);
void editorViewSetTabIndicatorColor(EditorView* view, const float* color);

// Drawing
void bufferDrawEditorView(OptimizedBuffer* buffer, EditorView* view, int32_t x, int32_t y);

// ============================================================================
// SyntaxStyle Functions
// ============================================================================

SyntaxStyle* createSyntaxStyle(void);
void destroySyntaxStyle(SyntaxStyle* style);
uint32_t syntaxStyleRegister(SyntaxStyle* style, const uint8_t* name, size_t nameLen, const float* fg, const float* bg, uint32_t attributes);
uint32_t syntaxStyleResolveByName(SyntaxStyle* style, const uint8_t* name, size_t nameLen);
size_t syntaxStyleGetStyleCount(SyntaxStyle* style);

// ============================================================================
// Unicode Encoding Functions
// ============================================================================

bool encodeUnicode(const uint8_t* text, size_t textLen, uint8_t widthMethod, EncodedChar** outChars, size_t* outLen);
void freeUnicode(const EncodedChar* chars, size_t len);

#ifdef __cplusplus
}
#endif

#endif // OPENTUI_H
