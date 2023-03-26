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

type GatewayEventName string

// Gateway Event Names - https://discord.com/developers/docs/topics/gateway-events
const (
	GatewayEventNameHello                               GatewayEventName = "HELLO"
	GatewayEventNameReady                               GatewayEventName = "READY"
	GatewayEventNameResumed                             GatewayEventName = "RESUMED"
	GatewayEventNameReconnect                           GatewayEventName = "RECONNECT"
	GatewayEventNameInvalidSession                      GatewayEventName = "INVALID_SESSION"
	GatewayEventNameApplicationCommandPermissionsUpdate GatewayEventName = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
	GatewayEventNameAutoModerationRuleCreate            GatewayEventName = "AUTO_MODERATION_RULE_CREATE"
	GatewayEventNameAutoModerationRuleUpdate            GatewayEventName = "AUTO_MODERATION_RULE_UPDATE"
	GatewayEventNameAutoModerationRuleDelete            GatewayEventName = "AUTO_MODERATION_RULE_DELETE"
	GatewayEventNameAutoModerationActionExecution       GatewayEventName = "AUTO_MODERATION_ACTION_EXECUTION"
	GatewayEventNameChannelCreate                       GatewayEventName = "CHANNEL_CREATE"
	GatewayEventNameChannelUpdate                       GatewayEventName = "CHANNEL_UPDATE"
	GatewayEventNameChannelDelete                       GatewayEventName = "CHANNEL_DELETE"
	GatewayEventNameChannelPinsUpdate                   GatewayEventName = "CHANNEL_PINS_UPDATE"
	GatewayEventNameThreadCreate                        GatewayEventName = "THREAD_CREATE"
	GatewayEventNameThreadUpdate                        GatewayEventName = "THREAD_UPDATE"
	GatewayEventNameThreadDelete                        GatewayEventName = "THREAD_DELETE"
	GatewayEventNameThreadListSync                      GatewayEventName = "THREAD_LIST_SYNC"
	GatewayEventNameThreadMemberUpdate                  GatewayEventName = "THREAD_MEMBER_UPDATE"
	GatewayEventNameThreadMembersUpdate                 GatewayEventName = "THREAD_MEMBERS_UPDATE"
	GatewayEventNameGuildCreate                         GatewayEventName = "GUILD_CREATE"
	GatewayEventNameGuildUpdate                         GatewayEventName = "GUILD_UPDATE"
	GatewayEventNameGuildDelete                         GatewayEventName = "GUILD_DELETE"
	GatewayEventNameGuildAuditLogEntryCreate            GatewayEventName = "GUILD_AUDIT_LOG_ENTRY_CREATE"
	GatewayEventNameGuildBanAdd                         GatewayEventName = "GUILD_BAN_ADD"
	GatewayEventNameGuildBanRemove                      GatewayEventName = "GUILD_BAN_REMOVE"
	GatewayEventNameGuildEmojisUpdate                   GatewayEventName = "GUILD_EMOJIS_UPDATE"
	GatewayEventNameGuildStickersUpdate                 GatewayEventName = "GUILD_STICKERS_UPDATE"
	GatewayEventNameGuildIntegrationsUpdate             GatewayEventName = "GUILD_INTEGRATIONS_UPDATE"
	GatewayEventNameGuildMemberAdd                      GatewayEventName = "GUILD_MEMBER_ADD"
	GatewayEventNameGuildMemberRemove                   GatewayEventName = "GUILD_MEMBER_REMOVE"
	GatewayEventNameGuildMemberUpdate                   GatewayEventName = "GUILD_MEMBER_UPDATE"
	GatewayEventNameGuildMembersChunk                   GatewayEventName = "GUILD_MEMBERS_CHUNK"
	GatewayEventNameGuildRoleCreate                     GatewayEventName = "GUILD_ROLE_CREATE"
	GatewayEventNameGuildRoleUpdate                     GatewayEventName = "GUILD_ROLE_UPDATE"
	GatewayEventNameGuildRoleDelete                     GatewayEventName = "GUILD_ROLE_DELETE"
	GatewayEventNameGuildScheduledEventCreate           GatewayEventName = "GUILD_SCHEDULED_EVENT_CREATE"
	GatewayEventNameGuildScheduledEventUpdate           GatewayEventName = "GUILD_SCHEDULED_EVENT_UPDATE"
	GatewayEventNameGuildScheduledEventDelete           GatewayEventName = "GUILD_SCHEDULED_EVENT_DELETE"
	GatewayEventNameGuildScheduledEventUserAdd          GatewayEventName = "GUILD_SCHEDULED_EVENT_USER_ADD"
	GatewayEventNameGuildScheduledEventUserRemove       GatewayEventName = "GUILD_SCHEDULED_EVENT_USER_REMOVE"
	GatewayEventNameIntegrationCreate                   GatewayEventName = "INTEGRATION_CREATE"
	GatewayEventNameIntegrationUpdate                   GatewayEventName = "INTEGRATION_UPDATE"
	GatewayEventNameIntegrationDelete                   GatewayEventName = "INTEGRATION_DELETE"
	GatewayEventNameInteractionCreate                   GatewayEventName = "INTERACTION_CREATE"
	GatewayEventNameInviteCreate                        GatewayEventName = "INVITE_CREATE"
	GatewayEventNameInviteDelete                        GatewayEventName = "INVITE_DELETE"
	GatewayEventNameMessageCreate                       GatewayEventName = "MESSAGE_CREATE"
	GatewayEventNameMessageUpdate                       GatewayEventName = "MESSAGE_UPDATE"
	GatewayEventNameMessageDelete                       GatewayEventName = "MESSAGE_DELETE"
	GatewayEventNameMessageDeleteBulk                   GatewayEventName = "MESSAGE_DELETE_BULK"
	GatewayEventNameMessageReactionAdd                  GatewayEventName = "MESSAGE_REACTION_ADD"
	GatewayEventNameMessageReactionRemove               GatewayEventName = "MESSAGE_REACTION_REMOVE"
	GatewayEventNameMessageReactionRemoveAll            GatewayEventName = "MESSAGE_REACTION_REMOVE_ALL"
	GatewayEventNameMessageReactionRemoveEmoji          GatewayEventName = "MESSAGE_REACTION_REMOVE_EMOJI"
	GatewayEventNamePresenceUpdate                      GatewayEventName = "PRESENCE_UPDATE"
	GatewayEventNameStageInstanceCreate                 GatewayEventName = "STAGE_INSTANCE_CREATE"
	GatewayEventNameStageInstanceDelete                 GatewayEventName = "STAGE_INSTANCE_DELETE"
	GatewayEventNameStageInstanceUpdate                 GatewayEventName = "STAGE_INSTANCE_UPDATE"
	GatewayEventNameTypingStart                         GatewayEventName = "TYPING_START"
	GatewayEventNameUserUpdate                          GatewayEventName = "USER_UPDATE"
	GatewayEventNameVoiceStateUpdate                    GatewayEventName = "VOICE_STATE_UPDATE"
	GatewayEventNameVoiceServerUpdate                   GatewayEventName = "VOICE_SERVER_UPDATE"
	GatewayEventNameWebhooksUpdate                      GatewayEventName = "WEBHOOKS_UPDATE"
)
