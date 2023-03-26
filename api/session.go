package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

var (
	gatewayURL        = "wss://gateway.discord.gg/"
	gatewayParameters = "?encoding=json&v=" + API_VERSION
	headers           = make(http.Header)
)

// User contains a self-bot's data.
type User struct {
	Discriminator string // Discriminator e.g. `0001`.
	ID            string // ID            e.g. `000000000000000000`.
	Locale        string // Locale        e.g. `en-GB`.
	Token         string // Token         e.g. ``.
	Username      string // Username      e.g. `user`.
}

// Gateway represents a Discord WebSocket connection.
type Gateway struct {
	CloseChan         chan struct{}   // CloseChan is used as a signal to stop for the gateway's goroutines.
	Conn              *websocket.Conn // Conn represents a connection to the Discord WebSocket.
	GatewayURL        string          // GatewayURL contains the URL used when resuming after a disconnect.
	HeartbeatInterval time.Duration   // The interval the client should wait between sending heartbeats.
	LastSeq           float64         // LastSeq contains the last sequence number received by the client.
	SessionID         string          // SessionID contains the ID of the gateway.
	User              *User           // User contains data relating to the user.
}

type GatewayClosure struct{}

func (err *GatewayClosure) Error() string {
	return "gateway closed"
}

func (user *User) GetData() error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.Header.Set("Authorization", user.Token)
	req.SetRequestURI("https://discord.com/api/v" + API_VERSION + "/users/@me")

	err := requestClient.Do(req, resp)

	if err != nil {
		return err
	} else if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return fmt.Errorf("invalid token (received status code %d)", resp.StatusCode())
	}

	var jsonMap = make(genericMap)

	if err = json.Unmarshal(resp.Body(), &jsonMap); err != nil {
		return err
	}

	user.Discriminator = jsonMap["discriminator"].(string)
	user.ID = jsonMap["id"].(string)
	user.Locale = jsonMap["locale"].(string)
	user.Username = jsonMap["username"].(string)

	return nil
}

func CreateGateway(user *User) *Gateway {
	return &Gateway{
		CloseChan: make(chan struct{}),
		User:      user,
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

func (gateway *Gateway) readMessage() (genericMap, error) {
	_, message, err := gateway.Conn.ReadMessage()

	if err != nil {
		if ce, ok := err.(*websocket.CloseError); ok {
			switch ce.Code {
			case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
				gateway.LastSeq = 0
				gateway.SessionID = ""

				go gateway.reconnect()
				return nil, &GatewayClosure{}
			}
		} else {
			return nil, err
		}
	}

	payload := make(genericMap)

	if err = json.Unmarshal(message, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (gateway *Gateway) sendMessage(jsonPayload genericMap, reconnect bool) error {
	payload, err := json.Marshal(jsonPayload)

	if err != nil {
		return err
	}

	err = gateway.Conn.WriteMessage(websocket.TextMessage, payload)

	if err != nil {
		if ce, ok := err.(*websocket.CloseError); ok {
			switch ce.Code {
			case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
				if reconnect {
					gateway.LastSeq = 0
					gateway.SessionID = ""

					go gateway.reconnect()
				}

				return &GatewayClosure{}
			}
		} else {
			return err
		}
	}

	return nil
}

func (gateway *Gateway) sendHeartbeat() error {
	var err error

	if gateway.LastSeq == 0 {
		err = gateway.sendMessage(genericMap{"op": GatewayOPCodeHeartbeat, "d": nil}, false)
	} else {
		err = gateway.sendMessage(genericMap{"op": GatewayOPCodeHeartbeat, "d": gateway.LastSeq}, false)
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
				if (err == &GatewayClosure{}) {
					return
				}

				continue
			}
		case <-gateway.CloseChan:
			return
		default:
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (gateway *Gateway) gatewayHello() error {
	payload, err := gateway.readMessage()

	if err != nil {
		if (err == &GatewayClosure{}) { // WebSocket closed.
			return gateway.reconnect()
		}

		return err
	}

	if payload["op"].(float64) != GatewayOPCodeHello {
		return fmt.Errorf("unexpected opcode when parsing hello event (expected %f, got %f)", GatewayOPCodeHello, payload["op"].(float64))
	}

	gateway.HeartbeatInterval = time.Duration(payload["d"].(genericMap)["heartbeat_interval"].(float64))

	go gateway.startHeartbeatSender()

	return nil
}

func (gateway *Gateway) gatewayIdentify() error {
	var err error

	if gateway.canReconnect() {
		err = gateway.sendMessage(genericMap{
			"op": GatewayOPCodeResume,
			"d": genericMap{
				"token":      gateway.User.Token,
				"session_id": gateway.SessionID,
				"seq":        int(gateway.LastSeq)},
		}, false)

		if err != nil {
			return err
		}

		err = gateway.sendMessage(genericMap{
			"op": GatewayOPCodePresenceUpdate,
			"d": genericMap{
				"status":     STATUS,
				"since":      0,
				"activities": []any{},
				"afk":        false},
		}, false)
	} else {
		err = gateway.sendMessage(genericMap{
			"op":       GatewayOPCodeIdentify,
			"compress": false,
			"d": genericMap{
				"token":       gateway.User.Token,
				"capabilties": CAPABILITIES,
				"properties": genericMap{
					"os":                       OS,
					"broswer":                  BROWSER,
					"device":                   DEVICE,
					"system_locale":            gateway.User.Locale,
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
		message, err := gateway.readMessage()

		if err != nil {
			if (err == &GatewayClosure{}) { // WebSocket closed.
				return
			}

			continue
		}

		fmt.Printf("%v\n\n", message)

		op := message["op"].(float64)

		if op == GatewayOPCodeHeartbeat { // Discord is asking for a hearbeat.
			gateway.sendHeartbeat()
		} else if op == GatewayOPCodeReconnect {
			gateway.reconnect()
			return
		} else if op == GatewayOPCodeHeartbeatACK { // Discord is acknowledging that we sent a heartbeat.
			continue
		}

		if message["s"] != nil { // Some payloads, for example the heartbeat ack, don't contain contribute to the sequence ID.
			gateway.LastSeq = message["s"].(float64)
		}
	}
}

func (gateway *Gateway) gatewayReady() error {
	payload, err := gateway.readMessage()

	if err != nil {
		if (err == &GatewayClosure{}) { // WebSocket closed.
			return gateway.reconnect()
		}

		return err
	}

	opcode := payload["op"].(float64)

	if opcode == GatewayOPCodeInvalidGateway { // Invalid session. Re-try the connection.
		<-gateway.CloseChan

		return gateway.reconnect()
	} else if opcode != GatewayOPCodeDispatch {
		return fmt.Errorf("unexpected opcode when parsing ready event (expected %f, got %f)", GatewayOPCodeDispatch, payload["op"].(float64))
	}

	gateway.GatewayURL = payload["d"].(genericMap)["resume_gateway_url"].(string)
	gateway.SessionID = payload["d"].(genericMap)["session_id"].(string)

	return nil
}

func (gateway *Gateway) reconnect() error {
	return gateway.Connect()
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
