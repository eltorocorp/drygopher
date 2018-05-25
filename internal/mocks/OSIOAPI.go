// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import os "os"

// OSIOAPI is an autogenerated mock type for the OSIOAPI type
type OSIOAPI struct {
	mock.Mock
}

// LookupEnv provides a mock function with given fields: key
func (_m *OSIOAPI) LookupEnv(key string) (string, bool) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ReadDir provides a mock function with given fields: dirname
func (_m *OSIOAPI) ReadDir(dirname string) ([]os.FileInfo, error) {
	ret := _m.Called(dirname)

	var r0 []os.FileInfo
	if rf, ok := ret.Get(0).(func(string) []os.FileInfo); ok {
		r0 = rf(dirname)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]os.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dirname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadFile provides a mock function with given fields: filename
func (_m *OSIOAPI) ReadFile(filename string) ([]byte, error) {
	ret := _m.Called(filename)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFile provides a mock function with given fields: filename, data, perm
func (_m *OSIOAPI) WriteFile(filename string, data []byte, perm os.FileMode) error {
	ret := _m.Called(filename, data, perm)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, os.FileMode) error); ok {
		r0 = rf(filename, data, perm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
