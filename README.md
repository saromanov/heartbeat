# heartbeat


# Usage

Adding a new endpoint for check

```

curl -X POST -H 'Content-type: application/json' --data '{"title":"Hello, World!", "url":"https://github.com"}' localhost:8100/api/checks

```

Getting the report of registered checks

```
 curl localhost:8100/api/status
```

