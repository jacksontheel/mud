package scheduler

import "container/heap"

type JobHeap []*Job

var _ heap.Interface = (*JobHeap)(nil)

func (h JobHeap) Len() int {
	return len(h)
}

// the job with the earliest NextRun time is considered "less".
func (h JobHeap) Less(i, j int) bool {
	return h[i].NextRun.Before(h[j].NextRun)
}

func (h JobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *JobHeap) Push(x interface{}) {
	job := x.(*Job)
	*h = append(*h, job)
}

func (h *JobHeap) Pop() interface{} {
	old := *h
	n := len(old)

	// Get the last element.
	job := old[n-1]

	// Shrink the slice.
	*h = old[:n-1]

	return job
}
