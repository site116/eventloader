# Event Loader

🧪 **Event Loader** — это утилита на Go для генерации фейковых JSON-сообщений и отправки их в Kafka (например, Redpanda). 

## 🚀 Возможности

- Генерация фейковых данных с помощью [`gofakeit`](https://github.com/brianvoe/gofakeit) и встроенного шаблонизатора.
- Генерация батчей событий и отправка их в заданный топик Kafka
- Поддержка конкуррентной отправки батчей.

## Параметры коммандной строки
```bash
go run main.go -template=template.json -config=config.yaml
```
 - template - путь к файлу шаблона по умолчанию `template.json`
 - config - путь к файлу конфигурации по умолчанию `config.yml`

## Параметры конфигурации config.yml
```yaml
broker:
  # The topic to which the producer will send messages
  topic: "topic"
  seeds:
    # The Kafka broker address
    - "localhost:9092"
  sasl:
    # Enable SASL authentication
    enabled: true
    # Username for SASL authentication
    username: 'emm'
    # Password for SASL authentication
    password: '123'
pool:
  # Number of worker threads to process messages
  num_workers: 5
  # Number of messages to process in each batch
  messagesPerBatch: 10
  # Number of batches to process before stopping
  batches: 100
  # Number of errors allowed before stopping the worker
  errorThreshold: 2
 ```
## Пример конфига template.json
```json
{
    "subjectType":"{{Email}}",
    "objectId":"{{RandomString (SliceString "obj1" "jobj2")}}",
    "objectType":"subsystem",
    "action":"readConfig",
    "result":"{{HTTPStatusCode}}",
    "endpoint":"{{URL}}",
    "requestQuery":"modified_at=2025-04-07T12%3A44%3A43.653470506+03%3A00",
    "httpMethod":"{{HTTPMethod}}",
    "gatewayId":"gw1",
    "eventTime":"{{Date.Format "2006-01-02T15:04:05Z07:00"}}",
    "tenantCode":"default"
}
```