// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	hostiface "github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
	mock "github.com/stretchr/testify/mock"
)

// ExecAPI is an autogenerated mock type for the ExecAPI type
type ExecAPI struct {
	mock.Mock
}

// Command provides a mock function with given fields: name, arg
func (_m *ExecAPI) Command(name string, arg ...string) hostiface.CommandAPI {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 hostiface.CommandAPI
	if rf, ok := ret.Get(0).(func(string, ...string) hostiface.CommandAPI); ok {
		r0 = rf(name, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(hostiface.CommandAPI)
		}
	}

	return r0
}
