# Оптимизация аллокаций памяти
## Агент
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

## Сервер
### Методика проверки
Сбор метрик выполнялся с помощью запуска приложения:
```shell
go run ./cmd/server -d 'postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
```

С последующим запуском нагрузочного скрипта
```shell
make test-memory
```
Снятие метрик производилось после завершения скрипта командой
```shell
curl -o result.profile  'http://localhost:8090/debug/pprof/heap?gc=0'
```

### Узкие места
Обнаружено, что основные аллокации приходят на декодирование json
```
18.03MB 23.60%  23.60% 18.53MB 24.26%  encoding/json.(*Decoder).refill
```

### Выполненные действия
Блок, кода по сжатию отправляемых данных в файле internal/server/handlers/metricupdatealljson.go:
```go
decoder := json.NewDecoder(r.Body) Оптимизировано инкремент #16
defer func() {
	_ = r.Body.Close()
}()

var mts []entity.Metrics
err := decoder.Decode(&mts)
```

Был заменен на:
```go
mts := mtsPool.Get().(entity.MetricsList)
mtsForRead := mts[0:0]
err := easyjson.UnmarshalFromReader(r.Body, &mtsForRead)
defer mtsPool.Put(mts)
```
Благодаря сгенерированному коду по объектам в пакете entity

### Результат
Разбор json теперь занимает меньше памяти, однако теперь само чтение из потока по занимает довольно значимую часть
```
9751.58kB 54.29%  54.29% 9751.58kB 54.29%  io.ReadAll
```

