**Project Overview**

Simple HTTP service for managing key-value pairs with Redis as the storage backend.

**Features**
1. CRUD operations for key-value pairs
2. HTTP REST API interface
3. Redis storage backend

**Quick Start**
1. Clone the repository
```bash
git clone https://github.com/ItserX/kvRedis.git
```  
2.Start the services
```bash  
cd kvRedis  
docker compose -f deployments/docker-compose.yml up --build
```

**Run Tests**  
```bash
$ go test -cover ./internal/handlers/ 
ok      kvRedis/internal/handlers     0.006s          coverage: 64.1% of statements  
```

**API Documentation** 
Create Key-Value Pair  
`POST /kv body: {"key": "key1", "value": {"v1":1, "v2": true, "v3": [1,2,3,4,5]}}`  

Get Value by Key  
`GET /kv/{id}`  

Update Value by key  
`PUT /kv/{id} body: {"value": {"new_value": 1}}`  

Delete Key  
`DELETE /kv/{id}`  

**Configuration**
```ini
APP_PORT=:8080                    #HTTP server port  
TARANTOOL_ADDRESS=localhost:6379  #DB host:port
```
