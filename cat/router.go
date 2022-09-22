package cat

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
)

const ( // Declared properties given by the router server.
	propertySample  = "sample"
	propertyRouters = "routers"
	propertyBlock   = "block"
)

type routerConfigXMLProperty struct {
	XMLName xml.Name `xml:"property"`
	Id      string   `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
}

type routerConfigXML struct {
	XMLName    xml.Name                  `xml:"property-config"`
	Properties []routerConfigXMLProperty `xml:"property"`
}

type Router struct {
	catServerVersion string
	env              string
	domain           string
	ip               string
	hostname         string
	routerServers    []string
	routers          []string
	current          string
	sample           float64
}

func (r *Router) updateRouterConfig() {
	var query = url.Values{}
	query.Add("env", r.env)
	query.Add("domain", r.domain)
	query.Add("ip", r.ip)
	query.Add("hostname", r.hostname)
	if r.catServerVersion == CatServerVersionV3 {
		query.Add("op", "xml")
	}

	u := url.URL{
		Scheme:   "http",
		Path:     "/cat/s/router",
		RawQuery: query.Encode(),
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	r.shuffleRouterServers()

	for _, server := range r.routerServers {
		u.Host = server
		log.Infof("getting router config from %s", u.String())

		resp, err := client.Get(u.String())
		if err != nil {
			log.Warnf("Error occurred while getting router config from url %s", u.String())
			continue
		}

		r.parse(resp.Body)
		return
	}

	log.Error("can't get router config from remote server.")
	return
}

func (r *Router) parse(reader io.ReadCloser) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	if r.catServerVersion == CatServerVersionV3 {
		t := new(routerConfigXML)
		if err := xml.Unmarshal(bytes, &t); err != nil {
			log.Warnf("error occurred while parsing router config xml content.\n%s", string(bytes))
		}

		for _, property := range t.Properties {
			switch property.Id {
			case propertySample:
				r.updateSample(property.Value)
			case propertyRouters:
				r.updateRouters(property.Value)
			case propertyBlock:
				r.updateBlock(property.Value)
			}
		}
	} else {
		r.updateRouters(string(bytes))
	}
}

func (r *Router) updateSample(v string) {
	sample, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Warnf("Sample should be a valid float, %s given", v)
	} else if math.Abs(sample-atomicx.LoadFloat64(&r.sample)) > 1e-9 {
		atomicx.StoreFloat64(&r.sample, sample)
		log.Infof("Sample rate has been set to %f%%", atomicx.LoadFloat64(&r.sample)*100)
	}
}

func (r *Router) updateRouters(router string) {
	newRouters := resolveServerAddresses(router)

	oldLen, newLen := len(r.routers), len(newRouters)

	if newLen == 0 {
		return
	} else if oldLen == 0 {
		log.Infof("routers has been initialized to: %s", newRouters)
		r.routers = newRouters
	} else if oldLen != newLen {
		log.Infof("routers has been changed to: %s", newRouters)
		r.routers = newRouters
	} else {
		for i := 0; i < oldLen; i++ {
			if r.routers[i] != newRouters[i] {
				log.Infof("routers has been changed to: %s", newRouters)
				r.routers = newRouters
				break
			}
		}
	}

	if len(newRouters) > 0 {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(len(newRouters))
		server := newRouters[randNum]

		if r.current == server {
			return
		}

		r.current = server
		log.Infof("Connected to %s.", server)
		return
	}

	log.Info("cannot established a connection to cat server.")
}

func (r *Router) updateBlock(v string) {
	if v == "false" {
		// enable()
	} else {
		// disable()
	}
}

func (r *Router) shuffleRouterServers() {
	rand.Seed(time.Now().UnixNano())
	length := len(r.routerServers)
	for i := 0; i < length; i++ {
		index := rand.Intn(length - i)
		r.routerServers[i], r.routerServers[index+i] = r.routerServers[index+i], r.routerServers[i]
	}
}
