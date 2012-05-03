// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
)

type LogType uint8

const (
	LogInfo LogType = iota
	LogError
)

// LogEntry is a single, pending log entry.
type LogEntry struct {
	data  []byte  // Message contents.
	ltype LogType // Type of message: info or error.
}

// Log writes log output to a given writer in a threadsafe fashion
type Log struct {
	queue chan *LogEntry
	info  io.WriteCloser
	error io.WriteCloser
}

// NewLog constructs a new Log instance for the given writers.
func NewLog(out, err io.WriteCloser) *Log {
	l := new(Log)
	l.info = out
	l.error = err
	l.queue = make(chan *LogEntry)
	go l.poll()
	return l
}

func (l *Log) Close() {
	close(l.queue)
	l.info.Close()
	l.error.Close()
}

// Printf writes a formatted informative message into the log.
func (l *Log) Printf(f string, argv ...interface{}) {
	l.queue <- &LogEntry{
		data:  []byte(fmt.Sprintf("%s\n", fmt.Sprintf(f, argv...))),
		ltype: LogInfo,
	}
}

// Errorf writes a formatted error message into the log.
func (l *Log) Errorf(f string, argv ...interface{}) {
	l.queue <- &LogEntry{
		data:  []byte(fmt.Sprintf("%s\n", fmt.Sprintf(f, argv...))),
		ltype: LogError,
	}
}

// Fatalf writes a formatted message into the log and exits the
// application immediately.
func (l *Log) Fatalf(f string, argv ...interface{}) {
	l.Errorf(f, argv...)
	os.Exit(1)
}

// poll checks the queue for pending messages and pushes them
// to the underlying writer.
func (l *Log) poll() {
	for {
		select {
		case entry := <-l.queue:
			if entry == nil {
				return
			}

			switch entry.ltype {
			case LogInfo:
				l.info.Write(entry.data)
			case LogError:
				l.error.Write(entry.data)
			}
		}
	}
}
