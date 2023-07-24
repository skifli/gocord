package api

import "github.com/switchupcb/dasgo/dasgo"

// Me represents the self-bot.
type SelfBot struct {
	AccentColor    float64       `mapstructure:"accent_color"`
	Avatar         string        `mapstructure:"avatar"`
	BannerColor    string        `mapstructure:"banner_color"`
	Bio            string        `mapstructure:"bio"`
	Desktop        bool          `mapstructure:"desktop"`
	Discriminator  string        `mapstructure:"discriminator"`
	DisplayName    string        `mapstructure:"display_name"`
	Email          string        `mapstructure:"email"`
	Flags          dasgo.BitFlag `mapstructure:"flags"`
	GlobalName     string        `mapstructure:"global_name"`
	ID             string        `mapstructure:"id"`
	MFAEnabled     bool          `mapstructure:"mfa_enabled"`
	Mobile         bool          `mapstructure:"mobile"`
	NSFWAllowed    bool          `mapstructure:"nsfw_allowed"`
	Phone          string        `mapstructure:"phone"`
	Premium        bool          `mapstructure:"premium"`
	PremiumType    float64       `mapstructure:"premium_type"`
	PublicFlags    dasgo.BitFlag `mapstructure:"public_flags"`
	PurchasedFlags dasgo.BitFlag `mapstructure:"purchased_flags"`
	Username       string        `mapstructure:"username"`
	Verified       bool          `mapstructure:"verified"`
	Token          string
}
