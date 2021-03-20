package async

import "context"

type Future interface {
	Await() interface{}
}

type future struct {
	await func(ctx context.Context) interface{}
}

func (f future) Await() interface{} {
	return f.await(context.Background())
}

// Executes the async function
func Invoke(f func() interface{}) Future {
	var result interface{}
	channel := make(chan struct{})

	go func() {
		defer close(channel)
		result = f()
	}()

	return future{
		await: func(ctx context.Context) interface{} {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-channel:
				return result
			}
		},
	}
}
