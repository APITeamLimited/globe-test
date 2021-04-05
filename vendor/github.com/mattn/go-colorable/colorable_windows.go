// +build windows
// +build !appengine

package colorable

import (
	"bytes"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/mattn/go-isatty"
)

const (
	foregroundBlue      = 0x1
	foregroundGreen     = 0x2
	foregroundRed       = 0x4
	foregroundIntensity = 0x8
	foregroundMask      = (foregroundRed | foregroundBlue | foregroundGreen | foregroundIntensity)
	backgroundBlue      = 0x10
	backgroundGreen     = 0x20
	backgroundRed       = 0x40
	backgroundIntensity = 0x80
	backgroundMask      = (backgroundRed | backgroundBlue | backgroundGreen | backgroundIntensity)
	commonLvbUnderscore = 0x8000

	cENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x4
)

const (
	genericRead  = 0x80000000
	genericWrite = 0x40000000
)

const (
	consoleTextmodeBuffer = 0x1
)

type wchar uint16
type short int16
type dword uint32
type word uint16

type coord struct ***REMOVED***
	x short
	y short
***REMOVED***

type smallRect struct ***REMOVED***
	left   short
	top    short
	right  short
	bottom short
***REMOVED***

type consoleScreenBufferInfo struct ***REMOVED***
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
***REMOVED***

type consoleCursorInfo struct ***REMOVED***
	size    dword
	visible int32
***REMOVED***

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procSetConsoleTextAttribute    = kernel32.NewProc("SetConsoleTextAttribute")
	procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")
	procGetConsoleCursorInfo       = kernel32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo       = kernel32.NewProc("SetConsoleCursorInfo")
	procSetConsoleTitle            = kernel32.NewProc("SetConsoleTitleW")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procCreateConsoleScreenBuffer  = kernel32.NewProc("CreateConsoleScreenBuffer")
)

// Writer provides colorable Writer to the console
type Writer struct ***REMOVED***
	out       io.Writer
	handle    syscall.Handle
	althandle syscall.Handle
	oldattr   word
	oldpos    coord
	rest      bytes.Buffer
	mutex     sync.Mutex
***REMOVED***

// NewColorable returns new instance of Writer which handles escape sequence from File.
func NewColorable(file *os.File) io.Writer ***REMOVED***
	if file == nil ***REMOVED***
		panic("nil passed instead of *os.File to NewColorable()")
	***REMOVED***

	if isatty.IsTerminal(file.Fd()) ***REMOVED***
		var mode uint32
		if r, _, _ := procGetConsoleMode.Call(file.Fd(), uintptr(unsafe.Pointer(&mode))); r != 0 && mode&cENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 ***REMOVED***
			return file
		***REMOVED***
		var csbi consoleScreenBufferInfo
		handle := syscall.Handle(file.Fd())
		procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
		return &Writer***REMOVED***out: file, handle: handle, oldattr: csbi.attributes, oldpos: coord***REMOVED***0, 0***REMOVED******REMOVED***
	***REMOVED***
	return file
***REMOVED***

// NewColorableStdout returns new instance of Writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer ***REMOVED***
	return NewColorable(os.Stdout)
***REMOVED***

// NewColorableStderr returns new instance of Writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer ***REMOVED***
	return NewColorable(os.Stderr)
***REMOVED***

var color256 = map[int]int***REMOVED***
	0:   0x000000,
	1:   0x800000,
	2:   0x008000,
	3:   0x808000,
	4:   0x000080,
	5:   0x800080,
	6:   0x008080,
	7:   0xc0c0c0,
	8:   0x808080,
	9:   0xff0000,
	10:  0x00ff00,
	11:  0xffff00,
	12:  0x0000ff,
	13:  0xff00ff,
	14:  0x00ffff,
	15:  0xffffff,
	16:  0x000000,
	17:  0x00005f,
	18:  0x000087,
	19:  0x0000af,
	20:  0x0000d7,
	21:  0x0000ff,
	22:  0x005f00,
	23:  0x005f5f,
	24:  0x005f87,
	25:  0x005faf,
	26:  0x005fd7,
	27:  0x005fff,
	28:  0x008700,
	29:  0x00875f,
	30:  0x008787,
	31:  0x0087af,
	32:  0x0087d7,
	33:  0x0087ff,
	34:  0x00af00,
	35:  0x00af5f,
	36:  0x00af87,
	37:  0x00afaf,
	38:  0x00afd7,
	39:  0x00afff,
	40:  0x00d700,
	41:  0x00d75f,
	42:  0x00d787,
	43:  0x00d7af,
	44:  0x00d7d7,
	45:  0x00d7ff,
	46:  0x00ff00,
	47:  0x00ff5f,
	48:  0x00ff87,
	49:  0x00ffaf,
	50:  0x00ffd7,
	51:  0x00ffff,
	52:  0x5f0000,
	53:  0x5f005f,
	54:  0x5f0087,
	55:  0x5f00af,
	56:  0x5f00d7,
	57:  0x5f00ff,
	58:  0x5f5f00,
	59:  0x5f5f5f,
	60:  0x5f5f87,
	61:  0x5f5faf,
	62:  0x5f5fd7,
	63:  0x5f5fff,
	64:  0x5f8700,
	65:  0x5f875f,
	66:  0x5f8787,
	67:  0x5f87af,
	68:  0x5f87d7,
	69:  0x5f87ff,
	70:  0x5faf00,
	71:  0x5faf5f,
	72:  0x5faf87,
	73:  0x5fafaf,
	74:  0x5fafd7,
	75:  0x5fafff,
	76:  0x5fd700,
	77:  0x5fd75f,
	78:  0x5fd787,
	79:  0x5fd7af,
	80:  0x5fd7d7,
	81:  0x5fd7ff,
	82:  0x5fff00,
	83:  0x5fff5f,
	84:  0x5fff87,
	85:  0x5fffaf,
	86:  0x5fffd7,
	87:  0x5fffff,
	88:  0x870000,
	89:  0x87005f,
	90:  0x870087,
	91:  0x8700af,
	92:  0x8700d7,
	93:  0x8700ff,
	94:  0x875f00,
	95:  0x875f5f,
	96:  0x875f87,
	97:  0x875faf,
	98:  0x875fd7,
	99:  0x875fff,
	100: 0x878700,
	101: 0x87875f,
	102: 0x878787,
	103: 0x8787af,
	104: 0x8787d7,
	105: 0x8787ff,
	106: 0x87af00,
	107: 0x87af5f,
	108: 0x87af87,
	109: 0x87afaf,
	110: 0x87afd7,
	111: 0x87afff,
	112: 0x87d700,
	113: 0x87d75f,
	114: 0x87d787,
	115: 0x87d7af,
	116: 0x87d7d7,
	117: 0x87d7ff,
	118: 0x87ff00,
	119: 0x87ff5f,
	120: 0x87ff87,
	121: 0x87ffaf,
	122: 0x87ffd7,
	123: 0x87ffff,
	124: 0xaf0000,
	125: 0xaf005f,
	126: 0xaf0087,
	127: 0xaf00af,
	128: 0xaf00d7,
	129: 0xaf00ff,
	130: 0xaf5f00,
	131: 0xaf5f5f,
	132: 0xaf5f87,
	133: 0xaf5faf,
	134: 0xaf5fd7,
	135: 0xaf5fff,
	136: 0xaf8700,
	137: 0xaf875f,
	138: 0xaf8787,
	139: 0xaf87af,
	140: 0xaf87d7,
	141: 0xaf87ff,
	142: 0xafaf00,
	143: 0xafaf5f,
	144: 0xafaf87,
	145: 0xafafaf,
	146: 0xafafd7,
	147: 0xafafff,
	148: 0xafd700,
	149: 0xafd75f,
	150: 0xafd787,
	151: 0xafd7af,
	152: 0xafd7d7,
	153: 0xafd7ff,
	154: 0xafff00,
	155: 0xafff5f,
	156: 0xafff87,
	157: 0xafffaf,
	158: 0xafffd7,
	159: 0xafffff,
	160: 0xd70000,
	161: 0xd7005f,
	162: 0xd70087,
	163: 0xd700af,
	164: 0xd700d7,
	165: 0xd700ff,
	166: 0xd75f00,
	167: 0xd75f5f,
	168: 0xd75f87,
	169: 0xd75faf,
	170: 0xd75fd7,
	171: 0xd75fff,
	172: 0xd78700,
	173: 0xd7875f,
	174: 0xd78787,
	175: 0xd787af,
	176: 0xd787d7,
	177: 0xd787ff,
	178: 0xd7af00,
	179: 0xd7af5f,
	180: 0xd7af87,
	181: 0xd7afaf,
	182: 0xd7afd7,
	183: 0xd7afff,
	184: 0xd7d700,
	185: 0xd7d75f,
	186: 0xd7d787,
	187: 0xd7d7af,
	188: 0xd7d7d7,
	189: 0xd7d7ff,
	190: 0xd7ff00,
	191: 0xd7ff5f,
	192: 0xd7ff87,
	193: 0xd7ffaf,
	194: 0xd7ffd7,
	195: 0xd7ffff,
	196: 0xff0000,
	197: 0xff005f,
	198: 0xff0087,
	199: 0xff00af,
	200: 0xff00d7,
	201: 0xff00ff,
	202: 0xff5f00,
	203: 0xff5f5f,
	204: 0xff5f87,
	205: 0xff5faf,
	206: 0xff5fd7,
	207: 0xff5fff,
	208: 0xff8700,
	209: 0xff875f,
	210: 0xff8787,
	211: 0xff87af,
	212: 0xff87d7,
	213: 0xff87ff,
	214: 0xffaf00,
	215: 0xffaf5f,
	216: 0xffaf87,
	217: 0xffafaf,
	218: 0xffafd7,
	219: 0xffafff,
	220: 0xffd700,
	221: 0xffd75f,
	222: 0xffd787,
	223: 0xffd7af,
	224: 0xffd7d7,
	225: 0xffd7ff,
	226: 0xffff00,
	227: 0xffff5f,
	228: 0xffff87,
	229: 0xffffaf,
	230: 0xffffd7,
	231: 0xffffff,
	232: 0x080808,
	233: 0x121212,
	234: 0x1c1c1c,
	235: 0x262626,
	236: 0x303030,
	237: 0x3a3a3a,
	238: 0x444444,
	239: 0x4e4e4e,
	240: 0x585858,
	241: 0x626262,
	242: 0x6c6c6c,
	243: 0x767676,
	244: 0x808080,
	245: 0x8a8a8a,
	246: 0x949494,
	247: 0x9e9e9e,
	248: 0xa8a8a8,
	249: 0xb2b2b2,
	250: 0xbcbcbc,
	251: 0xc6c6c6,
	252: 0xd0d0d0,
	253: 0xdadada,
	254: 0xe4e4e4,
	255: 0xeeeeee,
***REMOVED***

// `\033]0;TITLESTR\007`
func doTitleSequence(er *bytes.Reader) error ***REMOVED***
	var c byte
	var err error

	c, err = er.ReadByte()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if c != '0' && c != '2' ***REMOVED***
		return nil
	***REMOVED***
	c, err = er.ReadByte()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if c != ';' ***REMOVED***
		return nil
	***REMOVED***
	title := make([]byte, 0, 80)
	for ***REMOVED***
		c, err = er.ReadByte()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if c == 0x07 || c == '\n' ***REMOVED***
			break
		***REMOVED***
		title = append(title, c)
	***REMOVED***
	if len(title) > 0 ***REMOVED***
		title8, err := syscall.UTF16PtrFromString(string(title))
		if err == nil ***REMOVED***
			procSetConsoleTitle.Call(uintptr(unsafe.Pointer(title8)))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// returns Atoi(s) unless s == "" in which case it returns def
func atoiWithDefault(s string, def int) (int, error) ***REMOVED***
	if s == "" ***REMOVED***
		return def, nil
	***REMOVED***
	return strconv.Atoi(s)
***REMOVED***

// Write writes data on console
func (w *Writer) Write(data []byte) (n int, err error) ***REMOVED***
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(w.handle), uintptr(unsafe.Pointer(&csbi)))

	handle := w.handle

	var er *bytes.Reader
	if w.rest.Len() > 0 ***REMOVED***
		var rest bytes.Buffer
		w.rest.WriteTo(&rest)
		w.rest.Reset()
		rest.Write(data)
		er = bytes.NewReader(rest.Bytes())
	***REMOVED*** else ***REMOVED***
		er = bytes.NewReader(data)
	***REMOVED***
	var bw [1]byte
loop:
	for ***REMOVED***
		c1, err := er.ReadByte()
		if err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c1 != 0x1b ***REMOVED***
			bw[0] = c1
			w.out.Write(bw[:])
			continue
		***REMOVED***
		c2, err := er.ReadByte()
		if err != nil ***REMOVED***
			break loop
		***REMOVED***

		switch c2 ***REMOVED***
		case '>':
			continue
		case ']':
			w.rest.WriteByte(c1)
			w.rest.WriteByte(c2)
			er.WriteTo(&w.rest)
			if bytes.IndexByte(w.rest.Bytes(), 0x07) == -1 ***REMOVED***
				break loop
			***REMOVED***
			er = bytes.NewReader(w.rest.Bytes()[2:])
			err := doTitleSequence(er)
			if err != nil ***REMOVED***
				break loop
			***REMOVED***
			w.rest.Reset()
			continue
		// https://github.com/mattn/go-colorable/issues/27
		case '7':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			w.oldpos = csbi.cursorPosition
			continue
		case '8':
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&w.oldpos)))
			continue
		case 0x5b:
			// execute part after switch
		default:
			continue
		***REMOVED***

		w.rest.WriteByte(c1)
		w.rest.WriteByte(c2)
		er.WriteTo(&w.rest)

		var buf bytes.Buffer
		var m byte
		for i, c := range w.rest.Bytes()[2:] ***REMOVED***
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' ***REMOVED***
				m = c
				er = bytes.NewReader(w.rest.Bytes()[2+i+1:])
				w.rest.Reset()
				break
			***REMOVED***
			buf.Write([]byte(string(c)))
		***REMOVED***
		if m == 0 ***REMOVED***
			break loop
		***REMOVED***

		switch m ***REMOVED***
		case 'A':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.y -= short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'B':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.y += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'C':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'D':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x -= short(n)
			if csbi.cursorPosition.x < 0 ***REMOVED***
				csbi.cursorPosition.x = 0
			***REMOVED***
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'E':
			n, err = strconv.Atoi(buf.String())
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = 0
			csbi.cursorPosition.y += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'F':
			n, err = strconv.Atoi(buf.String())
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = 0
			csbi.cursorPosition.y -= short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'G':
			n, err = strconv.Atoi(buf.String())
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			if n < 1 ***REMOVED***
				n = 1
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = short(n - 1)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'H', 'f':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			if buf.Len() > 0 ***REMOVED***
				token := strings.Split(buf.String(), ";")
				switch len(token) ***REMOVED***
				case 1:
					n1, err := strconv.Atoi(token[0])
					if err != nil ***REMOVED***
						continue
					***REMOVED***
					csbi.cursorPosition.y = short(n1 - 1)
				case 2:
					n1, err := strconv.Atoi(token[0])
					if err != nil ***REMOVED***
						continue
					***REMOVED***
					n2, err := strconv.Atoi(token[1])
					if err != nil ***REMOVED***
						continue
					***REMOVED***
					csbi.cursorPosition.x = short(n2 - 1)
					csbi.cursorPosition.y = short(n1 - 1)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				csbi.cursorPosition.y = 0
			***REMOVED***
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'J':
			n := 0
			if buf.Len() > 0 ***REMOVED***
				n, err = strconv.Atoi(buf.String())
				if err != nil ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			var count, written dword
			var cursor coord
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			switch n ***REMOVED***
			case 0:
				cursor = coord***REMOVED***x: csbi.cursorPosition.x, y: csbi.cursorPosition.y***REMOVED***
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.size.y-csbi.cursorPosition.y)*dword(csbi.size.x)
			case 1:
				cursor = coord***REMOVED***x: csbi.window.left, y: csbi.window.top***REMOVED***
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.window.top-csbi.cursorPosition.y)*dword(csbi.size.x)
			case 2:
				cursor = coord***REMOVED***x: csbi.window.left, y: csbi.window.top***REMOVED***
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.size.y-csbi.cursorPosition.y)*dword(csbi.size.x)
			***REMOVED***
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
		case 'K':
			n := 0
			if buf.Len() > 0 ***REMOVED***
				n, err = strconv.Atoi(buf.String())
				if err != nil ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			var cursor coord
			var count, written dword
			switch n ***REMOVED***
			case 0:
				cursor = coord***REMOVED***x: csbi.cursorPosition.x, y: csbi.cursorPosition.y***REMOVED***
				count = dword(csbi.size.x - csbi.cursorPosition.x)
			case 1:
				cursor = coord***REMOVED***x: csbi.window.left, y: csbi.cursorPosition.y***REMOVED***
				count = dword(csbi.size.x - csbi.cursorPosition.x)
			case 2:
				cursor = coord***REMOVED***x: csbi.window.left, y: csbi.cursorPosition.y***REMOVED***
				count = dword(csbi.size.x)
			***REMOVED***
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
		case 'X':
			n := 0
			if buf.Len() > 0 ***REMOVED***
				n, err = strconv.Atoi(buf.String())
				if err != nil ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			var cursor coord
			var written dword
			cursor = coord***REMOVED***x: csbi.cursorPosition.x, y: csbi.cursorPosition.y***REMOVED***
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(n), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(n), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
		case 'm':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			attr := csbi.attributes
			cs := buf.String()
			if cs == "" ***REMOVED***
				procSetConsoleTextAttribute.Call(uintptr(handle), uintptr(w.oldattr))
				continue
			***REMOVED***
			token := strings.Split(cs, ";")
			for i := 0; i < len(token); i++ ***REMOVED***
				ns := token[i]
				if n, err = strconv.Atoi(ns); err == nil ***REMOVED***
					switch ***REMOVED***
					case n == 0 || n == 100:
						attr = w.oldattr
					case n == 4:
						attr |= commonLvbUnderscore
					case (1 <= n && n <= 3) || n == 5:
						attr |= foregroundIntensity
					case n == 7 || n == 27:
						attr =
							(attr &^ (foregroundMask | backgroundMask)) |
								((attr & foregroundMask) << 4) |
								((attr & backgroundMask) >> 4)
					case n == 22:
						attr &^= foregroundIntensity
					case n == 24:
						attr &^= commonLvbUnderscore
					case 30 <= n && n <= 37:
						attr &= backgroundMask
						if (n-30)&1 != 0 ***REMOVED***
							attr |= foregroundRed
						***REMOVED***
						if (n-30)&2 != 0 ***REMOVED***
							attr |= foregroundGreen
						***REMOVED***
						if (n-30)&4 != 0 ***REMOVED***
							attr |= foregroundBlue
						***REMOVED***
					case n == 38: // set foreground color.
						if i < len(token)-2 && (token[i+1] == "5" || token[i+1] == "05") ***REMOVED***
							if n256, err := strconv.Atoi(token[i+2]); err == nil ***REMOVED***
								if n256foreAttr == nil ***REMOVED***
									n256setup()
								***REMOVED***
								attr &= backgroundMask
								attr |= n256foreAttr[n256%len(n256foreAttr)]
								i += 2
							***REMOVED***
						***REMOVED*** else if len(token) == 5 && token[i+1] == "2" ***REMOVED***
							var r, g, b int
							r, _ = strconv.Atoi(token[i+2])
							g, _ = strconv.Atoi(token[i+3])
							b, _ = strconv.Atoi(token[i+4])
							i += 4
							if r > 127 ***REMOVED***
								attr |= foregroundRed
							***REMOVED***
							if g > 127 ***REMOVED***
								attr |= foregroundGreen
							***REMOVED***
							if b > 127 ***REMOVED***
								attr |= foregroundBlue
							***REMOVED***
						***REMOVED*** else ***REMOVED***
							attr = attr & (w.oldattr & backgroundMask)
						***REMOVED***
					case n == 39: // reset foreground color.
						attr &= backgroundMask
						attr |= w.oldattr & foregroundMask
					case 40 <= n && n <= 47:
						attr &= foregroundMask
						if (n-40)&1 != 0 ***REMOVED***
							attr |= backgroundRed
						***REMOVED***
						if (n-40)&2 != 0 ***REMOVED***
							attr |= backgroundGreen
						***REMOVED***
						if (n-40)&4 != 0 ***REMOVED***
							attr |= backgroundBlue
						***REMOVED***
					case n == 48: // set background color.
						if i < len(token)-2 && token[i+1] == "5" ***REMOVED***
							if n256, err := strconv.Atoi(token[i+2]); err == nil ***REMOVED***
								if n256backAttr == nil ***REMOVED***
									n256setup()
								***REMOVED***
								attr &= foregroundMask
								attr |= n256backAttr[n256%len(n256backAttr)]
								i += 2
							***REMOVED***
						***REMOVED*** else if len(token) == 5 && token[i+1] == "2" ***REMOVED***
							var r, g, b int
							r, _ = strconv.Atoi(token[i+2])
							g, _ = strconv.Atoi(token[i+3])
							b, _ = strconv.Atoi(token[i+4])
							i += 4
							if r > 127 ***REMOVED***
								attr |= backgroundRed
							***REMOVED***
							if g > 127 ***REMOVED***
								attr |= backgroundGreen
							***REMOVED***
							if b > 127 ***REMOVED***
								attr |= backgroundBlue
							***REMOVED***
						***REMOVED*** else ***REMOVED***
							attr = attr & (w.oldattr & foregroundMask)
						***REMOVED***
					case n == 49: // reset foreground color.
						attr &= foregroundMask
						attr |= w.oldattr & backgroundMask
					case 90 <= n && n <= 97:
						attr = (attr & backgroundMask)
						attr |= foregroundIntensity
						if (n-90)&1 != 0 ***REMOVED***
							attr |= foregroundRed
						***REMOVED***
						if (n-90)&2 != 0 ***REMOVED***
							attr |= foregroundGreen
						***REMOVED***
						if (n-90)&4 != 0 ***REMOVED***
							attr |= foregroundBlue
						***REMOVED***
					case 100 <= n && n <= 107:
						attr = (attr & foregroundMask)
						attr |= backgroundIntensity
						if (n-100)&1 != 0 ***REMOVED***
							attr |= backgroundRed
						***REMOVED***
						if (n-100)&2 != 0 ***REMOVED***
							attr |= backgroundGreen
						***REMOVED***
						if (n-100)&4 != 0 ***REMOVED***
							attr |= backgroundBlue
						***REMOVED***
					***REMOVED***
					procSetConsoleTextAttribute.Call(uintptr(handle), uintptr(attr))
				***REMOVED***
			***REMOVED***
		case 'h':
			var ci consoleCursorInfo
			cs := buf.String()
			if cs == "5>" ***REMOVED***
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 0
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			***REMOVED*** else if cs == "?25" ***REMOVED***
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 1
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			***REMOVED*** else if cs == "?1049" ***REMOVED***
				if w.althandle == 0 ***REMOVED***
					h, _, _ := procCreateConsoleScreenBuffer.Call(uintptr(genericRead|genericWrite), 0, 0, uintptr(consoleTextmodeBuffer), 0, 0)
					w.althandle = syscall.Handle(h)
					if w.althandle != 0 ***REMOVED***
						handle = w.althandle
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case 'l':
			var ci consoleCursorInfo
			cs := buf.String()
			if cs == "5>" ***REMOVED***
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 1
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			***REMOVED*** else if cs == "?25" ***REMOVED***
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 0
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			***REMOVED*** else if cs == "?1049" ***REMOVED***
				if w.althandle != 0 ***REMOVED***
					syscall.CloseHandle(w.althandle)
					w.althandle = 0
					handle = w.handle
				***REMOVED***
			***REMOVED***
		case 's':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			w.oldpos = csbi.cursorPosition
		case 'u':
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&w.oldpos)))
		***REMOVED***
	***REMOVED***

	return len(data), nil
***REMOVED***

type consoleColor struct ***REMOVED***
	rgb       int
	red       bool
	green     bool
	blue      bool
	intensity bool
***REMOVED***

func (c consoleColor) foregroundAttr() (attr word) ***REMOVED***
	if c.red ***REMOVED***
		attr |= foregroundRed
	***REMOVED***
	if c.green ***REMOVED***
		attr |= foregroundGreen
	***REMOVED***
	if c.blue ***REMOVED***
		attr |= foregroundBlue
	***REMOVED***
	if c.intensity ***REMOVED***
		attr |= foregroundIntensity
	***REMOVED***
	return
***REMOVED***

func (c consoleColor) backgroundAttr() (attr word) ***REMOVED***
	if c.red ***REMOVED***
		attr |= backgroundRed
	***REMOVED***
	if c.green ***REMOVED***
		attr |= backgroundGreen
	***REMOVED***
	if c.blue ***REMOVED***
		attr |= backgroundBlue
	***REMOVED***
	if c.intensity ***REMOVED***
		attr |= backgroundIntensity
	***REMOVED***
	return
***REMOVED***

var color16 = []consoleColor***REMOVED***
	***REMOVED***0x000000, false, false, false, false***REMOVED***,
	***REMOVED***0x000080, false, false, true, false***REMOVED***,
	***REMOVED***0x008000, false, true, false, false***REMOVED***,
	***REMOVED***0x008080, false, true, true, false***REMOVED***,
	***REMOVED***0x800000, true, false, false, false***REMOVED***,
	***REMOVED***0x800080, true, false, true, false***REMOVED***,
	***REMOVED***0x808000, true, true, false, false***REMOVED***,
	***REMOVED***0xc0c0c0, true, true, true, false***REMOVED***,
	***REMOVED***0x808080, false, false, false, true***REMOVED***,
	***REMOVED***0x0000ff, false, false, true, true***REMOVED***,
	***REMOVED***0x00ff00, false, true, false, true***REMOVED***,
	***REMOVED***0x00ffff, false, true, true, true***REMOVED***,
	***REMOVED***0xff0000, true, false, false, true***REMOVED***,
	***REMOVED***0xff00ff, true, false, true, true***REMOVED***,
	***REMOVED***0xffff00, true, true, false, true***REMOVED***,
	***REMOVED***0xffffff, true, true, true, true***REMOVED***,
***REMOVED***

type hsv struct ***REMOVED***
	h, s, v float32
***REMOVED***

func (a hsv) dist(b hsv) float32 ***REMOVED***
	dh := a.h - b.h
	switch ***REMOVED***
	case dh > 0.5:
		dh = 1 - dh
	case dh < -0.5:
		dh = -1 - dh
	***REMOVED***
	ds := a.s - b.s
	dv := a.v - b.v
	return float32(math.Sqrt(float64(dh*dh + ds*ds + dv*dv)))
***REMOVED***

func toHSV(rgb int) hsv ***REMOVED***
	r, g, b := float32((rgb&0xFF0000)>>16)/256.0,
		float32((rgb&0x00FF00)>>8)/256.0,
		float32(rgb&0x0000FF)/256.0
	min, max := minmax3f(r, g, b)
	h := max - min
	if h > 0 ***REMOVED***
		if max == r ***REMOVED***
			h = (g - b) / h
			if h < 0 ***REMOVED***
				h += 6
			***REMOVED***
		***REMOVED*** else if max == g ***REMOVED***
			h = 2 + (b-r)/h
		***REMOVED*** else ***REMOVED***
			h = 4 + (r-g)/h
		***REMOVED***
	***REMOVED***
	h /= 6.0
	s := max - min
	if max != 0 ***REMOVED***
		s /= max
	***REMOVED***
	v := max
	return hsv***REMOVED***h: h, s: s, v: v***REMOVED***
***REMOVED***

type hsvTable []hsv

func toHSVTable(rgbTable []consoleColor) hsvTable ***REMOVED***
	t := make(hsvTable, len(rgbTable))
	for i, c := range rgbTable ***REMOVED***
		t[i] = toHSV(c.rgb)
	***REMOVED***
	return t
***REMOVED***

func (t hsvTable) find(rgb int) consoleColor ***REMOVED***
	hsv := toHSV(rgb)
	n := 7
	l := float32(5.0)
	for i, p := range t ***REMOVED***
		d := hsv.dist(p)
		if d < l ***REMOVED***
			l, n = d, i
		***REMOVED***
	***REMOVED***
	return color16[n]
***REMOVED***

func minmax3f(a, b, c float32) (min, max float32) ***REMOVED***
	if a < b ***REMOVED***
		if b < c ***REMOVED***
			return a, c
		***REMOVED*** else if a < c ***REMOVED***
			return a, b
		***REMOVED*** else ***REMOVED***
			return c, b
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if a < c ***REMOVED***
			return b, c
		***REMOVED*** else if b < c ***REMOVED***
			return b, a
		***REMOVED*** else ***REMOVED***
			return c, a
		***REMOVED***
	***REMOVED***
***REMOVED***

var n256foreAttr []word
var n256backAttr []word

func n256setup() ***REMOVED***
	n256foreAttr = make([]word, 256)
	n256backAttr = make([]word, 256)
	t := toHSVTable(color16)
	for i, rgb := range color256 ***REMOVED***
		c := t.find(rgb)
		n256foreAttr[i] = c.foregroundAttr()
		n256backAttr[i] = c.backgroundAttr()
	***REMOVED***
***REMOVED***

// EnableColorsStdout enable colors if possible.
func EnableColorsStdout(enabled *bool) func() ***REMOVED***
	var mode uint32
	h := os.Stdout.Fd()
	if r, _, _ := procGetConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode))); r != 0 ***REMOVED***
		if r, _, _ = procSetConsoleMode.Call(h, uintptr(mode|cENABLE_VIRTUAL_TERMINAL_PROCESSING)); r != 0 ***REMOVED***
			if enabled != nil ***REMOVED***
				*enabled = true
			***REMOVED***
			return func() ***REMOVED***
				procSetConsoleMode.Call(h, uintptr(mode))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if enabled != nil ***REMOVED***
		*enabled = true
	***REMOVED***
	return func() ***REMOVED******REMOVED***
***REMOVED***
