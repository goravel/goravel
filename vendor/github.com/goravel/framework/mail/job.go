package mail

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
)

type SendMailJob struct {
	config config.Config
}

func NewSendMailJob(config config.Config) *SendMailJob {
	return &SendMailJob{
		config: config,
	}
}

// Signature The name and signature of the job.
func (r *SendMailJob) Signature() string {
	return "goravel_send_mail_job"
}

// Handle Execute the job.
func (r *SendMailJob) Handle(args ...any) error {
	if len(args) != 10 {
		return fmt.Errorf("expected 10 arguments, got %d", len(args))
	}

	subject, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("SUBJECT should be of type string")
	}

	html, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("HTML BODY should be of type string")
	}

	text, ok := args[2].(string)
	if !ok {
		return fmt.Errorf("TEXT BODY should be of type string")
	}

	fromAddress, ok := args[3].(string)
	if !ok {
		return fmt.Errorf("FROM should be of type string")
	}

	fromName, ok := args[4].(string)
	if !ok {
		return fmt.Errorf("FROM NAME should be of type string")
	}

	to, ok := args[5].([]string)
	if !ok {
		return fmt.Errorf("TO should be of type []string")
	}

	cc, ok := args[6].([]string)
	if !ok {
		return fmt.Errorf("CC should be of type []string")
	}

	bcc, ok := args[7].([]string)
	if !ok {
		return fmt.Errorf("BCC should be of type []string")
	}

	attachments, ok := args[8].([]string)
	if !ok {
		return fmt.Errorf("ATTACHMENTS should be of type []string")
	}

	headerSlice, ok := args[9].([]string)
	if !ok {
		return fmt.Errorf("HEADERS should be of type []string")
	}

	params := Params{
		Subject:     subject,
		HTML:        html,
		Text:        text,
		FromAddress: fromAddress,
		FromName:    fromName,
		To:          to,
		CC:          cc,
		BCC:         bcc,
		Attachments: attachments,
		Headers:     convertSliceHeadersToMap(headerSlice),
	}

	return SendMail(r.config, params)
}
