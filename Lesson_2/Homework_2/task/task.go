package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrEmptyTitle  = errors.New("пустой заголовок")
	ErrAlreadyDone = errors.New("задача уже выполнена")
)

type Timestamps struct {
	CreatedAt time.Time
}

func (ts *Timestamps) TouchCreated() {
	ts.CreatedAt = time.Now()
}

type Task struct {
	Timestamps 
	title      string
	done       bool
}

func New(title string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTitle
	}
	t := &Task{title: title}
	t.TouchCreated()
	return t, nil
}

func (t Task) Label() string { return t.title }

func (t Task) IsDone() bool { return t.done }

func (t Task) CreatedAtString() string {
	return t.CreatedAt.Format("02.01.2006 15:04")
}

func (t *Task) Complete() error {
	if t.done {
		return ErrAlreadyDone
	}
	t.done = true
	return nil
}

func (t Task) String() string {
	mark := " "
	if t.done {
		mark = "x"
	}
	return fmt.Sprintf("[%s] %s", mark, t.title)
}

type Completer interface {
	Complete() error
}

func MarkDone(c Completer) error {
	return c.Complete()
}

var _ Completer = (*Task)(nil)
