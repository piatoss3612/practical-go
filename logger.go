package main

import (
	"sync"

	"go.uber.org/zap"
)

var (
	l        *zap.Logger
	s        *zap.SugaredLogger
	syncOnce sync.Once
)

func SetupLogger(development bool, withArgs ...interface{}) error {
	var err error

	syncOnce.Do(func() {
		err = setupLogger(development, withArgs...)
	})

	return err
}

func setupLogger(development bool, withArgs ...interface{}) error {
	var err error

	if development {
		l, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	} else {
		l, err = zap.NewProduction(zap.AddCallerSkip(1))
	}
	if err != nil {
		return err
	}

	s = l.Sugar()

	s = s.With(withArgs...)

	zap.ReplaceGlobals(l)

	return nil
}

func Sync() error {
	if l == nil {
		return nil
	}

	return l.Sync()
}

func Info(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Infow(msg, args...)
		return
	}
	s.Infow(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Debugw(msg, args...)
		return
	}
	s.Debugw(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Warnw(msg, args...)
		return
	}
	s.Warnw(msg, args...)
}

func Error(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Errorw(msg, args...)
		return
	}
	s.Errorw(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Panicw(msg, args...)
		return
	}
	s.Panicw(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Fatalw(msg, args...)
		return
	}
	s.Fatalw(msg, args...)
}
