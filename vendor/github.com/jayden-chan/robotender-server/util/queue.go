package util

// QueueItem represents one drink in the queue
type QueueItem struct {
	Name      string                 `json:"name"`
	QueueID   string                 `json:"queueID"`
	Glass     string                 `json:"glass"`
	Tall      bool                   `json:"tall"`
	Spec      map[string]interface{} `json:"spec"`
	SpecClean map[string]interface{} `json:"spec_clean"`
	Note      string                 `json:"note"`
}

type queue struct {
	Items  []QueueItem `json:"queue"`
	Status string      `json:"status"`
}

// push adds a drink to the queue
func (q *queue) push(drink QueueItem) {
	q.Items = append(q.Items, drink)
}

// pop returns the first item in the queue
func (q *queue) pop() (drink QueueItem) {
	drink = q.Items[0]
	q.Items = q.Items[1:]
	return
}

// get returns the item at the specified index
func (q *queue) get(index int) QueueItem {
	return q.Items[index]
}

// delete removes the item at the specified index
func (q *queue) remove(index int) {
	q.Items = append(q.Items[:index], q.Items[index+1:]...)
}

// size returns the size of the queue
func (q *queue) size() int {
	return len(q.Items)
}

// isEmpty returns whether there are any items in the queue
func (q *queue) isEmpty() bool {
	return len(q.Items) > 0
}

// clear removes all items from the queue
func (q *queue) clear() {
	q.Items = nil
}

// setStatus sets the status
func (q *queue) setStatus(status string) {
	q.Status = status
}
