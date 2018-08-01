package main

import "os"

type FileManager struct {}

func (f *FileManager) GetMap(filename string) (*os.File, error) {
	mapPath := f.GetPath("maps", filename, "map")
	return f.GetContents(mapPath)
}

func (f *FileManager) GetArea(filename string) (*os.File, error) {
	areaPath := f.GetPath("maps", filename, "area")
	return f.GetContents(areaPath)
}

func (f *FileManager) GetContents(path string) (*os.File, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return fp, nil
}

func (f *FileManager) GetImagePath(filename string) string {
	return f.GetPath("assets", filename, "png")
}

func (f *FileManager) GetTilesetPath(filename string) string {
	return f.GetPath("tilesets", filename, "png")
}

func (f *FileManager) GetPath(dir string, filename string, fileExtension string) string {
	wd, _ := os.Getwd()
	return wd + "/" + dir + "/" + filename + "." + fileExtension
}