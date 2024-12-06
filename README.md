# Extra Credit Assignment: Sensor Network

This project simulates a distributed sensor network using Go. It demonstrates:
- **Concurrency**: Goroutines for simulating sensors.
- **Synchronization**: `sync.Mutex` for thread safety and `sync.Once` for one-time initialization.
- **Graceful Shutdown**: Proper cleanup on program termination using signal handling.

## Features
- A sensor network that collects random data from sensors.
- Real-time data monitoring.
- Automatic shutdown after 1 minute or upon manual interruption (Ctrl+C).

## Running the Program
```bash
go run main.go
