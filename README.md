# AB Platform

Simple A/B testing backend in Go.

Flow:

1. Create experiment
2. Assign user to A/B
3. Send events (impression, click)
4. Worker updates metrics from Kafka
5. Read results with uplift and significance

## Stack

- Go + Gin
- PostgreSQL
- Kafka

## Quick Start

1. Start docker compose

```bash
docker compose up -d
```

2. Create tables (run once in Postgres)

```sql
CREATE TABLE IF NOT EXISTS experiments (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'running',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assignments (
    id SERIAL PRIMARY KEY,
    experiment_id INT NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    variant TEXT NOT NULL CHECK (variant IN ('A', 'B')),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (experiment_id, user_id)
);

CREATE TABLE IF NOT EXISTS metrics (
    experiment_id INT PRIMARY KEY REFERENCES experiments(id) ON DELETE CASCADE,
    impressions_a INT NOT NULL DEFAULT 0,
    impressions_b INT NOT NULL DEFAULT 0,
    conversions_a INT NOT NULL DEFAULT 0,
    conversions_b INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_events (
    id SERIAL PRIMARY KEY,
    experiment_id INT NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    event_name TEXT NOT NULL CHECK (event_name IN ('impression', 'click')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (experiment_id, user_id, event_name)
);
```

3. Run services in separate terminals

```bash
go run ./cmd/api
go run ./cmd/worker
```

4. Optional load test (simulator)

```bash
go run ./cmd/simulator
```

## API

Create experiment:

```bash
curl -X POST http://localhost:8080/experiments -H "Content-Type: application/json" -d '{"name":"CTA test"}'
```


Assign user:

```bash
curl "http://localhost:8080/assign?experiment_id=1&user_id=user123"
```

Send event:

```bash
curl -X POST http://localhost:8080/events -H "Content-Type: application/json" -d '{"experiment_id":1,"user_id":"user123","event_name":"impression"}'
```

Read results:

```bash
curl "http://localhost:8080/results?experiment_id=1"
```
