# Dock Runner

`Dock Runner` is a small orchestrator written in **Go** that manages
**self-hosted GitHub Actions runners** on a single machine using **Docker**.

The goal is to provide a simple, observable and extensible way to:

- Define the desired state of containers (replicas, image, config).
- Continuously reconcile this desired state with the real Docker state.
- React to configuration changes without restarting the service.

---

## High-level Architecture

The system is composed of three core components:

- **ControlPlane**
- **ConfigAgent**
- **SchedulerAgent**

Plus supporting pieces:

- **In-memory database** (cluster state)
- **Docker Engine**
- **Config file** (source of truth for desired state)

```text
+---------------------------------------------------------+
|                        Host                             |
|                                                         |
|   +-------------+            +----------------------+   |
|   | ControlPlane| <--------> |   ConfigAgent        |   |
|   +-------------+   gRPC     +----------------------+   |
|         |                                               |
|         |                                               |
|         v                                               |
|   +-------------------+       +---------------------+   |
|   |  SchedulerAgent   |<----->| In-memory database  |   |
|   +-------------------+       +---------------------+   |
|      ^           |                                      |
|      |           |                                      |
| Watches status   | Increase/decrease containers         |
|      |           v                                      |
|   +--------------------+                                |
|   |      Docker        |  (Container-1, Container-2..)  |
|   +--------------------+                                |
+--------------------------------------------------------+

Filesystem:
+------------------------+
|      Config file       |
+------------------------+
        ^
        |
   watched by
   ConfigAgent
````

### Components

#### ControlPlane

* Exposes the main **gRPC API** for external clients (CLI, UI, other services).
* Validates and stores the **desired state** for managed workloads:

  * Desired number of replicas per service.
  * Container image, environment variables, restart policy, etc.
* Triggers scaling procedures (increase/decrease of containers) by calling the **SchedulerAgent**.
* Provides basic introspection endpoints (list services, list containers, get status, etc).

#### ConfigAgent

* Watches the **configuration file** in the filesystem.
* On file changes:

  * Parses and validates the new configuration.
  * Sends update commands to the **ControlPlane** using gRPC.
* Makes it possible to manage `Dock Runner` declaratively just by editing a file.

#### SchedulerAgent

* Responsible for the **reconciliation loop** between desired state and real Docker state.
* Functions:

  * Read the **desired state** from the in-memory database.
  * Read container information from Docker (running containers, status, etc).
  * Start/stop containers to match the desired number of replicas.
  * Watch container status and report changes back to the in-memory database / ControlPlane.
* Interacts directly with **Docker** using the Docker API/SDK for Go.

#### In-memory database

* Lightweight state store (e.g. based on Go structures or an embedded KV) containing:

  * Registered services.
  * Desired replicas.
  * Currently running containers and their metadata.
* Used by:

  * ControlPlane (to read/write state).
  * SchedulerAgent (to read/write state during reconciliation).

---

## Technology Stack

* **Language**: Go
* **Communication**: gRPC (with Protocol Buffers)
* **Containers**: Docker Engine
* **Config**: Toml file
* **Build/Distribution**: Go modules, Makefile (planned), Docker images (planned)
* **Testing**: Go test
