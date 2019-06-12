package peertask

import (
	"time"

	pq "github.com/ipfs/go-ipfs-pq"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// FIFOCompare is a basic task comparator that returns tasks in the order created.
var FIFOCompare = func(a, b *TaskBlock) bool {
	return a.created.Before(b.created)
}

// PriorityCompare respects the target peer's task priority. For tasks involving
// different peers, the oldest task is prioritized.
var PriorityCompare = func(a, b *TaskBlock) bool {
	if a.Target == b.Target {
		return a.Priority > b.Priority
	}
	return FIFOCompare(a, b)
}

// WrapCompare wraps a TaskBlock comparison function so it can be used as
// comparison for a priority queue
func WrapCompare(f func(a, b *TaskBlock) bool) func(a, b pq.Elem) bool {
	return func(a, b pq.Elem) bool {
		return f(a.(*TaskBlock), b.(*TaskBlock))
	}
}

// Identifier is a unique identifier for a task. It's used by the client library
// to act on a task once it exits the queue.
type Identifier interface{}

// Task is a single task to be executed as part of a task block.
type Task struct {
	Identifier Identifier
	Priority   int
}

// TaskBlock is a block of tasks to execute on a single peer.
type TaskBlock struct {
	Tasks    []Task
	Priority int
	Target   peer.ID

	// A callback to signal that this task block has been completed
	Done func([]Task)

	// toPrune are the tasks that have already been taken care of as part of
	// a different task block which can be removed from the task block.
	toPrune map[Identifier]struct{}
	created time.Time // created marks the time that the task was added to the queue
	index   int       // book-keeping field used by the pq container
}

// NewTaskBlock creates a new task block with the given tasks, priority, target
// peer, and task completion function.
func NewTaskBlock(tasks []Task, priority int, target peer.ID, done func([]Task)) *TaskBlock {
	return &TaskBlock{
		Tasks:    tasks,
		Priority: priority,
		Target:   target,
		Done:     done,
		toPrune:  make(map[Identifier]struct{}, len(tasks)),
		created:  time.Now(),
	}
}

// MarkPrunable marks any tasks with the given identifier as prunable at the time
// the task block is pulled of the queue to execute (because they've already been removed).
func (pt *TaskBlock) MarkPrunable(identifier Identifier) {
	pt.toPrune[identifier] = struct{}{}
}

// PruneTasks removes all tasks previously marked as prunable from the lists of
// tasks in the block
func (pt *TaskBlock) PruneTasks() {
	newTasks := make([]Task, 0, len(pt.Tasks)-len(pt.toPrune))
	for _, task := range pt.Tasks {
		if _, ok := pt.toPrune[task.Identifier]; !ok {
			newTasks = append(newTasks, task)
		}
	}
	pt.Tasks = newTasks
}

// Index implements pq.Elem.
func (pt *TaskBlock) Index() int {
	return pt.index
}

// SetIndex implements pq.Elem.
func (pt *TaskBlock) SetIndex(i int) {
	pt.index = i
}
