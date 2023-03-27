package api

import (
	"errors"
	"sync"

	"github.com/mitchellh/mapstructure"
)

// Handlers contains handlers for gateway events.
type Handlers struct {
	OnHello []func(*GatewayEventHello)
	mutex   sync.Mutex // Used to prevents concurrent writes to the handlers.
}

func (handlers *Handlers) Add(event GatewayEventName, function any) error {
	handlers.mutex.Lock()
	defer handlers.mutex.Unlock()

	failed := false

	switch event {
	case GatewayEventNameHello:
		if function, ok := function.(func(*GatewayEventHello)); ok {
			handlers.OnHello = append(handlers.OnHello, function)
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

func createGatewayEvent(payload genericMap, container any) error {
	return mapstructure.Decode(payload, &container)
}

// Hello Structure - https://discord.com/developers/docs/topics/gateway-events#hello-hello-structure
type GatewayEventHello struct {
	HeartbeatInterval float64 `mapstructure:"heartbeat_interval"` // In milliseconds
}

// Ready Event Fields - https://discord.com/developers/docs/topics/gateway-events#ready-ready-event-fields
type GatewayEventReady struct {
	ResumeGatewayURL string  `mapstructure:"resume_gateway_url"`
	SessionID        string  `mapstructure:"session_id"`
	Version          float64 `mapstructure:"v"`
}
