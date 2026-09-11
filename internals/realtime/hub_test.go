package realtime

import (
	"testing"
)

func TestHub(t *testing.T) {
	hub := NewHub()

	// Register client 1
	sendChan1 := make(chan Event, 10)
	client1 := &Client{
		UserID: "user-1",
		Role:   "rider",
		Send:   sendChan1,
	}
	hub.Register(client1)

	// Register client 2
	sendChan2 := make(chan Event, 10)
	client2 := &Client{
		UserID: "driver-1",
		Role:   "driver",
		Send:   sendChan2,
	}
	hub.Register(client2)

	if !hub.IsUserOnline("user-1") {
		t.Fatal("expected user-1 to be online")
	}
	if !hub.IsUserOnline("driver-1") {
		t.Fatal("expected driver-1 to be online")
	}

	// Test SendToUser
	testEvent := Event{
		Type: "test.event",
		Data: map[string]interface{}{"msg": "hello"},
	}
	hub.SendToUser("user-1", testEvent)

	select {
	case received := <-sendChan1:
		if received.Type != "test.event" {
			t.Errorf("expected type test.event, got %s", received.Type)
		}
	default:
		t.Fatal("expected to receive event on sendChan1")
	}

	// Test SendToRole
	roleEvent := Event{
		Type: "driver.notification",
		Data: "important alert",
	}
	hub.SendToRole("driver", roleEvent)

	select {
	case received := <-sendChan2:
		if received.Type != "driver.notification" {
			t.Errorf("expected driver.notification, got %s", received.Type)
		}
	default:
		t.Fatal("expected to receive role event on sendChan2")
	}

	// Test Unregister
	hub.Unregister("user-1")
	if hub.IsUserOnline("user-1") {
		t.Fatal("expected user-1 to be offline after unregister")
	}
}
