package conf

const (
	ErrNotFound = confError("config file not found")
	ErrNotOpen  = confError("yml file not open")
)

// Implementation of excellent Dave Chaney idea about constant errors.
// https://dave.cheney.net/2016/04/07/constant-errors
type confError string

func (e confError) Error() string {
	return string(e)
}
