// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sse

import (
	"bytes"
	"io"
	"io/ioutil"
)

type decoder struct ***REMOVED***
	events []Event
***REMOVED***

func Decode(r io.Reader) ([]Event, error) ***REMOVED***
	var dec decoder
	return dec.decode(r)
***REMOVED***

func (d *decoder) dispatchEvent(event Event, data string) ***REMOVED***
	dataLength := len(data)
	if dataLength > 0 ***REMOVED***
		//If the data buffer's last character is a U+000A LINE FEED (LF) character, then remove the last character from the data buffer.
		data = data[:dataLength-1]
		dataLength--
	***REMOVED***
	if dataLength == 0 && event.Event == "" ***REMOVED***
		return
	***REMOVED***
	if event.Event == "" ***REMOVED***
		event.Event = "message"
	***REMOVED***
	event.Data = data
	d.events = append(d.events, event)
***REMOVED***

func (d *decoder) decode(r io.Reader) ([]Event, error) ***REMOVED***
	buf, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var currentEvent Event
	var dataBuffer *bytes.Buffer = new(bytes.Buffer)
	// TODO (and unit tests)
	// Lines must be separated by either a U+000D CARRIAGE RETURN U+000A LINE FEED (CRLF) character pair,
	// a single U+000A LINE FEED (LF) character,
	// or a single U+000D CARRIAGE RETURN (CR) character.
	lines := bytes.Split(buf, []byte***REMOVED***'\n'***REMOVED***)
	for _, line := range lines ***REMOVED***
		if len(line) == 0 ***REMOVED***
			// If the line is empty (a blank line). Dispatch the event.
			d.dispatchEvent(currentEvent, dataBuffer.String())

			// reset current event and data buffer
			currentEvent = Event***REMOVED******REMOVED***
			dataBuffer.Reset()
			continue
		***REMOVED***
		if line[0] == byte(':') ***REMOVED***
			// If the line starts with a U+003A COLON character (:), ignore the line.
			continue
		***REMOVED***

		var field, value []byte
		colonIndex := bytes.IndexRune(line, ':')
		if colonIndex != -1 ***REMOVED***
			// If the line contains a U+003A COLON character character (:)
			// Collect the characters on the line before the first U+003A COLON character (:),
			// and let field be that string.
			field = line[:colonIndex]
			// Collect the characters on the line after the first U+003A COLON character (:),
			// and let value be that string.
			value = line[colonIndex+1:]
			// If value starts with a single U+0020 SPACE character, remove it from value.
			if len(value) > 0 && value[0] == ' ' ***REMOVED***
				value = value[1:]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Otherwise, the string is not empty but does not contain a U+003A COLON character character (:)
			// Use the whole line as the field name, and the empty string as the field value.
			field = line
			value = []byte***REMOVED******REMOVED***
		***REMOVED***
		// The steps to process the field given a field name and a field value depend on the field name,
		// as given in the following list. Field names must be compared literally,
		// with no case folding performed.
		switch string(field) ***REMOVED***
		case "event":
			// Set the event name buffer to field value.
			currentEvent.Event = string(value)
		case "id":
			// Set the event stream's last event ID to the field value.
			currentEvent.Id = string(value)
		case "retry":
			// If the field value consists of only characters in the range U+0030 DIGIT ZERO (0) to U+0039 DIGIT NINE (9),
			// then interpret the field value as an integer in base ten, and set the event stream's reconnection time to that integer.
			// Otherwise, ignore the field.
			currentEvent.Id = string(value)
		case "data":
			// Append the field value to the data buffer,
			dataBuffer.Write(value)
			// then append a single U+000A LINE FEED (LF) character to the data buffer.
			dataBuffer.WriteString("\n")
		default:
			//Otherwise. The field is ignored.
			continue
		***REMOVED***
	***REMOVED***
	// Once the end of the file is reached, the user agent must dispatch the event one final time.
	d.dispatchEvent(currentEvent, dataBuffer.String())

	return d.events, nil
***REMOVED***
