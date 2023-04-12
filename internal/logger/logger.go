package logger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wazofski/storz/utils"
)

type Logger interface {
	Printf(string, ...interface{})
	Object(string, interface{})
	Fatal(error)
}

type _Logger struct {
	Module string
}

func Factory(module string) Logger {
	return &_Logger{
		Module: module,
	}
}

type _Msg struct {
	Who  *string      `json:"who,omitempty"`
	What *interface{} `json:"what,omitempty"`
	When *string      `json:"when,omitempty"`
}

func jsonify(module string, msg interface{}) string {
	ts := utils.Timestamp()
	_msg := _Msg{
		Who:  &module,
		What: &msg,
		When: &ts,
	}

	data, _ := json.MarshalIndent(_msg, "", " ")

	return string(data)
}

func (l *_Logger) Printf(msg string, params ...interface{}) {
	fmt.Println(jsonify(l.Module, fmt.Sprintf(msg, params...)))
}

func (l *_Logger) Object(title string, obj interface{}) {
	fmt.Println(jsonify(l.Module,
		_Msg{
			Who:  &title,
			What: &obj,
		}))
}

func (l *_Logger) Fatal(msg error) {
	log.Panicf(jsonify(l.Module, msg.Error()))
}
