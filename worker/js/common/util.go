// Package common contains helpers for interacting with the JavaScript runtime.
package common

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

func findLocationError(errorString string) string {
	// Find " at " in the error string
	if !strings.Contains(errorString, " at ") {
		return ""
	}

	// Split on the first occurrence of " at "
	splitString := strings.SplitN(errorString, " at ", 2)
	// If we have two elements, the second is the location
	if len(splitString) != 2 {
		return ""
	}

	// Get second element
	location := splitString[1]

	// Look for at *.js:*
	if !strings.Contains(location, ".js:") {
		return ""
	}

	// Split on the first occurrence of ".js:"
	splitString = strings.SplitN(location, ".js:", 2)

	if len(splitString) != 2 {
		return ""
	}

	// Get first element
	filename := fmt.Sprintf("%s.js", splitString[0])

	// Get second element
	lineNumber := splitString[1]

	return fmt.Sprintf("%s:%s", filename, lineNumber)
}

// Throw a JS error; avoids re-wrapping GoErrors.
func Throw(rt *goja.Runtime, err error) {
	if e, ok := err.(*goja.Exception); ok { //nolint:errorlint // we don't really want to unwrap here
		panic(e)
	}
	panic(rt.ToValue(err))
}

// GetReader tries to return an io.Reader value from an exported goja value.
func GetReader(data interface{}) (io.Reader, error) {
	switch r := data.(type) {
	case string:
		return bytes.NewBufferString(r), nil
	case []byte:
		return bytes.NewBuffer(r), nil
	case io.Reader:
		return r, nil
	case goja.ArrayBuffer:
		return bytes.NewBuffer(r.Bytes()), nil
	default:
		return nil, fmt.Errorf("invalid type %T, it needs to be a string, byte array or an ArrayBuffer", data)
	}
}

// ToBytes tries to return a byte slice from compatible types.
func ToBytes(data interface{}) ([]byte, error) {
	switch dt := data.(type) {
	case []byte:
		return dt, nil
	case string:
		return []byte(dt), nil
	case goja.ArrayBuffer:
		return dt.Bytes(), nil
	default:
		return nil, fmt.Errorf("invalid type %T, expected string, []byte or ArrayBuffer", data)
	}
}

// ToString tries to return a string from compatible types.
func ToString(data interface{}) (string, error) {
	switch dt := data.(type) {
	case []byte:
		return string(dt), nil
	case string:
		return dt, nil
	case goja.ArrayBuffer:
		return string(dt.Bytes()), nil
	default:
		return "", fmt.Errorf("invalid type %T, expected string, []byte or ArrayBuffer", data)
	}
}

// RunWithPanicCatching catches panic and converts into an InterruptError error that should abort a script
func RunWithPanicCatching(logger logrus.FieldLogger, rt *goja.Runtime, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()

	return fn()
}
