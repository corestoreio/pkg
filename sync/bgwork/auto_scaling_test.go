package bgwork_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

func TestRun(t *testing.T) {
	chanEvents := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	var lastStatistics string
	go func() {
		bgwork.AutoScaling(ctx, chanEvents, func(event interface{}) {
			time.Sleep(randInt(3, 9) * 100 * time.Millisecond)
			// fmt.Printf("%#v\n", event)
		}, bgwork.ScalingOptions{
			WorkerCheckInterval: 1000 * time.Millisecond,
			GetStatistics: func(s string) {
				println(s)
				lastStatistics = s
			},
		})
	}()
	for i := 0; i < 100; i++ {
		chanEvents <- fmt.Sprintf("Event_%03d", i)
		if i%20 == 0 {
			time.Sleep(300 * time.Millisecond)
		}
	}
	time.Sleep(4 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
	assert.Regexp(t, `Total Init 10
Total WIP 100
Total Waiting 110
Total Terminate 9
Total Terminated 8
WrkrID:[0-9] State: idle
WrkrID:[0-9] State: idle`+"\n", lastStatistics)
}

func randInt(min int, max int) time.Duration {
	return time.Duration(min) + time.Duration(rand.Intn(max-min))
}
