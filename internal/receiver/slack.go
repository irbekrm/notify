package receiver

type slack struct {
	webhookUrl    string
	messageHeader string
}

type slackOption func(*slack)

func (s *slack) applyOptions(opts ...slackOption) {
	for _, o := range opts {
		o(s)
	}
}

func MessageHeader(msg string) slackOption {
	return func(s *slack) {
		s.messageHeader = msg
	}
}

func NewSlackReceiver(webhookUrl string, opts ...slackOption) (Notifier, error) {
	s := &slack{}
	s.messageHeader = "New GitHub issue"
	s.webhookUrl = webhookUrl
	s.applyOptions(opts...)
	return s, nil
}

func (s *slack) Notify(msg string) error {
	return nil
}
