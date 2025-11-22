package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeAuditStorageWriter struct {
	locker sync.Mutex
	file   *os.File
	writer *bufio.Writer
}

func NewIncomeAuditStorageWriter(storagePath string) (*IncomeAuditStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &IncomeAuditStorageWriter{
		file:   storage,
		writer: bufio.NewWriter(storage),
	}, nil
}

// Closer interface

func (w *IncomeAuditStorageWriter) Close() error {
	w.locker.Lock()
	defer w.locker.Unlock()

	err := w.writer.Flush()
	if err != nil {
		return err
	}

	return w.file.Close()
}

// implementation

func (w *IncomeAuditStorageWriter) SaveData(ctx context.Context, dto *dto.IncomeAuditDto) error {
	if dto == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		{
			jsonBytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}

			err = w.writeLine(jsonBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *IncomeAuditStorageWriter) writeLine(jsonBytes []byte) error {
	w.locker.Lock()
	defer w.locker.Unlock()

	_, err := w.writer.Write(jsonBytes)
	if err != nil {
		return err
	}
	_, err = w.writer.Write([]byte("\n"))
	if err != nil {
		return err
	}

	err = w.writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
