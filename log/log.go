package log

import (
	"github.com/gin-gonic/gin"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel     Level = zap.FatalLevel // 5
	DebugLevel     Level = zap.DebugLevel // -1
	logTmFmtWithMS       = "2006-01-02 15:04:05.000"
	timeLoc              = "Asia/Shanghai"
)

type Field = zap.Field

func (l *Logger) Debug(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Error(msg, fields...)
}

func (l *Logger) DPanic(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panic(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(c *gin.Context, msg string, fields ...Field) {
	getID(c)
	l.l.Fatal(msg, fields...)
}

// function variables for all field types
// in github.com/uber-go/zap/field.go

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any

	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

// ResetDefault not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

var std = New(os.Stderr, InfoLevel)

var traceID string
var spanID string
var FileName = "test.log"

func getID(c *gin.Context) {
	if oteltrace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
		traceID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
		spanID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()
	}
}

func Default() *Logger {
	return std
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}

	// 自定义日志文件写入格式
	opts := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   FileName,
			MaxSize:    1024,
			MaxBackups: 20,
			MaxAge:     28,
			Compress:   true,
		}),
	}
	opts = append(opts, zapcore.AddSync(os.Stdout))
	syncWriter := zapcore.NewMultiWriteSyncer(opts...)

	// 自定义时间输出格式
	loc, _ := time.LoadLocation(timeLoc)
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.In(loc).Format(logTmFmtWithMS) + "]")
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[TraceID " + traceID + "]")
		enc.AppendString("[SpanID " + spanID + "]")
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	cfg := zapcore.EncoderConfig{
		CallerKey:      "caller_line", // 打印文件名和行数
		LevelKey:       "level_name",
		MessageKey:     "msg",
		TimeKey:        "ts",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,   // 自定义时间格式
		EncodeLevel:    customLevelEncoder,  // 小写编码器
		EncodeCaller:   customCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		syncWriter,
		zap.NewAtomicLevelAt(level),
	)
	logger := &Logger{
		l:     zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)),
		level: level,
	}
	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
