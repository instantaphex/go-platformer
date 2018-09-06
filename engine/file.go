package engine

import (
	"fmt"
	"os"
	"io/ioutil"
)

type FileManager struct {}

func (f *FileManager) GetMap(filename string) (*os.File, error) {
	mapPath := f.GetPath("maps", filename, "map")
	return f.GetContents(mapPath)
}

func (f *FileManager) GetFontPath(filename string) string {
	return f.GetPath("fonts", filename, "ttf")
}

func (f *FileManager) GetArea(filename string) (*os.File, error) {
	areaPath := f.GetPath("maps", filename, "area")
	return f.GetContents(areaPath)
}

func (f *FileManager) GetAudioPath(filename string) string {
	return f.GetDirectoryPath("audio") + "/" + filename
}

func (f *FileManager) GetContents(path string) (*os.File, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return fp, nil
}

func (f *FileManager) GetDirectoryContents(dir string) ([]string, error) {
	var fileNames []string
	files, err := ioutil.ReadDir(f.GetDirectoryPath(dir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Filed to read directory: %s\n", err)
		return fileNames, err
	}
	for _, file := range files {
		fileNames = append(fileNames, f.GetDirectoryPath("audio") + "/" + file.Name())
	}
	return fileNames, nil
}

func (f *FileManager) GetDirectoryPath(dir string) string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get working directory: %s", err)
	}

	return path + "/" + dir
}

func (f *FileManager) GetImagePath(filename string) string {
	return f.GetPath("assets", filename, "png")
}

func (f *FileManager) GetTilesetPath(filename string) string {
	return f.GetPath("tilesets", filename, "png")
}

func (f *FileManager) GetPath(dir string, filename string, fileExtension string) string {
	return f.GetDirectoryPath(dir) + "/" + filename + "." + fileExtension
}
