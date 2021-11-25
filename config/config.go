package config

import "github.com/byteintellect/go_commons/config"

type OtpServiceConfig struct {
	BaseConfig config.BaseConfig `yaml:"base_config" json:"base_config"`
	OtpLength  int               `yaml:"otp_length" json:"otp_length"`
}
