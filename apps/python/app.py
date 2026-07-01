import random
import threading
import time

from flask import Flask
from prometheus_client import CONTENT_TYPE_LATEST, Counter, Gauge, generate_latest

LANG = "python"
app = Flask(__name__)

# Identity metric: unmistakable "this is the python app" signal.
# Named "demo_app" (no "_info" suffix) to stay consistent with the Spring Boot
# app, whose Prometheus client strips the reserved "_info" suffix from gauges.
app_identity = Gauge("demo_app", "Demo app identity", ["language", "app"])
app_identity.labels(LANG, "app-python").set(1)

# Simulated workload so the graphs move during the demo.
requests_total = Counter("demo_requests", "Simulated processed requests", ["language"])
inflight = Gauge("demo_inflight_requests", "Simulated in-flight requests", ["language"])


def _workload():
    while True:
        requests_total.labels(LANG).inc(random.randint(1, 5))
        inflight.labels(LANG).set(random.randint(0, 20))
        time.sleep(2)


@app.get("/")
def root():
    return f"Demo app: {LANG}\n"


@app.get("/metrics")
def metrics():
    # Default registry also exposes python_gc_*, python_info, process_* -> obviously Python.
    return generate_latest(), 200, {"Content-Type": CONTENT_TYPE_LATEST}


if __name__ == "__main__":
    threading.Thread(target=_workload, daemon=True).start()
    app.run(host="0.0.0.0", port=8000)
