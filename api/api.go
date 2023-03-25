package api

import (
	"time"

	"github.com/valyala/fasthttp"
)

type genericMap = map[string]any

var (
	clientBuildNumber = mustGetLatestBuild()
	requestClient     = fasthttp.Client{
		ReadBufferSize:                8192,
		ReadTimeout:                   time.Duration(time.Second),
		WriteTimeout:                  time.Duration(time.Second),
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
)

const (
	API_VERSION     = "9"
	BROWSER         = "Firefox"
	BROWSER_VERSION = "111.0"
	CAPABILITIES    = 4093
	DEVICE          = "" // Discord's official client sends an empty string.
	OS              = "Windows"
	OS_VERSION      = "10"
	STATUS          = "offline" // https://discord.com/developers/docs/topics/gateway-events#update-presence-status-types
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/" + BROWSER_VERSION
)

// Gateway OPCodes - https://discord.com/developers/docs/topics/OPCodes-and-status-codes#gateway-gateway-OPCodes
const (
	GatewayOPCodeDispatch            float64 = 0
	GatewayOPCodeHeartbeat           float64 = 1
	GatewayOPCodeIdentify            float64 = 2
	GatewayOPCodePresenceUpdate      float64 = 3
	GatewayOPCodeVoiceStateUpdate    float64 = 4
	GatewayOPCodeResume              float64 = 6
	GatewayOPCodeReconnect           float64 = 7
	GatewayOPCodeRequestGuildMembers float64 = 8
	GatewayOPCodeInvalidGateway      float64 = 9
	GatewayOPCodeHello               float64 = 10
	GatewayOPCodeHeartbeatACK        float64 = 11
)
