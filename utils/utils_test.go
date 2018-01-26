package utils

import (
	"reflect"
	"testing"
)

func TestGetProgressBar(t *testing.T) {
	pb := GetProgressBar(10)
	expectedType := "*pb.ProgressBar"
	returnType := reflect.TypeOf(pb).String()
	if returnType != expectedType {
		t.Errorf("GetProgressBar returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}

func TestWordWrap(t *testing.T) {
	inputStr := `This is a test string for testing word wrap function`
	wrappedString := WordWrap(inputStr, " ", 5)
	expectedString := `This is a test string
 for testing word wrap function`
	if wrappedString != expectedString {
		t.Errorf("WordWrap returned incorrect wrapped string, got: %s, want: %s", wrappedString, expectedString)
	}
}
