package config

import (
	"fmt"
	"time"
)

type ClientAPI struct {
	Matrix  *Global  `yaml:"-"`
	Derived *Derived `yaml:"-"` // TODO: Nuke Derived from orbit

	Listen Address `yaml:"listen"`
	Bind   Address `yaml:"bind"`

	// If set, allows registration by anyone who also has the shared
	// secret, even if registration is otherwise disabled.
	RegistrationSharedSecret string `yaml:"registration_shared_secret"`
	// This Home Server's ReCAPTCHA public key.
	RecaptchaPublicKey string `yaml:"recaptcha_public_key"`
	// This Home Server's ReCAPTCHA private key.
	RecaptchaPrivateKey string `yaml:"recaptcha_private_key"`
	// Boolean stating whether catpcha registration is enabled
	// and required
	RecaptchaEnabled bool `yaml:"enable_registration_captcha"`
	// Secret used to bypass the captcha registration entirely
	RecaptchaBypassSecret string `yaml:"captcha_bypass_secret"`
	// HTTP API endpoint used to verify whether the captcha response
	// was successful
	RecaptchaSiteVerifyAPI string `yaml:"recaptcha_siteverify_api"`
	// If set disables new users from registering (except via shared
	// secrets)
	RegistrationDisabled bool `yaml:"registration_disabled"`

	// TURN options
	TURN TURN `yaml:"turn"`
}

func (c *ClientAPI) Defaults() {
	c.Listen = "localhost:7771"
	c.Bind = "localhost:7771"
	c.RegistrationSharedSecret = ""
	c.RecaptchaPublicKey = ""
	c.RecaptchaPrivateKey = ""
	c.RecaptchaEnabled = false
	c.RecaptchaBypassSecret = ""
	c.RecaptchaSiteVerifyAPI = ""
	c.RegistrationDisabled = false
}

func (c *ClientAPI) Verify(configErrs *configErrors) {
	checkNotEmpty(configErrs, "client_api.listen", string(c.Listen))
	checkNotEmpty(configErrs, "client_api.bind", string(c.Bind))
	if c.RecaptchaEnabled {
		checkNotEmpty(configErrs, "client_api.recaptcha_public_key", string(c.RecaptchaPublicKey))
		checkNotEmpty(configErrs, "client_api.recaptcha_private_key", string(c.RecaptchaPrivateKey))
		checkNotEmpty(configErrs, "client_api.recaptcha_siteverify_api", string(c.RecaptchaSiteVerifyAPI))
	}
	c.TURN.Verify(configErrs)
}

type TURN struct {
	// TODO Guest Support
	// Whether or not guests can request TURN credentials
	// AllowGuests bool `yaml:"turn_allow_guests"`
	// How long the authorization should last
	UserLifetime string `yaml:"turn_user_lifetime"`
	// The list of TURN URIs to pass to clients
	URIs []string `yaml:"turn_uris"`

	// Authorization via Shared Secret
	// The shared secret from coturn
	SharedSecret string `yaml:"turn_shared_secret"`

	// Authorization via Static Username & Password
	// Hardcoded Username and Password
	Username string `yaml:"turn_username"`
	Password string `yaml:"turn_password"`
}

func (c *TURN) Verify(configErrs *configErrors) {
	value := c.UserLifetime
	if value != "" {
		if _, err := time.ParseDuration(value); err != nil {
			configErrs.Add(fmt.Sprintf("invalid duration for config key %q: %s", "client_api.turn.turn_user_lifetime", value))
		}
	}
}
