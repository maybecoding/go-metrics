package backupstorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
)

type BackupStorage struct {
	interval      int64
	path          string
	isRestoreOnUp bool
}

func NewBackupStorage(interval int64, path string, isRestoreOnUp bool) *BackupStorage {
	return &BackupStorage{
		interval:      interval,
		path:          path,
		isRestoreOnUp: isRestoreOnUp,
	}
}

func (bs *BackupStorage) Save(metrics []*sapp.Metric) error {
	logger.Log.Info().Msg("start metrics saving")
	if bs.path == "" || len(metrics) == 0 {
		return nil
	}

	// Открываем файл на перезапись
	outFile, err := os.OpenFile(bs.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return fmt.Errorf("error due open file for writing: %w", err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	// По каждой метрике выполняем сохранение в файл
	for _, metric := range metrics {
		mj, err := json.Marshal(metric)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due marshal metric")
			continue
		}
		_, err = outFile.Write(append(mj, '\n'))
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due write metric")
			continue
		}
	}

	return nil
}

func (bs *BackupStorage) Restore() ([]*sapp.Metric, error) {
	logger.Log.Info().Msg("start metrics restoring")
	if bs.path == "" {
		return nil, nil
	}
	// Открываем файл на чтение
	inFile, err := os.OpenFile(bs.path, os.O_RDONLY|os.O_CREATE, 0660)
	if err != nil {
		return nil, fmt.Errorf("error due open file for reading: %w", err)
	}
	defer func() {
		_ = inFile.Close()
	}()

	metrics := make([]*sapp.Metric, 0)
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		data := scanner.Bytes()
		m := sapp.Metric{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due read metric from line")
			continue
		}
		metrics = append(metrics, &m)
	}
	return metrics, nil
}

func (bs *BackupStorage) GetBackupInterval() int64 {
	return bs.interval
}

func (bs *BackupStorage) GetIsNeedRestore() bool {
	return bs.isRestoreOnUp
}
