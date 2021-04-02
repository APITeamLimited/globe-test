package pflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// BytesHex adapts []byte for use as a flag. Value of flag is HEX encoded
type bytesHexValue []byte

// String implements pflag.Value.String.
func (bytesHex bytesHexValue) String() string ***REMOVED***
	return fmt.Sprintf("%X", []byte(bytesHex))
***REMOVED***

// Set implements pflag.Value.Set.
func (bytesHex *bytesHexValue) Set(value string) error ***REMOVED***
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*bytesHex = bin

	return nil
***REMOVED***

// Type implements pflag.Value.Type.
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

// BytesBase64 adapts []byte for use as a flag. Value of flag is Base64 encoded
type bytesBase64Value []byte

// String implements pflag.Value.String.
func (bytesBase64 bytesBase64Value) String() string ***REMOVED***
	return base64.StdEncoding.EncodeToString([]byte(bytesBase64))
***REMOVED***

// Set implements pflag.Value.Set.
func (bytesBase64 *bytesBase64Value) Set(value string) error ***REMOVED***
	bin, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*bytesBase64 = bin

	return nil
***REMOVED***

// Type implements pflag.Value.Type.
func (*bytesBase64Value) Type() string ***REMOVED***
	return "bytesBase64"
***REMOVED***

func newBytesBase64Value(val []byte, p *[]byte) *bytesBase64Value ***REMOVED***
	*p = val
	return (*bytesBase64Value)(p)
***REMOVED***

func bytesBase64ValueConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***

	bin, err := base64.StdEncoding.DecodeString(sval)
	if err == nil ***REMOVED***
		return bin, nil
	***REMOVED***

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
***REMOVED***

// GetBytesBase64 return the []byte value of a flag with the given name
func (f *FlagSet) GetBytesBase64(name string) ([]byte, error) ***REMOVED***
	val, err := f.getFlagType(name, "bytesBase64", bytesBase64ValueConv)

	if err != nil ***REMOVED***
		return []byte***REMOVED******REMOVED***, err
	***REMOVED***

	return val.([]byte), nil
***REMOVED***

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (f *FlagSet) BytesBase64Var(p *[]byte, name string, value []byte, usage string) ***REMOVED***
	f.VarP(newBytesBase64Value(value, p), name, "", usage)
***REMOVED***

// BytesBase64VarP is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) ***REMOVED***
	f.VarP(newBytesBase64Value(value, p), name, shorthand, usage)
***REMOVED***

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesBase64Var(p *[]byte, name string, value []byte, usage string) ***REMOVED***
	CommandLine.VarP(newBytesBase64Value(value, p), name, "", usage)
***REMOVED***

// BytesBase64VarP is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash.
func BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) ***REMOVED***
	CommandLine.VarP(newBytesBase64Value(value, p), name, shorthand, usage)
***REMOVED***

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (f *FlagSet) BytesBase64(name string, value []byte, usage string) *[]byte ***REMOVED***
	p := new([]byte)
	f.BytesBase64VarP(p, name, "", value, usage)
	return p
***REMOVED***

// BytesBase64P is like BytesBase64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte ***REMOVED***
	p := new([]byte)
	f.BytesBase64VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesBase64(name string, value []byte, usage string) *[]byte ***REMOVED***
	return CommandLine.BytesBase64P(name, "", value, usage)
***REMOVED***

// BytesBase64P is like BytesBase64, but accepts a shorthand letter that can be used after a single dash.
func BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte ***REMOVED***
	return CommandLine.BytesBase64P(name, shorthand, value, usage)
***REMOVED***
