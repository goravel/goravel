package mail

type Mail interface {
	// Attach attaches files to the Mail.
	Attach(files []string) Mail
	// Bcc adds a "blind carbon copy" address to the Mail.
	Bcc(addresses []string) Mail
	// Cc adds a "carbon copy" address to the Mail.
	Cc(addresses []string) Mail
	// Content set the content of Mail.
	Content(content Content) Mail
	// From set the sender of Mail.
	From(address Address) Mail
	// Headers adds custom headers to the Mail.
	Headers(headers map[string]string) Mail
	// Queue a given Mail
	Queue(mailable ...Mailable) error
	// Send the Mail
	Send(mailable ...Mailable) error
	// Subject set the subject of Mail.
	Subject(subject string) Mail
	// To set the recipients of Mail.
	To(addresses []string) Mail
}

type Mailable interface {
	// Attachments set the attachments of Mailable.
	Attachments() []string
	// Content set the content of Mailable.
	Content() *Content
	// Envelope set the envelope of Mailable.
	Envelope() *Envelope
	// Headers adds custom headers to the Mail.
	Headers() map[string]string
	// Queue set the queue of Mailable.
	Queue() *Queue
}

type Content struct {
	Html string
	Text string
	View string
	With map[string]any
}

type Queue struct {
	Connection string
	Queue      string
}

type Address struct {
	Address string
	Name    string
}

type Envelope struct {
	Bcc     []string
	Cc      []string
	From    Address
	Subject string
	To      []string
}
