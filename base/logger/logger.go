package logger

import (
    "io"
    "os"
    "fmt"
    "time"
    "github.com/op/go-logging"
    "github.com/lestrrat/go-file-rotatelogs"
)


var logger *logging.Logger = nil

var format = logging.MustStringFormatter(
    `%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x} %{message}`,
)

func createBackend(w io.Writer, level logging.Level) logging.Backend {
    backend := logging.NewLogBackend(w, "", 0)
    backendLeveled := logging.AddModuleLevel(logging.NewBackendFormatter(backend, format))
    backendLeveled.SetLevel(level, "");
    return backendLeveled;
}


func New(logLevel string, logPath string, module string) error {
    level, err := logging.LogLevel(logLevel);
    if err != nil {
        level = logging.WARNING;
    }
    if logger == nil {
        logger = logging.MustGetLogger(module);
    }
    consoleBackend :=  createBackend(os.Stdout, level);

    if len(logPath) != 0 {
        _, err = os.Stat(logPath)
        if err != nil && os.IsNotExist(err)  {
            // create log dir
            err = os.Mkdir(logPath, 0777)
        }
        
        if err == nil {
            fileWriter, err := rotatelogs.New(
                            logPath + string(os.PathSeparator)+"logview-%Y%m%d%H%M.log",
                            // generate soft link, point to latest log file
                            //rotatelogs.WithLinkName(logSoftLink),
                            // maximum time to save log files
                            rotatelogs.WithMaxAge(7*24*time.Hour),
                            // time period of log file switching
                            rotatelogs.WithRotationTime(24*time.Hour),
                            //rotatelogs.WithRotationTime(2*time.Minute),
                        )
            if err == nil {
                fileBackend := createBackend(fileWriter, level);
                logging.SetBackend(consoleBackend, fileBackend);
                return nil;
            } else {
                return err;
            }
        } else {
            fmt.Printf("LogPath(%s) Not Exist, err:%v \n", logPath, err);
            return err
        }
    }

    // log to console only.
    logging.SetBackend(consoleBackend)
    return nil
    
}

func Debug(i ...interface{}) {
    logger.Debug(i...)
}

func Debugf(format string, args ...interface{}) {
    logger.Debugf(format, args...)
}

func Info(i ...interface{}) {
    logger.Info(i...)
}

func Infof(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

func Warn(i ...interface{}) {
    logger.Warning(i...)
}

func Warnf(format string, args ...interface{}) {
    logger.Warningf(format, args...)
}

func Error(i ...interface{}) {
    logger.Error(i...)
}

func Errorf(format string, args ...interface{}) {
    logger.Errorf(format, args...)
}

func Fatal(i ...interface{}) {
    logger.Fatal(i...)
}

func Fatalf(format string, args ...interface{}) {
    logger.Fatalf(format, args...)
}

func Panic(i ...interface{}) {
    logger.Panic(i...)
}

func Panicf(format string, args ...interface{}) {
    logger.Panicf(format, args...)
}