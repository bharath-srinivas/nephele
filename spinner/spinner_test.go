package spinner

import (
	"reflect"
	"testing"
	"time"
)

var sp = NewSpinner("testing ", 32, 100*time.Millisecond)

func TestDefault(t *testing.T) {
	s := Default("testing ")
	expectedType := "*spinner.Spinner"
	returnType := reflect.TypeOf(s).String()
	if returnType != expectedType {
		t.Errorf("Default returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}

func TestNewSpinner(t *testing.T) {
	s := NewSpinner("", 0, 1*time.Second)
	expectedType := "*spinner.Spinner"
	returnType := reflect.TypeOf(s).String()
	if returnType != expectedType {
		t.Errorf("NewSpinner returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}

func TestStart(t *testing.T) {
	sp.Start()
	time.Sleep(2 * time.Second)

	if !sp.active {
		t.Error("Start did not start the spinner")
	}
}

func TestStop(t *testing.T) {
	sp.Stop()
	time.Sleep(100 * time.Millisecond)

	if sp.active {
		t.Error("Stop did not stop the spinner")
	}
}