package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type ShortURLStorageReader struct {
	file   *os.File
	reader *bufio.Scanner
	log    logger.Logger
}

func NewShortURLStorageReader(storagePath string, log logger.Logger) (*ShortURLStorageReader, error) {
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURLStorageReader{
		file:   storage,
		reader: bufio.NewScanner(storage),
		log:    log.GetLogger("shortURI storage reader"),
	}, nil
}

type ShortURLStorageWriter struct {
	file   *os.File
	writer *bufio.Writer
	log    logger.Logger
}

func NewShortURLStorageWriter(storagePath string, log logger.Logger) (*ShortURLStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURLStorageWriter{
		file:   storage,
		writer: bufio.NewWriter(storage),
		log:    log.GetLogger("shortURI storage writer"),
	}, nil
}

func (r *ShortURLStorageReader) Close() error {
	return r.file.Close()
}

func (r *ShortURLStorageReader) LoadData(cache map[string]*model.ShortURI) error {
	clear(cache)
	for r.reader.Scan() {
		data := r.reader.Bytes()
		var shortURL model.ShortURI
		if err := json.Unmarshal(data, &shortURL); err != nil {
			r.log.Warn("Failed to unmarshal short URI")
			continue
		}
		r.log.Infof("Add short URI: [%s]", shortURL.ID)
		cache[shortURL.ID] = &shortURL
	}
	if err := r.reader.Err(); err != nil {
		return err
	}

	return nil
}

func (w *ShortURLStorageWriter) Close() error {
	return w.file.Close()
}

func (w *ShortURLStorageWriter) SaveData(cache map[string]*model.ShortURI) error {
	w.log.Infof("Short URI length [%d]", len(cache))
	for id, shortURL := range cache {
		w.log.Infof("Saving short URI [%s] to [%s]", id, w.file.Name())
		data, err := json.Marshal(shortURL)
		if err != nil {
			w.log.Warnf("Failed to marshal short URI id [%s]", id)
			return err
		}
		if err := w.writeLine(data, shortURL.ID); err != nil {
			w.log.Warnf("Failed to write short URI id [%s]", id)
		}
	}

	return nil
}

func (w *ShortURLStorageWriter) writeLine(data []byte, id string) error {
	if _, err := w.writer.Write(data); err != nil {
		w.log.Warnf("Failed to write short URI id [%s]", id)

		return err
	}
	if err := w.writer.WriteByte('\n'); err != nil {
		w.log.Warnf("Failed to write term symbol short URI id [%s]", id)

		return err
	}

	if err := w.writer.Flush(); err != nil {
		w.log.Warnf("Failed to flush short URI id [%s]", id)

		return err
	}

	return nil
}
