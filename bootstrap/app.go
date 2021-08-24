package bootstrap

import (
	"errors"
	"fmt"
	"go.uber.org/zap/zapcore"
	"hpv/bootstrap/lib"
	"os"

	"github.com/jinzhu/configor"
	"go.uber.org/zap"

	"gopkg.in/natefinch/lumberjack.v2"
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

// initConfig 初始化配置文件
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
	err = configor.Load(config, file)
	if err != nil {
		return
	}
	app.Ctx.SetAppConfig(config)
	return
}

// initLogger 初始化日志
func (app App) initLogger() (err error) {
	rootDir := app.Ctx.Get("root_dir").(string)
	file := rootDir + "/logs/runtime.log"
	conf := app.Ctx.GetAppConfig()
	if conf == nil {
		err = errors.New("config is nil")
		return
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}

	zapConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller_line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志多个输出流
	opts := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    128,
			MaxAge:     7,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		}),
		zapcore.AddSync(os.Stdout),
	}

	// 设置日志级别
	zapLevel := zap.NewAtomicLevel()
	if err = zapLevel.UnmarshalText([]byte(conf.Logger.Level)); err != nil {
		return
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapConfig), zapcore.NewMultiWriteSyncer(opts...), zapLevel)

	// 开启堆栈跟踪和文件及行号
	logger := zap.New(core, zap.AddCaller(), zap.Development())
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return
}
