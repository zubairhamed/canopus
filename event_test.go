package canopus

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEvents(t *testing.T) {
	events := NewEvents()

	assert.NotNil(t, events)

	// OnNotify
	evtOnNotifyCalled := false
	events.OnNotify(func(string, interface{}, *Message) {
		evtOnNotifyCalled = true
	})
	events.Notify("/test", "", nil)

	// OnStarted
	evtOnStartedCalled := false
	events.OnStart(func(CoapServer) {
		evtOnStartedCalled = true
	})
	events.Started(nil)

	// OnClosed
	evtOnClosedCalled := false
	events.OnClose(func(CoapServer) {
		evtOnClosedCalled = true
	})
	events.Closed(nil)

	// OnDiscover
	evtOnDiscoverCalled := false
	events.OnDiscover(func() {
		evtOnDiscoverCalled = true
	})
	events.Discover()

	// OnError
	evtOnErrorCalled := false
	events.OnError(func(error) {
		evtOnErrorCalled = true
	})
	events.Error(errors.New("An error occured"))

	// OnObserve
	evtOnObserveCalled := false
	events.OnObserve(func(string, *Message) {
		evtOnObserveCalled = true
	})
	events.Observe("/test", nil)

	// OnObserveCancelled
	evtOnObserveCancelledCalled := false
	events.OnObserveCancel(func(string, *Message) {
		evtOnObserveCancelledCalled = true
	})
	events.ObserveCancelled("/test", nil)

	// OnMessage
	evtOnMessageCalled := false
	events.OnMessage(func(*Message, bool) {
		evtOnMessageCalled = true
	})
	events.Message(nil, true)

	time.Sleep(3000)

	assert.True(t, evtOnNotifyCalled)
	assert.True(t, evtOnStartedCalled)
	assert.True(t, evtOnClosedCalled)
	assert.True(t, evtOnDiscoverCalled)
	assert.True(t, evtOnErrorCalled)
	assert.True(t, evtOnObserveCalled)
	assert.True(t, evtOnObserveCancelledCalled)
	assert.True(t, evtOnMessageCalled)
}
