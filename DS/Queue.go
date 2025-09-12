package DataStructures

// ~ File Description = Implements a basic fixed size Queue Data Structure.

type Queue[T any] struct {
	buffer []T
	head   int
	tail   int
	size   int
	cap    int
}

func NewQueue[T any](n int) *Queue[T] {
	return &Queue[T]{
		buffer: make([]T, n),
		cap:    n,
	}
}

func (q *Queue[T]) Enqueue(item T) bool {
	if q.size == q.cap {
		return false // Queue is full
	}
	q.buffer[q.tail] = item
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return true
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false // Queue is empty
	}
	item := q.buffer[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.buffer[q.head], true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) IsFull() bool {
	return q.size == q.cap
}

func (q *Queue[T]) Len() int {
	return q.size
}
