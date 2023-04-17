package api

// User represents a Discord user.
type User struct {
	AccentColor   int     `mapstructure:"accent_color"`
	Avatar        string  `mapstructure:"avatar"`
	Banner        string  `mapstructure:"banner"`
	Bot           bool    `mapstructure:"bot"`
	Discriminator string  `mapstructure:"discriminator"`
	Email         string  `mapstructure:"email"`
	Flags         BitFlag `mapstructure:"flag"`
	ID            string  `mapstructure:"id"`
	Locale        string  `mapstructure:"locale"`
	MFAEnabled    bool    `mapstructure:"mfa_enabled"`
	PremiumType   Flag    `mapstructure:"premium_type"`
	PublicFlags   BitFlag `mapstructure:"public_flag"`
	System        bool    `mapstructure:"system"`
	Username      string  `mapstructure:"username"`
	Verified      bool    `mapstructure:"verified"`
}

// Me represents the self-bot.
type SelfBot struct {
	AccentColor    float64 `mapstructure:"accent_color"`
	Avatar         string  `mapstructure:"avatar"`
	BannerColor    string  `mapstructure:"banner_color"`
	Bio            string  `mapstructure:"bio"`
	Desktop        bool    `mapstructure:"desktop"`
	Discriminator  string  `mapstructure:"discriminator"`
	DisplayName    string  `mapstructure:"display_name"`
	Email          string  `mapstructure:"email"`
	Flags          BitFlag `mapstructure:"flags"`
	GlobalName     string  `mapstructure:"global_name"`
	ID             string  `mapstructure:"id"`
	MFAEnabled     bool    `mapstructure:"mfa_enabled"`
	Mobile         bool    `mapstructure:"mobile"`
	NSFWAllowed    bool    `mapstructure:"nsfw_allowed"`
	Phone          string  `mapstructure:"phone"`
	Premium        bool    `mapstructure:"premium"`
	PremiumType    float64 `mapstructure:"premium_type"`
	PublicFlags    BitFlag `mapstructure:"public_flags"`
	PurchasedFlags BitFlag `mapstructure:"purchased_flags"`
	Username       string  `mapstructure:"username"`
	Verified       bool    `mapstructure:"verified"`
	Token          string
}
