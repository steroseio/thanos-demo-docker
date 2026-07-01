const express = require("express");
const client = require("prom-client");

const LANG = "nodejs";
const register = client.register;

// Default metrics expose nodejs_*, process_* -> obviously Node.js.
client.collectDefaultMetrics();

// Identity metric. Named "demo_app" (no "_info" suffix) to stay consistent with
// the Spring Boot app, whose Prometheus client strips the reserved "_info" suffix.
const appIdentity = new client.Gauge({
  name: "demo_app",
  help: "Demo app identity",
  labelNames: ["language", "app"],
});
appIdentity.set({ language: LANG, app: "app-node" }, 1);

// Simulated workload so the graphs move during the demo.
const requestsTotal = new client.Counter({
  name: "demo_requests_total",
  help: "Simulated processed requests",
  labelNames: ["language"],
});
const inflight = new client.Gauge({
  name: "demo_inflight_requests",
  help: "Simulated in-flight requests",
  labelNames: ["language"],
});

setInterval(() => {
  requestsTotal.inc({ language: LANG }, Math.floor(Math.random() * 5) + 1);
  inflight.set({ language: LANG }, Math.floor(Math.random() * 20));
}, 2000);

const app = express();
app.get("/", (_req, res) => res.type("text").send(`Demo app: ${LANG}\n`));
app.get("/metrics", async (_req, res) => {
  res.set("Content-Type", register.contentType);
  res.end(await register.metrics());
});

app.listen(3000, () => console.log("node demo app listening on :3000"));
