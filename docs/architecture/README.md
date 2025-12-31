# Observability Hub Architecture

This document serves as the entry point for the system's architecture.

## System Context

The hub integrates standard observability tools with custom Go services.

```mermaid
graph TD
    subgraph "Host Environment"
        Hardware[Host Hardware]
        Docker[Docker Containers]
    end

    subgraph "External Data"
        Mongo[(MongoDB)]
    end

    subgraph "Observability Hub"
        direction TB
        
        subgraph "Collection"
            Metrics[System Metrics Collector]
            Promtail[Promtail Agent]
            Proxy[Proxy Service / ETL]
        end

        subgraph "Storage"
            PG[(PostgreSQL)]
            Loki[(Loki)]
        end

        subgraph "Visualization"
            Grafana[Grafana Dashboard]
        end
    end

    %% Data Flows
    Hardware -->|Stats: CPU, RAM, Disk, Memory| Metrics
    Metrics -->|Writes Metrics| PG
    
    Mongo -->|Reads 'ingested' docs| Proxy
    Proxy -->|Writes Analytics| PG
    
    Docker -->|Logs| Promtail
    Promtail -->|Pushes Logs| Loki
    
    PG -->|Query Metrics| Grafana
    Loki -->|Query Logs| Grafana
```

## Detailed Architecture Documents

| Component | Description |
| :----------- | :------------- |
| **[Proxy Service](./proxy-service.md)** | Architecture of the Go-based API Gateway and ETL Engine. |
| **[System Metrics](./system-metrics.md)** | Details on the custom host telemetry collector (`gopsutil`). |
| **[Infrastructure](./infrastructure.md)** | Deployment (Docker), Storage (Postgres/Loki), and Security config. |

## Related Documentation

- **[Decisions](../decisions/)**: Architecture Decision Records (ADRs).
- **[Planning](../planning/)**: Future trends and RFC drafts.
