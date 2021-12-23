package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetReportCaller(true)
	Log.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			arr := strings.Split(f.File, "Go-GraphQL-Mongo-Server")
			//There can be a parent folder named Go-GraphQL-Mongo-Server, so taking the last element
			fileName := arr[len(arr)-1]
			// Finding next / after Go-GraphQL-Mongo-Server, adding +1 to remove the /
			fileName = fileName[strings.Index(fileName, "/")+1:]
			return funcName, fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	}
}
