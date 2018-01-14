package pool

import (
	"errors"
	"testing"
	"time"
)

func TestPoolSuccess(t *testing.T) {
	job := Job(func() error {
		time.Sleep(time.Second)
		return nil
	})

	pool := NewPool(2)

	pool.Add(job)
	pool.Add(job)
	pool.Add(job)
	pool.Add(job)

	startTime := time.Now()
	errs := pool.Run()

	duration := time.Since(startTime)

	if duration.Seconds() > 3 {
		t.Errorf(
			"Expected test to take less than 3 seconds (2 threads, 2 seconds per thread).  Took %v",
			duration.Seconds(),
		)
	}
	if len(errs) != 0 {
		t.Error("Expected no errors returned")
	}
}

func TestPoolErrors(t *testing.T) {
	errMessage := "Something went wrong"
	goodJob := Job(func() error {
		time.Sleep(time.Second)
		return nil
	})
	badJob := Job(func() error {
		time.Sleep(time.Second)
		return errors.New(errMessage)
	})

	pool := NewPool(2)

	pool.Add(goodJob)
	pool.Add(badJob)
	pool.Add(badJob)
	pool.Add(goodJob)

	startTime := time.Now()
	errs := pool.Run()

	duration := time.Since(startTime)

	if duration.Seconds() > 3 {
		t.Errorf(
			"Expected test to take less than 3 seconds (2 threads, 2 seconds per thread).  Took %v",
			duration.Seconds(),
		)
	}

	if len(errs) != 2 {
		t.Error("Incorrect number of errors returned")
	}
	for _, err := range errs {
		if err.Error() != errMessage {
			t.Error("Incorrect error message returned")
		}
	}
}
