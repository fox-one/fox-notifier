# fox-notifier

## Notify Message

```bash
POST /message
```

**Params:**

```javascript
{
    "message_id": "191b3672-1014-4f54-aeed-48caf0a8d0af", // optional
    "topic": "test", // required
    "message": "just a test", // required
}
```

**Response:**

```javascript
{
  "code": 0,
  "data": {
    "id": 3,
    "topic": "test",
    "conversation_id": "3fed03e2-9799-3cf1-9a11-c90574e99209",
    "message_id": "191b3672-1014-4f54-aeed-48caf0a8d0af",
    "message": "just a test",
    "created_at": "2019-10-09T21:28:08.525309+08:00",
    "updated_at": "2019-10-09T21:28:08.525309+08:00"
  }
}
```
