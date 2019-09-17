// +build !go1.13

/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package cmd

import (
	"errors"
	"os"
	"runtime"
)

// TODO(cuonglm): remove this when last supported version bump to go1.13
func configDir() (string, error) ***REMOVED***
	var dir string

	switch runtime.GOOS ***REMOVED***
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" ***REMOVED***
			return "", errors.New("%AppData% is not defined")
		***REMOVED***

	case "darwin":
		dir = os.Getenv("HOME")
		if dir == "" ***REMOVED***
			return "", errors.New("$HOME is not defined")
		***REMOVED***
		dir += "/Library/Application Support"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" ***REMOVED***
			return "", errors.New("$home is not defined")
		***REMOVED***
		dir += "/lib"

	default: // Unix
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" ***REMOVED***
			dir = os.Getenv("HOME")
			if dir == "" ***REMOVED***
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			***REMOVED***
			dir += "/.config"
		***REMOVED***
	***REMOVED***

	return dir, nil
***REMOVED***
