# Storage
## Environment Variables
* POSTGRES_ADDR - hostname or ip
* POSTGRES_PORT - port (standard 5432) 
* POSTGRES_USER - postgres user
* POSTGRES_PASS - postgres user password
* POSTGRES_DB - database name
## Build
```bash
docker build -t storage .
```
## Run
```bash
docker-compose up -d
```
