package receiver

//go:generate mockgen -source=notifier.go -destination=../../mocks/mock_notifier.go -package=mocks
type Notifier interface {
	Notify(string) error
}
