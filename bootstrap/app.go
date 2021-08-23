package bootstrap

import (
	"errors"
	"fmt"
	"go.uber.org/zap/zapcore"
	"hpv/bootstrap/lib"
	"os"

	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type App struct {
	ConsoleApp
}

// NewApp _
func NewApp() (app *App) {
	app = new(App)
	app.initCommon()
	if err := app.initConfig(); err != nil {
		zap.L().Panic("app config init fail", zap.Error(err))
	}
	if err := app.initLogger(); err != nil {
		zap.L().Panic("app logger init fail", zap.Error(err))
	}
	return
}

func (app App) initConfig() (err error) {
	// config init
	rootDir := app.Ctx.Get("root_dir").(string)
	file := rootDir + "/config/config.yml"
	stat, err := os.Stat(file)
	if err != nil {
		return
	}
	if stat.IsDir() {
		err = fmt.Errorf("[%s] is not a valid config file", file)
		return
	}
	config := new(lib.YamlConfig)
	err = configor.New(&configor.Config{Debug: true}).Load(config, file)
	app.Ctx.SetAppConfig(config)
	return
}

func (app App) initLogger() (err error) {
	rootDir := app.Ctx.Get("root_dir").(string)
	file := rootDir + "/logs/runtime.log"
	conf := app.Ctx.GetAppConfig()
	if conf == nil {
		err = errors.New("config is nil")
		return
	}

	var zapConf zap.Config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                      // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}
	if conf.Logger.Debug {
		zapConf = zap.NewDevelopmentConfig()
	} else {
		zapConf = zap.NewProductionConfig()
	}
	zapConf.EncoderConfig = encoderConfig

	zapLevel := zap.NewAtomicLevel()
	if err = zapLevel.UnmarshalText([]byte(conf.Logger.Level)); err != nil {
		return
	}
	zapConf.Level = zapLevel
	zapConf.Encoding = "json"
	zapConf.OutputPaths = []string{file}
	zapConf.ErrorOutputPaths = []string{file}

	logger, _ := zapConf.Build()

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return
}
