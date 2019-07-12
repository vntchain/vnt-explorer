package common

import (
	"github.com/astaxie/beego"
	"fmt"
)

var LOG_LEVEL = map[string]int {
	"ERROR": beego.LevelError,
	"WARN": beego.LevelWarning,
	"INFO": beego.LevelInformational,
	"DEBUG": beego.LevelDebug,
}

func InitLogLevel() {
	level := beego.AppConfig.String("log::level")
	fmt.Printf("Will set log level to: %s=%d \n", level, LOG_LEVEL[level])
	beego.SetLevel(LOG_LEVEL[level])
}
