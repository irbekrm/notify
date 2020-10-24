package receiver

type Notifier interface {
	Notify(msg string) error
}
