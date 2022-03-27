package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := ""

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logName := "my-log"
	logger := client.Logger(logName).StandardLogger(logging.Info)
	var log *zap.Logger
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logging.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	writeStdout := zapcore.AddSync(CustomWriter{})
	gcpLogger := zapcore.AddSync(logger.Writer())
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(w, writeStdout, gcpLogger),
		zap.InfoLevel)

	log = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	fmt.Println(111)
	log.Info("Hello Zap!")
	log.Warn("Beware of getting Zapped! (Pun)")
	log.Error("I'm out of Zap joke!")
}
