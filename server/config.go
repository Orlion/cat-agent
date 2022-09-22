package server

type Config struct {
	Addr               string
	ReadTimeoutMillis  int
	WriteTimeoutMillis int
	IdleTimeoutMillis  int
}
