# Thanos Unification Demo

A self-contained, laptop-runnable demo for the "how I unified multiple Prometheus with Thanos" talk.

Four apps in four languages each expose `/metrics`. Two Prometheus servers each scrape a different pair of apps (so no single Prometheus sees everything). Thanos Query fans out across both and presents one unified, deduplicated view that Grafana plugs into. MinIO stands in for AWS S3.

For brevity this omits the Compactor and Store Gateway. The demo is about **unifying live metrics under one source**, not historical read-back.

## What's in the stack

| Container | Language / role | Recognisable metrics |
|---|---|---|
| `app-python` | Python (Flask) | `python_gc_*`, `python_info`, `demo_app{language="python"}` |
| `app-node` | Node.js (Express) | `nodejs_*`, `process_*`, `demo_app{language="nodejs"}` |
| `app-go` | Go (standard library only) | `go_info`, `go_goroutines`, `go_memstats_*`, `demo_app{language="go"}` |
| `app-springboot` | Java (Spring Boot + Micrometer) | `jvm_*`, `tomcat_*`, `demo_app{language="java"}` |
| `app-shared` | Python (minimal) | `shared_heartbeat_total` — scraped by **both** Prometheus, for the dedup demo |
| `prometheus-1` | replica **A** | scrapes python + springboot + shared |
| `prometheus-2` | replica **B** | scrapes node + go + shared |
| `thanos-sidecar-1/2` | Thanos sidecars | expose each Prometheus over gRPC; ship blocks to MinIO |
| `thanos-query` | Thanos Query | fan-out + `--query.replica-label=replica` |
| `minio` + `minio-init` | S3-compatible object storage | bucket `thanos` created on startup |
| `grafana` | Grafana | Thanos pre-wired as default datasource + starter dashboard |

## Run it

Requires Docker Desktop (or any Docker Engine) with Compose v2. Works on Apple Silicon.

```bash
cd thanos-unification-demo
docker compose up --build -d
```

First build takes a couple of minutes (Spring Boot pulls Maven deps; Python/Node pull pip/npm packages). The Go app is standard-library only, so its build downloads nothing. Watch it come up:

```bash
docker compose logs -f thanos-query
```

Tear down (removes volumes too):

```bash
docker compose down -v
```

## Ports

| URL | What |
|---|---|
| http://localhost:3001 | **Grafana** — dashboard "Thanos Unification Demo" (`admin` / `admin`; anonymous access is read-only) |
| http://localhost:9090 | **Thanos Query** — unified UI, Stores page, dedup toggle |
| http://localhost:9091 | **Prometheus A** — only sees python + springboot (+ shared) |
| http://localhost:9092 | **Prometheus B** — only sees node + go (+ shared) |
| http://localhost:9001 | **MinIO console** — `minioadmin` / `minioadmin` |
| http://localhost:8001/metrics | app-python |
| http://localhost:8002/metrics | app-node |
| http://localhost:8003/metrics | app-go |
| http://localhost:8004/actuator/prometheus | app-springboot |
| http://localhost:8005/metrics | app-shared |

If any host port clashes with something on your laptop, change the left-hand side of the `ports:` mapping in `docker-compose.yml`.