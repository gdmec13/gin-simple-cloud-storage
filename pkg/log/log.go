package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	oplogging "github.com/op/go-logging"
	"io"
	"os"
	"simple-cloud-storage/app/config"
	"simple-cloud-storage/app/global"
	"simple-cloud-storage/pkg/util"
	"strings"
	"time"
)

const (
	logDir = "log"
	module = "simple-cloud-storage"
)

var (
	defaultFormatter = `%{time:2006/01/02 - 15:04:05.000} %{longfile} %{color:bold}▶ [%{level:.6s}] %{message}%{color:reset}`
)

func InitLogger() {
	c := global.APP_CONFIG.Log
	if c.Prefix == "" {
		_ = fmt.Errorf("logger prefix not found")
	}
	logger := oplogging.MustGetLogger(module)
	var backends []oplogging.Backend
	backends = registerStdout(c, backends)
	backends = registerFile(c, backends)
	oplogging.SetBackend(backends...)
	global.APP_LOG = logger
}

func registerStdout(c config.Log, backends []oplogging.Backend) []oplogging.Backend {
	if c.Stdout != "" {
		level, err := oplogging.LogLevel(c.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		backends = append(backends, createBackend(os.Stdout, c, level))
	}
	return backends
}

func registerFile(c config.Log, backends []oplogging.Backend) []oplogging.Backend {
	if c.File != "" {
		if ok, _ := util.PathExists(logDir); !ok {
			// directory not exist
			fmt.Println("create log directory")
			_ = os.Mkdir(logDir, os.ModePerm)
		}
		fileWriter, err := rotatelogs.New(
			logDir+string(os.PathSeparator)+"%Y-%m-%d-%H-%M.log",
			// 设置文件清理前的最长保存时间
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 设置日志分割的时间，隔多久分割一次
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
			return backends
		}
		level, err := oplogging.LogLevel(c.File)
		if err != nil {
			fmt.Println(err)
		}
		backends = append(backends, createBackend(fileWriter, c, level))
	}
	return backends
}

func createBackend(w io.Writer, c config.Log, level oplogging.Level) oplogging.Backend {
	backend := oplogging.NewLogBackend(w, c.Prefix, 0)
	stdoutWriter := false
	if w == os.Stdout {
		stdoutWriter = true
	}
	format := getLogFormatter(c, stdoutWriter)
	backendLeveled := oplogging.AddModuleLevel(oplogging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, module)
	return backendLeveled
}

func getLogFormatter(c config.Log, stdoutWriter bool) oplogging.Formatter {
	pattern := defaultFormatter
	if !stdoutWriter {
		// Color is only required for console output
		// Other writers don't need %{color} tag
		pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}
	if !c.LogFile {
		// Remove %{logfile} tag
		pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	}
	return oplogging.MustStringFormatter(pattern)
}
