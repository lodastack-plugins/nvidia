package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/lodastack/models"
)

var ms []models.Metric

func main() {
	nvml.Init()
	defer nvml.Shutdown()

	count, err := nvml.GetDeviceCount()
	if err != nil {
		fmt.Printf("Error getting device count:", err)
		return
	}

	var devices []*nvml.Device
	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			fmt.Printf("Error getting device %d: %v\n", i, err)
			return
		}
		devices = append(devices, device)
	}

	valueMap := make(map[string]int)
	for i, device := range devices {
		st, err := device.Status()
		if err != nil {
			fmt.Printf("Error getting device %d status: %v\n", i, err)
			return
		}
		valueMap["power_W"] = *st.Power
		valueMap["temperature_C"] = *st.Temperature
		valueMap["utilization.GPU"] = *st.Utilization.GPU
		valueMap["utilization.memory"] = *st.Utilization.Memory
		valueMap["utilization.encoder"] = *st.Utilization.Encoder
		valueMap["utilization.decoder"] = *st.Utilization.Decoder
		valueMap["clocks.memory_MHz"] = *st.Clocks.Memory
		valueMap["clocks.cores_MHz"] = *st.Clocks.Cores

		tags := make(map[string]string)
		tags["idx"] = i

		ts := time.Now().UnixNano()
		for k, v := range valueMap {
			var m models.Metric
			m.Name = k
			m.Value = v
			m.Tags = tags
			m.Timestamp = ts
			ms = append(ms, m)
		}
	}

	data, err := json.Marshal(ms)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", data)
}
