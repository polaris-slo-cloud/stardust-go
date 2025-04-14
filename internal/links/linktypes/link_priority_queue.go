package linktypes

import (
	"container/heap"
)

// LinkPriorityQueue is a min-heap of IslLinks based on Distance
type LinkPriorityQueue struct {
	items []*linkItem
	index map[*IslLink]int
}

type linkItem struct {
	link     *IslLink
	priority float64
}

func NewLinkPriorityQueue() *LinkPriorityQueue {
	return &LinkPriorityQueue{
		items: []*linkItem{},
		index: make(map[*IslLink]int),
	}
}

func (pq *LinkPriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *LinkPriorityQueue) Less(i, j int) bool {
	return pq.items[i].priority < pq.items[j].priority
}

func (pq *LinkPriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.index[pq.items[i].link] = i
	pq.index[pq.items[j].link] = j
}

func (pq *LinkPriorityQueue) Push(x interface{}) {
	item := x.(*linkItem)
	pq.index[item.link] = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *LinkPriorityQueue) Pop() interface{} {
	n := len(pq.items)
	item := pq.items[n-1]
	pq.items[n-1] = nil
	pq.items = pq.items[:n-1]
	delete(pq.index, item.link)
	return item
}

// Public API

// Push adds a link with its priority (distance)
func (pq *LinkPriorityQueue) Enqueue(link *IslLink, priority float64) {
	if _, exists := pq.index[link]; exists {
		return // already present
	}
	heap.Push(pq, &linkItem{link: link, priority: priority})
}

// Pop removes and returns the link with the lowest priority (distance)
func (pq *LinkPriorityQueue) Dequeue() *IslLink {
	if pq.Len() == 0 {
		return nil
	}
	item := heap.Pop(pq).(*linkItem)
	return item.link
}

// Clear resets the queue
func (pq *LinkPriorityQueue) Clear() {
	pq.items = []*linkItem{}
	pq.index = make(map[*IslLink]int)
}

// Init prepares the heap for use
func (pq *LinkPriorityQueue) Init() {
	heap.Init(pq)
}
