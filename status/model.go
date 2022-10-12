package status

import "encoding/xml"

type StatusExtension interface {
	GetId() string
	GetDesc() string
	GetProperties() map[string]string
}

type ExtensionDetail struct {
	Id    string `xml:"id,attr"`
	Value string `xml:"value,attr"`
}

type Extension struct {
	Id      string            `xml:"id,attr"`
	Desc    string            `xml:"description"`
	Details []ExtensionDetail `xml:"extensionDetail"`
}

type CustomInfo struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type Status struct {
	XMLName     xml.Name     `xml:"status"`
	Extensions  []Extension  `xml:"extension"`
	CustomInfos []CustomInfo `xml:"customInfo"`
}
