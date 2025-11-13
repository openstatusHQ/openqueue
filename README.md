# openqueue

A simple HTTP message queue


If you want to run it locally:
```
curl \
    --header 'Content-Type: application/json' \
    --data '{"queue_name": "my_queue", "task":{"url":"http://localhost:1337/ping", "method":"GET"}}' \
    http://localhost:8080/api.v1.QueueService/CreateTask
```
