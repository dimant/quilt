package mocks

import "github.com/stretchr/testify/mock"

// IPTables is an autogenerated mock type for the IPTables type
type IPTables struct {
	mock.Mock
}

// Append provides a mock function with given fields: _a0, _a1, _a2
func (_m *IPTables) Append(_a0 string, _a1 string, _a2 ...string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...string) error); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppendUnique provides a mock function with given fields: _a0, _a1, _a2
func (_m *IPTables) AppendUnique(_a0 string, _a1 string, _a2 ...string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...string) error); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *IPTables) Delete(_a0 string, _a1 string, _a2 ...string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...string) error); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: _a0, _a1
func (_m *IPTables) List(_a0 string, _a1 string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}