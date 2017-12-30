/**
Package logging implements a logging infrastructure. It supports different
MessageHandler like StreamHandler for console and FileHandler for file output;
Every MessageHandler can have different MessageFormatter and MessageFilter.
*/
package logging

import (
	"fmt"
	"os"
)

type MessageLevel int

const (
	// For StreamHandler, there are different colors that correspond
	// to different levels of the log. Details as follows:
	_            = iota + 30 // black
	ColorRed                 // red
	ColorGreen               // green
	ColorYellow              // yellow
	ColorBlue                // blue
	ColorMagenta             // magenta
	_                        // cyan
	_                        // white
)

const LevelColorSeqClear = "\033[0m"

// LevelColorFlag, MessageLevel color flag.
var LevelColorFlag = []string{
	DEBUG:    levelColorSeq(ColorBlue, 0),
	INFO:     levelColorSeq(ColorGreen, 0),
	WARNING:  levelColorSeq(ColorYellow, 0),
	ERROR:    levelColorSeq(ColorRed, 1),
	CRITICAL: levelColorSeq(ColorMagenta, 1),
}

// LevelString, MessageLevel string.
var LevelString = map[MessageLevel]string{
	DEBUG:    "DEBUG",
	INFO:     "INFO",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
}

// MessageLevel
const (
	NOTSET   = iota
	DEBUG    = MessageLevel(10 * iota) // DEBUG = 10
	INFO     = MessageLevel(10 * iota) // INFO = 20
	WARNING  = MessageLevel(10 * iota) // WARNING = 30
	ERROR    = MessageLevel(10 * iota) // ERROR = 40
	CRITICAL = MessageLevel(10 * iota) // CRITICAL = 50
)

func levelColorSeq(l MessageLevel, way int) string {
	return fmt.Sprintf("\033[%d;%dm", way, MessageLevel(l))
}

// Logger, define logger entity.
type Logger struct {
	Level         MessageLevel          // continue only message level gte Level
	Filter        MessageFilter         // logger message filter, you can define it as your will.
	Record        *MessageRecord        // message entity, you must not instance it.
	StreamHandler *StreamMessageHandler // StreamMessageHandler
	FileHandler   *FileMessageHandler   // FileMessageHandler
}

// logging.GetDefaultLogger, return a default logger object.
func GetDefaultLogger() *Logger {

	return &Logger{
		Level: DEBUG,
		StreamHandler: &StreamMessageHandler{
			Level: DEBUG,
			Formatter: &MessageFormatter{
				Format:     `{{.Color}}[{{.Time}}] {{.LevelString | printf "%8s"}}  {{.FuncName}} {{.ShortFileName}} {{.Line}} {{.ColorClear}} {{.Message}}`,
				TimeFormat: "2006-01-02 15:04:05",
			},
			Destination: os.Stdout,
		},
	}
}

// Logger.Log, sed message to different handler.
func (l *Logger) log(level MessageLevel, format string, a ...interface{}) {

	if level >= l.Level {

		l.Record = GetMessageRecord(level, format, a...)

		if l.Filter == nil || (l.Filter != nil && l.Filter(l)) {

			if l.StreamHandler != nil && level >= l.StreamHandler.Level {
				if l.StreamHandler.Filter == nil || (l.StreamHandler.Filter != nil && l.StreamHandler.Filter(l)) {
					l.StreamHandler.Write([]byte(l.StreamHandler.Formatter.GetMessage(l)))
				}
			}

			if l.FileHandler != nil && level >= l.FileHandler.Level {
				if l.FileHandler.Filter == nil || (l.FileHandler.Filter != nil && l.FileHandler.Filter(l)) {
					l.FileHandler.Write([]byte(l.FileHandler.Formatter.GetMessage(l)))
				}
			}

		}
	}
}

// Logger.Debug, record DEBUG message.
func (l *Logger) DEBUG(format string, a ...interface{}) {
	l.log(DEBUG, format, a...)
}

// Logger.Debug, record INFO message.
func (l *Logger) INFO(format string, a ...interface{}) {
	l.log(INFO, format, a...)
}

// Logger.Debug, record WARNING message.
func (l *Logger) WARNING(format string, a ...interface{}) {
	l.log(WARNING, format, a...)
}

// Logger.Debug, record ERROR message.
func (l *Logger) ERROR(format string, a ...interface{}) {
	l.log(ERROR, format, a...)
}

// Logger.Debug, record CRITICAL message.
func (l *Logger) CRITICAL(format string, a ...interface{}) {
	l.log(CRITICAL, format, a...)
}
