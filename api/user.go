package api

// User represents a Discord user.
type User struct {
	PublicFlags   BitFlag `mapstructure:"public_flag,omitempty"`
	PremiumType   Flag    `mapstructure:"premium_type,omitempty"`
	Flags         BitFlag `mapstructure:"flag,omitempty"`
	Avatar        string  `mapstructure:"avatar"`
	Bot           bool    `mapstructure:"bot,omitempty"`
	System        bool    `mapstructure:"system,omitempty"`
	MFAEnabled    bool    `mapstructure:"mfa_enabled,omitempty"`
	Banner        string  `mapstructure:"banner,omitempty"`
	AccentColor   int     `mapstructure:"accent_color,omitempty"`
	Locale        string  `mapstructure:"locale,omitempty"`
	Verified      bool    `mapstructure:"verified,omitempty"`
	Email         string  `mapstructure:"email,omitempty"`
	Discriminator string  `mapstructure:"discriminator"`
	Username      string  `mapstructure:"username"`
	ID            string  `mapstructure:"id"`
}

// Me represents the self-bot.
type SelfBot struct {
	Token string
}
