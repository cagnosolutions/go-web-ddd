package db

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type FileDB struct {
	fd *os.File
	r  *csv.Reader
	w  *csv.Writer
}

func OpenFileDB(path string) (*FileDB, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fd, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		err = fd.Close()
		if err != nil {
			return nil, err
		}
	}
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FileDB{
		fd: fd,
		r:  csv.NewReader(fd),
		w:  csv.NewWriter(fd),
	}, nil
}

func (f *FileDB) Get(id int) (interface{}, error) {
	var record interface{}
	for {
		data, err := f.r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if data[0] == strconv.Itoa(id) {
			record = data
			break
		}
	}
	if record == nil {
		return nil, errors.New("not found")
	}
	return record, nil
}

func (f *FileDB) GetAll() ([]interface{}, error) {
	data, err := f.r.ReadAll()
	if err != nil {
		return nil, err
	}
	var records []interface{}
	for i := range data {
		records = append(records, data[i])
	}
	return records, nil
}

func (f *FileDB) Close() error {
	err := f.fd.Sync()
	if err != nil {
		return err
	}
	err = f.fd.Close()
	if err != nil {
		return err
	}
	return nil
}
