/*
Copyright 2022 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package sessionlogger

import "os"
import "io"
import "io/ioutil"

type logLevel int

// Logger levels for use with the config builder functions.
const (
	Info = logLevel(iota)
	Warn
	Err
)

// LoggerConfig contains the configuration for the current logger. You can either fill it out manually,
// or use the provided helper functions. The zero value is a valid config that writes log messages to
// Stdout and Stderr and has all log levels enabled.
//
// Note that changes to the config will not effect loggers created before the changes were made.
type LoggerConfig struct {
	Disabled [3]bool      // Info, Warn, Err
	Writers  [3]io.Writer // If nil, use the default for this level.
}

// Disable is a convenience method that makes a specific log level as disabled. Will panic if the level is invalid.
func (lc *LoggerConfig) Disable(l logLevel) *LoggerConfig {
	if l < 0 || l > 3 {
		panic("Log level out of range. Use the constants dumdum.")
	}

	lc.Disabled[l] = true
	return lc
}

// Writer is a convenience method that combines all the given writers and uses them as the output for the
// given log level.
func (lc *LoggerConfig) Writer(l logLevel, w ...io.Writer) *LoggerConfig {
	if l < 0 || l > 3 {
		panic("Log level out of range. Use the constants dumdum.")
	}

	lc.Writers[l] = io.MultiWriter(w...)
	return lc
}

var defaultWriters = []io.Writer{
	os.Stdout,
	os.Stdout,
	os.Stderr,
}

// GetWriter gets a writer for the given log level. No matter what, a valid writer will be
// returned (assuming no invalid logger was manually set in the config).
func (lc *LoggerConfig) GetWriter(l logLevel) io.Writer {
	if l < 0 || l > 3 {
		return os.Stdout
	}
	if lc.Disabled[l] {
		return ioutil.Discard
	}
	if lc.Writers[l] == nil {
		return defaultWriters[l]
	}
	return lc.Writers[l]
}
