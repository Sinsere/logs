package logs

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

var errorLogger *zap.SugaredLogger

func init() {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "call",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	infoWriter := getWriter("console_info.log")

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	)

	//AddCallerSkip增加了调用者注释跳过的调用者数量
	//大白话：就是显示调用打印日志的是哪一行的code行数，你可以改为0看看效果，哈哈
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = log.Sugar()
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y-%m-%d.log",
		// WithMaxAge和WithRotationCount二者只能设置一个
		// WithMaxAge设置文件清理前的最长保存时间
		// WithRotationCount设置文件清理前最多保存的个数
		//文件最大寿命30天
		rotatelogs.WithMaxAge(time.Hour*24*30),
		//每隔1天进行文件分割
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
