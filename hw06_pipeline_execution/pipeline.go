package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		if stage != nil {
			in = process(stage(in), done)
		}
	}
	return in
}

func process(in In, done In) Out {
	out := make(Bi, 1) // буфер для избежания блокировок

	go func() {
		defer func() {
			for range in {
				// читаем остатки горутин при резком сигнале <-done
				_ = struct{}{}
			}
		}()
		defer close(out)

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- val:
				case <-done:
					// прерываем передачу, если уже закрыт done
					return
				}
			}
		}
	}()

	return out
}
