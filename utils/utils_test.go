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
	input := "This-is-a-test-string-with-separator-for-testing-the-word-wrap-function"
	inputStr := "This is a test string for testing word wrap function"

	if got := WordWrap(input, '-', 0); got != input {
		t.Errorf("WordWrap did not handle parts: 0 properly, got: %s, want: %s", got, input)
	}

	if got := WordWrap(input, '-', -1); got != input {
		t.Errorf("WordWrap did not handle parts: -1 properly, got: %s, want: %s", got, input)
	}

	if got, want := WordWrap(input, '-', 2), `This-is-a-test-string-with-separator-
for-testing-the-word-wrap-function`; got != want {
		t.Errorf("WordWrap returned incorrect wrapped string, got: %s, want: %s", got, want)
	}

	if got, want := WordWrap(inputStr, ' ', 2), inputStr; got != want {
		t.Errorf("WordWrap returned incorrect output, got: %s, want: %s", got, want)
	}
}
