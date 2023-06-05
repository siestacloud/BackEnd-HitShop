package pkg

import (
	"github.com/sirupsen/logrus"
)

var (
	transport  = "transport"
	service    = "service"
	repository = "repository"
)

func InfoPrint(layer, status interface{}, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Info(message...)
}

func WarnPrint(layer, status interface{}, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": layer, "status": status}).Warn(message...)
}

func ErrPrintT(uri string, status interface{}, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": transport, "uri": uri, "status": status}).Error(message...)
}
func ErrPrintR(status interface{}, message ...interface{}) {
	logrus.WithFields(logrus.Fields{"layer": repository, "status": status}).Error(message...)
}
