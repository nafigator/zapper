package conf

const (
	ErrNotFound = Error("config file not found")
	ErrNotOpen  = Error("yml file not open")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
