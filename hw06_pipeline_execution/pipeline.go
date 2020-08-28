package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
	I   interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		bi := make(Bi)
		go func(done In, out Out) {
			defer close(bi)
			for {
				select {
				case <-done:
					return
				case v, ok := <-out:
					if !ok {
						return
					}
					bi <- v
				default:
				}
			}
		}(done, out)
		out = stage(bi)
	}
	return out
}
