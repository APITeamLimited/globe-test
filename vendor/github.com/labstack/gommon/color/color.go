package color

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

type (
	inner func(interface***REMOVED******REMOVED***, []string, *Color) string
)

// Color styles
const (
	// Blk Black text style
	Blk = "30"
	// Rd red text style
	Rd = "31"
	// Grn green text style
	Grn = "32"
	// Yel yellow text style
	Yel = "33"
	// Blu blue text style
	Blu = "34"
	// Mgn magenta text style
	Mgn = "35"
	// Cyn cyan text style
	Cyn = "36"
	// Wht white text style
	Wht = "37"
	// Gry grey text style
	Gry = "90"

	// BlkBg black background style
	BlkBg = "40"
	// RdBg red background style
	RdBg = "41"
	// GrnBg green background style
	GrnBg = "42"
	// YelBg yellow background style
	YelBg = "43"
	// BluBg blue background style
	BluBg = "44"
	// MgnBg magenta background style
	MgnBg = "45"
	// CynBg cyan background style
	CynBg = "46"
	// WhtBg white background style
	WhtBg = "47"

	// R reset emphasis style
	R = "0"
	// B bold emphasis style
	B = "1"
	// D dim emphasis style
	D = "2"
	// I italic emphasis style
	I = "3"
	// U underline emphasis style
	U = "4"
	// In inverse emphasis style
	In = "7"
	// H hidden emphasis style
	H = "8"
	// S strikeout emphasis style
	S = "9"
)

var (
	black   = outer(Blk)
	red     = outer(Rd)
	green   = outer(Grn)
	yellow  = outer(Yel)
	blue    = outer(Blu)
	magenta = outer(Mgn)
	cyan    = outer(Cyn)
	white   = outer(Wht)
	grey    = outer(Gry)

	blackBg   = outer(BlkBg)
	redBg     = outer(RdBg)
	greenBg   = outer(GrnBg)
	yellowBg  = outer(YelBg)
	blueBg    = outer(BluBg)
	magentaBg = outer(MgnBg)
	cyanBg    = outer(CynBg)
	whiteBg   = outer(WhtBg)

	reset     = outer(R)
	bold      = outer(B)
	dim       = outer(D)
	italic    = outer(I)
	underline = outer(U)
	inverse   = outer(In)
	hidden    = outer(H)
	strikeout = outer(S)

	global = New()
)

func outer(n string) inner ***REMOVED***
	return func(msg interface***REMOVED******REMOVED***, styles []string, c *Color) string ***REMOVED***
		// TODO: Drop fmt to boost performance?
		if c.disabled ***REMOVED***
			return fmt.Sprintf("%v", msg)
		***REMOVED***

		b := new(bytes.Buffer)
		b.WriteString("\x1b[")
		b.WriteString(n)
		for _, s := range styles ***REMOVED***
			b.WriteString(";")
			b.WriteString(s)
		***REMOVED***
		b.WriteString("m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), msg)
	***REMOVED***
***REMOVED***

type (
	Color struct ***REMOVED***
		output   io.Writer
		disabled bool
	***REMOVED***
)

// New creates a Color instance.
func New() (c *Color) ***REMOVED***
	c = new(Color)
	c.SetOutput(colorable.NewColorableStdout())
	return
***REMOVED***

// Output returns the output.
func (c *Color) Output() io.Writer ***REMOVED***
	return c.output
***REMOVED***

// SetOutput sets the output.
func (c *Color) SetOutput(w io.Writer) ***REMOVED***
	c.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) ***REMOVED***
		c.disabled = true
	***REMOVED***
***REMOVED***

// Disable disables the colors and styles.
func (c *Color) Disable() ***REMOVED***
	c.disabled = true
***REMOVED***

// Enable enables the colors and styles.
func (c *Color) Enable() ***REMOVED***
	c.disabled = false
***REMOVED***

// Print is analogous to `fmt.Print` with termial detection.
func (c *Color) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprint(c.output, args...)
***REMOVED***

// Println is analogous to `fmt.Println` with termial detection.
func (c *Color) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprintln(c.output, args...)
***REMOVED***

// Printf is analogous to `fmt.Printf` with termial detection.
func (c *Color) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprintf(c.output, format, args...)
***REMOVED***

func (c *Color) Black(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return black(msg, styles, c)
***REMOVED***

func (c *Color) Red(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return red(msg, styles, c)
***REMOVED***

func (c *Color) Green(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return green(msg, styles, c)
***REMOVED***

func (c *Color) Yellow(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return yellow(msg, styles, c)
***REMOVED***

func (c *Color) Blue(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return blue(msg, styles, c)
***REMOVED***

func (c *Color) Magenta(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return magenta(msg, styles, c)
***REMOVED***

func (c *Color) Cyan(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return cyan(msg, styles, c)
***REMOVED***

func (c *Color) White(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return white(msg, styles, c)
***REMOVED***

func (c *Color) Grey(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return grey(msg, styles, c)
***REMOVED***

func (c *Color) BlackBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return blackBg(msg, styles, c)
***REMOVED***

func (c *Color) RedBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return redBg(msg, styles, c)
***REMOVED***

func (c *Color) GreenBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return greenBg(msg, styles, c)
***REMOVED***

func (c *Color) YellowBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return yellowBg(msg, styles, c)
***REMOVED***

func (c *Color) BlueBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return blueBg(msg, styles, c)
***REMOVED***

func (c *Color) MagentaBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return magentaBg(msg, styles, c)
***REMOVED***

func (c *Color) CyanBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return cyanBg(msg, styles, c)
***REMOVED***

func (c *Color) WhiteBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return whiteBg(msg, styles, c)
***REMOVED***

func (c *Color) Reset(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return reset(msg, styles, c)
***REMOVED***

func (c *Color) Bold(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return bold(msg, styles, c)
***REMOVED***

func (c *Color) Dim(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return dim(msg, styles, c)
***REMOVED***

func (c *Color) Italic(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return italic(msg, styles, c)
***REMOVED***

func (c *Color) Underline(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return underline(msg, styles, c)
***REMOVED***

func (c *Color) Inverse(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return inverse(msg, styles, c)
***REMOVED***

func (c *Color) Hidden(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return hidden(msg, styles, c)
***REMOVED***

func (c *Color) Strikeout(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return strikeout(msg, styles, c)
***REMOVED***

// Output returns the output.
func Output() io.Writer ***REMOVED***
	return global.output
***REMOVED***

// SetOutput sets the output.
func SetOutput(w io.Writer) ***REMOVED***
	global.SetOutput(w)
***REMOVED***

func Disable() ***REMOVED***
	global.Disable()
***REMOVED***

func Enable() ***REMOVED***
	global.Enable()
***REMOVED***

// Print is analogous to `fmt.Print` with termial detection.
func Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Print(args...)
***REMOVED***

// Println is analogous to `fmt.Println` with termial detection.
func Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Println(args...)
***REMOVED***

// Printf is analogous to `fmt.Printf` with termial detection.
func Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Printf(format, args...)
***REMOVED***

func Black(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Black(msg, styles...)
***REMOVED***

func Red(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Red(msg, styles...)
***REMOVED***

func Green(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Green(msg, styles...)
***REMOVED***

func Yellow(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Yellow(msg, styles...)
***REMOVED***

func Blue(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Blue(msg, styles...)
***REMOVED***

func Magenta(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Magenta(msg, styles...)
***REMOVED***

func Cyan(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Cyan(msg, styles...)
***REMOVED***

func White(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.White(msg, styles...)
***REMOVED***

func Grey(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Grey(msg, styles...)
***REMOVED***

func BlackBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.BlackBg(msg, styles...)
***REMOVED***

func RedBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.RedBg(msg, styles...)
***REMOVED***

func GreenBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.GreenBg(msg, styles...)
***REMOVED***

func YellowBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.YellowBg(msg, styles...)
***REMOVED***

func BlueBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.BlueBg(msg, styles...)
***REMOVED***

func MagentaBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.MagentaBg(msg, styles...)
***REMOVED***

func CyanBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.CyanBg(msg, styles...)
***REMOVED***

func WhiteBg(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.WhiteBg(msg, styles...)
***REMOVED***

func Reset(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Reset(msg, styles...)
***REMOVED***

func Bold(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Bold(msg, styles...)
***REMOVED***

func Dim(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Dim(msg, styles...)
***REMOVED***

func Italic(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Italic(msg, styles...)
***REMOVED***

func Underline(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Underline(msg, styles...)
***REMOVED***

func Inverse(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Inverse(msg, styles...)
***REMOVED***

func Hidden(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Hidden(msg, styles...)
***REMOVED***

func Strikeout(msg interface***REMOVED******REMOVED***, styles ...string) string ***REMOVED***
	return global.Strikeout(msg, styles...)
***REMOVED***
