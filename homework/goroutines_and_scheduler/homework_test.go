package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks   []Task
	posByID map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{
		tasks:   []Task{},
		posByID: map[int]int{},
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)

	i := len(s.tasks) - 1
	s.posByID[s.tasks[i].Identifier] = i

	_ = s.sieveUp(i)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	idx, ok := s.posByID[taskID]
	if !ok {
		return
	}

	s.tasks[idx].Priority = newPriority

	i := idx
	i = s.sieveUp(i)

	s.sieveDown(i)
}

func (s *Scheduler) GetTask() Task {
	if len(s.tasks) == 0 {
		return Task{}
	}

	result := s.tasks[0]

	s.tasks[0] = s.tasks[len(s.tasks)-1]
	s.tasks = s.tasks[:len(s.tasks)-1]

	i := 0
	s.sieveDown(i)

	return result
}

func (s *Scheduler) sieveUp(i int) int {
	for i > 0 {
		parent := (i - 1) / 2
		if s.tasks[parent].Priority < s.tasks[i].Priority {
			s.posByID[s.tasks[i].Identifier] = parent
			s.posByID[s.tasks[parent].Identifier] = i
			s.tasks[parent], s.tasks[i] = s.tasks[i], s.tasks[parent]
		}

		i = parent
	}

	return i
}

func (s *Scheduler) sieveDown(i int) {
	l, r := 2*i+1, 2*i+2

	for l < len(s.tasks) {
		var c int
		if l < len(s.tasks) {
			c = l

			if r < len(s.tasks) {
				if s.tasks[l].Priority < s.tasks[r].Priority {
					c = r
				}
			}
		}

		if s.tasks[i].Priority < s.tasks[c].Priority {
			s.posByID[s.tasks[i].Identifier] = c
			s.posByID[s.tasks[c].Identifier] = i
			s.tasks[i], s.tasks[c] = s.tasks[c], s.tasks[i]
		}

		i = c
		l, r = 2*i+1, 2*i+2
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1.Identifier, task.Identifier) // кажется ждать старый приоритет, изменив его выше, не вполне правильно

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
