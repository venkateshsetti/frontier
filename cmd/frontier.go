package main

import (
	"fmt"
	"frontier/api"
	"frontier/configs"
	"frontier/internal/app/blockQuery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var LOGGER *zap.SugaredLogger

// ControllerRegistry mapping the http method string to its routed controller
var ControllerRegistry = map[string]api.HTTPController{}

/*This is the  run function which gets called by Execute function , Basically it initiates the every controller instance
and manager instance through out the project and assigning the required parameters by reading from config

 */
func run(cmd *cobra.Command, args []string) {
	//Logger initialization
	setUpLogger()
	LOGGER.Info("Booting Frontier")
	LOGGER.Infof("HTTP_HOST : %v:%v", configs.HTTP_HOST, configs.HTTP_PORT)

	//intializing the BlockQuery Manager and controller instances and adding the block query controller to the controller registry
	blockQueryMgr := blockQuery.NewManager(LOGGER)
	ControllerRegistry["block_query"] = blockQuery.NewBlockQueryController(LOGGER,blockQueryMgr)

	//initializing the http server
	initHTTP()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	LOGGER.Info("Press Ctrl+C to exit")
	<-sigs
}

//cobra package Command  model to config the project details and main run method
var rootCmd = &cobra.Command{
	Use:   "frontier",
	Short: "frontier ",
	Run:   run, // function from oms.go
}

// Execute executes the command config
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

//Main function
func main(){
  Execute()
}

//initializing the host and port of api NewFrontierWebServer method and staring the webserver start function
func initHTTP() {
	webServer := api.NewFrontierWebServer(LOGGER, configs.HTTP_HOST, configs.HTTP_PORT)
	for key, value := range ControllerRegistry {
		LOGGER.Infof("Initializing routes for %v", key)
		webServer.SetRoute(value)
	}

	go func() {
		err := webServer.Start()
		if err != nil {
			LOGGER.Errorf("Error in web server %v", err)
		}
		LOGGER.Infof("Web server closed")
	}()
}


//init function lets us to read the config variables before the main function get called
//assigning the config path  where yaml file is located
func init() {
	cobra.OnInitialize(initConfig)
	viper.SetConfigName("frontier")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf(
			"No config file found %s,"+
				" Default values from environment variables. \n",
			err,
		)
	}
	err = viper.SafeWriteConfig()
	if err != nil {
		_ = fmt.Errorf("failed to write config %v", err)
	}
}
//assigning the config variables
func initConfig() {
	configs.LOG_LEVEL = viper.GetString("LOG_LEVEL")
	configs.HTTP_HOST = viper.GetString("HTTP_HOST")
	configs.HTTP_PORT = viper.GetInt("HTTP_PORT")

}

//setting up the logger with frontier repo and all internal packages and its file path
//based upon log type such as INFO,DEBUG etc
func setUpLogger() {
	fileName := "frontier" + time.Now().Format("20060102150405") + ".log"
	writer := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "/tmp/" + fileName,
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		},
	)
	opt := zap.ErrorOutput(writer)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(getLogLevel(configs.LOG_LEVEL)),
		OutputPaths:      []string{"stderr", "/tmp/" + fileName},
		ErrorOutputPaths: []string{"stderr", "/tmp/" + fileName},
		EncoderConfig:    encoderConfig,
	}
	_logger, _ := cfg.Build(opt)
	LOGGER = _logger.Sugar()
}

//setting the different log types
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	case "WARN", "warn":
		return zapcore.WarnLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}