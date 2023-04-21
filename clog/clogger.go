package clog

import (
	"log"
)

var (
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
	DebugLog   *log.Logger
)
