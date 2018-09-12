package util

import (
	"fmt"
)

// DrinkQueue stores the drink queue
type DrinkQueue struct {
	lanes [4]queue
}

// NewDrinkQueue returns a new DrinkQueue
func NewDrinkQueue() *DrinkQueue {
	ret := DrinkQueue{}
	ret.init()
	return &ret
}

// init initializes the queue
func (d *DrinkQueue) init() {
	for i := 0; i < 4; i++ {
		d.lanes[i].Items = []QueueItem{}
	}
	d.SetAllStatus("READY")
}

// Push adds a drink to the lane with the least items
func (d *DrinkQueue) Push(drink QueueItem) {
	lowest := 0
	for i, q := range d.lanes {
		if q.size() < d.lanes[lowest].size() {
			lowest = i
		}
	}

	d.lanes[lowest].push(drink)
}

// GetAll returns a list of all queue items
func (d *DrinkQueue) GetAll() (list queue) {
	for i := 0; i < len(d.lanes); i++ {
		for j := 0; j < d.lanes[i].size(); j++ {
			list.push(d.lanes[i].get(j))
		}
	}

	return
}

// Get returns the queue for the specified lane
func (d *DrinkQueue) Get(lane int) (list queue, err error) {
	if lane < 4 {
		return d.lanes[lane], nil
	}
	return queue{}, fmt.Errorf("Index out of bounds")
}

// SetStatus sets the lane status
func (d *DrinkQueue) SetStatus(lane int, status string) error {
	if lane < 4 {
		d.lanes[lane].setStatus(status)
		return nil
	}
	return fmt.Errorf("Index out of bounds")
}

// SetAllStatus sets all lane statuses
func (d *DrinkQueue) SetAllStatus(status string) {
	for i := 0; i < len(d.lanes); i++ {
		d.lanes[i].setStatus(status)
	}
}

// Remove removes the specified item from the correct lane
func (d *DrinkQueue) Remove(queueID string) bool {
	for i := 0; i < len(d.lanes); i++ {
		for j := 0; j < d.lanes[i].size(); j++ {
			if d.lanes[i].get(j).QueueID == queueID {
				d.lanes[i].remove(j)
				return true
			}
		}
	}

	return false
}

// Clear removes all items from the queue
func (d *DrinkQueue) Clear() {
	for i := 0; i < len(d.lanes); i++ {
		for j := d.lanes[i].size() - 1; j >= 0; j-- {
			d.lanes[i].remove(j)
		}
	}
}

// Size returns the size of the whole drink queue
func (d *DrinkQueue) Size() (sum int) {
	for _, q := range d.lanes {
		sum += q.size()
	}
	return
}

// LaneSize returns the size of an individual lane queue
func (d *DrinkQueue) LaneSize(lane int) (size int) {
	if lane >= 0 && lane < 4 {
		return d.lanes[lane].size()
	}
	return 0
}
