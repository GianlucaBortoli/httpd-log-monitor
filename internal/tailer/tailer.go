package tailer

import (
	"fmt"

	"github.com/hpcloud/tail"
)

// Tailer the file tailer
type Tailer struct {
	fileName string
	tailConf tail.Config
	tail     *tail.Tail
	started  bool
}

// New returns a tailer for the given file
func New(fileName string) *Tailer {
	return &Tailer{
		fileName: fileName,
		tailConf: tail.Config{
			MustExist: true, // Fail early if the file does not exist
			Follow:    true, // Continue looking for new lines (tail -f)
			ReOpen:    true, // Reopen recreated/truncated files (tail -F)
		},
	}
}

// Start starts the tailing process in a separate goroutine.
// Returns the lines channel and an error
func (t *Tailer) Start() (<-chan *tail.Line, error) {
	if t.started {
		return nil, fmt.Errorf("tailer can be started only once")
	}

	tf, err := tail.TailFile(t.fileName, t.tailConf)
	if err != nil {
		return nil, err
	}
	t.started = true
	t.tail = tf
	return tf.Lines, nil
}

// Stop stops the tailing process, gracefully exiting the background goroutines
func (t *Tailer) Stop() error {
	if !t.started {
		return fmt.Errorf("tailer can be stopped only after start")
	}
	return t.tail.Stop()
}

// Wait blocks until the tailer goroutine is in a dead state.
// Returns the reason for its death.
func (t *Tailer) Wait() error {
	if t.tail != nil && t.started {
		return t.tail.Wait()
	}
	return fmt.Errorf("tailer cannot wait if not started")
}
