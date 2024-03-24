package mainosexit

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	// Функция analysistest.Run применяет тестируемый анализатор mainosexit.Analyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
