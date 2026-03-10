package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/mail"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/mail/template"
)

// Params represents all parameters needed for sending mail
type Params struct {
	Subject     string            `json:"subject"`
	HTML        string            `json:"html"`
	Text        string            `json:"text"`
	FromAddress string            `json:"from_address"`
	FromName    string            `json:"from_name"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	Attachments []string          `json:"attachments"`
	Headers     map[string]string `json:"headers"`
}

type Application struct {
	config   config.Config
	queue    contractsqueue.Queue
	template mail.Template
	params   Params
	clone    int

	view string
	text string
	with map[string]any
}

func NewApplication(config config.Config, queue contractsqueue.Queue) (*Application, error) {
	templateEngine, err := template.Get(config)
	if err != nil {
		return nil, err
	}

	return &Application{
		config:   config,
		queue:    queue,
		template: templateEngine,
	}, nil
}

func (r *Application) Attach(attachments []string) mail.Mail {
	instance := r.instance()
	instance.params.Attachments = attachments

	return instance
}

func (r *Application) Bcc(bcc []string) mail.Mail {
	instance := r.instance()
	instance.params.BCC = bcc

	return instance
}

func (r *Application) Cc(cc []string) mail.Mail {
	instance := r.instance()
	instance.params.CC = cc

	return instance
}

func (r *Application) Content(content mail.Content) mail.Mail {
	instance := r.instance()
	instance.params.HTML = content.Html
	instance.view = content.View
	instance.text = content.Text
	instance.with = content.With

	return instance
}

func (r *Application) From(address mail.Address) mail.Mail {
	instance := r.instance()
	instance.params.FromAddress = address.Address
	instance.params.FromName = address.Name

	return instance
}

func (r *Application) Headers(headers map[string]string) mail.Mail {
	instance := r.instance()
	instance.params.Headers = headers

	return instance
}

func (r *Application) Queue(mailable ...mail.Mailable) error {
	if len(mailable) > 0 {
		r.setUsingMailable(mailable[0])
	}

	if err := r.renderViewTemplate(); err != nil {
		return err
	}

	job := r.queue.Job(NewSendMailJob(r.config), []contractsqueue.Arg{
		{
			Type:  "string",
			Value: r.params.Subject,
		},
		{
			Type:  "string",
			Value: r.params.HTML,
		},
		{
			Type:  "string",
			Value: r.params.Text,
		},
		{
			Type:  "string",
			Value: r.params.FromAddress,
		},
		{
			Type:  "string",
			Value: r.params.FromName,
		},
		{
			Type:  "[]string",
			Value: r.params.To,
		},
		{
			Type:  "[]string",
			Value: r.params.CC,
		},
		{
			Type:  "[]string",
			Value: r.params.BCC,
		},
		{
			Type:  "[]string",
			Value: r.params.Attachments,
		},
		{
			Type:  "[]string",
			Value: convertMapHeadersToSlice(r.params.Headers),
		},
	})

	if len(mailable) > 0 {
		if queue := mailable[0].Queue(); queue != nil {
			if queue.Connection != "" {
				job.OnConnection(queue.Connection)
			}
			if queue.Queue != "" {
				job.OnQueue(queue.Queue)
			}
		}
	}

	return job.Dispatch()
}

func (r *Application) Send(mailable ...mail.Mailable) error {
	if len(mailable) > 0 {
		r.setUsingMailable(mailable[0])
	}

	if err := r.renderViewTemplate(); err != nil {
		return err
	}

	return SendMail(r.config, r.params)
}

func (r *Application) Subject(subject string) mail.Mail {
	instance := r.instance()
	instance.params.Subject = subject

	return instance
}

func (r *Application) To(to []string) mail.Mail {
	instance := r.instance()
	instance.params.To = to

	return instance
}

func (r *Application) instance() *Application {
	if r.clone == 0 {
		return &Application{
			clone:    1,
			config:   r.config,
			queue:    r.queue,
			template: r.template,
		}
	}

	return r
}

func (r *Application) setUsingMailable(mailable mail.Mailable) {
	if content := mailable.Content(); content != nil {
		if content.Html != "" {
			r.params.HTML = content.Html
		}
		r.view = content.View
		r.text = content.Text
		r.with = content.With
	}

	if attachments := mailable.Attachments(); len(attachments) > 0 {
		r.params.Attachments = attachments
	}

	if headers := mailable.Headers(); len(headers) > 0 {
		r.params.Headers = headers
	}

	if envelope := mailable.Envelope(); envelope != nil {
		if envelope.From.Address != "" {
			r.params.FromAddress = envelope.From.Address
			r.params.FromName = envelope.From.Name
		}
		if len(envelope.To) > 0 {
			r.params.To = envelope.To
		}
		if len(envelope.Cc) > 0 {
			r.params.CC = envelope.Cc
		}
		if len(envelope.Bcc) > 0 {
			r.params.BCC = envelope.Bcc
		}
		if envelope.Subject != "" {
			r.params.Subject = envelope.Subject
		}
	}
}

func (r *Application) renderViewTemplate() error {
	if r.view != "" && r.template != nil {
		html, err := r.template.Render(r.view, r.with)
		if err != nil {
			return err
		}
		r.params.HTML = html
	}

	if r.text != "" && r.template != nil {
		text, err := r.template.Render(r.text, r.with)
		if err != nil {
			return err
		}
		r.params.Text = text
	}

	return nil
}

func SendMail(config config.Config, params Params) error {
	e := NewEmail()
	fromAddress, fromName := params.FromAddress, params.FromName
	if fromAddress == "" {
		fromAddress, fromName = config.GetString("mail.from.address"), config.GetString("mail.from.name")
	}

	e.From = fmt.Sprintf("%s <%s>", fromName, fromAddress)
	e.To = params.To
	if len(params.BCC) > 0 {
		e.Bcc = params.BCC
	}
	if len(params.CC) > 0 {
		e.Cc = params.CC
	}
	e.Subject = params.Subject

	if len(params.HTML) > 0 {
		e.HTML = []byte(params.HTML)
	}

	if len(params.Text) > 0 {
		e.Text = []byte(params.Text)
	}

	for _, attach := range params.Attachments {
		if _, err := e.AttachFile(attach); err != nil {
			return err
		}
	}

	for key, val := range params.Headers {
		e.Headers.Add(key, val)
	}

	port := config.GetInt("mail.port")
	switch port {
	case 465:
		return e.SendWithTLS(fmt.Sprintf("%s:%d", config.GetString("mail.host"), config.GetInt("mail.port")),
			LoginAuth(config.GetString("mail.username"), config.GetString("mail.password")),
			&tls.Config{ServerName: config.GetString("mail.host")})
	case 587:
		return e.SendWithStartTLS(fmt.Sprintf("%s:%d", config.GetString("mail.host"), config.GetInt("mail.port")),
			LoginAuth(config.GetString("mail.username"), config.GetString("mail.password")),
			&tls.Config{ServerName: config.GetString("mail.host")})
	default:
		return e.Send(fmt.Sprintf("%s:%d", config.GetString("mail.host"), port),
			LoginAuth(config.GetString("mail.username"), config.GetString("mail.password")))
	}
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(*smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
