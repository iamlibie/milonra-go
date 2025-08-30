package api_test

import (
	"encoding/json"
	"testing"

	"github.com/iamlibie/milonra-go/api"
)

func TestNewMessage(t *testing.T) {
	msg := api.NewMessage()
	if msg == nil {
		t.Error("NewMessage returned nil")
	}
}

func TestTextMessage(t *testing.T) {
	msg := api.NewMessage().Text("Hello World")
	segments := msg.Build()

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Type != "text" {
		t.Errorf("Expected type 'text', got %s", segments[0].Type)
	}

	if segments[0].Data["text"] != "Hello World" {
		t.Errorf("Expected 'Hello World', got %s", segments[0].Data["text"])
	}
}

func TestAtMessage(t *testing.T) {
	msg := api.NewMessage().At(123456789)
	segments := msg.Build()

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Type != "at" {
		t.Errorf("Expected type 'at', got %s", segments[0].Type)
	}

	if segments[0].Data["qq"] != "123456789" {
		t.Errorf("Expected '123456789', got %s", segments[0].Data["qq"])
	}
}

func TestChainedMessage(t *testing.T) {
	msg := api.NewMessage().
		Text("Hello ").
		At(123456789).
		Text("!")

	segments := msg.Build()

	if len(segments) != 3 {
		t.Errorf("Expected 3 segments, got %d", len(segments))
	}
}

func TestToCQCode(t *testing.T) {
	msg := api.NewMessage().
		Text("Hello ").
		At(123456789)

	cqCode := msg.ToCQCode()
	expected := "Hello [CQ:at,qq=123456789]"

	if cqCode != expected {
		t.Errorf("Expected %s, got %s", expected, cqCode)
	}
}

func TestBuildJSON(t *testing.T) {
	msg := api.NewMessage().Text("test")
	jsonStr, err := msg.BuildJSON()

	if err != nil {
		t.Errorf("BuildJSON failed: %v", err)
	}

	var segments []api.MessageSegment
	err = json.Unmarshal([]byte(jsonStr), &segments)
	if err != nil {
		t.Errorf("JSON unmarshal failed: %v", err)
	}

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment in JSON, got %d", len(segments))
	}
}

func TestParseCQCode(t *testing.T) {
	cqCode := "Hello [CQ:at,qq=123456789] World"
	msg := api.ParseCQCode(cqCode)
	segments := msg.Build()

	if len(segments) < 2 {
		t.Errorf("Expected at least 2 segments, got %d", len(segments))
	}
}

func TestImageMessage(t *testing.T) {
	msg := api.NewMessage().Image("test.jpg", "http://example.com/test.jpg")
	segments := msg.Build()

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Type != "image" {
		t.Errorf("Expected type 'image', got %s", segments[0].Type)
	}

	if segments[0].Data["file"] != "test.jpg" {
		t.Errorf("Expected 'test.jpg', got %s", segments[0].Data["file"])
	}
}

func TestFaceMessage(t *testing.T) {
	msg := api.NewMessage().Face(1)
	segments := msg.Build()

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Type != "face" {
		t.Errorf("Expected type 'face', got %s", segments[0].Type)
	}
}
