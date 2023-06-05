# **Hit Shop**



---
docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres

docker-compose up -d


select * from pg_available_extensions;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


curl -X 'POST'   'http://localhost:8080/api/session' -H 'Authorization: Bearer ...'
