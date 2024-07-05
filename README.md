# Process Manager Server

Simple process manager written on Go

## Feature 

1. Spawn processes
2. Stop processes
3. Rerun processes
4. Return list process

## Endpoints and request examples
### `/processes`
`POST` (Spawn new process)
#### Request
``` json
{
    "command": "ping",
    "args": ["1.1.1.1"]
}
```
#### Response
``` json
{
    "id": 3,
    "command": "ping",
    "args": [
        "1.1.1.1"
    ],
    "running": true,
    "pid": 79376
}
```
`GET` (Get processes list)
#### Request
``` 
empty 
```
#### Response
``` json
[
    {
    "id": 3,
    "command": "ping",
    "args": [
        "1.1.1.1"
    ],
    "running": true,
    "pid": 79376
    },
...
]
```
### `processes/{id}`
`DELETE` (Stop process and remove from list)
#### Request
``` 
empty 
```
#### Response
```
204 No Content (Successful)
```
### `processes/rerun/{id}`
`POST` (Start new proc after kill/normal exit)
#### Request
``` 
empty
```
#### Response
```json
{
    "id": 0,
    "command": "ping",
    "args": [
        "1.1.1.1"
    ],
    "running": true,
    "pid": 8228
}
# New PID and "running: true"
```
```
400 Bad Request
Plain text: 
"Process with id 0 is already running" (If process already running)

Plain text: 
"Process with id 0 not found" (If process has been removed)
```
## Plans
1. CLI for setting HOST and PORT parameters
2. API improvement
3. Normal error messages
4. Uniform responses
5. State storage
6. Recovery and restart process from stage storage after restart server

## Related projects
UI for that server [Remote Process Manager](https://github.com/RomanCherkasov/process_manager_ui?tab=readme-ov-file)