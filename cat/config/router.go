package config

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Orlion/cat-agent/log"
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

func (c *ConfigService) pullRouters() error {
	var query = url.Values{}
	query.Add("env", c.config.Env)
	query.Add("domain", c.config.Domain)
	query.Add("ip", c.config.Ip)
	query.Add("hostname", c.config.Hostname)
	query.Add("op", "xml")

	u := url.URL{
		Scheme:   "http",
		Path:     "/cat/s/router",
		RawQuery: query.Encode(),
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	c.shuffleRouterServers()

	for _, server := range c.config.Servers {
		u.Host = server
		log.Infof("getting router config from %s", u.String())

		resp, err := client.Get(u.String())
		if err != nil {
			log.Warnf("Error occurred while getting router config from url %s : %s", u.String(), err.Error())
			continue
		}

		c.parseRouterConfig(resp.Body)
		return nil
	}

	return errors.New("can't get router config from remote server")
}

func (c *ConfigService) parseRouterConfig(reader io.ReadCloser) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	t := new(routerConfigXML)
	if err := xml.Unmarshal(bytes, &t); err != nil {
		log.Warnf("error occurred while parsing router config xml content.\n%s", string(bytes))
	}

	for _, property := range t.Properties {
		switch property.Id {
		case propertySample:
			c.updateSample(property.Value)
		case propertyRouters:
			c.updateRouters(property.Value)
		case propertyBlock:
			c.updateBlock(property.Value)
		}
	}
}

func (c *ConfigService) shuffleRouterServers() {
	rand.Seed(time.Now().UnixNano())
	length := len(c.config.Servers)
	for i := 0; i < length; i++ {
		index := rand.Intn(length - i)
		c.config.Servers[i], c.config.Servers[index+i] = c.config.Servers[index+i], c.config.Servers[i]
	}
}

func resolveServerAddresses(router string) (addresses []string) {
	for _, segment := range strings.Split(router, ";") {
		if len(segment) == 0 {
			continue
		}
		fragments := strings.Split(segment, ":")
		if len(fragments) != 2 {
			continue
		}

		addresses = append(addresses, segment)
	}

	return
}
