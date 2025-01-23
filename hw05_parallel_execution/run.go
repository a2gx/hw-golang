package hw05parallelexecution

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 {
		return errors.New("n must be greater than 0")
	}

	// канал с задачами
	taskChan := make(chan Task)

	// обработка ошибок, при m <= 0 игнорируем ошибки в принципе
	var errCount int32
	var mInt32 int32

	if m > math.MaxInt32 {
		m = math.MaxInt32 // литер ругается, тут можно кинуть ошибку
	}

	checkErrorLimit := func() bool {
		return m > 0 && atomic.LoadInt32(&errCount) >= mInt32
	}

	// запускаем воркеры
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for task := range taskChan {
				// пропускаем пустую задачу
				if task == nil {
					continue
				}
				// обрабатываем ошибку
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	// передаем задачи в канал
	for _, task := range tasks {
		// останавливаем передачу задач при превышении лимита ошибок
		if checkErrorLimit() {
			break
		}
		taskChan <- task
	}

	// закрываем канал и дожидаемся завершения всех воркеров
	close(taskChan)
	wg.Wait()

	// проверяем превышение лимита ошибок
	if checkErrorLimit() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
