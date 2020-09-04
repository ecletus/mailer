package mailer

import (
	"fmt"
	"net/mail"

	"github.com/ecletus/render"
	"github.com/moisespsena-go/assetfs"

	"github.com/ecletus/core"
)

// Mailer mailer struct
type Mailer struct {
	*Config
}

// Config mailer config
type Config struct {
	DefaultEmailTemplate *Email
	AssetFS              assetfs.Interface
	Sender               SenderInterface
	From                 *mail.Address
	*render.Render
}

// New initialize mailer
func New(config *Config) *Mailer {
	if config == nil {
		config = &Config{}
	}

	if config.Render == nil {
		config.Render = render.New(nil)
		config.Render.SetAssetFS(config.AssetFS)
	}

	return &Mailer{config}
}

// Send send email
func (mailer Mailer) Send(site *core.Site, email *Email, templates ...Template) (err error) {
	copy := *email
	if copy.From == nil {
		copy.From = mailer.From
	}

	var formatAddr = func(addr ...mail.Address) (result []mail.Address) {
		result = make([]mail.Address, len(addr), len(addr))
		for i, addr := range addr {
			if addr.Name != "" {
				if addr.Name, err = site.TextRender(addr.Name); err != nil {
					return
				}
			}
			if addr.Address, err = site.TextRender(addr.Address); err != nil {
				return
			}

			result[i] = addr
		}
		return
	}

	*copy.From = formatAddr(*copy.From)[0]
	if err != nil {
		return
	}
	copy.TO = formatAddr(copy.TO...)
	if err != nil {
		return
	}

	if mailer.DefaultEmailTemplate != nil {
		copy = *mailer.DefaultEmailTemplate.Merge(&copy)
	}

	if len(templates) == 0 {
		return mailer.Sender.Send(&copy)
	}

	var langs = copy.Lang

	for _, template := range templates {
		Email, err := mailer.Render(template, langs...)
		if err == nil {
			if err := mailer.Sender.Send(Email.Merge(&copy)); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// WithSender set sender now
func (mailer Mailer) WithSender(Sender SenderInterface) Mailer {
	if Sender != nil {
		mailer.Sender = Sender
	}
	return mailer
}

type SiteMailer struct {
	Site *core.Site
}

func (this SiteMailer) Mailer() *Mailer {
	return MustGet(this.Site.Data)
}

func (this SiteMailer) Send(email *Email, templates ...Template) (err error) {
	mailer := this.Mailer()
	if mailer == nil {
		return fmt.Errorf("no mailer for site %q", this.Site.Name())
	}
	return mailer.Send(this.Site, email, templates...)
}
