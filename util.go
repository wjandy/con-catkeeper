package main

import (
	"encoding/xml"
)

type VNCinfo struct {
	VNCPort string `xml:"port,attr"`
}

type DiskSource struct {
	Path string `xml:"file,attr"`
	Name string `xml:"name,attr"`
}

type DiskTarget struct {
	Dev string `xml:"dev,attr"`
}

type Disk struct {
	Source DiskSource `xml:"source"`
	Target DiskTarget `xml:"target"`
}

type Devices struct {
	Graphics VNCinfo `xml:"graphics"`
	Disks    []Disk  `xml:"disk"`
}

type xmlParseResult struct {
	Name    string  `xml:"name"`
	UUID    string  `xml:"uuid"`
	Devices Devices `xml:"devices"`
}

func ParseDomainXML(xmlData string) (*xmlParseResult, error) {
	var v = xmlParseResult{}
	err := xml.Unmarshal([]byte(xmlData), &v)
	return &v, err
}
