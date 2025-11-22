package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type ShortURIUserStorageReader struct {
	file   *os.File
	reader *bufio.Scanner
	log    logger.Logger
}

func NewShortURIUserStorageReader(storagePath string, log logger.Logger) (*ShortURIUserStorageReader, error) {
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURIUserStorageReader{
		file:   storage,
		reader: bufio.NewScanner(storage),
		log:    log.GetLogger("shortURIUser storage reader"),
	}, nil
}

func (r *ShortURIUserStorageReader) Close() error {
	return r.file.Close()
}

type ShortURIUserStorageWriter struct {
	file   *os.File
	writer *bufio.Writer
	log    logger.Logger
}

func NewShortURIUserStorageWriter(storagePath string, log logger.Logger) (*ShortURIUserStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURIUserStorageWriter{
		file:   storage,
		writer: bufio.NewWriter(storage),
		log:    log.GetLogger("shortURIUser storage writer"),
	}, nil
}

func (w *ShortURIUserStorageWriter) Close() error {
	return w.file.Close()
}

func (r *ShortURIUserStorageReader) LoadData(cache map[string]*model.ShortURIUser) error {
	clear(cache)
	for r.reader.Scan() {
		r.log.Infof("Loading data: [%s]", r.reader.Text())
		data := r.reader.Bytes()
		var entity model.ShortURIUser
		if err := json.Unmarshal(data, &entity); err != nil {
			r.log.Errorf("failed to unmarshal entity: [%v]", err)
			continue
		}
		cache[entity.ID] = &entity
	}
	if err := r.reader.Err(); err != nil {
		return err
	}

	return nil
}

func (w *ShortURIUserStorageWriter) SaveData(cache map[string]*model.ShortURIUser) error {
	for id, entity := range cache {
		w.log.Infof("Saving short URI [%s] to [%s]", id, w.file.Name())
		data, err := json.Marshal(entity)
		if err != nil {
			w.log.Warnf("Failed to marshal short URI id [%s]", id)
			return err
		}
		if err := w.writeLine(data, entity.ID); err != nil {
			w.log.Warnf("Failed to write short URI id [%s]", id)
		}
	}

	return nil
}

func (w *ShortURIUserStorageWriter) writeLine(data []byte, id string) error {
	if _, err := w.writer.Write(data); err != nil {
		w.log.Warnf("Failed to write short URI user with id [%s]", id)

		return err
	}
	if err := w.writer.WriteByte('\n'); err != nil {
		w.log.Warnf("Failed to write term symbol short URI user with id [%s]", id)

		return err
	}

	if err := w.writer.Flush(); err != nil {
		w.log.Warnf("Failed to flush short URI user with id [%s]", id)

		return err
	}

	return nil
}
