package api

type BitFlag uint64          // BitFlag represents an alias for a Discord API Bitwise Flag denoted by 1 << x.
type Flag uint8              // Flag represents an alias for a Discord API Flag ranging from 0 - 255.
type GatewayEventName string // GatewayEventName is used internally in event handler manipulation functions.

// Gateway OPCodes - https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
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

// Gateway Close Event - https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
type GatewayCloseEvent struct {
	Code        int
	Description string
	Explanation string
	Reconnect   bool
}

func (gatewayCloseEvent *GatewayCloseEvent) Error() string {
	return gatewayCloseEvent.Explanation
}

var (
	FlagGatewayCloseEventUnknownError = &GatewayCloseEvent{
		Code:        4000,
		Description: "Unknown error",
		Explanation: "We're not sure what went wrong. Try reconnecting?",
		Reconnect:   true,
	}

	FlagGatewayCloseEventUnknownOpcode = &GatewayCloseEvent{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid Gateway opcode or an invalid payload for an opcode. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventDecodeError = &GatewayCloseEvent{
		Code:        4002,
		Description: "Decode error",
		Explanation: "You sent an invalid payload to Discord. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventNotAuthenticated = &GatewayCloseEvent{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent us a payload prior to identifying.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventAuthenticationFailed = &GatewayCloseEvent{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The account token sent with your identify payload is incorrect.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventAlreadyAuthenticated = &GatewayCloseEvent{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventInvalidSeq = &GatewayCloseEvent{
		Code:        4007,
		Description: "Invalid seq",
		Explanation: "The sequence sent when resuming the session was invalid. Reconnect and start a new session.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventRateLimited = &GatewayCloseEvent{
		Code:        4008,
		Description: "Rate limited",
		Explanation: "Woah nelly! You're sending payloads to us too quickly. Slow it down! You will be disconnected on receiving this.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventSessionTimed = &GatewayCloseEvent{
		Code:        4009,
		Description: "Session timed out",
		Explanation: "Your session timed out. Reconnect and start a new one.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventInvalidShard = &GatewayCloseEvent{
		Code:        4010,
		Description: "Invalid shard",
		Explanation: "You sent us an invalid shard when identifying.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventShardingRequired = &GatewayCloseEvent{
		Code:        4011,
		Description: "Sharding required",
		Explanation: "The session would have handled too many guilds - you are required to shard your connection in order to connect.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventInvalidAPIVersion = &GatewayCloseEvent{
		Code:        4012,
		Description: "Invalid API version",
		Explanation: "You sent an invalid version for the gateway.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventInvalidIntent = &GatewayCloseEvent{
		Code:        4013,
		Description: "Invalid intent(s)",
		Explanation: "You sent an invalid intent for a Gateway Intent. You may have incorrectly calculated the bitwise value.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventDisallowedIntent = &GatewayCloseEvent{
		Code:        4014,
		Description: "Disallowed intent(s)",
		Explanation: "You sent a disallowed intent for a Gateway Intent. You may have tried to specify an intent that you have not enabled or are not approved for.",
		Reconnect:   false,
	}

	GatewayCloseEvents = map[int]*GatewayCloseEvent{
		FlagGatewayCloseEventUnknownError.Code:         FlagGatewayCloseEventUnknownError,
		FlagGatewayCloseEventUnknownOpcode.Code:        FlagGatewayCloseEventUnknownOpcode,
		FlagGatewayCloseEventDecodeError.Code:          FlagGatewayCloseEventDecodeError,
		FlagGatewayCloseEventNotAuthenticated.Code:     FlagGatewayCloseEventNotAuthenticated,
		FlagGatewayCloseEventAuthenticationFailed.Code: FlagGatewayCloseEventAuthenticationFailed,
		FlagGatewayCloseEventAlreadyAuthenticated.Code: FlagGatewayCloseEventAlreadyAuthenticated,
		FlagGatewayCloseEventInvalidSeq.Code:           FlagGatewayCloseEventInvalidSeq,
		FlagGatewayCloseEventRateLimited.Code:          FlagGatewayCloseEventRateLimited,
		FlagGatewayCloseEventSessionTimed.Code:         FlagGatewayCloseEventSessionTimed,
		FlagGatewayCloseEventInvalidShard.Code:         FlagGatewayCloseEventInvalidShard,
		FlagGatewayCloseEventInvalidAPIVersion.Code:    FlagGatewayCloseEventInvalidAPIVersion,
		FlagGatewayCloseEventInvalidIntent.Code:        FlagGatewayCloseEventInvalidIntent,
		FlagGatewayCloseEventDisallowedIntent.Code:     FlagGatewayCloseEventDisallowedIntent,
	}
)

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

// User Flags - https://discord.com/developers/docs/resources/user#user-object-user-flags
const (
	FlagUserNONE                         BitFlag = 0
	FlagUserSTAFF                        BitFlag = 1 << 0
	FlagUserPARTNER                      BitFlag = 1 << 1
	FlagUserHYPESQUAD                    BitFlag = 1 << 2
	FlagUserBUG_HUNTER_LEVEL_1           BitFlag = 1 << 3
	FlagUserHYPESQUAD_ONLINE_HOUSE_ONE   BitFlag = 1 << 6
	FlagUserHYPESQUAD_ONLINE_HOUSE_TWO   BitFlag = 1 << 7
	FlagUserHYPESQUAD_ONLINE_HOUSE_THREE BitFlag = 1 << 8
	FlagUserPREMIUM_EARLY_SUPPORTER      BitFlag = 1 << 9
	FlagUserTEAM_PSEUDO_USER             BitFlag = 1 << 10
	FlagUserBUG_HUNTER_LEVEL_2           BitFlag = 1 << 14
	FlagUserVERIFIED_BOT                 BitFlag = 1 << 16
	FlagUserVERIFIED_DEVELOPER           BitFlag = 1 << 17
	FlagUserCERTIFIED_MODERATOR          BitFlag = 1 << 18
	FlagUserBOT_HTTP_INTERACTIONS        BitFlag = 1 << 19
	FlagUserACTIVE_DEVELOPER             BitFlag = 1 << 22
)
