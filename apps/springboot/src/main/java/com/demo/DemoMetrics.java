package com.demo;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.Random;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Registers the identity metric and a simulated workload.
 * Micrometer's Prometheus registry also exposes jvm_*, tomcat_*, system_*
 * metrics automatically -> obviously a JVM / Spring Boot app.
 */
@Component
public class DemoMetrics {

    private static final String LANG = "java";

    private final Counter requests;
    private final AtomicInteger inflight = new AtomicInteger(0);
    private final Random random = new Random();

    public DemoMetrics(MeterRegistry registry) {
        // Identity metric -> exposed as demo_app{language="java",app="app-springboot"} 1.0
        // NB: the metric is named "demo_app", not "demo_app_info". Spring Boot's
        // Prometheus client treats "_info" as a reserved suffix and strips it from
        // a plain gauge, so we drop the suffix to stay consistent with the other
        // three apps. strongReference keeps the supplier from being GC'd.
        Gauge.builder("demo_app", () -> 1.0)
                .description("Demo app identity")
                .tag("language", LANG)
                .tag("app", "app-springboot")
                .strongReference(true)
                .register(registry);

        // Micrometer appends _total, so this is exposed as demo_requests_total.
        this.requests = Counter.builder("demo_requests")
                .tag("language", LANG)
                .register(registry);

        Gauge.builder("demo_inflight_requests", inflight, AtomicInteger::doubleValue)
                .tag("language", LANG)
                .strongReference(true)
                .register(registry);
    }

    @Scheduled(fixedRate = 2000)
    public void tick() {
        requests.increment(random.nextInt(5) + 1);
        inflight.set(random.nextInt(20));
    }
}
