package owl

import (
	"fmt"
	"strconv"
)

type TaskStatus byte

const (
	JobDone    TaskStatus = 0
	JobRunning TaskStatus = 1
	JobFailed  TaskStatus = 2
	JobTimeout TaskStatus = 3
)

func (sha TaskStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(sha))), nil
}
func (sha *TaskStatus) UnmarshalJSON(json []byte) error {
	v, err := strconv.Atoi(string(json))
	if err != nil {
		return err
	}

	statusValue := TaskStatus(v)

	switch statusValue {
	case JobDone, JobRunning, JobFailed, JobTimeout:
		*sha = statusValue
		return nil
	}

	return fmt.Errorf(`Cannot convert [%d] to "TaskStatus"`, v)
}
