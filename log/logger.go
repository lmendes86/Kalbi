package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

//Log global log object
var Log = logrus.New()

func init() {
	Log.Level = logrus.FatalLevel
	Log.Out = os.Stdout
}
