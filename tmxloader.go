package main

import (
	"encoding/xml"
	"strings"
	"strconv"
	"io/ioutil"
	"fmt"
	"os"
	"errors"
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

func (t *TmxMap) GetLayerByName(name string) (*TmxLayer, error) {
	var layer *TmxLayer
	var err error
	for _, v := range t.Layers {
		if v.Name == name {
			layer = &v
		}
	}
	if layer == nil {
		err = errors.New("No layer with named " + name)
	}
	return layer, err
}

func (t *TmxMap) GetObjGroupByName(name string) (*TmxObjectGroup, error) {
	var group *TmxObjectGroup
	var err error
	for _, v := range t.ObjectGroups {
		if v.Name == name {
			group = &v
		}
	}
	if group == nil {
		err = errors.New("No layer with named " + name)
	}
	return group, err
}

func parseTmxMap(b []byte) (TmxMap, error) {
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

func NewTmxMap(path string) (TmxMap, error) {
	f, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read tmx file: %s\n", err)
	}

	parsed, err := parseTmxMap(f)
	return parsed, err
}