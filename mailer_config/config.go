package mailer_config

import (
	"github.com/moisespsena-go/getters"

	mailer_ "unapu.com/mailer"
)

type keyType uint8

const Key keyType = iota

type SMTP struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Config = mailer_.ConfigMixed

func Get(getter getters.Getter) (cfg *Config, ok bool) {
	if v, ok := getter.Get(Key); ok {
		return v.(*Config), ok
	}
	return
}
