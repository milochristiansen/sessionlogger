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

// This package adds a "session logger", a log system that allows you to create loggers with three log levels, a
// prefix, and a random unique ID that will be used to prefix every message. This makes it easy to tell where the
// messages are coming from, if a request generates multiple log messages, etc.
//
// All the session loggers for a given program share a log file, and also log to stdout and stderr.
//
// This is not intended to be the One True Logging Solution™ rather I made this to simplify basic logging in
// simple server applications, specifically in REST endpoints for low traffic server apps, chat bots, and other
// endpoint or callback based microservices.
//
// Under the covers, logging is done by the standard library log package. No attempt is made to control log file
// size outside of creating a new one every time the program is started.
package sessionlogger

import "os"
import "io"
import "log"
import "time"

import "github.com/teris-io/shortid"

var (
	infoLog  io.Writer
	warnLog  io.Writer
	errorLog io.Writer
)

var logIDService <-chan string

// MustInitalizeLogger is just InitalizeLogger that panics on error.
func MustInitalizeLogger(logdir string) {
	panic("Logger initialization failed. *shrug* Guess I'll die.\n" + InitalizeLogger(logdir).Error())
}

var loggerInitted = false

// ErrAlreadyInitalized is the error returned by InitalizeLogger if called more than once.
var ErrAlreadyInitalized = errors.New("Logger system was already initialized.")

// InitalizeLogger sets up the logging system (creates files, starts ID service, etc). logdir should
// be a path to the directory you want your log files to be placed in. If this path does not exist it
// will be created. If logdir is set to the empty string, no log file will be created and logging will
// be only to stderr (error level) and stdout (info and warning levels).
func InitalizeLogger(logdir string) error {
	// This may only be set up once.
	if loggerInitted {
		return ErrLoggerAlreadyInitalized
	}

	if logdir != "" {
		err := os.MkdirAll(logdir, 0775)
		if err != nil {
			return err
		}

		f, err := os.Create(logdir + "/" + time.Now().UTC().Format("m01-d02-t150405") + ".log")
		if err != nil {
			return err
		}

		infoLog = io.MultiWriter(f, os.Stdout)
		warnLog = io.MultiWriter(f, os.Stdout)
		errorLog = io.MultiWriter(f, os.Stderr)
	} else {
		infoLog = os.Stdout
		warnLog = os.Stdout
		errorLog = os.Stderr
	}

	go func() {
		c := make(chan string)
		logIDService = c

		idsource := shortid.MustNew(16, shortid.DefaultABC, uint64(time.Now().UnixNano()))

		for {
			c <- idsource.MustGenerate()
		}
	}()

	loggerInitted = true
	return nil
}

// Logger is a logger instance. Possibly with a prefix and unique instance ID.
type Logger struct {
	// Info, Warning, and Error log levels.
	I, W, E *log.Logger
}

// NewMasterLogger creates a new Logger without prefix or instance ID.
func NewMasterLogger() *Logger {
	return newLogger("")
}

// NewSessionLogger creates a Logger that prefixes messages with the endpoint being logged and a unique
// ID individual to that particular Logger.
func NewSessionLogger(endpoint string) *Logger {
	return newLogger("@" + endpoint + ":" + <-logIDService)
}

func newLogger(prefix string) *Logger {
	return &Logger{
		I: log.New(infoLog, "INFO"+prefix+": ", log.Ldate|log.Ltime|log.Lshortfile),
		W: log.New(warnLog, "WARN"+prefix+": ", log.Ldate|log.Ltime|log.Lshortfile),
		E: log.New(errorLog, " ERR"+prefix+": ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
