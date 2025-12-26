# Self-Hosted Observability Platform

A personal telemetry system that collects **system metrics** (CPU, memory, disk, network) and **application events** (like reading pipeline status), stores them in **PostgreSQL with TimescaleDB**, and visualizes everything in **Grafana**—so I can explore, correlate, and understand system behavior over time.

---

## 🔍 What It Does

- **Collects system metrics** (CPU, memory, disk, network) from host machines via a lightweight Go collector (`gopsutil`)
- **Tracks application events** from [personal reading analytics dashboard](https://github.com/victoriacheng15/personal-reading-analytics-dashboard)
- **Stores all data in TimescaleDB (PostgreSQL)** using a flexible JSONB schema for long-term retention and cross-type querying
- **Visualizes in Grafana** with two focused dashboards:
  - **Infrastructure Health**: per-host metrics over time
  - **Application Telemetry**: event counts, success rates, trends
- **Ensures durability** via automated volume backups and idempotent collectors

---

## 🏗 Architecture Overview

The system is designed for simplicity and reliability:

- **Go collectors**: Cron-driven, stateless, and idempotent  
- **TimescaleDB**: Single table with JSONB for unified metric/event storage  
- **Grafana**: Dashboards separated by concern (infra vs. app)  
- **Extensible observability**: Includes **Loki** for log aggregation and a **Go proxy** for ETL sync (e.g., pulling reading pipeline events from MongoDB), enabling full telemetry coverage: metrics, events, and logs.

For full architecture details, data model, and Mermaid diagrams:  
→ [docs/architecture/observability-system.md](./docs/architecture/observability-system.md)

---

## 💡 Why I Built This

I built this to understand how observability can reveal system behavior over time—not just through logs or one-off checks, but through continuous, contextual telemetry.

It started as a personal experiment.  
Now, it’s how I monitor host health, debug pipeline failures, and connect infrastructure metrics to application events—all without cloud dependencies.
