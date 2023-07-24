package api

import (
	"errors"
	"sync"

	"github.com/goccy/go-json"
	"github.com/switchupcb/dasgo/dasgo"
)

// Handlers contains handlers for gateway events.
type Handlers struct {
	OnHello         []func(*dasgo.Hello)
	OnMessageCreate []func(*dasgo.Message)
	OnReady         []func(*dasgo.Ready)
	OnReconnect     []func(*dasgo.Reconnect)
	mutex           sync.Mutex // Used to prevents concurrent writes to the handlers.
}

func (handlers *Handlers) Add(event string, function any) error {
	handlers.mutex.Lock()
	defer handlers.mutex.Unlock()

	failed := false

	switch event {
	case dasgo.FlagGatewayEventNameHello:
		if function, ok := function.(func(*dasgo.Hello)); ok {
			handlers.OnHello = append(handlers.OnHello, function)
		} else {
			failed = true
		}
	case dasgo.FlagGatewayEventNameMessageCreate:
		if function, ok := function.(func(*dasgo.Message)); ok {
			handlers.OnMessageCreate = append(handlers.OnMessageCreate, function)
		} else {
			failed = true
		}
	case dasgo.FlagGatewayEventNameReady:
		if function, ok := function.(func(*dasgo.Ready)); ok {
			handlers.OnReady = append(handlers.OnReady, function)
		} else {
			failed = true
		}
	case dasgo.FlagGatewayEventNameReconnect:
		if function, ok := function.(func(*dasgo.Reconnect)); ok {
			handlers.OnReconnect = append(handlers.OnReconnect, function)
		} else {
			failed = true
		}
	default:
		return errors.New("failed to match event to gateway event")
	}

	if failed {
		return errors.New("function signature was not correct for the specified event")
	}

	return nil
}

func createGatewayEvent(genericMap []byte, container any) error {
	payload := new(dasgo.GatewayPayload)

	if err := json.Unmarshal(genericMap, payload); err != nil {
		return err
	}

	if err := json.Unmarshal(payload.Data, container); err != nil {
		return err
	}

	return nil
}
