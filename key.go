package mailer

import (
	"github.com/ecletus/core"
	"github.com/moisespsena-go/maps"
)

type keyType uint8

const Key keyType = iota

func Set(to *maps.Map, value *Mailer) {
	to.Set(Key, value)
}

func Get(m maps.Map) (value *Mailer, ok bool) {
	if v, ok := m.Get(Key); ok {
		return v.(*Mailer), true
	}
	return
}

func MustGet(m maps.Map) *Mailer {
	v, _ := Get(m)
	return v
}

func FromSite(site *core.Site) *SiteMailer {
	return &SiteMailer{Site: site}
}

type GetSiteMailer func(site *core.Site) *SiteMailer
