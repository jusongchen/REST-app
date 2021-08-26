package postgres

// Error is a custom error type for errors returned by envconfig.
type Error string

// Error implements error.
func (e Error) Error() string {
	return string(e)
}

const (
	//ErrConnPoolFail return this error when creating connection pool failed
	ErrConnPoolFail = Error("creating pgx connection pool failed")
)
