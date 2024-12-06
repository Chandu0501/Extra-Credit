// File: main_test.go
package main

import (
	"sync"
	"testing"
)

func TestSensorNetworkInitialization(t *testing.T) {
	network1 := GetNetworkInstance()
	network2 := GetNetworkInstance()

	if network1 != network2 {
		t.Error("Expected the same instance for the sensor network.")
	}
}

func TestAddAndUpdateSensors(t *testing.T) {
	var once sync.Once
	once.Do(InitializeNetwork)

	network := GetNetworkInstance()

	network.AddSensor(1)
	network.UpdateData(1, 42)

	data := 0
	network.mu.Lock()
	if sensor, exists := network.sensors[1]; exists {
		data = sensor.data
	}
	network.mu.Unlock()

	if data != 42 {
		t.Errorf("Expected sensor data 42, got %d", data)
	}
}
