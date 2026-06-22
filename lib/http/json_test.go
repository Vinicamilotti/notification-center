package http

import (
	"encoding/json"
	"io"
	"testing"
)

type sampleStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestMarshalJsonToBody_ValidStruct(t *testing.T) {
	input := sampleStruct{Name: "test", Value: 42}

	reader, err := MarshalJsonToBody(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reader == nil {
		t.Fatal("expected non-nil reader")
	}

	data, _ := io.ReadAll(reader)
	var result sampleStruct
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if result.Name != "test" || result.Value != 42 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestMarshalJsonToBody_Nil(t *testing.T) {
	reader, err := MarshalJsonToBody(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reader != nil {
		t.Error("expected nil reader for nil input")
	}
}

func TestMarshalJsonToBody_Unmarshalable(t *testing.T) {
	ch := make(chan int)
	_, err := MarshalJsonToBody(ch)
	if err == nil {
		t.Error("expected error for unmarshalable type")
	}
}

func TestMarshalJsonToString_ValidStruct(t *testing.T) {
	input := sampleStruct{Name: "hello", Value: 7}

	result, err := MarshalJsonToString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty string")
	}

	var parsed sampleStruct
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to unmarshal result string: %v", err)
	}
	if parsed.Name != "hello" || parsed.Value != 7 {
		t.Errorf("unexpected parsed result: %+v", parsed)
	}
}

func TestMarshalJsonToString_Nil(t *testing.T) {
	result, err := MarshalJsonToString(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for nil input, got %q", result)
	}
}

func TestMarshalJsonToString_Unmarshalable(t *testing.T) {
	ch := make(chan int)
	_, err := MarshalJsonToString(ch)
	if err == nil {
		t.Error("expected error for unmarshalable type")
	}
}
