package mail

import (
	"github.com/goravel/framework/contracts/mail"
)

func Address(address, name string) mail.Address {
	return mail.Address{
		Address: address,
		Name:    name,
	}
}

func Html(html string) mail.Content {
	return mail.Content{
		Html: html,
	}
}

type QueueMail struct {
	queue *mail.Queue
}

func Queue() *QueueMail {
	return &QueueMail{
		queue: &mail.Queue{},
	}
}

// Attachments attach files to the mail
func (r *QueueMail) Attachments() []string {
	return []string{}
}

// Content set the content of the mail
func (r *QueueMail) Content() *mail.Content {
	return &mail.Content{}
}

// Envelope set the envelope of the mail
func (r *QueueMail) Envelope() *mail.Envelope {
	return &mail.Envelope{}
}

// Headers add custom headers to the mail.
func (r *QueueMail) Headers() map[string]string {
	return map[string]string{}
}

// Queue set the queue of the mail
func (r *QueueMail) Queue() *mail.Queue {
	return r.queue
}

func (r *QueueMail) OnConnection(connection string) *QueueMail {
	r.queue.Connection = connection

	return r
}

func (r *QueueMail) OnQueue(queue string) *QueueMail {
	r.queue.Queue = queue

	return r
}
