package pistaged

// Server .
type Server interface {
	Serve() error
	Close()
	Reload() error
	Exit() <-chan struct{}
}
