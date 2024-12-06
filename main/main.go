package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Sensor struct {
	id   int
	data int
}

type SensorNetwork struct {
	mu            sync.Mutex
	sensors       map[int]*Sensor
	collectedData []int
}

var network *SensorNetwork
var once sync.Once

// InitializeNetwork initializes the sensor network (singleton).
func InitializeNetwork() {
	network = &SensorNetwork{
		sensors:       make(map[int]*Sensor),
		collectedData: []int{},
	}
	fmt.Println("Sensor network initialized.")
}

// GetNetworkInstance retrieves the singleton instance of the sensor network.
func GetNetworkInstance() *SensorNetwork {
	once.Do(func() {
		network = &SensorNetwork{
			sensors:       make(map[int]*Sensor),
			collectedData: []int{},
		}
	})
	return network
}

// AddSensor adds a new sensor to the network.
func (sn *SensorNetwork) AddSensor(id int) {
	sn.mu.Lock()
	defer sn.mu.Unlock()
	if _, exists := sn.sensors[id]; !exists {
		sn.sensors[id] = &Sensor{id: id}
		fmt.Printf("Sensor %d added to the network.\n", id)
	}
}

// UpdateData allows sensors to push data to the network.
func (sn *SensorNetwork) UpdateData(sensorID, data int) {
	sn.mu.Lock()
	defer sn.mu.Unlock()
	if sensor, exists := sn.sensors[sensorID]; exists {
		sensor.data = data
		sn.collectedData = append(sn.collectedData, data)
		fmt.Printf("Sensor %d reported data: %d\n", sensorID, data)
	}
}

// Monitor continuously collects data from all sensors until a stop signal is received.
func (sn *SensorNetwork) Monitor(stop chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stop:
			fmt.Println("Stopping monitoring...")
			return
		default:
			time.Sleep(2 * time.Second)
			sn.mu.Lock()
			for id, sensor := range sn.sensors {
				fmt.Printf("Sensor %d Current data: %d\n", id, sensor.data)
			}
			sn.mu.Unlock()
		}
	}
}

func main() {
	once.Do(InitializeNetwork)
	network := GetNetworkInstance()

	wg := sync.WaitGroup{}
	stop := make(chan struct{})

	// Capture OS signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Goroutine to handle the interrupt signal
	go func() {
		select {
		case <-sig:
			fmt.Println("\nGraceful shutdown initiated...")
			close(stop) // Signal all goroutines to stop
		case <-time.After(1 * time.Minute): // Automatically stop after 1 minute
			fmt.Println("\n1-minute timeout reached. Shutting down gracefully...")
			close(stop)
		}
	}()

	// Simulate 5 sensors
	for i := 1; i <= 5; i++ {
		network.AddSensor(i)
		wg.Add(1)
		go func(sensorID int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					fmt.Printf("Sensor %d stopping...\n", sensorID)
					return
				default:
					data := rand.Intn(100) // Random sensor data
					network.UpdateData(sensorID, data)
					time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
				}
			}
		}(i)
	}

	// Start monitoring data
	wg.Add(1)
	go network.Monitor(stop, &wg)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Printf("Final Collected Data: %v\n", network.collectedData)
	fmt.Println("Program exited cleanly.")
}
