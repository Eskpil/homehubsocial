package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	HeartbeatEntityId              string = "sensor.heartbeat"
	BloodPressureSystolicEntityId  string = "sensor.bloodpressure_systolic"
	BloodPressureDiastolicEntityId string = "sensor.bloodpressure_diastolic"

	AccessToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3ZjgzMjZiYTMwMGI0NGRiYmQyZmVhMmExOTc1YjE0MSIsImlhdCI6MTY3NjM4OTQxMSwiZXhwIjoxOTkxNzQ5NDExfQ.CR2yJ-sg3NJIKY2kE7ESGKMQk_4fVFf6PqkPwBMCGsE"

	HeartbeatMin = 0.5
	HeartbeatMax = 2.5
)

type HeartBeatRequest struct {
	Value      string            `json:"state"`
	Attributes map[string]string `json:"attributes"`
}

type BloodPressureRequest struct {
	Value      string            `json:"state"`
	Attributes map[string]string `json:"attributes"`
}

type LastReportedHeartbeat struct {
	Lock sync.Mutex
	Bpm  uint64
}

func simulateBloodPressure(heartBeat uint64) (float64, float64) {
	cardiacOutput := 5.0         // liters per minute
	peripheralResistance := 20.0 // units unknown
	bloodVolume := 5.0           // liters

	MAP := cardiacOutput * peripheralResistance

	// Generate random variation in MAP
	rand.Seed(time.Now().UnixNano())
	variation := rand.Float64() * 10.0 // varies between 0 and 10 mmHg
	MAP += variation

	// Calculate systolic and diastolic pressures
	systolic := MAP + cardiacOutput/3.0 + rand.Float64()*10.0
	diastolic := MAP - cardiacOutput/3.0 - rand.Float64()*10.0

	cardiacOutput = bloodVolume * float64(heartBeat) / 1000.0

	_ = diastolic

	// Print out the simulated blood pressure readings
	return math.Round(systolic*100) / 100, math.Round(diastolic*100) / 100
}

func main() {
	client := resty.New()
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", AccessToken))

	quit := make(chan struct{})

	var wg sync.WaitGroup

	lastReported := &LastReportedHeartbeat{
		Lock: sync.Mutex{},
		Bpm:  90,
	}

	wg.Add(1)
	go func(group *sync.WaitGroup) {
		defer group.Done()
		ticker := time.NewTicker(time.Second)

		count := 0
		beats := 0.0

		for {
			select {
			case <-ticker.C:
				count += 1

				rand.Seed(time.Now().UnixNano())
				beats += HeartbeatMin + rand.Float64()*(HeartbeatMax-HeartbeatMin)

				// beats += 1.16

				if count == 15 {
					beats = beats * 4

					heartBeat := HeartBeatRequest{
						Value:      fmt.Sprintf("%d", uint64(math.Floor(beats))),
						Attributes: map[string]string{},
					}
					heartBeat.Attributes["friendly_name"] = "Slag i minuttet"

					_, err := client.R().
						SetBody(heartBeat).
						Post(fmt.Sprintf("http://localhost:8123/api/states/%s", HeartbeatEntityId))
					if err != nil {
						log.Fatalf("Failed to send hearbeat: %v", err)
					}

					log.Infof("Reported %v beats per minute", uint64(math.Floor(beats)))

					lastReported.Lock.Lock()
					lastReported.Bpm = uint64(beats)
					lastReported.Lock.Unlock()

					count = 0
					beats = 0
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(&wg)

	wg.Add(2)
	go func(group *sync.WaitGroup) {
		defer group.Done()
		ticker := time.NewTicker(20 * time.Second)

		for {
			select {
			case <-ticker.C:
				lastReported.Lock.Lock()
				bpm := lastReported.Bpm
				lastReported.Lock.Unlock()

				systolic, diastolic := simulateBloodPressure(bpm)
				systolicRequest := BloodPressureRequest{
					Value:      fmt.Sprintf("%2.f", systolic),
					Attributes: map[string]string{},
				}
				systolicRequest.Attributes["unit_of_measurement"] = "%"
				systolicRequest.Attributes["friendly_name"] = "Systolisk"

				_, err := client.R().
					SetBody(systolicRequest).
					Post(fmt.Sprintf("http://localhost:8123/api/states/%s", BloodPressureSystolicEntityId))
				if err != nil {
					log.Fatalf("Failed to send systolic: %v", err)
				}

				diastolicRequest := BloodPressureRequest{
					Value:      fmt.Sprintf("%2.f", diastolic),
					Attributes: map[string]string{},
				}
				diastolicRequest.Attributes["unit_of_measurement"] = "%"
				diastolicRequest.Attributes["friendly_name"] = "Diastolisk"

				_, err = client.R().
					SetBody(diastolicRequest).
					Post(fmt.Sprintf("http://localhost:8123/api/states/%s", BloodPressureDiastolicEntityId))
				if err != nil {
					log.Fatalf("Failed to send diastolic: %v", err)
				}

				log.Infof("Reported blood pressure: %.2f/%2.f mmHg", systolic, diastolic)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(&wg)

	wg.Wait()
}
