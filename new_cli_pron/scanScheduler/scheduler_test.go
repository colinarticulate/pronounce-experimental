package scanScheduler

import (
	"testing"
)

func Test_scheduler(t *testing.T) {
	// GIVEN
	batchResults := make(map[batchId]*batchScan)
	pending := make(map[batchId][]PsScan)
	newScan := make(chan PsScan)
	closing := make(chan chan error)

	sch := Scheduler{
		"", "", "", "",
		batchResults,
		pending,
		[]PsScan{},
		newScan,
		4,
		make(chan bool),
		closing,
		make(chan string),
	}
	go sch.loop()

	// WHEN
	for i := 0; i < 10; i++ {
	}

	// THEN

}
