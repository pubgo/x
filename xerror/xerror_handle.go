package xerror

import (
	"errors"
	"os"
	"reflect"
)

type errF1 = func() (err error)
type errF2 = func(...interface{}) (err error)
type errF3 = func(...interface{}) func(...interface{}) error
type errF4 = func(...interface{}) func(...interface{}) func(...interface{}) error

func _handle(err error) IErr {
	var m IErr
	switch _e := err.(type) {
	case IErr:
		m = _e
	case error:
		m = &_Err{err: _e}
	default:
		m = &_Err{err: ErrUnknownType.Err("type %T", _e)}
	}

	return m
}

// Assert errors
func Assert() {
	ErrHandle(recover(), func(err IErr) {
		err.P()
		os.Exit(1)
	})
}

// Resp errors
func Resp(fn func(err IErr)) {
	ErrHandle(recover(), func(err IErr) {
		err.Caller(_caller.FromFunc(reflect.ValueOf(fn)))
		fn(err)
	})
}

// RespErr errors
func RespErr(err *error) {
	ErrHandle(recover(), func(_err IErr) {
		*err = _err
	})
}

// ErrHandle errors
func ErrHandle(err interface{}, fn ...func(err IErr)) {
	if _isNone(err) {
		return
	}

	if _m := _handle(err); !_isNone(_m) {
		if len(fn) != 0 {
			fn[0](_m)
		}
	}
}
