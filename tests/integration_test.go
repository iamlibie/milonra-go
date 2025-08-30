package integration_test

import (
	"testing"
	"time"

	"github.com/iamlibie/milonra-go/bot"
	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/plugin"
)

// Mock WebSocket connection for testing
type mockConn struct {
	messages []interface{}
}

func (m *mockConn) WriteJSON(v interface{}) error {
	m.messages = append(m.messages, v)
	return nil
}

func (m *mockConn) ReadJSON(v interface{}) error {
	// Mock implementation
	return nil
}

func (m *mockConn) Close() error {
	return nil
}

// Integration test for bot message handling
func TestBotMessageHandling(t *testing.T) {
	// Register a test plugin
	testPlugin := func(bot plugin.Bot, e *event.MessageEvent) string {
		if e.Message == "integration_test" {
			return "integration_response"
		}
		return ""
	}
	plugin.Register("integration_test", testPlugin)

	// Create bot instance
	botInstance := &bot.Bot{
		SelfID: 123456789,
	}

	// Mock message data
	messageData := map[string]interface{}{
		"post_type":    "message",
		"message_type": "group",
		"group_id":     float64(987654321),
		"user_id":      float64(111222333),
		"message":      "integration_test",
		"raw_message":  "integration_test",
		"time":         float64(time.Now().Unix()),
	}

	// Test message handling
	botInstance.HandleMessage(messageData)

	// Note: In a real integration test, you would check if the response
	// was sent correctly. This requires more complex mocking of the
	// WebSocket connection and API calls.
}

// Test plugin chain execution
func TestPluginChainExecution(t *testing.T) {
	executed := make([]string, 0)

	// Register multiple test plugins
	plugin1 := func(bot plugin.Bot, e *event.MessageEvent) string {
		executed = append(executed, "plugin1")
		return ""
	}

	plugin2 := func(bot plugin.Bot, e *event.MessageEvent) string {
		executed = append(executed, "plugin2")
		if e.Message == "test_chain" {
			return "chain_response"
		}
		return ""
	}

	plugin.Register("chain_test_1", plugin1)
	plugin.Register("chain_test_2", plugin2)

	// Create test event
	e := &event.MessageEvent{
		UserID:  123456,
		Message: "test_chain",
	}

	// Execute all plugins
	plugins := plugin.GetPlugins()
	responses := make([]string, 0)

	for _, pluginFunc := range plugins {
		response := pluginFunc(nil, e) // nil bot for this test
		if response != "" {
			responses = append(responses, response)
		}
	}

	// Verify execution
	if len(executed) < 2 {
		t.Errorf("Expected at least 2 plugins to execute, got %d", len(executed))
	}
}

// Benchmark plugin execution
func BenchmarkPluginExecution(b *testing.B) {
	testPlugin := func(bot plugin.Bot, e *event.MessageEvent) string {
		if e.Message == "benchmark" {
			return "response"
		}
		return ""
	}

	e := &event.MessageEvent{
		UserID:  123456,
		Message: "benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPlugin(nil, e)
	}
}
