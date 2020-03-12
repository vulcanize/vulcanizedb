package core

type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}
