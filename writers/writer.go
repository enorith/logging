package writers

type CleanupWriter interface {
	Cleanup() error
}
