package main

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"strconv"
)

type TmxMap struct {
	Version      string           `xml:"version,attr"`
	Orientation  string           `xml:"orientation,attr"`
	Width        int              `xml:"width,attr"`
	Height       int              `xml:"height,attr"`
	TileWidth    int              `xml:"tilewidth,attr"`
	TileHeight   int              `xml:"tileheight,attr"`
	Properties   []TmxProperties  `xml:"properties"`
	Tilesets     []TmxTileset     `xml:"tileset"`
	Layers       []TmxLayer       `xml:"layer"`
	ObjectGroups []TmxObjectGroup `xml:"objectgroup"`
}

type TmxProperties struct {
	Property []TmxProperty `xml:"property"`
}

type TmxProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TmxTileset struct {
	FirstGid   int        `xml:"firstgid,attr"`
	Name       string     `xml:"name,attr"`
	TileWidth  int        `xml:"tilewidth,attr"`
	TileHeight int        `xml:"tileheight,attr"`
	Images     []TmxImage `xml:"image"`
}

type TmxImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type TmxLayer struct {
	Name   string  `xml:"name,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Data   TmxData `xml:"data"`
}

type TmxData struct {
	Encoding string `xml:"encoding,attr"`
	Value    string `xml:",chardata"`
	ParsedData []int
}

type TmxObjectGroup struct {
	Name    string      `xml:"name,attr"`
	Width   int         `xml:"width,attr"`
	Height  int         `xml:"height,attr"`
	Objects []TmxObject `xml:"object"`
}

type TmxObject struct {
	Type   string `xml:"type,attr"`
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// Parse reads the TMX-encoded data and converts it into a TmxMap object
// Returns an error if the TMX-encoded data is malformed
func Parse(b []byte) (TmxMap, error) {
	var parsed TmxMap
	err := xml.Unmarshal(b, &parsed)

	for i, v := range parsed.Layers {
		str := strings.Replace(v.Data.Value, "\n", "", -1)
		arr := strings.Split(str, ",")
		var converted []int
		for _, v := range arr {
			num, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			converted = append(converted, num)
		}
		parsed.Layers[i].Data.ParsedData = converted
	}

	return parsed, err
}

func openTMX(filename string) TmxMap {
	path := fileManager.GetPath("tmx", filename, "tmx")
	f, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	parsed, _ := Parse(f)
	return parsed
}
