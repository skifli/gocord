package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/switchupcb/dasgo/dasgo"
	"github.com/valyala/fasthttp"
)

var (
	gatewayURL        = "wss://gateway.discord.gg/"
	gatewayParameters = "?encoding=json&v=" + API_VERSION
	headers           = make(http.Header)
)

// Gateway represents a Discord WebSocket connection.
type Gateway struct {
	CloseChan         chan struct{}   // CloseChan is used as a signal to stop for the gateway's goroutines.
	Conn              *websocket.Conn // Conn represents a connection to the Discord WebSocket.
	GatewayURL        string          // GatewayURL contains the URL used when resuming after a disconnect.
	Handlers          *Handlers       // Handles for gateway events
	HeartbeatInterval time.Duration   // The interval the client should wait between sending heartbeats.
	LastSeq           float64         // LastSeq contains the last sequence number received by the client.
	SelfBot           *SelfBot        // SelfBot contains data relating to the self-bot.
	SessionID         string          // SessionID contains the ID of the gateway.
}

func CreateGateway(selfBot *SelfBot) *Gateway {
	return &Gateway{
		CloseChan: make(chan struct{}),
		Handlers:  new(Handlers),
		SelfBot:   selfBot,
	}
}

func getGatewayURL() error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("https://discord.com/api/v" + API_VERSION + "/gateway")

	if err := requestClient.Do(req, resp); err != nil {
		return err
	}

	jsonMap := make(map[string]string)

	if err := json.Unmarshal(req.Body(), &jsonMap); err != nil {
		return err
	}

	gatewayURL = jsonMap["url"]

	return nil
}

func (gateway *Gateway) canReconnect() bool {
	return gateway.SessionID != "" && gateway.GatewayURL != "" && gateway.LastSeq != 0
}

func (gateway *Gateway) readMessage() ([]byte, genericMap, error) {
	_, message, err := gateway.Conn.ReadMessage()

	if err != nil {
		closeError := err.(*websocket.CloseError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
			go gateway.reset()
			return nil, nil, err
		default:
			if closeEvent, ok := dasgo.GatewayCloseEventCodes[closeError.Code]; ok {
				if closeEvent.Reconnect { // If the session is reconnectable.
					go gateway.reconnect()
				}

				return nil, nil, fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
			} else {
				return nil, nil, err
			}
		}
	}

	payload := make(genericMap)

	if err = json.Unmarshal(message, &payload); err != nil {
		return nil, nil, err
	}

	return message, payload, nil
}

func (gateway *Gateway) sendMessage(jsonPayload genericMap, reconnect bool) error {
	payload, err := json.Marshal(jsonPayload)

	if err != nil {
		return err
	}

	err = gateway.Conn.WriteMessage(websocket.TextMessage, payload)

	if err != nil {
		closeError := err.(*websocket.CloseError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
			go gateway.reset()
			return err
		default:
			if closeEvent, ok := dasgo.GatewayCloseEventCodes[closeError.Code]; ok {
				if closeEvent.Reconnect { // If the session is reconnectable.
					go gateway.reconnect()
				}

				return fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
			} else {
				return err
			}
		}
	}

	return nil
}

func (gateway *Gateway) sendHeartbeat() error {
	var err error

	if gateway.LastSeq == 0 {
		err = gateway.sendMessage(genericMap{"op": dasgo.FlagGatewayOpcodeHeartbeat, "d": nil}, false)
	} else {
		err = gateway.sendMessage(genericMap{"op": dasgo.FlagGatewayOpcodeHeartbeat, "d": gateway.LastSeq}, false)
	}

	if err != nil {
		return err
	}

	return nil
}

func (gateway *Gateway) startHeartbeatSender() {
	ticker := time.NewTicker(gateway.HeartbeatInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := gateway.sendHeartbeat(); err != nil {
				return
			}
		case <-gateway.CloseChan:
			return
		default:
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (gateway *Gateway) gatewayHello() error {
	payloadBytes, payload, err := gateway.readMessage()

	if err != nil {
		return err
	}

	if payload["op"].(float64) != dasgo.FlagGatewayOpcodeHello {
		return fmt.Errorf("unexpected opcode when parsing hello event (expected %d, got %f)", dasgo.FlagGatewayOpcodeHello, payload["op"].(float64))
	}

	gateway.HeartbeatInterval = time.Duration(payload["d"].(genericMap)["heartbeat_interval"].(float64))

	go gateway.startHeartbeatSender()

	helloEvent := new(dasgo.Hello)

	if err = createGatewayEvent(payloadBytes, helloEvent); err != nil {
		return err
	}

	for _, handler := range gateway.Handlers.OnHello {
		go handler(helloEvent)
	}

	return nil
}

func (gateway *Gateway) gatewayIdentify() error {
	var err error

	if gateway.canReconnect() {
		err = gateway.sendMessage(genericMap{
			"op": dasgo.FlagGatewayOpcodeResume,
			"d": genericMap{
				"token":      gateway.SelfBot.Token,
				"session_id": gateway.SessionID,
				"seq":        int(gateway.LastSeq)},
		}, false)

		if err != nil {
			return err
		}

		err = gateway.sendMessage(genericMap{
			"op": dasgo.FlagGatewayOpcodePresenceUpdate,
			"d": genericMap{
				"status":     STATUS,
				"since":      0,
				"activities": []any{},
				"afk":        false},
		}, false)
	} else {
		err = gateway.sendMessage(genericMap{
			"op":       dasgo.FlagGatewayOpcodeIdentify,
			"compress": false,
			"d": genericMap{
				"token":       gateway.SelfBot.Token,
				"capabilties": CAPABILITIES,
				"properties": genericMap{
					"os":                       OS,
					"broswer":                  BROWSER,
					"device":                   DEVICE,
					"system_locale":            systemLocale,
					"browser_user_agent":       USER_AGENT,
					"browser_version":          BROWSER_VERSION,
					"os_version":               OS_VERSION,
					"referrer":                 "",
					"referring_domain":         "",
					"referrer_current":         "",
					"referring_domain_current": "",
					"release_channel":          "stable",
					"client_build_number":      clientBuildNumber,
					"client_event_source":      nil,
				},
				"presence": genericMap{
					"status":     STATUS,
					"since":      0,
					"activities": []any{},
					"afk":        false,
				},
				"compress": false,
				"client_state": genericMap{
					"guild_versions":              genericMap{},
					"highest_last_message_id":     "0",
					"read_state_version":          0,
					"user_guild_settings_version": -1,
					"user_settings_version":       -1,
					"private_channels_version":    "0",
					"api_code_version":            0},
			},
		}, false)
	}

	if err != nil {
		return err
	}

	return nil
}

func (gateway *Gateway) startMessageHandler() {
	for {
		messageBytes, message, err := gateway.readMessage()

		if err != nil { // Error occurred, assume readMessage handled it.
			return
		}

		op := message["op"].(float64)

		switch op {
		case dasgo.FlagGatewayOpcodeDispatch: // Dispatch event.
			eventName := message["t"].(string)

			switch eventName {
			case string(dasgo.FlagGatewayEventNameMessageCreate):
				messageEvent := new(dasgo.Message)

				check(createGatewayEvent(messageBytes, messageEvent))

				for _, handler := range gateway.Handlers.OnMessageCreate {
					go handler(messageEvent)
				}
			}
		case dasgo.FlagGatewayOpcodeHeartbeat: // Discord is asking for a hearbeat.
			gateway.sendHeartbeat()
		case dasgo.FlagGatewayOpcodeHeartbeatACK: // Discord is acknowledging that we sent a heartbeat.
			continue
		case dasgo.FlagGatewayOpcodeReconnect:
			recconectEvent := new(dasgo.Reconnect)

			for _, handler := range gateway.Handlers.OnReconnect {
				go handler(recconectEvent)
			}

			gateway.reconnect()
			return
		}

		if message["s"] != nil { // Some payloads, for example the heartbeat ack, don't contribute to the sequence ID.
			gateway.LastSeq = message["s"].(float64)
		}
	}
}

func (gateway *Gateway) gatewayReady() error {
	payloadBytes, payload, err := gateway.readMessage()

	if err != nil {
		return err
	}

	opcode := payload["op"].(float64)

	if opcode == dasgo.FlagGatewayOpcodeInvalidSession { // Invalid session. Re-try the connection.
		<-gateway.CloseChan

		return gateway.reconnect()
	} else if opcode != dasgo.FlagGatewayOpcodeDispatch {
		return fmt.Errorf("unexpected opcode when parsing ready event (expected %d, got %f)", dasgo.FlagGatewayOpcodeDispatch, payload["op"].(float64))
	}

	check(mapstructure.Decode(payload["d"].(genericMap)["user"].(genericMap), gateway.SelfBot))

	gateway.GatewayURL = payload["d"].(genericMap)["resume_gateway_url"].(string)
	gateway.SessionID = payload["d"].(genericMap)["session_id"].(string)

	readyEvent := new(dasgo.Ready)

	if err = createGatewayEvent(payloadBytes, readyEvent); err != nil {
		return err
	}

	for _, handler := range gateway.Handlers.OnReady {
		go handler(readyEvent)
	}

	return nil
}

func (gateway *Gateway) reconnect() error {
	return gateway.Connect()
}

func (gateway *Gateway) reset() error {
	gateway.LastSeq = 0
	gateway.SessionID = ""

	return gateway.reconnect()
}

func (gateway *Gateway) Connect() error {
	if gateway.GatewayURL == "" {
		gateway.GatewayURL = gatewayURL + gatewayParameters
	}

	if len(headers) == 0 {
		headers.Set("Host", "gateway.discord.gg")
		headers.Set("User-Agent", USER_AGENT)
	}

	conn, resp, err := websocket.DefaultDialer.Dial(gateway.GatewayURL, headers)

	if resp.StatusCode == 404 { // WebSocket URL was invalid, try getting the latest from the API.
		if err = getGatewayURL(); err != nil {
			return err
		}

		return gateway.Connect()
	} else if err != nil {
		return err
	}

	gateway.Conn = conn

	if err = gateway.gatewayHello(); err != nil {
		return err
	} else if err = gateway.gatewayIdentify(); err != nil {
		return err
	} else if err = gateway.gatewayReady(); err != nil {
		return err
	}

	gateway.startMessageHandler()

	return nil
}
