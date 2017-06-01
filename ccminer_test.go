package ccminer

import (
	"fmt"
	"testing"
)

// WARNING: These tests are currently terrible, and require a setup such as mine
// (I'm connecting to a real cgminer instance running on a different machine).
// Once I figure out how to mock things out, these tests should improve substantially.
// For now they're more of just a convenient scratch area for manual testing.

func Test_Summary(t *testing.T) {
	miner := New("127.0.0.1", 4068)
	summary, err := miner.Summary()
	if err != nil {
		t.Error(err)
		return
	}
	if summary == nil {
		t.Error("Summary returned nil")
		return
	}
	// TODO: Make some assertions. Need to mock out the data source first?
}

func Test_Devs(t *testing.T) {
	miner := New("127.0.0.1", 4068)
	devs, err := miner.Devs()
	if err != nil {
		t.Error(err)
		return
	}
	if devs == nil {
		t.Error("Summary returned nil")
		return
	}
	for _, dev := range *devs {
		fmt.Printf("Dev %d temp: %f\n", dev.Id, dev.Temp)
	}
}

func Test_Pools(t *testing.T) {
	miner := New("127.0.0.1", 4068)
	pools, err := miner.Pools()
	if err != nil {
		t.Error(err)
		return
	}
	for _, pool := range pools {
		fmt.Printf("Pool %d: %s\n", pool.Name, pool.Url)
	}
}
