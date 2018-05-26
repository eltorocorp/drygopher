// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import pckg "github.com/eltorocorp/drygopher/drygopher/coverage/pckg"

// AnalysisAPI is an autogenerated mock type for the AnalysisAPI type
type AnalysisAPI struct {
	mock.Mock
}

// GetCoverageStatistics provides a mock function with given fields: packages
func (_m *AnalysisAPI) GetCoverageStatistics(packages []string) (pckg.Group, pckg.Group, error) {
	ret := _m.Called(packages)

	var r0 pckg.Group
	if rf, ok := ret.Get(0).(func([]string) pckg.Group); ok {
		r0 = rf(packages)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pckg.Group)
		}
	}

	var r1 pckg.Group
	if rf, ok := ret.Get(1).(func([]string) pckg.Group); ok {
		r1 = rf(packages)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(pckg.Group)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func([]string) error); ok {
		r2 = rf(packages)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}