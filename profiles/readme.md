## Оптимизация аллокаций памяти
### Диагностика аллокаций памяти
Сбор метрик выполнялся на основе теста agent_benchmark_test.go из каталога cmd/agent с помощью команды
```shell
 go test -bench=BenchmarkMain -memprofile=base.pprof 
```

### Результаты
Было выявлено, что основная аллокация памяти приходится на flate.NewWriter
```
690.16MB  74.09% 74.09%  832.01MB 89.32%  compress/flate.NewWriter (inline)
```

### Выполненные действия
Блок, кода по сжатию отправляемых данных в файле internal/agent/sender/sendmetric.go:
```go
//Создаем сжатый поток (до оптимизации по инкременту 18)
buf := bytes.NewBuffer(nil)
zw := gzip.NewWriter(buf)

// И записываем в него данные
_, err = zw.Write(rd)
if err != nil {
    logger.Error().Err(err).Msg("can't write into gzip writer")
    return
}
err = zw.Close()
if err != nil {
    logger.Error().Err(err).Msg("can't close gzip writer")
    return
}
rdGz := buf.Bytes()
```

Был заменен на:
```go
rdGz, err := zipper.ZippedBytes(rd)
if err != nil {
    logger.Error().Err(err).Msg("sender - sendMetric - zipper.ZippedBytes")
    return
}
```

В котором используется функция *ZippedBytes* нового пакета zipper, в котором реализовано
переиспользование *bytes.Buffer и *gzip.Writer при помощи sync.Pool 

### Проверка эффективности оптимизации
После применения оптимизации был запущен повторный тест командой
```shell
 go test -bench=BenchmarkMain -memprofile=result.pprof 
```
Затем командой 
```shell
go tool pprof -top -diff_base=base.pprof result.pprof
```
Установлено что потребление памяти снизилось на ~550 MB
```
      flat  flat%   sum%        cum   cum%
 -556.18MB 59.71% 59.71%  -670.26MB 71.95%  compress/flate.NewWriter (inline)

```