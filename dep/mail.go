package dep

import (
	"github.com/go-mail/mail"
	"github.com/online-bnsp/backend/util/mailer"
	"github.com/spf13/viper"
)

var mailerObj *mailer.Mailer

func (di *DI) GetMailer() *mailer.Mailer {
	if mailerObj == nil {
		d := mail.NewDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.user"), viper.GetString("smtp.pass"))
		mailerObj = mailer.New(d, viper.GetString("mail_sender"))
	}
	return mailerObj
}
