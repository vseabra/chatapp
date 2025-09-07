# Bot Integration Guide

This defines the spec (AMPQ/ message broker) bots need to implement. In short:

1) User sends `/stock aapl` → backend publishes `bot.requested` 
2) Bot consumes, performs work, renders reply.
3) Bot publishes to `bot.response.submit` → backend persists as bot message and broadcasts.

## Routing keys
- **bot.requested** (incoming work): user invoked a command (e.g., `/stock aapl`). Consume this.
- **bot.response.submit** (outgoing reply): publish your reply to this key. Backend persists as `type="bot"` and broadcasts.

## Payload contracts

### BotRequested (consume from `bot.requested`)
```json
{
  "command": "stock",
  "args": "aapl.us",
  "roomId": "<ROOM_ID>",
  "requestUserId": "<USER_ID>",
  "requestedAt": "2025-01-01T12:00:00Z"
}
```

### BotResponseSubmit (publish to `bot.response.submit`)
```json
{
  "roomId": "<ROOM_ID>",
  "text": "AAPL quote is 187.43 USD"
}
```
## Test via RabbitMQ HTTP API
- Publish `bot.requested`:
```bash
curl -u guest:guest -H "content-type: application/json" -X POST \
  -d '{"properties":{"content_type":"application/json"},"routing_key":"bot.requested","payload":"{\"command\":\"stock\",\"args\":\"aapl.us\",\"roomId\":\"<ROOM_ID>\",\"requestUserId\":\"<USER_ID>\",\"messageId\":\"\",\"requestedAt\":\"2025-01-01T12:00:00Z\"}","payload_encoding":"string"}' \
  http://localhost:15672/api/exchanges/%2F/chat.events/publish
```

- Publish `bot.response.submit`:
```bash
curl -u guest:guest -H "content-type: application/json" -X POST \
  -d '{"properties":{"content_type":"application/json"},"routing_key":"bot.response.submit","payload":"{\"roomId\":\"<ROOM_ID>\",\"text\":\"BOT: test reply\"}","payload_encoding":"string"}' \
  http://localhost:15672/api/exchanges/%2F/chat.events/publish
```

