package filereader

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type FileReaderWriter struct {
	fileName string
	file     *os.File
}

type FileURL struct {
	UUID      string `json:"uuid"`
	ShortURL  string `json:"short_url"`
	OriginURL string `json:"origin_url"`
}

type FileURLs []FileURL

func New(filename string) (*FileReaderWriter, error) {
	if filename == "" {
		return &FileReaderWriter{
			fileName: filename,
			file:     nil,
		}, nil
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileReaderWriter{
		fileName: filename,
		file:     file,
	}, nil
}

func (f *FileReaderWriter) ReadFileEvents() (FileURLs, error) {
	if f.fileName == "" {
		return nil, nil
	}
	scan := bufio.NewScanner(f.file)
	urls := make(FileURLs, 0)
	for scan.Scan() {
		var url FileURL
		data := scan.Bytes()
		err := json.Unmarshal(data, &url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (f FileReaderWriter) WriteFileEvent(shortURL, originURL string) error {
	if f.fileName == "" {
		return nil
	}
	event := FileURL{
		UUID:      uuid.New().String(),
		ShortURL:  shortURL,
		OriginURL: originURL,
	}
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	_, err = f.file.Write(data)
	if err != nil {
		return err
	}
	if _, err := f.file.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}
