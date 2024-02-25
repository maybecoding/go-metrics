package metricmemstorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
)

type Dumper struct {
	path string
}

func NewDumper(path string) *Dumper {
	return &Dumper{
		path: path,
	}
}

func (bs *Dumper) Save(metrics []*metric.Metrics) error {
	logger.Info().Msg("start metric saving")
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
	for _, m := range metrics {
		mj, err := json.Marshal(m)
		if err != nil {
			logger.Error().Err(err).Msg("error due marshal metric")
			continue
		}
		_, err = outFile.Write(append(mj, '\n'))
		if err != nil {
			logger.Error().Err(err).Msg("error due write metric")
			continue
		}
	}

	return nil
}

func (bs *Dumper) Restore() ([]metric.Metrics, error) {
	logger.Info().Msg("start metric restoring")
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

	metrics := make([]metric.Metrics, 0)
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		data := scanner.Bytes()
		m := metric.Metrics{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			logger.Error().Err(err).Msg("error due read metric from line")
			continue
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
