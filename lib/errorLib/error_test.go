package errorlib

import (
	"errors"
	"testing"
)

func TestExecMultipleCanError_AllSucceed(t *testing.T) {
	called := make([]bool, 3)
	funcs := []CanErrorFuncNoReturn{
		func() error { called[0] = true; return nil },
		func() error { called[1] = true; return nil },
		func() error { called[2] = true; return nil },
	}

	err := ExecMultipleCanError(funcs...)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	for i, c := range called {
		if !c {
			t.Errorf("expected func[%d] to be called", i)
		}
	}
}

func TestExecMultipleCanError_StopsOnFirstError(t *testing.T) {
	sentinel := errors.New("stop here")
	thirdCalled := false

	err := ExecMultipleCanError(
		func() error { return nil },
		func() error { return sentinel },
		func() error { thirdCalled = true; return nil },
	)

	if err != sentinel {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if thirdCalled {
		t.Error("expected third func to not be called after an error")
	}
}

func TestExecMultipleCanError_NoFuncs(t *testing.T) {
	err := ExecMultipleCanError()
	if err != nil {
		t.Errorf("expected nil error for empty input, got %v", err)
	}
}

func TestExecMultipleCanError_FirstFuncErrors(t *testing.T) {
	sentinel := errors.New("first fails")
	secondCalled := false

	err := ExecMultipleCanError(
		func() error { return sentinel },
		func() error { secondCalled = true; return nil },
	)

	if err != sentinel {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if secondCalled {
		t.Error("expected second func to not be called")
	}
}
