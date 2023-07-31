## Работа с клиентом

### 1. подготовка контекста

#### Заполним основные поля
1. BufferBatch - кеш готовых батчей
2. BufferInput - кеш входящего потока
3. BufferResort - кеш буффера для ресорта
4. Endpoint - адрес сервиса
5. CheckStatusInterval - время опроса сервиса на доступность
6. Service - сервис должен реализовывать интерфейс, как в syncer/interfaces -> Service

```go
  client := (&Client{
    BufferBatch:  10, 
    BufferInput:  50,
    BufferResort: 3,
    Endpoint:     "https://external-service.local",
    CheckStatusInterval:  5 * time.Second,
    Service:      &service.ExternalBadService{ProcessLimit: 5, Duration: 10 * time.Second},
  }).Run()

```
### 2. После клиент готов для приёма новых items

```go
client.AddNewItem(interfaces.Item{})
```


### 3. Подпишемся на события от клиента

```go

// получить информацию о последней ошибке 
client.ExceptionCallback = func(err error) {
    
}

// Перехватить отправленный batch
client.SuccessCallback = func(batch interfaces.Batch) {
    
}

```
### 4. Получить дополнительную инфу о процессе/прогрессе 

```go

// сколько успешно отправленных бачей, а сколько ошибочных
success, exceptions := client.GetStats() 


// возвращает информацию о заполненности каналов, где 
// 0 - количество входящих в очереди
// 1 - количество батчей на пересортировку
// 2 - количество бачей на синхронизацию
chanStats := client.GetChanStats() // [2, 0, 5]


```