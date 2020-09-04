package mailer_plugin

import (
	"context"

	"github.com/ecletus/db"
	"github.com/ecletus/plug"
	"github.com/moisespsena-go/logging"
	"github.com/moisespsena-go/maps"
	path_helpers "github.com/moisespsena-go/path-helpers"
	mailer_ "unapu.com/mailer"
	"unapu.com/mailer/db-report/aorm_report"

	"github.com/ecletus/mailer"
	"github.com/ecletus/mailer/mailer_config"
	"github.com/moisespsena-go/aorm"

	"github.com/moisespsena-go/assetfs/assetfsapi"

	"github.com/ecletus/core"
)

var log = logging.GetOrCreateLogger(path_helpers.GetCalledDir())

type Plugin struct {
	plug.EventDispatcher
	db.DBNames

	SitesRegisterKey,
	MailerKey,
	ConfigDirKey,
	AssetFsKey string
	DbLog    bool
	GetDbLog func(site *core.Site, ctx context.Context) (db *aorm.DB)
}

func (this *Plugin) OnRegister() {
	if this.DbLog {
		if this.GetDbLog == nil {
			this.GetDbLog = func(site *core.Site, ctx context.Context) (db *aorm.DB) {
				return site.GetSystemDB().DB
			}
		}

		db.Events(this).DBOnMigrate(func(e *db.DBEvent) error {
			db := this.GetDbLog(core.GetSiteFromDB(e.DB.DB), context.Background())
			return db.AutoMigrate(&aorm_report.Log{}).Error
		})
	}
}

func (this Plugin) RequireOptions() []string {
	return []string{this.SitesRegisterKey, this.ConfigDirKey, this.AssetFsKey}
}

func (this Plugin) ProvideOptions() []string {
	return []string{this.MailerKey}
}

func (this Plugin) ProvidesOptions(options *plug.Options) {
	options.Set(this.MailerKey, mailer.GetSiteMailer(mailer.FromSite))
}

func (this *Plugin) Init(options *plug.Options) {
	sitesRegister := options.GetInterface(this.SitesRegisterKey).(*core.SitesRegister)
	if sitesRegister == nil {
		panic("sites register is nil")
	}
	sitesRegister.OnAdd(func(site *core.Site) {
		var (
			ok         bool
			cfg        *mailer_config.Config
			fs         = options.GetInterface(this.AssetFsKey).(assetfsapi.Interface).NameSpace("mailer")
			siteConfig = site.Config()
		)
		if cfg, ok = mailer_config.Get(siteConfig); !ok {
			if cfgm, ok := maps.GetMapSI(siteConfig, "mailer"); ok {
				cfg = &mailer_config.Config{}
				if err := cfgm.CopyTo(cfg); err != nil {
					log.Errorf("unmarshall config from %q failed: $v", site.Name(), err)
					return
				}
			} else {
				log.Warningf("site %q does not have mailer config", site.Name())
				return
			}
		}
		mixer, err := cfg.Build()
		if err != nil {
			log.Error("create sender failed:", err.Error())
			return
		}

		Sender := mailer.Fallback(mixer.CreateSender())

		if this.DbLog {
			Sender = mailer.SendhandlerFunc(func(next mailer_.Sender) mailer.SenderInterface {
				return mailer.SenderFunc(func(email *mailer.Email) error {
					reporter := aorm_report.NewReporter(this.GetDbLog(site, email.Context))
					email.Context = mailer.WithReporter(email.Context, reporter, false)
					return next.Send(email)
				})
			}).Handle(Sender)
		}

		Mailer := mailer.New(&mailer.Config{
			From:    cfg.From,
			Sender:  Sender,
			AssetFS: fs,
		})
		mailer.Set(&site.Data, Mailer)
	})
}
