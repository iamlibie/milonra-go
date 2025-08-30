package plugin_test

import (
	"testing"

	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/plugin"
)

// Mock Bot for testing
type mockBot struct {
	selfID int64
}

func (m *mockBot) WriteJSON(v interface{}) error {
	return nil
}

func (m *mockBot) GetSelfID() int64 {
	return m.selfID
}

// Test plugin function
func testPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if e.Message == "test" {
		return "test response"
	}
	return ""
}

func TestRegister(t *testing.T) {
	// Clear plugins before test
	plugin.Register("test-plugin", testPlugin)

	plugins := plugin.GetPlugins()
	if _, exists := plugins["test-plugin"]; !exists {
		t.Error("Plugin registration failed")
	}
}

func TestPluginExecution(t *testing.T) {
	bot := &mockBot{selfID: 123456}
	event := &event.MessageEvent{
		UserID:  789,
		Message: "test",
	}

	result := testPlugin(bot, event)
	expected := "test response"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestPluginNoResponse(t *testing.T) {
	bot := &mockBot{selfID: 123456}
	event := &event.MessageEvent{
		UserID:  789,
		Message: "no match",
	}

	result := testPlugin(bot, event)

	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestGetPlugins(t *testing.T) {
	plugins := plugin.GetPlugins()
	if plugins == nil {
		t.Error("GetPlugins returned nil")
	}
}
