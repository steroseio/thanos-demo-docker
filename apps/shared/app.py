import time

from prometheus_client import Counter, Info, start_http_server

# This target is scraped by BOTH Prometheus servers (replica A and replica B).
# Its series are therefore identical except for the replica label, so Thanos
# Query deduplicates them via --query.replica-label=replica.
Info("shared_app", "Shared demo target scraped by both Prometheus servers").info(
    {"role": "shared", "note": "scraped-by-both-prometheus"}
)

heartbeat = Counter("shared_heartbeat", "Heartbeat from the shared target")

if __name__ == "__main__":
    start_http_server(9000)  # serves /metrics on :9000
    while True:
        heartbeat.inc(1)
        time.sleep(2)
