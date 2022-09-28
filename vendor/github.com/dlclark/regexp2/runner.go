package regexp2

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dlclark/regexp2/syntax"
)

type runner struct ***REMOVED***
	re   *Regexp
	code *syntax.Code

	runtextstart int // starting point for search

	runtext    []rune // text to search
	runtextpos int    // current position in text
	runtextend int

	// The backtracking stack.  Opcodes use this to store data regarding
	// what they have matched and where to backtrack to.  Each "frame" on
	// the stack takes the form of [CodePosition Data1 Data2...], where
	// CodePosition is the position of the current opcode and
	// the data values are all optional.  The CodePosition can be negative, and
	// these values (also called "back2") are used by the BranchMark family of opcodes
	// to indicate whether they are backtracking after a successful or failed
	// match.
	// When we backtrack, we pop the CodePosition off the stack, set the current
	// instruction pointer to that code position, and mark the opcode
	// with a backtracking flag ("Back").  Each opcode then knows how to
	// handle its own data.
	runtrack    []int
	runtrackpos int

	// This stack is used to track text positions across different opcodes.
	// For example, in /(a*b)+/, the parentheses result in a SetMark/CaptureMark
	// pair. SetMark records the text position before we match a*b.  Then
	// CaptureMark uses that position to figure out where the capture starts.
	// Opcodes which push onto this stack are always paired with other opcodes
	// which will pop the value from it later.  A successful match should mean
	// that this stack is empty.
	runstack    []int
	runstackpos int

	// The crawl stack is used to keep track of captures.  Every time a group
	// has a capture, we push its group number onto the runcrawl stack.  In
	// the case of a balanced match, we push BOTH groups onto the stack.
	runcrawl    []int
	runcrawlpos int

	runtrackcount int // count of states that may do backtracking

	runmatch *Match // result object

	ignoreTimeout       bool
	timeout             time.Duration // timeout in milliseconds (needed for actual)
	timeoutChecksToSkip int
	timeoutAt           time.Time

	operator        syntax.InstOp
	codepos         int
	rightToLeft     bool
	caseInsensitive bool
***REMOVED***

// run searches for matches and can continue from the previous match
//
// quick is usually false, but can be true to not return matches, just put it in caches
// textstart is -1 to start at the "beginning" (depending on Right-To-Left), otherwise an index in input
// input is the string to search for our regex pattern
func (re *Regexp) run(quick bool, textstart int, input []rune) (*Match, error) ***REMOVED***

	// get a cached runner
	runner := re.getRunner()
	defer re.putRunner(runner)

	if textstart < 0 ***REMOVED***
		if re.RightToLeft() ***REMOVED***
			textstart = len(input)
		***REMOVED*** else ***REMOVED***
			textstart = 0
		***REMOVED***
	***REMOVED***

	return runner.scan(input, textstart, quick, re.MatchTimeout)
***REMOVED***

// Scans the string to find the first match. Uses the Match object
// both to feed text in and as a place to store matches that come out.
//
// All the action is in the Go() method. Our
// responsibility is to load up the class members before
// calling Go.
//
// The optimizer can compute a set of candidate starting characters,
// and we could use a separate method Skip() that will quickly scan past
// any characters that we know can't match.
func (r *runner) scan(rt []rune, textstart int, quick bool, timeout time.Duration) (*Match, error) ***REMOVED***
	r.timeout = timeout
	r.ignoreTimeout = (time.Duration(math.MaxInt64) == timeout)
	r.runtextstart = textstart
	r.runtext = rt
	r.runtextend = len(rt)

	stoppos := r.runtextend
	bump := 1

	if r.re.RightToLeft() ***REMOVED***
		bump = -1
		stoppos = 0
	***REMOVED***

	r.runtextpos = textstart
	initted := false

	r.startTimeoutWatch()
	for ***REMOVED***
		if r.re.Debug() ***REMOVED***
			//fmt.Printf("\nSearch content: %v\n", string(r.runtext))
			fmt.Printf("\nSearch range: from 0 to %v\n", r.runtextend)
			fmt.Printf("Firstchar search starting at %v stopping at %v\n", r.runtextpos, stoppos)
		***REMOVED***

		if r.findFirstChar() ***REMOVED***
			if err := r.checkTimeout(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if !initted ***REMOVED***
				r.initMatch()
				initted = true
			***REMOVED***

			if r.re.Debug() ***REMOVED***
				fmt.Printf("Executing engine starting at %v\n\n", r.runtextpos)
			***REMOVED***

			if err := r.execute(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if r.runmatch.matchcount[0] > 0 ***REMOVED***
				// We'll return a match even if it touches a previous empty match
				return r.tidyMatch(quick), nil
			***REMOVED***

			// reset state for another go
			r.runtrackpos = len(r.runtrack)
			r.runstackpos = len(r.runstack)
			r.runcrawlpos = len(r.runcrawl)
		***REMOVED***

		// failure!

		if r.runtextpos == stoppos ***REMOVED***
			r.tidyMatch(true)
			return nil, nil
		***REMOVED***

		// Recognize leading []* and various anchors, and bump on failure accordingly

		// r.bump by one and start again

		r.runtextpos += bump
	***REMOVED***
	// We never get here
***REMOVED***

func (r *runner) execute() error ***REMOVED***

	r.goTo(0)

	for ***REMOVED***

		if r.re.Debug() ***REMOVED***
			r.dumpState()
		***REMOVED***

		if err := r.checkTimeout(); err != nil ***REMOVED***
			return err
		***REMOVED***

		switch r.operator ***REMOVED***
		case syntax.Stop:
			return nil

		case syntax.Nothing:
			break

		case syntax.Goto:
			r.goTo(r.operand(0))
			continue

		case syntax.Testref:
			if !r.runmatch.isMatched(r.operand(0)) ***REMOVED***
				break
			***REMOVED***
			r.advance(1)
			continue

		case syntax.Lazybranch:
			r.trackPush1(r.textPos())
			r.advance(1)
			continue

		case syntax.Lazybranch | syntax.Back:
			r.trackPop()
			r.textto(r.trackPeek())
			r.goTo(r.operand(0))
			continue

		case syntax.Setmark:
			r.stackPush(r.textPos())
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Nullmark:
			r.stackPush(-1)
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Setmark | syntax.Back, syntax.Nullmark | syntax.Back:
			r.stackPop()
			break

		case syntax.Getmark:
			r.stackPop()
			r.trackPush1(r.stackPeek())
			r.textto(r.stackPeek())
			r.advance(0)
			continue

		case syntax.Getmark | syntax.Back:
			r.trackPop()
			r.stackPush(r.trackPeek())
			break

		case syntax.Capturemark:
			if r.operand(1) != -1 && !r.runmatch.isMatched(r.operand(1)) ***REMOVED***
				break
			***REMOVED***
			r.stackPop()
			if r.operand(1) != -1 ***REMOVED***
				r.transferCapture(r.operand(0), r.operand(1), r.stackPeek(), r.textPos())
			***REMOVED*** else ***REMOVED***
				r.capture(r.operand(0), r.stackPeek(), r.textPos())
			***REMOVED***
			r.trackPush1(r.stackPeek())

			r.advance(2)

			continue

		case syntax.Capturemark | syntax.Back:
			r.trackPop()
			r.stackPush(r.trackPeek())
			r.uncapture()
			if r.operand(0) != -1 && r.operand(1) != -1 ***REMOVED***
				r.uncapture()
			***REMOVED***

			break

		case syntax.Branchmark:
			r.stackPop()

			matched := r.textPos() - r.stackPeek()

			if matched != 0 ***REMOVED*** // Nonempty match -> loop now
				r.trackPush2(r.stackPeek(), r.textPos()) // Save old mark, textpos
				r.stackPush(r.textPos())                 // Make new mark
				r.goTo(r.operand(0))                     // Loop
			***REMOVED*** else ***REMOVED*** // Empty match -> straight now
				r.trackPushNeg1(r.stackPeek()) // Save old mark
				r.advance(1)                   // Straight
			***REMOVED***
			continue

		case syntax.Branchmark | syntax.Back:
			r.trackPopN(2)
			r.stackPop()
			r.textto(r.trackPeekN(1))      // Recall position
			r.trackPushNeg1(r.trackPeek()) // Save old mark
			r.advance(1)                   // Straight
			continue

		case syntax.Branchmark | syntax.Back2:
			r.trackPop()
			r.stackPush(r.trackPeek()) // Recall old mark
			break                      // Backtrack

		case syntax.Lazybranchmark:
			***REMOVED***
				// We hit this the first time through a lazy loop and after each
				// successful match of the inner expression.  It simply continues
				// on and doesn't loop.
				r.stackPop()

				oldMarkPos := r.stackPeek()

				if r.textPos() != oldMarkPos ***REMOVED*** // Nonempty match -> try to loop again by going to 'back' state
					if oldMarkPos != -1 ***REMOVED***
						r.trackPush2(oldMarkPos, r.textPos()) // Save old mark, textpos
					***REMOVED*** else ***REMOVED***
						r.trackPush2(r.textPos(), r.textPos())
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// The inner expression found an empty match, so we'll go directly to 'back2' if we
					// backtrack.  In this case, we need to push something on the stack, since back2 pops.
					// However, in the case of ()+? or similar, this empty match may be legitimate, so push the text
					// position associated with that empty match.
					r.stackPush(oldMarkPos)

					r.trackPushNeg1(r.stackPeek()) // Save old mark
				***REMOVED***
				r.advance(1)
				continue
			***REMOVED***

		case syntax.Lazybranchmark | syntax.Back:

			// After the first time, Lazybranchmark | syntax.Back occurs
			// with each iteration of the loop, and therefore with every attempted
			// match of the inner expression.  We'll try to match the inner expression,
			// then go back to Lazybranchmark if successful.  If the inner expression
			// fails, we go to Lazybranchmark | syntax.Back2

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.trackPushNeg1(r.trackPeek()) // Save old mark
			r.stackPush(pos)               // Make new mark
			r.textto(pos)                  // Recall position
			r.goTo(r.operand(0))           // Loop
			continue

		case syntax.Lazybranchmark | syntax.Back2:
			// The lazy loop has failed.  We'll do a true backtrack and
			// start over before the lazy loop.
			r.stackPop()
			r.trackPop()
			r.stackPush(r.trackPeek()) // Recall old mark
			break

		case syntax.Setcount:
			r.stackPush2(r.textPos(), r.operand(0))
			r.trackPush()
			r.advance(1)
			continue

		case syntax.Nullcount:
			r.stackPush2(-1, r.operand(0))
			r.trackPush()
			r.advance(1)
			continue

		case syntax.Setcount | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Nullcount | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Branchcount:
			// r.stackPush:
			//  0: Mark
			//  1: Count

			r.stackPopN(2)
			mark := r.stackPeek()
			count := r.stackPeekN(1)
			matched := r.textPos() - mark

			if count >= r.operand(1) || (matched == 0 && count >= 0) ***REMOVED*** // Max loops or empty match -> straight now
				r.trackPushNeg2(mark, count) // Save old mark, count
				r.advance(2)                 // Straight
			***REMOVED*** else ***REMOVED*** // Nonempty match -> count+loop now
				r.trackPush1(mark)                 // remember mark
				r.stackPush2(r.textPos(), count+1) // Make new mark, incr count
				r.goTo(r.operand(0))               // Loop
			***REMOVED***
			continue

		case syntax.Branchcount | syntax.Back:
			// r.trackPush:
			//  0: Previous mark
			// r.stackPush:
			//  0: Mark (= current pos, discarded)
			//  1: Count
			r.trackPop()
			r.stackPopN(2)
			if r.stackPeekN(1) > 0 ***REMOVED*** // Positive -> can go straight
				r.textto(r.stackPeek())                           // Zap to mark
				r.trackPushNeg2(r.trackPeek(), r.stackPeekN(1)-1) // Save old mark, old count
				r.advance(2)                                      // Straight
				continue
			***REMOVED***
			r.stackPush2(r.trackPeek(), r.stackPeekN(1)-1) // recall old mark, old count
			break

		case syntax.Branchcount | syntax.Back2:
			// r.trackPush:
			//  0: Previous mark
			//  1: Previous count
			r.trackPopN(2)
			r.stackPush2(r.trackPeek(), r.trackPeekN(1)) // Recall old mark, old count
			break                                        // Backtrack

		case syntax.Lazybranchcount:
			// r.stackPush:
			//  0: Mark
			//  1: Count

			r.stackPopN(2)
			mark := r.stackPeek()
			count := r.stackPeekN(1)

			if count < 0 ***REMOVED*** // Negative count -> loop now
				r.trackPushNeg1(mark)              // Save old mark
				r.stackPush2(r.textPos(), count+1) // Make new mark, incr count
				r.goTo(r.operand(0))               // Loop
			***REMOVED*** else ***REMOVED*** // Nonneg count -> straight now
				r.trackPush3(mark, count, r.textPos()) // Save mark, count, position
				r.advance(2)                           // Straight
			***REMOVED***
			continue

		case syntax.Lazybranchcount | syntax.Back:
			// r.trackPush:
			//  0: Mark
			//  1: Count
			//  2: r.textPos

			r.trackPopN(3)
			mark := r.trackPeek()
			textpos := r.trackPeekN(2)

			if r.trackPeekN(1) < r.operand(1) && textpos != mark ***REMOVED*** // Under limit and not empty match -> loop
				r.textto(textpos)                        // Recall position
				r.stackPush2(textpos, r.trackPeekN(1)+1) // Make new mark, incr count
				r.trackPushNeg1(mark)                    // Save old mark
				r.goTo(r.operand(0))                     // Loop
				continue
			***REMOVED*** else ***REMOVED*** // Max loops or empty match -> backtrack
				r.stackPush2(r.trackPeek(), r.trackPeekN(1)) // Recall old mark, count
				break                                        // backtrack
			***REMOVED***

		case syntax.Lazybranchcount | syntax.Back2:
			// r.trackPush:
			//  0: Previous mark
			// r.stackPush:
			//  0: Mark (== current pos, discarded)
			//  1: Count
			r.trackPop()
			r.stackPopN(2)
			r.stackPush2(r.trackPeek(), r.stackPeekN(1)-1) // Recall old mark, count
			break                                          // Backtrack

		case syntax.Setjump:
			r.stackPush2(r.trackpos(), r.crawlpos())
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Setjump | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Backjump:
			// r.stackPush:
			//  0: Saved trackpos
			//  1: r.crawlpos
			r.stackPopN(2)
			r.trackto(r.stackPeek())

			for r.crawlpos() != r.stackPeekN(1) ***REMOVED***
				r.uncapture()
			***REMOVED***

			break

		case syntax.Forejump:
			// r.stackPush:
			//  0: Saved trackpos
			//  1: r.crawlpos
			r.stackPopN(2)
			r.trackto(r.stackPeek())
			r.trackPush1(r.stackPeekN(1))
			r.advance(0)
			continue

		case syntax.Forejump | syntax.Back:
			// r.trackPush:
			//  0: r.crawlpos
			r.trackPop()

			for r.crawlpos() != r.trackPeek() ***REMOVED***
				r.uncapture()
			***REMOVED***

			break

		case syntax.Bol:
			if r.leftchars() > 0 && r.charAt(r.textPos()-1) != '\n' ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.Eol:
			if r.rightchars() > 0 && r.charAt(r.textPos()) != '\n' ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.Boundary:
			if !r.isBoundary(r.textPos(), 0, r.runtextend) ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.Nonboundary:
			if r.isBoundary(r.textPos(), 0, r.runtextend) ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.ECMABoundary:
			if !r.isECMABoundary(r.textPos(), 0, r.runtextend) ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.NonECMABoundary:
			if r.isECMABoundary(r.textPos(), 0, r.runtextend) ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.Beginning:
			if r.leftchars() > 0 ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.Start:
			if r.textPos() != r.textstart() ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.EndZ:
			rchars := r.rightchars()
			if rchars > 1 ***REMOVED***
				break
			***REMOVED***
			// RE2 and EcmaScript define $ as "asserts position at the end of the string"
			// PCRE/.NET adds "or before the line terminator right at the end of the string (if any)"
			if (r.re.options & (RE2 | ECMAScript)) != 0 ***REMOVED***
				// RE2/Ecmascript mode
				if rchars > 0 ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if rchars == 1 && r.charAt(r.textPos()) != '\n' ***REMOVED***
				// "regular" mode
				break
			***REMOVED***

			r.advance(0)
			continue

		case syntax.End:
			if r.rightchars() > 0 ***REMOVED***
				break
			***REMOVED***
			r.advance(0)
			continue

		case syntax.One:
			if r.forwardchars() < 1 || r.forwardcharnext() != rune(r.operand(0)) ***REMOVED***
				break
			***REMOVED***

			r.advance(1)
			continue

		case syntax.Notone:
			if r.forwardchars() < 1 || r.forwardcharnext() == rune(r.operand(0)) ***REMOVED***
				break
			***REMOVED***

			r.advance(1)
			continue

		case syntax.Set:

			if r.forwardchars() < 1 || !r.code.Sets[r.operand(0)].CharIn(r.forwardcharnext()) ***REMOVED***
				break
			***REMOVED***

			r.advance(1)
			continue

		case syntax.Multi:
			if !r.runematch(r.code.Strings[r.operand(0)]) ***REMOVED***
				break
			***REMOVED***

			r.advance(1)
			continue

		case syntax.Ref:

			capnum := r.operand(0)

			if r.runmatch.isMatched(capnum) ***REMOVED***
				if !r.refmatch(r.runmatch.matchIndex(capnum), r.runmatch.matchLength(capnum)) ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if (r.re.options & ECMAScript) == 0 ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			r.advance(1)
			continue

		case syntax.Onerep:

			c := r.operand(1)

			if r.forwardchars() < c ***REMOVED***
				break
			***REMOVED***

			ch := rune(r.operand(0))

			for c > 0 ***REMOVED***
				if r.forwardcharnext() != ch ***REMOVED***
					goto BreakBackward
				***REMOVED***
				c--
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Notonerep:

			c := r.operand(1)

			if r.forwardchars() < c ***REMOVED***
				break
			***REMOVED***
			ch := rune(r.operand(0))

			for c > 0 ***REMOVED***
				if r.forwardcharnext() == ch ***REMOVED***
					goto BreakBackward
				***REMOVED***
				c--
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Setrep:

			c := r.operand(1)

			if r.forwardchars() < c ***REMOVED***
				break
			***REMOVED***

			set := r.code.Sets[r.operand(0)]

			for c > 0 ***REMOVED***
				if !set.CharIn(r.forwardcharnext()) ***REMOVED***
					goto BreakBackward
				***REMOVED***
				c--
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Oneloop:

			c := r.operand(1)

			if c > r.forwardchars() ***REMOVED***
				c = r.forwardchars()
			***REMOVED***

			ch := rune(r.operand(0))
			i := c

			for ; i > 0; i-- ***REMOVED***
				if r.forwardcharnext() != ch ***REMOVED***
					r.backwardnext()
					break
				***REMOVED***
			***REMOVED***

			if c > i ***REMOVED***
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Notoneloop:

			c := r.operand(1)

			if c > r.forwardchars() ***REMOVED***
				c = r.forwardchars()
			***REMOVED***

			ch := rune(r.operand(0))
			i := c

			for ; i > 0; i-- ***REMOVED***
				if r.forwardcharnext() == ch ***REMOVED***
					r.backwardnext()
					break
				***REMOVED***
			***REMOVED***

			if c > i ***REMOVED***
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Setloop:

			c := r.operand(1)

			if c > r.forwardchars() ***REMOVED***
				c = r.forwardchars()
			***REMOVED***

			set := r.code.Sets[r.operand(0)]
			i := c

			for ; i > 0; i-- ***REMOVED***
				if !set.CharIn(r.forwardcharnext()) ***REMOVED***
					r.backwardnext()
					break
				***REMOVED***
			***REMOVED***

			if c > i ***REMOVED***
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Oneloop | syntax.Back, syntax.Notoneloop | syntax.Back:

			r.trackPopN(2)
			i := r.trackPeek()
			pos := r.trackPeekN(1)

			r.textto(pos)

			if i > 0 ***REMOVED***
				r.trackPush2(i-1, pos-r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Setloop | syntax.Back:

			r.trackPopN(2)
			i := r.trackPeek()
			pos := r.trackPeekN(1)

			r.textto(pos)

			if i > 0 ***REMOVED***
				r.trackPush2(i-1, pos-r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Onelazy, syntax.Notonelazy:

			c := r.operand(1)

			if c > r.forwardchars() ***REMOVED***
				c = r.forwardchars()
			***REMOVED***

			if c > 0 ***REMOVED***
				r.trackPush2(c-1, r.textPos())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Setlazy:

			c := r.operand(1)

			if c > r.forwardchars() ***REMOVED***
				c = r.forwardchars()
			***REMOVED***

			if c > 0 ***REMOVED***
				r.trackPush2(c-1, r.textPos())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Onelazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if r.forwardcharnext() != rune(r.operand(0)) ***REMOVED***
				break
			***REMOVED***

			i := r.trackPeek()

			if i > 0 ***REMOVED***
				r.trackPush2(i-1, pos+r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Notonelazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if r.forwardcharnext() == rune(r.operand(0)) ***REMOVED***
				break
			***REMOVED***

			i := r.trackPeek()

			if i > 0 ***REMOVED***
				r.trackPush2(i-1, pos+r.bump())
			***REMOVED***

			r.advance(2)
			continue

		case syntax.Setlazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if !r.code.Sets[r.operand(0)].CharIn(r.forwardcharnext()) ***REMOVED***
				break
			***REMOVED***

			i := r.trackPeek()

			if i > 0 ***REMOVED***
				r.trackPush2(i-1, pos+r.bump())
			***REMOVED***

			r.advance(2)
			continue

		default:
			return errors.New("unknown state in regex runner")
		***REMOVED***

	BreakBackward:
		;

		// "break Backward" comes here:
		r.backtrack()
	***REMOVED***
***REMOVED***

// increase the size of stack and track storage
func (r *runner) ensureStorage() ***REMOVED***
	if r.runstackpos < r.runtrackcount*4 ***REMOVED***
		doubleIntSlice(&r.runstack, &r.runstackpos)
	***REMOVED***
	if r.runtrackpos < r.runtrackcount*4 ***REMOVED***
		doubleIntSlice(&r.runtrack, &r.runtrackpos)
	***REMOVED***
***REMOVED***

func doubleIntSlice(s *[]int, pos *int) ***REMOVED***
	oldLen := len(*s)
	newS := make([]int, oldLen*2)

	copy(newS[oldLen:], *s)
	*pos += oldLen
	*s = newS
***REMOVED***

// Save a number on the longjump unrolling stack
func (r *runner) crawl(i int) ***REMOVED***
	if r.runcrawlpos == 0 ***REMOVED***
		doubleIntSlice(&r.runcrawl, &r.runcrawlpos)
	***REMOVED***
	r.runcrawlpos--
	r.runcrawl[r.runcrawlpos] = i
***REMOVED***

// Remove a number from the longjump unrolling stack
func (r *runner) popcrawl() int ***REMOVED***
	val := r.runcrawl[r.runcrawlpos]
	r.runcrawlpos++
	return val
***REMOVED***

// Get the height of the stack
func (r *runner) crawlpos() int ***REMOVED***
	return len(r.runcrawl) - r.runcrawlpos
***REMOVED***

func (r *runner) advance(i int) ***REMOVED***
	r.codepos += (i + 1)
	r.setOperator(r.code.Codes[r.codepos])
***REMOVED***

func (r *runner) goTo(newpos int) ***REMOVED***
	// when branching backward or in place, ensure storage
	if newpos <= r.codepos ***REMOVED***
		r.ensureStorage()
	***REMOVED***

	r.setOperator(r.code.Codes[newpos])
	r.codepos = newpos
***REMOVED***

func (r *runner) textto(newpos int) ***REMOVED***
	r.runtextpos = newpos
***REMOVED***

func (r *runner) trackto(newpos int) ***REMOVED***
	r.runtrackpos = len(r.runtrack) - newpos
***REMOVED***

func (r *runner) textstart() int ***REMOVED***
	return r.runtextstart
***REMOVED***

func (r *runner) textPos() int ***REMOVED***
	return r.runtextpos
***REMOVED***

// push onto the backtracking stack
func (r *runner) trackpos() int ***REMOVED***
	return len(r.runtrack) - r.runtrackpos
***REMOVED***

func (r *runner) trackPush() ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
***REMOVED***

func (r *runner) trackPush1(I1 int) ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
***REMOVED***

func (r *runner) trackPush2(I1, I2 int) ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
***REMOVED***

func (r *runner) trackPush3(I1, I2, I3 int) ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I3
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
***REMOVED***

func (r *runner) trackPushNeg1(I1 int) ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = -r.codepos
***REMOVED***

func (r *runner) trackPushNeg2(I1, I2 int) ***REMOVED***
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = -r.codepos
***REMOVED***

func (r *runner) backtrack() ***REMOVED***
	newpos := r.runtrack[r.runtrackpos]
	r.runtrackpos++

	if r.re.Debug() ***REMOVED***
		if newpos < 0 ***REMOVED***
			fmt.Printf("       Backtracking (back2) to code position %v\n", -newpos)
		***REMOVED*** else ***REMOVED***
			fmt.Printf("       Backtracking to code position %v\n", newpos)
		***REMOVED***
	***REMOVED***

	if newpos < 0 ***REMOVED***
		newpos = -newpos
		r.setOperator(r.code.Codes[newpos] | syntax.Back2)
	***REMOVED*** else ***REMOVED***
		r.setOperator(r.code.Codes[newpos] | syntax.Back)
	***REMOVED***

	// When branching backward, ensure storage
	if newpos < r.codepos ***REMOVED***
		r.ensureStorage()
	***REMOVED***

	r.codepos = newpos
***REMOVED***

func (r *runner) setOperator(op int) ***REMOVED***
	r.caseInsensitive = (0 != (op & syntax.Ci))
	r.rightToLeft = (0 != (op & syntax.Rtl))
	r.operator = syntax.InstOp(op & ^(syntax.Rtl | syntax.Ci))
***REMOVED***

func (r *runner) trackPop() ***REMOVED***
	r.runtrackpos++
***REMOVED***

// pop framesize items from the backtracking stack
func (r *runner) trackPopN(framesize int) ***REMOVED***
	r.runtrackpos += framesize
***REMOVED***

// Technically we are actually peeking at items already popped.  So if you want to
// get and pop the top item from the stack, you do
// r.trackPop();
// r.trackPeek();
func (r *runner) trackPeek() int ***REMOVED***
	return r.runtrack[r.runtrackpos-1]
***REMOVED***

// get the ith element down on the backtracking stack
func (r *runner) trackPeekN(i int) int ***REMOVED***
	return r.runtrack[r.runtrackpos-i-1]
***REMOVED***

// Push onto the grouping stack
func (r *runner) stackPush(I1 int) ***REMOVED***
	r.runstackpos--
	r.runstack[r.runstackpos] = I1
***REMOVED***

func (r *runner) stackPush2(I1, I2 int) ***REMOVED***
	r.runstackpos--
	r.runstack[r.runstackpos] = I1
	r.runstackpos--
	r.runstack[r.runstackpos] = I2
***REMOVED***

func (r *runner) stackPop() ***REMOVED***
	r.runstackpos++
***REMOVED***

// pop framesize items from the grouping stack
func (r *runner) stackPopN(framesize int) ***REMOVED***
	r.runstackpos += framesize
***REMOVED***

// Technically we are actually peeking at items already popped.  So if you want to
// get and pop the top item from the stack, you do
// r.stackPop();
// r.stackPeek();
func (r *runner) stackPeek() int ***REMOVED***
	return r.runstack[r.runstackpos-1]
***REMOVED***

// get the ith element down on the grouping stack
func (r *runner) stackPeekN(i int) int ***REMOVED***
	return r.runstack[r.runstackpos-i-1]
***REMOVED***

func (r *runner) operand(i int) int ***REMOVED***
	return r.code.Codes[r.codepos+i+1]
***REMOVED***

func (r *runner) leftchars() int ***REMOVED***
	return r.runtextpos
***REMOVED***

func (r *runner) rightchars() int ***REMOVED***
	return r.runtextend - r.runtextpos
***REMOVED***

func (r *runner) bump() int ***REMOVED***
	if r.rightToLeft ***REMOVED***
		return -1
	***REMOVED***
	return 1
***REMOVED***

func (r *runner) forwardchars() int ***REMOVED***
	if r.rightToLeft ***REMOVED***
		return r.runtextpos
	***REMOVED***
	return r.runtextend - r.runtextpos
***REMOVED***

func (r *runner) forwardcharnext() rune ***REMOVED***
	var ch rune
	if r.rightToLeft ***REMOVED***
		r.runtextpos--
		ch = r.runtext[r.runtextpos]
	***REMOVED*** else ***REMOVED***
		ch = r.runtext[r.runtextpos]
		r.runtextpos++
	***REMOVED***

	if r.caseInsensitive ***REMOVED***
		return unicode.ToLower(ch)
	***REMOVED***
	return ch
***REMOVED***

func (r *runner) runematch(str []rune) bool ***REMOVED***
	var pos int

	c := len(str)
	if !r.rightToLeft ***REMOVED***
		if r.runtextend-r.runtextpos < c ***REMOVED***
			return false
		***REMOVED***

		pos = r.runtextpos + c
	***REMOVED*** else ***REMOVED***
		if r.runtextpos-0 < c ***REMOVED***
			return false
		***REMOVED***

		pos = r.runtextpos
	***REMOVED***

	if !r.caseInsensitive ***REMOVED***
		for c != 0 ***REMOVED***
			c--
			pos--
			if str[c] != r.runtext[pos] ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for c != 0 ***REMOVED***
			c--
			pos--
			if str[c] != unicode.ToLower(r.runtext[pos]) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !r.rightToLeft ***REMOVED***
		pos += len(str)
	***REMOVED***

	r.runtextpos = pos

	return true
***REMOVED***

func (r *runner) refmatch(index, len int) bool ***REMOVED***
	var c, pos, cmpos int

	if !r.rightToLeft ***REMOVED***
		if r.runtextend-r.runtextpos < len ***REMOVED***
			return false
		***REMOVED***

		pos = r.runtextpos + len
	***REMOVED*** else ***REMOVED***
		if r.runtextpos-0 < len ***REMOVED***
			return false
		***REMOVED***

		pos = r.runtextpos
	***REMOVED***
	cmpos = index + len

	c = len

	if !r.caseInsensitive ***REMOVED***
		for c != 0 ***REMOVED***
			c--
			cmpos--
			pos--
			if r.runtext[cmpos] != r.runtext[pos] ***REMOVED***
				return false
			***REMOVED***

		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for c != 0 ***REMOVED***
			c--
			cmpos--
			pos--

			if unicode.ToLower(r.runtext[cmpos]) != unicode.ToLower(r.runtext[pos]) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !r.rightToLeft ***REMOVED***
		pos += len
	***REMOVED***

	r.runtextpos = pos

	return true
***REMOVED***

func (r *runner) backwardnext() ***REMOVED***
	if r.rightToLeft ***REMOVED***
		r.runtextpos++
	***REMOVED*** else ***REMOVED***
		r.runtextpos--
	***REMOVED***
***REMOVED***

func (r *runner) charAt(j int) rune ***REMOVED***
	return r.runtext[j]
***REMOVED***

func (r *runner) findFirstChar() bool ***REMOVED***

	if 0 != (r.code.Anchors & (syntax.AnchorBeginning | syntax.AnchorStart | syntax.AnchorEndZ | syntax.AnchorEnd)) ***REMOVED***
		if !r.code.RightToLeft ***REMOVED***
			if (0 != (r.code.Anchors&syntax.AnchorBeginning) && r.runtextpos > 0) ||
				(0 != (r.code.Anchors&syntax.AnchorStart) && r.runtextpos > r.runtextstart) ***REMOVED***
				r.runtextpos = r.runtextend
				return false
			***REMOVED***
			if 0 != (r.code.Anchors&syntax.AnchorEndZ) && r.runtextpos < r.runtextend-1 ***REMOVED***
				r.runtextpos = r.runtextend - 1
			***REMOVED*** else if 0 != (r.code.Anchors&syntax.AnchorEnd) && r.runtextpos < r.runtextend ***REMOVED***
				r.runtextpos = r.runtextend
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if (0 != (r.code.Anchors&syntax.AnchorEnd) && r.runtextpos < r.runtextend) ||
				(0 != (r.code.Anchors&syntax.AnchorEndZ) && (r.runtextpos < r.runtextend-1 ||
					(r.runtextpos == r.runtextend-1 && r.charAt(r.runtextpos) != '\n'))) ||
				(0 != (r.code.Anchors&syntax.AnchorStart) && r.runtextpos < r.runtextstart) ***REMOVED***
				r.runtextpos = 0
				return false
			***REMOVED***
			if 0 != (r.code.Anchors&syntax.AnchorBeginning) && r.runtextpos > 0 ***REMOVED***
				r.runtextpos = 0
			***REMOVED***
		***REMOVED***

		if r.code.BmPrefix != nil ***REMOVED***
			return r.code.BmPrefix.IsMatch(r.runtext, r.runtextpos, 0, r.runtextend)
		***REMOVED***

		return true // found a valid start or end anchor
	***REMOVED*** else if r.code.BmPrefix != nil ***REMOVED***
		r.runtextpos = r.code.BmPrefix.Scan(r.runtext, r.runtextpos, 0, r.runtextend)

		if r.runtextpos == -1 ***REMOVED***
			if r.code.RightToLeft ***REMOVED***
				r.runtextpos = 0
			***REMOVED*** else ***REMOVED***
				r.runtextpos = r.runtextend
			***REMOVED***
			return false
		***REMOVED***

		return true
	***REMOVED*** else if r.code.FcPrefix == nil ***REMOVED***
		return true
	***REMOVED***

	r.rightToLeft = r.code.RightToLeft
	r.caseInsensitive = r.code.FcPrefix.CaseInsensitive

	set := r.code.FcPrefix.PrefixSet
	if set.IsSingleton() ***REMOVED***
		ch := set.SingletonChar()
		for i := r.forwardchars(); i > 0; i-- ***REMOVED***
			if ch == r.forwardcharnext() ***REMOVED***
				r.backwardnext()
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := r.forwardchars(); i > 0; i-- ***REMOVED***
			n := r.forwardcharnext()
			//fmt.Printf("%v in %v: %v\n", string(n), set.String(), set.CharIn(n))
			if set.CharIn(n) ***REMOVED***
				r.backwardnext()
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (r *runner) initMatch() ***REMOVED***
	// Use a hashtable'ed Match object if the capture numbers are sparse

	if r.runmatch == nil ***REMOVED***
		if r.re.caps != nil ***REMOVED***
			r.runmatch = newMatchSparse(r.re, r.re.caps, r.re.capsize, r.runtext, r.runtextstart)
		***REMOVED*** else ***REMOVED***
			r.runmatch = newMatch(r.re, r.re.capsize, r.runtext, r.runtextstart)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.runmatch.reset(r.runtext, r.runtextstart)
	***REMOVED***

	// note we test runcrawl, because it is the last one to be allocated
	// If there is an alloc failure in the middle of the three allocations,
	// we may still return to reuse this instance, and we want to behave
	// as if the allocations didn't occur. (we used to test _trackcount != 0)

	if r.runcrawl != nil ***REMOVED***
		r.runtrackpos = len(r.runtrack)
		r.runstackpos = len(r.runstack)
		r.runcrawlpos = len(r.runcrawl)
		return
	***REMOVED***

	r.initTrackCount()

	tracksize := r.runtrackcount * 8
	stacksize := r.runtrackcount * 8

	if tracksize < 32 ***REMOVED***
		tracksize = 32
	***REMOVED***
	if stacksize < 16 ***REMOVED***
		stacksize = 16
	***REMOVED***

	r.runtrack = make([]int, tracksize)
	r.runtrackpos = tracksize

	r.runstack = make([]int, stacksize)
	r.runstackpos = stacksize

	r.runcrawl = make([]int, 32)
	r.runcrawlpos = 32
***REMOVED***

func (r *runner) tidyMatch(quick bool) *Match ***REMOVED***
	if !quick ***REMOVED***
		match := r.runmatch

		r.runmatch = nil

		match.tidy(r.runtextpos)
		return match
	***REMOVED*** else ***REMOVED***
		// send back our match -- it's not leaving the package, so it's safe to not clean it up
		// this reduces allocs for frequent calls to the "IsMatch" bool-only functions
		return r.runmatch
	***REMOVED***
***REMOVED***

// capture captures a subexpression. Note that the
// capnum used here has already been mapped to a non-sparse
// index (by the code generator RegexWriter).
func (r *runner) capture(capnum, start, end int) ***REMOVED***
	if end < start ***REMOVED***
		T := end
		end = start
		start = T
	***REMOVED***

	r.crawl(capnum)
	r.runmatch.addMatch(capnum, start, end-start)
***REMOVED***

// transferCapture captures a subexpression. Note that the
// capnum used here has already been mapped to a non-sparse
// index (by the code generator RegexWriter).
func (r *runner) transferCapture(capnum, uncapnum, start, end int) ***REMOVED***
	var start2, end2 int

	// these are the two intervals that are cancelling each other

	if end < start ***REMOVED***
		T := end
		end = start
		start = T
	***REMOVED***

	start2 = r.runmatch.matchIndex(uncapnum)
	end2 = start2 + r.runmatch.matchLength(uncapnum)

	// The new capture gets the innermost defined interval

	if start >= end2 ***REMOVED***
		end = start
		start = end2
	***REMOVED*** else if end <= start2 ***REMOVED***
		start = start2
	***REMOVED*** else ***REMOVED***
		if end > end2 ***REMOVED***
			end = end2
		***REMOVED***
		if start2 > start ***REMOVED***
			start = start2
		***REMOVED***
	***REMOVED***

	r.crawl(uncapnum)
	r.runmatch.balanceMatch(uncapnum)

	if capnum != -1 ***REMOVED***
		r.crawl(capnum)
		r.runmatch.addMatch(capnum, start, end-start)
	***REMOVED***
***REMOVED***

// revert the last capture
func (r *runner) uncapture() ***REMOVED***
	capnum := r.popcrawl()
	r.runmatch.removeMatch(capnum)
***REMOVED***

//debug

func (r *runner) dumpState() ***REMOVED***
	back := ""
	if r.operator&syntax.Back != 0 ***REMOVED***
		back = " Back"
	***REMOVED***
	if r.operator&syntax.Back2 != 0 ***REMOVED***
		back += " Back2"
	***REMOVED***
	fmt.Printf("Text:  %v\nTrack: %v\nStack: %v\n       %s%s\n\n",
		r.textposDescription(),
		r.stackDescription(r.runtrack, r.runtrackpos),
		r.stackDescription(r.runstack, r.runstackpos),
		r.code.OpcodeDescription(r.codepos),
		back)
***REMOVED***

func (r *runner) stackDescription(a []int, index int) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***

	fmt.Fprintf(buf, "%v/%v", len(a)-index, len(a))
	if buf.Len() < 8 ***REMOVED***
		buf.WriteString(strings.Repeat(" ", 8-buf.Len()))
	***REMOVED***

	buf.WriteRune('(')
	for i := index; i < len(a); i++ ***REMOVED***
		if i > index ***REMOVED***
			buf.WriteRune(' ')
		***REMOVED***

		buf.WriteString(strconv.Itoa(a[i]))
	***REMOVED***

	buf.WriteRune(')')

	return buf.String()
***REMOVED***

func (r *runner) textposDescription() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***

	buf.WriteString(strconv.Itoa(r.runtextpos))

	if buf.Len() < 8 ***REMOVED***
		buf.WriteString(strings.Repeat(" ", 8-buf.Len()))
	***REMOVED***

	if r.runtextpos > 0 ***REMOVED***
		buf.WriteString(syntax.CharDescription(r.runtext[r.runtextpos-1]))
	***REMOVED*** else ***REMOVED***
		buf.WriteRune('^')
	***REMOVED***

	buf.WriteRune('>')

	for i := r.runtextpos; i < r.runtextend; i++ ***REMOVED***
		buf.WriteString(syntax.CharDescription(r.runtext[i]))
	***REMOVED***
	if buf.Len() >= 64 ***REMOVED***
		buf.Truncate(61)
		buf.WriteString("...")
	***REMOVED*** else ***REMOVED***
		buf.WriteRune('$')
	***REMOVED***

	return buf.String()
***REMOVED***

// decide whether the pos
// at the specified index is a boundary or not. It's just not worth
// emitting inline code for this logic.
func (r *runner) isBoundary(index, startpos, endpos int) bool ***REMOVED***
	return (index > startpos && syntax.IsWordChar(r.runtext[index-1])) !=
		(index < endpos && syntax.IsWordChar(r.runtext[index]))
***REMOVED***

func (r *runner) isECMABoundary(index, startpos, endpos int) bool ***REMOVED***
	return (index > startpos && syntax.IsECMAWordChar(r.runtext[index-1])) !=
		(index < endpos && syntax.IsECMAWordChar(r.runtext[index]))
***REMOVED***

// this seems like a comment to justify randomly picking 1000 :-P
// We have determined this value in a series of experiments where x86 retail
// builds (ono-lab-optimized) were run on different pattern/input pairs. Larger values
// of TimeoutCheckFrequency did not tend to increase performance; smaller values
// of TimeoutCheckFrequency tended to slow down the execution.
const timeoutCheckFrequency int = 1000

func (r *runner) startTimeoutWatch() ***REMOVED***
	if r.ignoreTimeout ***REMOVED***
		return
	***REMOVED***

	r.timeoutChecksToSkip = timeoutCheckFrequency
	r.timeoutAt = time.Now().Add(r.timeout)
***REMOVED***

func (r *runner) checkTimeout() error ***REMOVED***
	if r.ignoreTimeout ***REMOVED***
		return nil
	***REMOVED***
	r.timeoutChecksToSkip--
	if r.timeoutChecksToSkip != 0 ***REMOVED***
		return nil
	***REMOVED***

	r.timeoutChecksToSkip = timeoutCheckFrequency
	return r.doCheckTimeout()
***REMOVED***

func (r *runner) doCheckTimeout() error ***REMOVED***
	current := time.Now()

	if current.Before(r.timeoutAt) ***REMOVED***
		return nil
	***REMOVED***

	if r.re.Debug() ***REMOVED***
		//Debug.WriteLine("")
		//Debug.WriteLine("RegEx match timeout occurred!")
		//Debug.WriteLine("Specified timeout:       " + TimeSpan.FromMilliseconds(_timeout).ToString())
		//Debug.WriteLine("Timeout check frequency: " + TimeoutCheckFrequency)
		//Debug.WriteLine("Search pattern:          " + _runregex._pattern)
		//Debug.WriteLine("Input:                   " + r.runtext)
		//Debug.WriteLine("About to throw RegexMatchTimeoutException.")
	***REMOVED***

	return fmt.Errorf("match timeout after %v on input `%v`", r.timeout, string(r.runtext))
***REMOVED***

func (r *runner) initTrackCount() ***REMOVED***
	r.runtrackcount = r.code.TrackCount
***REMOVED***

// getRunner returns a run to use for matching re.
// It uses the re's runner cache if possible, to avoid
// unnecessary allocation.
func (re *Regexp) getRunner() *runner ***REMOVED***
	re.muRun.Lock()
	if n := len(re.runner); n > 0 ***REMOVED***
		z := re.runner[n-1]
		re.runner = re.runner[:n-1]
		re.muRun.Unlock()
		return z
	***REMOVED***
	re.muRun.Unlock()
	z := &runner***REMOVED***
		re:   re,
		code: re.code,
	***REMOVED***
	return z
***REMOVED***

// putRunner returns a runner to the re's cache.
// There is no attempt to limit the size of the cache, so it will
// grow to the maximum number of simultaneous matches
// run using re.  (The cache empties when re gets garbage collected.)
func (re *Regexp) putRunner(r *runner) ***REMOVED***
	re.muRun.Lock()
	re.runner = append(re.runner, r)
	re.muRun.Unlock()
***REMOVED***
