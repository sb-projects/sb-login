package client

type (
	Client interface {
	}
	impl struct {
	}
)

func NewClient() Client {
	return impl{}
}
