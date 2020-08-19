package main

import (
	//	"time"

	"github.com/astaxie/beego/logs"
)

type TGLog struct {
}

func (TGLog *TGLog) Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)

}
func (TGLog *TGLog) Emergency(f interface{}, v ...interface{}) {

	logs.Emergency(f, v...)

}

//Alert(f interface{}, v ...interface{})
func (TGLog *TGLog) Alert(f interface{}, v ...interface{}) {
	logs.Alert(f, v...)

}
func (TGLog *TGLog) Critical(f interface{}, v ...interface{}) {
	logs.Critical(f, v...)

}
func (TGLog *TGLog) Error(f interface{}, v ...interface{}) {
	logs.Error(f, v...)

}
func (TGLog *TGLog) Warning(f interface{}, v ...interface{}) {
	logs.Warning(f, v...)

}
func (TGLog *TGLog) Warn(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)

}
func (TGLog *TGLog) Notice(f interface{}, v ...interface{}) {
	logs.Notice(f, v...)

}
func (TGLog *TGLog) Informational(f interface{}, v ...interface{}) {
	logs.Informational(f, v...)

}
func (TGLog *TGLog) Info(f interface{}, v ...interface{}) {
	logs.Info(f, v...)

}

func (TGLog *TGLog) Trace(f interface{}, v ...interface{}) {
	logs.Trace(f, v...)

}
