package testing

import (
	"context"
	"testing"
	"time"

	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"goravel/bootstrap"
)

type MailTestSuite struct {
	suite.Suite
}

func TestMailTestSuite(t *testing.T) {
	bootstrap.Boot()
	suite.Run(t, new(MailTestSuite))
}

func (s *MailTestSuite) SetupTest() {

}

func (s *MailTestSuite) TestSendMail() {
	t := s.T()
	assert.Nil(t, facades.Mail.To([]string{facades.Config.Env("MAIL_TO").(string)}).
		Cc([]string{facades.Config.Env("MAIL_CC").(string)}).
		Bcc([]string{facades.Config.Env("MAIL_BCC").(string)}).
		Attach([]string{"./resources/logo.png"}).
		Content(mail.Content{Subject: "Goravel Test", Html: "<h1>Hello Goravel</h1>"}).
		Send())
}

func (s *MailTestSuite) TestSendMailWithFrom() {
	t := s.T()
	assert.Nil(t, facades.Mail.From(mail.From{Address: facades.Config.GetString("mail.from.address"), Name: facades.Config.GetString("mail.from.name")}).
		To([]string{facades.Config.Env("MAIL_TO").(string)}).
		Cc([]string{facades.Config.Env("MAIL_CC").(string)}).
		Bcc([]string{facades.Config.Env("MAIL_BCC").(string)}).
		Attach([]string{"./resources/logo.png"}).
		Content(mail.Content{Subject: "Goravel Test With From", Html: "<h1>Hello Goravel</h1>"}).
		Send())
}

func (s *MailTestSuite) TestQueueMail() {
	t := s.T()
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	go func(ctx context.Context) {
		if err := facades.Queue.Worker(nil).Run(); err != nil {
			facades.Log.Errorf("Queue run error: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	assert.Nil(t, facades.Mail.To([]string{facades.Config.Env("MAIL_TO").(string)}).
		Cc([]string{facades.Config.Env("MAIL_CC").(string)}).
		Bcc([]string{facades.Config.Env("MAIL_BCC").(string)}).
		Attach([]string{"./resources/logo.png"}).
		Content(mail.Content{Subject: "Goravel Test Queue", Html: "<h1>Hello Goravel</h1>"}).
		Queue(nil))

	time.Sleep(3 * time.Second)
}
