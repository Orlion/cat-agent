package server

type Request struct {
	Service string
	Method  string
	Body    []byte
}

func readRequest() (req *Request, err error) {
	return
}
