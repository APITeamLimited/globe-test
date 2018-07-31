package pflag

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// BytesHex adapts []byte for use as a flag. Value of flag is HEX encoded
type bytesHexValue []byte

func (bytesHex bytesHexValue) String() string ***REMOVED***
	return fmt.Sprintf("%X", []byte(bytesHex))
***REMOVED***

func (bytesHex *bytesHexValue) Set(value string) error ***REMOVED***
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*bytesHex = bin

	return nil
***REMOVED***

func (*bytesHexValue) Type() string ***REMOVED***
	return "bytesHex"
***REMOVED***

func newBytesHexValue(val []byte, p *[]byte) *bytesHexValue ***REMOVED***
	*p = val
	return (*bytesHexValue)(p)
***REMOVED***

func bytesHexConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***

	bin, err := hex.DecodeString(sval)

	if err == nil ***REMOVED***
		return bin, nil
	***REMOVED***

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
***REMOVED***

// GetBytesHex return the []byte value of a flag with the given name
func (f *FlagSet) GetBytesHex(name string) ([]byte, error) ***REMOVED***
	val, err := f.getFlagType(name, "bytesHex", bytesHexConv)

	if err != nil ***REMOVED***
		return []byte***REMOVED******REMOVED***, err
	***REMOVED***

	return val.([]byte), nil
***REMOVED***

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (f *FlagSet) BytesHexVar(p *[]byte, name string, value []byte, usage string) ***REMOVED***
	f.VarP(newBytesHexValue(value, p), name, "", usage)
***REMOVED***

// BytesHexVarP is like BytesHexVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) ***REMOVED***
	f.VarP(newBytesHexValue(value, p), name, shorthand, usage)
***REMOVED***

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesHexVar(p *[]byte, name string, value []byte, usage string) ***REMOVED***
	CommandLine.VarP(newBytesHexValue(value, p), name, "", usage)
***REMOVED***

// BytesHexVarP is like BytesHexVar, but accepts a shorthand letter that can be used after a single dash.
func BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) ***REMOVED***
	CommandLine.VarP(newBytesHexValue(value, p), name, shorthand, usage)
***REMOVED***

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (f *FlagSet) BytesHex(name string, value []byte, usage string) *[]byte ***REMOVED***
	p := new([]byte)
	f.BytesHexVarP(p, name, "", value, usage)
	return p
***REMOVED***

// BytesHexP is like BytesHex, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesHexP(name, shorthand string, value []byte, usage string) *[]byte ***REMOVED***
	p := new([]byte)
	f.BytesHexVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesHex(name string, value []byte, usage string) *[]byte ***REMOVED***
	return CommandLine.BytesHexP(name, "", value, usage)
***REMOVED***

// BytesHexP is like BytesHex, but accepts a shorthand letter that can be used after a single dash.
func BytesHexP(name, shorthand string, value []byte, usage string) *[]byte ***REMOVED***
	return CommandLine.BytesHexP(name, shorthand, value, usage)
***REMOVED***
