package cat

type Config struct {
	CatServerVersion string
	domain           string
	hostname         string
	env              string
	ip               string
	ipHex            string

	httpServerPort      int
	httpServerAddresses []serverAddress

	serverAddress []serverAddress
}
