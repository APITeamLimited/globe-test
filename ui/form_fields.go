/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Verify that the fields implement the interface
var _ Field = StringField***REMOVED******REMOVED***
var _ Field = PasswordField***REMOVED******REMOVED***

// StringField is just a simple field for reading cleartext strings
type StringField struct ***REMOVED***
	Key     string
	Label   string
	Default string

	// Length constraints.
	Min, Max int
***REMOVED***

// GetKey returns the field's key
func (f StringField) GetKey() string ***REMOVED***
	return f.Key
***REMOVED***

// GetLabel returns the field's label
func (f StringField) GetLabel() string ***REMOVED***
	return f.Label
***REMOVED***

// GetLabelExtra returns the field's default value
func (f StringField) GetLabelExtra() string ***REMOVED***
	return f.Default
***REMOVED***

// GetContents simply reads a string in cleartext from the supplied reader
// It's compllicated and doesn't use t he bufio utils because we can't read ahead
// of the newline and consume more of the stdin, because we'll mess up the next form field
func (f StringField) GetContents(r io.Reader) (string, error) ***REMOVED***
	result := make([]byte, 0, 20)
	buf := make([]byte, 1)
	for ***REMOVED***
		n, err := io.ReadAtLeast(r, buf, 1)
		if err != nil ***REMOVED***
			return string(result), err
		***REMOVED*** else if n != 1 ***REMOVED***
			// Shouldn't happen, but just in case
			return string(result), errors.New("unexpected input when reading string field")
		***REMOVED*** else if buf[0] == '\n' ***REMOVED***
			return string(result), nil
		***REMOVED***
		result = append(result, buf[0])
	***REMOVED***
***REMOVED***

// Clean trims the spaces in the string and checks for min and max length
func (f StringField) Clean(s string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	s = strings.TrimSpace(s)
	if f.Min != 0 && len(s) < f.Min ***REMOVED***
		return nil, fmt.Errorf("invalid input, min length is %d", f.Min)
	***REMOVED***
	if f.Max != 0 && len(s) > f.Max ***REMOVED***
		return nil, fmt.Errorf("invalid input, max length is %d", f.Max)
	***REMOVED***
	if s == "" ***REMOVED***
		s = f.Default
	***REMOVED***
	return s, nil
***REMOVED***

// PasswordField masks password input
type PasswordField struct ***REMOVED***
	Key   string
	Label string
	Min   int
***REMOVED***

// GetKey returns the field's key
func (f PasswordField) GetKey() string ***REMOVED***
	return f.Key
***REMOVED***

// GetLabel returns the field's label
func (f PasswordField) GetLabel() string ***REMOVED***
	return f.Label
***REMOVED***

// GetLabelExtra doesn't return anything so we don't expose the current password
func (f PasswordField) GetLabelExtra() string ***REMOVED***
	return ""
***REMOVED***

// GetContents simply reads a string in cleartext from the supplied reader
func (f PasswordField) GetContents(r io.Reader) (string, error) ***REMOVED***
	stdin, ok := r.(*os.File)
	if !ok ***REMOVED***
		return "", errors.New("cannot read password from the supplied terminal")
	***REMOVED***
	password, err := term.ReadPassword(int(stdin.Fd()))
	if err != nil ***REMOVED***
		// Possibly running on Cygwin/mintty which doesn't emulate
		// pseudo terminals properly, so fallback to plain text input.
		// Note that passwords will be echoed if this is the case.
		// See https://github.com/mintty/mintty/issues/56
		// A workaround is to use winpty or mintty compiled with
		// Cygwin >=3.1.0 which supports the new ConPTY Windows API.
		bufR := bufio.NewReader(r)
		password, err = bufR.ReadBytes('\n')
	***REMOVED***
	return string(password), err
***REMOVED***

// Clean just checks if the minimum length is exceeded, it doesn't trim the string!
func (f PasswordField) Clean(s string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if f.Min != 0 && len(s) < f.Min ***REMOVED***
		return nil, fmt.Errorf("invalid input, min length is %d", f.Min)
	***REMOVED***
	return s, nil
***REMOVED***
