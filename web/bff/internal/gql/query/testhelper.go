package query

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type mockEncodedValue struct {
	value interface{}
}

func (v *mockEncodedValue) HasValue() bool {
	return v.value != nil
}
func (v *mockEncodedValue) Get(ptr interface{}) error {
	if !v.HasValue() {
		return fmt.Errorf("no value present")
	}
	if reflect.TypeOf(ptr) != reflect.TypeOf(v.value) {
		return fmt.Errorf("wrong type of value. received %T but got %T", ptr, v.value)
	}
	result := reflect.ValueOf(v.value).Elem()
	reflect.ValueOf(ptr).Elem().Set(result)
	return nil
}

type mockAuth struct {
	mock.Mock
}

func (m *mockAuth) SessionID() string {
	a := m.Called()
	return a.String(0)
}
func (m *mockAuth) Token() string {
	a := m.Called()
	return a.String(0)
}
