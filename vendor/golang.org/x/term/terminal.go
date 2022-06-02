// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"bytes"
	"io"
	"runtime"
	"strconv"
	"sync"
	"unicode/utf8"
)

// EscapeCodes contains escape sequences that can be written to the terminal in
// order to achieve different styles of text.
type EscapeCodes struct ***REMOVED***
	// Foreground colors
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White []byte

	// Reset all attributes
	Reset []byte
***REMOVED***

var vt100EscapeCodes = EscapeCodes***REMOVED***
	Black:   []byte***REMOVED***keyEscape, '[', '3', '0', 'm'***REMOVED***,
	Red:     []byte***REMOVED***keyEscape, '[', '3', '1', 'm'***REMOVED***,
	Green:   []byte***REMOVED***keyEscape, '[', '3', '2', 'm'***REMOVED***,
	Yellow:  []byte***REMOVED***keyEscape, '[', '3', '3', 'm'***REMOVED***,
	Blue:    []byte***REMOVED***keyEscape, '[', '3', '4', 'm'***REMOVED***,
	Magenta: []byte***REMOVED***keyEscape, '[', '3', '5', 'm'***REMOVED***,
	Cyan:    []byte***REMOVED***keyEscape, '[', '3', '6', 'm'***REMOVED***,
	White:   []byte***REMOVED***keyEscape, '[', '3', '7', 'm'***REMOVED***,

	Reset: []byte***REMOVED***keyEscape, '[', '0', 'm'***REMOVED***,
***REMOVED***

// Terminal contains the state for running a VT100 terminal that is capable of
// reading lines of input.
type Terminal struct ***REMOVED***
	// AutoCompleteCallback, if non-null, is called for each keypress with
	// the full input line and the current position of the cursor (in
	// bytes, as an index into |line|). If it returns ok=false, the key
	// press is processed normally. Otherwise it returns a replacement line
	// and the new cursor position.
	AutoCompleteCallback func(line string, pos int, key rune) (newLine string, newPos int, ok bool)

	// Escape contains a pointer to the escape codes for this terminal.
	// It's always a valid pointer, although the escape codes themselves
	// may be empty if the terminal doesn't support them.
	Escape *EscapeCodes

	// lock protects the terminal and the state in this object from
	// concurrent processing of a key press and a Write() call.
	lock sync.Mutex

	c      io.ReadWriter
	prompt []rune

	// line is the current line being entered.
	line []rune
	// pos is the logical position of the cursor in line
	pos int
	// echo is true if local echo is enabled
	echo bool
	// pasteActive is true iff there is a bracketed paste operation in
	// progress.
	pasteActive bool

	// cursorX contains the current X value of the cursor where the left
	// edge is 0. cursorY contains the row number where the first row of
	// the current line is 0.
	cursorX, cursorY int
	// maxLine is the greatest value of cursorY so far.
	maxLine int

	termWidth, termHeight int

	// outBuf contains the terminal data to be sent.
	outBuf []byte
	// remainder contains the remainder of any partial key sequences after
	// a read. It aliases into inBuf.
	remainder []byte
	inBuf     [256]byte

	// history contains previously entered commands so that they can be
	// accessed with the up and down keys.
	history stRingBuffer
	// historyIndex stores the currently accessed history entry, where zero
	// means the immediately previous entry.
	historyIndex int
	// When navigating up and down the history it's possible to return to
	// the incomplete, initial line. That value is stored in
	// historyPending.
	historyPending string
***REMOVED***

// NewTerminal runs a VT100 terminal on the given ReadWriter. If the ReadWriter is
// a local terminal, that terminal must first have been put into raw mode.
// prompt is a string that is written at the start of each input line (i.e.
// "> ").
func NewTerminal(c io.ReadWriter, prompt string) *Terminal ***REMOVED***
	return &Terminal***REMOVED***
		Escape:       &vt100EscapeCodes,
		c:            c,
		prompt:       []rune(prompt),
		termWidth:    80,
		termHeight:   24,
		echo:         true,
		historyIndex: -1,
	***REMOVED***
***REMOVED***

const (
	keyCtrlC     = 3
	keyCtrlD     = 4
	keyCtrlU     = 21
	keyEnter     = '\r'
	keyEscape    = 27
	keyBackspace = 127
	keyUnknown   = 0xd800 /* UTF-16 surrogate area */ + iota
	keyUp
	keyDown
	keyLeft
	keyRight
	keyAltLeft
	keyAltRight
	keyHome
	keyEnd
	keyDeleteWord
	keyDeleteLine
	keyClearScreen
	keyPasteStart
	keyPasteEnd
)

var (
	crlf       = []byte***REMOVED***'\r', '\n'***REMOVED***
	pasteStart = []byte***REMOVED***keyEscape, '[', '2', '0', '0', '~'***REMOVED***
	pasteEnd   = []byte***REMOVED***keyEscape, '[', '2', '0', '1', '~'***REMOVED***
)

// bytesToKey tries to parse a key sequence from b. If successful, it returns
// the key and the remainder of the input. Otherwise it returns utf8.RuneError.
func bytesToKey(b []byte, pasteActive bool) (rune, []byte) ***REMOVED***
	if len(b) == 0 ***REMOVED***
		return utf8.RuneError, nil
	***REMOVED***

	if !pasteActive ***REMOVED***
		switch b[0] ***REMOVED***
		case 1: // ^A
			return keyHome, b[1:]
		case 2: // ^B
			return keyLeft, b[1:]
		case 5: // ^E
			return keyEnd, b[1:]
		case 6: // ^F
			return keyRight, b[1:]
		case 8: // ^H
			return keyBackspace, b[1:]
		case 11: // ^K
			return keyDeleteLine, b[1:]
		case 12: // ^L
			return keyClearScreen, b[1:]
		case 23: // ^W
			return keyDeleteWord, b[1:]
		case 14: // ^N
			return keyDown, b[1:]
		case 16: // ^P
			return keyUp, b[1:]
		***REMOVED***
	***REMOVED***

	if b[0] != keyEscape ***REMOVED***
		if !utf8.FullRune(b) ***REMOVED***
			return utf8.RuneError, b
		***REMOVED***
		r, l := utf8.DecodeRune(b)
		return r, b[l:]
	***REMOVED***

	if !pasteActive && len(b) >= 3 && b[0] == keyEscape && b[1] == '[' ***REMOVED***
		switch b[2] ***REMOVED***
		case 'A':
			return keyUp, b[3:]
		case 'B':
			return keyDown, b[3:]
		case 'C':
			return keyRight, b[3:]
		case 'D':
			return keyLeft, b[3:]
		case 'H':
			return keyHome, b[3:]
		case 'F':
			return keyEnd, b[3:]
		***REMOVED***
	***REMOVED***

	if !pasteActive && len(b) >= 6 && b[0] == keyEscape && b[1] == '[' && b[2] == '1' && b[3] == ';' && b[4] == '3' ***REMOVED***
		switch b[5] ***REMOVED***
		case 'C':
			return keyAltRight, b[6:]
		case 'D':
			return keyAltLeft, b[6:]
		***REMOVED***
	***REMOVED***

	if !pasteActive && len(b) >= 6 && bytes.Equal(b[:6], pasteStart) ***REMOVED***
		return keyPasteStart, b[6:]
	***REMOVED***

	if pasteActive && len(b) >= 6 && bytes.Equal(b[:6], pasteEnd) ***REMOVED***
		return keyPasteEnd, b[6:]
	***REMOVED***

	// If we get here then we have a key that we don't recognise, or a
	// partial sequence. It's not clear how one should find the end of a
	// sequence without knowing them all, but it seems that [a-zA-Z~] only
	// appears at the end of a sequence.
	for i, c := range b[0:] ***REMOVED***
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '~' ***REMOVED***
			return keyUnknown, b[i+1:]
		***REMOVED***
	***REMOVED***

	return utf8.RuneError, b
***REMOVED***

// queue appends data to the end of t.outBuf
func (t *Terminal) queue(data []rune) ***REMOVED***
	t.outBuf = append(t.outBuf, []byte(string(data))...)
***REMOVED***

var eraseUnderCursor = []rune***REMOVED***' ', keyEscape, '[', 'D'***REMOVED***
var space = []rune***REMOVED***' '***REMOVED***

func isPrintable(key rune) bool ***REMOVED***
	isInSurrogateArea := key >= 0xd800 && key <= 0xdbff
	return key >= 32 && !isInSurrogateArea
***REMOVED***

// moveCursorToPos appends data to t.outBuf which will move the cursor to the
// given, logical position in the text.
func (t *Terminal) moveCursorToPos(pos int) ***REMOVED***
	if !t.echo ***REMOVED***
		return
	***REMOVED***

	x := visualLength(t.prompt) + pos
	y := x / t.termWidth
	x = x % t.termWidth

	up := 0
	if y < t.cursorY ***REMOVED***
		up = t.cursorY - y
	***REMOVED***

	down := 0
	if y > t.cursorY ***REMOVED***
		down = y - t.cursorY
	***REMOVED***

	left := 0
	if x < t.cursorX ***REMOVED***
		left = t.cursorX - x
	***REMOVED***

	right := 0
	if x > t.cursorX ***REMOVED***
		right = x - t.cursorX
	***REMOVED***

	t.cursorX = x
	t.cursorY = y
	t.move(up, down, left, right)
***REMOVED***

func (t *Terminal) move(up, down, left, right int) ***REMOVED***
	m := []rune***REMOVED******REMOVED***

	// 1 unit up can be expressed as ^[[A or ^[A
	// 5 units up can be expressed as ^[[5A

	if up == 1 ***REMOVED***
		m = append(m, keyEscape, '[', 'A')
	***REMOVED*** else if up > 1 ***REMOVED***
		m = append(m, keyEscape, '[')
		m = append(m, []rune(strconv.Itoa(up))...)
		m = append(m, 'A')
	***REMOVED***

	if down == 1 ***REMOVED***
		m = append(m, keyEscape, '[', 'B')
	***REMOVED*** else if down > 1 ***REMOVED***
		m = append(m, keyEscape, '[')
		m = append(m, []rune(strconv.Itoa(down))...)
		m = append(m, 'B')
	***REMOVED***

	if right == 1 ***REMOVED***
		m = append(m, keyEscape, '[', 'C')
	***REMOVED*** else if right > 1 ***REMOVED***
		m = append(m, keyEscape, '[')
		m = append(m, []rune(strconv.Itoa(right))...)
		m = append(m, 'C')
	***REMOVED***

	if left == 1 ***REMOVED***
		m = append(m, keyEscape, '[', 'D')
	***REMOVED*** else if left > 1 ***REMOVED***
		m = append(m, keyEscape, '[')
		m = append(m, []rune(strconv.Itoa(left))...)
		m = append(m, 'D')
	***REMOVED***

	t.queue(m)
***REMOVED***

func (t *Terminal) clearLineToRight() ***REMOVED***
	op := []rune***REMOVED***keyEscape, '[', 'K'***REMOVED***
	t.queue(op)
***REMOVED***

const maxLineLength = 4096

func (t *Terminal) setLine(newLine []rune, newPos int) ***REMOVED***
	if t.echo ***REMOVED***
		t.moveCursorToPos(0)
		t.writeLine(newLine)
		for i := len(newLine); i < len(t.line); i++ ***REMOVED***
			t.writeLine(space)
		***REMOVED***
		t.moveCursorToPos(newPos)
	***REMOVED***
	t.line = newLine
	t.pos = newPos
***REMOVED***

func (t *Terminal) advanceCursor(places int) ***REMOVED***
	t.cursorX += places
	t.cursorY += t.cursorX / t.termWidth
	if t.cursorY > t.maxLine ***REMOVED***
		t.maxLine = t.cursorY
	***REMOVED***
	t.cursorX = t.cursorX % t.termWidth

	if places > 0 && t.cursorX == 0 ***REMOVED***
		// Normally terminals will advance the current position
		// when writing a character. But that doesn't happen
		// for the last character in a line. However, when
		// writing a character (except a new line) that causes
		// a line wrap, the position will be advanced two
		// places.
		//
		// So, if we are stopping at the end of a line, we
		// need to write a newline so that our cursor can be
		// advanced to the next line.
		t.outBuf = append(t.outBuf, '\r', '\n')
	***REMOVED***
***REMOVED***

func (t *Terminal) eraseNPreviousChars(n int) ***REMOVED***
	if n == 0 ***REMOVED***
		return
	***REMOVED***

	if t.pos < n ***REMOVED***
		n = t.pos
	***REMOVED***
	t.pos -= n
	t.moveCursorToPos(t.pos)

	copy(t.line[t.pos:], t.line[n+t.pos:])
	t.line = t.line[:len(t.line)-n]
	if t.echo ***REMOVED***
		t.writeLine(t.line[t.pos:])
		for i := 0; i < n; i++ ***REMOVED***
			t.queue(space)
		***REMOVED***
		t.advanceCursor(n)
		t.moveCursorToPos(t.pos)
	***REMOVED***
***REMOVED***

// countToLeftWord returns then number of characters from the cursor to the
// start of the previous word.
func (t *Terminal) countToLeftWord() int ***REMOVED***
	if t.pos == 0 ***REMOVED***
		return 0
	***REMOVED***

	pos := t.pos - 1
	for pos > 0 ***REMOVED***
		if t.line[pos] != ' ' ***REMOVED***
			break
		***REMOVED***
		pos--
	***REMOVED***
	for pos > 0 ***REMOVED***
		if t.line[pos] == ' ' ***REMOVED***
			pos++
			break
		***REMOVED***
		pos--
	***REMOVED***

	return t.pos - pos
***REMOVED***

// countToRightWord returns then number of characters from the cursor to the
// start of the next word.
func (t *Terminal) countToRightWord() int ***REMOVED***
	pos := t.pos
	for pos < len(t.line) ***REMOVED***
		if t.line[pos] == ' ' ***REMOVED***
			break
		***REMOVED***
		pos++
	***REMOVED***
	for pos < len(t.line) ***REMOVED***
		if t.line[pos] != ' ' ***REMOVED***
			break
		***REMOVED***
		pos++
	***REMOVED***
	return pos - t.pos
***REMOVED***

// visualLength returns the number of visible glyphs in s.
func visualLength(runes []rune) int ***REMOVED***
	inEscapeSeq := false
	length := 0

	for _, r := range runes ***REMOVED***
		switch ***REMOVED***
		case inEscapeSeq:
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ***REMOVED***
				inEscapeSeq = false
			***REMOVED***
		case r == '\x1b':
			inEscapeSeq = true
		default:
			length++
		***REMOVED***
	***REMOVED***

	return length
***REMOVED***

// handleKey processes the given key and, optionally, returns a line of text
// that the user has entered.
func (t *Terminal) handleKey(key rune) (line string, ok bool) ***REMOVED***
	if t.pasteActive && key != keyEnter ***REMOVED***
		t.addKeyToLine(key)
		return
	***REMOVED***

	switch key ***REMOVED***
	case keyBackspace:
		if t.pos == 0 ***REMOVED***
			return
		***REMOVED***
		t.eraseNPreviousChars(1)
	case keyAltLeft:
		// move left by a word.
		t.pos -= t.countToLeftWord()
		t.moveCursorToPos(t.pos)
	case keyAltRight:
		// move right by a word.
		t.pos += t.countToRightWord()
		t.moveCursorToPos(t.pos)
	case keyLeft:
		if t.pos == 0 ***REMOVED***
			return
		***REMOVED***
		t.pos--
		t.moveCursorToPos(t.pos)
	case keyRight:
		if t.pos == len(t.line) ***REMOVED***
			return
		***REMOVED***
		t.pos++
		t.moveCursorToPos(t.pos)
	case keyHome:
		if t.pos == 0 ***REMOVED***
			return
		***REMOVED***
		t.pos = 0
		t.moveCursorToPos(t.pos)
	case keyEnd:
		if t.pos == len(t.line) ***REMOVED***
			return
		***REMOVED***
		t.pos = len(t.line)
		t.moveCursorToPos(t.pos)
	case keyUp:
		entry, ok := t.history.NthPreviousEntry(t.historyIndex + 1)
		if !ok ***REMOVED***
			return "", false
		***REMOVED***
		if t.historyIndex == -1 ***REMOVED***
			t.historyPending = string(t.line)
		***REMOVED***
		t.historyIndex++
		runes := []rune(entry)
		t.setLine(runes, len(runes))
	case keyDown:
		switch t.historyIndex ***REMOVED***
		case -1:
			return
		case 0:
			runes := []rune(t.historyPending)
			t.setLine(runes, len(runes))
			t.historyIndex--
		default:
			entry, ok := t.history.NthPreviousEntry(t.historyIndex - 1)
			if ok ***REMOVED***
				t.historyIndex--
				runes := []rune(entry)
				t.setLine(runes, len(runes))
			***REMOVED***
		***REMOVED***
	case keyEnter:
		t.moveCursorToPos(len(t.line))
		t.queue([]rune("\r\n"))
		line = string(t.line)
		ok = true
		t.line = t.line[:0]
		t.pos = 0
		t.cursorX = 0
		t.cursorY = 0
		t.maxLine = 0
	case keyDeleteWord:
		// Delete zero or more spaces and then one or more characters.
		t.eraseNPreviousChars(t.countToLeftWord())
	case keyDeleteLine:
		// Delete everything from the current cursor position to the
		// end of line.
		for i := t.pos; i < len(t.line); i++ ***REMOVED***
			t.queue(space)
			t.advanceCursor(1)
		***REMOVED***
		t.line = t.line[:t.pos]
		t.moveCursorToPos(t.pos)
	case keyCtrlD:
		// Erase the character under the current position.
		// The EOF case when the line is empty is handled in
		// readLine().
		if t.pos < len(t.line) ***REMOVED***
			t.pos++
			t.eraseNPreviousChars(1)
		***REMOVED***
	case keyCtrlU:
		t.eraseNPreviousChars(t.pos)
	case keyClearScreen:
		// Erases the screen and moves the cursor to the home position.
		t.queue([]rune("\x1b[2J\x1b[H"))
		t.queue(t.prompt)
		t.cursorX, t.cursorY = 0, 0
		t.advanceCursor(visualLength(t.prompt))
		t.setLine(t.line, t.pos)
	default:
		if t.AutoCompleteCallback != nil ***REMOVED***
			prefix := string(t.line[:t.pos])
			suffix := string(t.line[t.pos:])

			t.lock.Unlock()
			newLine, newPos, completeOk := t.AutoCompleteCallback(prefix+suffix, len(prefix), key)
			t.lock.Lock()

			if completeOk ***REMOVED***
				t.setLine([]rune(newLine), utf8.RuneCount([]byte(newLine)[:newPos]))
				return
			***REMOVED***
		***REMOVED***
		if !isPrintable(key) ***REMOVED***
			return
		***REMOVED***
		if len(t.line) == maxLineLength ***REMOVED***
			return
		***REMOVED***
		t.addKeyToLine(key)
	***REMOVED***
	return
***REMOVED***

// addKeyToLine inserts the given key at the current position in the current
// line.
func (t *Terminal) addKeyToLine(key rune) ***REMOVED***
	if len(t.line) == cap(t.line) ***REMOVED***
		newLine := make([]rune, len(t.line), 2*(1+len(t.line)))
		copy(newLine, t.line)
		t.line = newLine
	***REMOVED***
	t.line = t.line[:len(t.line)+1]
	copy(t.line[t.pos+1:], t.line[t.pos:])
	t.line[t.pos] = key
	if t.echo ***REMOVED***
		t.writeLine(t.line[t.pos:])
	***REMOVED***
	t.pos++
	t.moveCursorToPos(t.pos)
***REMOVED***

func (t *Terminal) writeLine(line []rune) ***REMOVED***
	for len(line) != 0 ***REMOVED***
		remainingOnLine := t.termWidth - t.cursorX
		todo := len(line)
		if todo > remainingOnLine ***REMOVED***
			todo = remainingOnLine
		***REMOVED***
		t.queue(line[:todo])
		t.advanceCursor(visualLength(line[:todo]))
		line = line[todo:]
	***REMOVED***
***REMOVED***

// writeWithCRLF writes buf to w but replaces all occurrences of \n with \r\n.
func writeWithCRLF(w io.Writer, buf []byte) (n int, err error) ***REMOVED***
	for len(buf) > 0 ***REMOVED***
		i := bytes.IndexByte(buf, '\n')
		todo := len(buf)
		if i >= 0 ***REMOVED***
			todo = i
		***REMOVED***

		var nn int
		nn, err = w.Write(buf[:todo])
		n += nn
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		buf = buf[todo:]

		if i >= 0 ***REMOVED***
			if _, err = w.Write(crlf); err != nil ***REMOVED***
				return n, err
			***REMOVED***
			n++
			buf = buf[1:]
		***REMOVED***
	***REMOVED***

	return n, nil
***REMOVED***

func (t *Terminal) Write(buf []byte) (n int, err error) ***REMOVED***
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.cursorX == 0 && t.cursorY == 0 ***REMOVED***
		// This is the easy case: there's nothing on the screen that we
		// have to move out of the way.
		return writeWithCRLF(t.c, buf)
	***REMOVED***

	// We have a prompt and possibly user input on the screen. We
	// have to clear it first.
	t.move(0 /* up */, 0 /* down */, t.cursorX /* left */, 0 /* right */)
	t.cursorX = 0
	t.clearLineToRight()

	for t.cursorY > 0 ***REMOVED***
		t.move(1 /* up */, 0, 0, 0)
		t.cursorY--
		t.clearLineToRight()
	***REMOVED***

	if _, err = t.c.Write(t.outBuf); err != nil ***REMOVED***
		return
	***REMOVED***
	t.outBuf = t.outBuf[:0]

	if n, err = writeWithCRLF(t.c, buf); err != nil ***REMOVED***
		return
	***REMOVED***

	t.writeLine(t.prompt)
	if t.echo ***REMOVED***
		t.writeLine(t.line)
	***REMOVED***

	t.moveCursorToPos(t.pos)

	if _, err = t.c.Write(t.outBuf); err != nil ***REMOVED***
		return
	***REMOVED***
	t.outBuf = t.outBuf[:0]
	return
***REMOVED***

// ReadPassword temporarily changes the prompt and reads a password, without
// echo, from the terminal.
func (t *Terminal) ReadPassword(prompt string) (line string, err error) ***REMOVED***
	t.lock.Lock()
	defer t.lock.Unlock()

	oldPrompt := t.prompt
	t.prompt = []rune(prompt)
	t.echo = false

	line, err = t.readLine()

	t.prompt = oldPrompt
	t.echo = true

	return
***REMOVED***

// ReadLine returns a line of input from the terminal.
func (t *Terminal) ReadLine() (line string, err error) ***REMOVED***
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.readLine()
***REMOVED***

func (t *Terminal) readLine() (line string, err error) ***REMOVED***
	// t.lock must be held at this point

	if t.cursorX == 0 && t.cursorY == 0 ***REMOVED***
		t.writeLine(t.prompt)
		t.c.Write(t.outBuf)
		t.outBuf = t.outBuf[:0]
	***REMOVED***

	lineIsPasted := t.pasteActive

	for ***REMOVED***
		rest := t.remainder
		lineOk := false
		for !lineOk ***REMOVED***
			var key rune
			key, rest = bytesToKey(rest, t.pasteActive)
			if key == utf8.RuneError ***REMOVED***
				break
			***REMOVED***
			if !t.pasteActive ***REMOVED***
				if key == keyCtrlD ***REMOVED***
					if len(t.line) == 0 ***REMOVED***
						return "", io.EOF
					***REMOVED***
				***REMOVED***
				if key == keyCtrlC ***REMOVED***
					return "", io.EOF
				***REMOVED***
				if key == keyPasteStart ***REMOVED***
					t.pasteActive = true
					if len(t.line) == 0 ***REMOVED***
						lineIsPasted = true
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else if key == keyPasteEnd ***REMOVED***
				t.pasteActive = false
				continue
			***REMOVED***
			if !t.pasteActive ***REMOVED***
				lineIsPasted = false
			***REMOVED***
			line, lineOk = t.handleKey(key)
		***REMOVED***
		if len(rest) > 0 ***REMOVED***
			n := copy(t.inBuf[:], rest)
			t.remainder = t.inBuf[:n]
		***REMOVED*** else ***REMOVED***
			t.remainder = nil
		***REMOVED***
		t.c.Write(t.outBuf)
		t.outBuf = t.outBuf[:0]
		if lineOk ***REMOVED***
			if t.echo ***REMOVED***
				t.historyIndex = -1
				t.history.Add(line)
			***REMOVED***
			if lineIsPasted ***REMOVED***
				err = ErrPasteIndicator
			***REMOVED***
			return
		***REMOVED***

		// t.remainder is a slice at the beginning of t.inBuf
		// containing a partial key sequence
		readBuf := t.inBuf[len(t.remainder):]
		var n int

		t.lock.Unlock()
		n, err = t.c.Read(readBuf)
		t.lock.Lock()

		if err != nil ***REMOVED***
			return
		***REMOVED***

		t.remainder = t.inBuf[:n+len(t.remainder)]
	***REMOVED***
***REMOVED***

// SetPrompt sets the prompt to be used when reading subsequent lines.
func (t *Terminal) SetPrompt(prompt string) ***REMOVED***
	t.lock.Lock()
	defer t.lock.Unlock()

	t.prompt = []rune(prompt)
***REMOVED***

func (t *Terminal) clearAndRepaintLinePlusNPrevious(numPrevLines int) ***REMOVED***
	// Move cursor to column zero at the start of the line.
	t.move(t.cursorY, 0, t.cursorX, 0)
	t.cursorX, t.cursorY = 0, 0
	t.clearLineToRight()
	for t.cursorY < numPrevLines ***REMOVED***
		// Move down a line
		t.move(0, 1, 0, 0)
		t.cursorY++
		t.clearLineToRight()
	***REMOVED***
	// Move back to beginning.
	t.move(t.cursorY, 0, 0, 0)
	t.cursorX, t.cursorY = 0, 0

	t.queue(t.prompt)
	t.advanceCursor(visualLength(t.prompt))
	t.writeLine(t.line)
	t.moveCursorToPos(t.pos)
***REMOVED***

func (t *Terminal) SetSize(width, height int) error ***REMOVED***
	t.lock.Lock()
	defer t.lock.Unlock()

	if width == 0 ***REMOVED***
		width = 1
	***REMOVED***

	oldWidth := t.termWidth
	t.termWidth, t.termHeight = width, height

	switch ***REMOVED***
	case width == oldWidth:
		// If the width didn't change then nothing else needs to be
		// done.
		return nil
	case len(t.line) == 0 && t.cursorX == 0 && t.cursorY == 0:
		// If there is nothing on current line and no prompt printed,
		// just do nothing
		return nil
	case width < oldWidth:
		// Some terminals (e.g. xterm) will truncate lines that were
		// too long when shinking. Others, (e.g. gnome-terminal) will
		// attempt to wrap them. For the former, repainting t.maxLine
		// works great, but that behaviour goes badly wrong in the case
		// of the latter because they have doubled every full line.

		// We assume that we are working on a terminal that wraps lines
		// and adjust the cursor position based on every previous line
		// wrapping and turning into two. This causes the prompt on
		// xterms to move upwards, which isn't great, but it avoids a
		// huge mess with gnome-terminal.
		if t.cursorX >= t.termWidth ***REMOVED***
			t.cursorX = t.termWidth - 1
		***REMOVED***
		t.cursorY *= 2
		t.clearAndRepaintLinePlusNPrevious(t.maxLine * 2)
	case width > oldWidth:
		// If the terminal expands then our position calculations will
		// be wrong in the future because we think the cursor is
		// |t.pos| chars into the string, but there will be a gap at
		// the end of any wrapped line.
		//
		// But the position will actually be correct until we move, so
		// we can move back to the beginning and repaint everything.
		t.clearAndRepaintLinePlusNPrevious(t.maxLine)
	***REMOVED***

	_, err := t.c.Write(t.outBuf)
	t.outBuf = t.outBuf[:0]
	return err
***REMOVED***

type pasteIndicatorError struct***REMOVED******REMOVED***

func (pasteIndicatorError) Error() string ***REMOVED***
	return "terminal: ErrPasteIndicator not correctly handled"
***REMOVED***

// ErrPasteIndicator may be returned from ReadLine as the error, in addition
// to valid line data. It indicates that bracketed paste mode is enabled and
// that the returned line consists only of pasted data. Programs may wish to
// interpret pasted data more literally than typed data.
var ErrPasteIndicator = pasteIndicatorError***REMOVED******REMOVED***

// SetBracketedPasteMode requests that the terminal bracket paste operations
// with markers. Not all terminals support this but, if it is supported, then
// enabling this mode will stop any autocomplete callback from running due to
// pastes. Additionally, any lines that are completely pasted will be returned
// from ReadLine with the error set to ErrPasteIndicator.
func (t *Terminal) SetBracketedPasteMode(on bool) ***REMOVED***
	if on ***REMOVED***
		io.WriteString(t.c, "\x1b[?2004h")
	***REMOVED*** else ***REMOVED***
		io.WriteString(t.c, "\x1b[?2004l")
	***REMOVED***
***REMOVED***

// stRingBuffer is a ring buffer of strings.
type stRingBuffer struct ***REMOVED***
	// entries contains max elements.
	entries []string
	max     int
	// head contains the index of the element most recently added to the ring.
	head int
	// size contains the number of elements in the ring.
	size int
***REMOVED***

func (s *stRingBuffer) Add(a string) ***REMOVED***
	if s.entries == nil ***REMOVED***
		const defaultNumEntries = 100
		s.entries = make([]string, defaultNumEntries)
		s.max = defaultNumEntries
	***REMOVED***

	s.head = (s.head + 1) % s.max
	s.entries[s.head] = a
	if s.size < s.max ***REMOVED***
		s.size++
	***REMOVED***
***REMOVED***

// NthPreviousEntry returns the value passed to the nth previous call to Add.
// If n is zero then the immediately prior value is returned, if one, then the
// next most recent, and so on. If such an element doesn't exist then ok is
// false.
func (s *stRingBuffer) NthPreviousEntry(n int) (value string, ok bool) ***REMOVED***
	if n < 0 || n >= s.size ***REMOVED***
		return "", false
	***REMOVED***
	index := s.head - n
	if index < 0 ***REMOVED***
		index += s.max
	***REMOVED***
	return s.entries[index], true
***REMOVED***

// readPasswordLine reads from reader until it finds \n or io.EOF.
// The slice returned does not include the \n.
// readPasswordLine also ignores any \r it finds.
// Windows uses \r as end of line. So, on Windows, readPasswordLine
// reads until it finds \r and ignores any \n it finds during processing.
func readPasswordLine(reader io.Reader) ([]byte, error) ***REMOVED***
	var buf [1]byte
	var ret []byte

	for ***REMOVED***
		n, err := reader.Read(buf[:])
		if n > 0 ***REMOVED***
			switch buf[0] ***REMOVED***
			case '\b':
				if len(ret) > 0 ***REMOVED***
					ret = ret[:len(ret)-1]
				***REMOVED***
			case '\n':
				if runtime.GOOS != "windows" ***REMOVED***
					return ret, nil
				***REMOVED***
				// otherwise ignore \n
			case '\r':
				if runtime.GOOS == "windows" ***REMOVED***
					return ret, nil
				***REMOVED***
				// otherwise ignore \r
			default:
				ret = append(ret, buf[0])
			***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			if err == io.EOF && len(ret) > 0 ***REMOVED***
				return ret, nil
			***REMOVED***
			return ret, err
		***REMOVED***
	***REMOVED***
***REMOVED***
