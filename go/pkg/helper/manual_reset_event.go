package helper

// ManualResetEvent mimics the behavior of C#'s ManualResetEvent.
type ManualResetEvent struct {
	ch    chan struct{}
	isSet bool
}

// NewManualResetEvent creates a new ManualResetEvent.
// If initialState is true, the event starts in the set state.
func NewManualResetEvent(initialState bool) *ManualResetEvent {
	mre := &ManualResetEvent{
		ch:    make(chan struct{}),
		isSet: initialState,
	}
	// If initialState is true, close the channel to unblock Wait()
	if initialState {
		close(mre.ch)
	}
	return mre
}

// Set signals the event, allowing any waiting goroutines to proceed.
func (mre *ManualResetEvent) Set() {
	if !mre.isSet {
		mre.isSet = true
		close(mre.ch) // Closing the channel unblocks all Wait() calls
	}
}

// Wait blocks until the event is set.
func (mre *ManualResetEvent) Wait() {
	// If the event is set, the channel is closed and Wait() returns immediately.
	// Otherwise, it blocks until the channel is closed.
	_, ok := <-mre.ch
	if !ok {
		// Channel is closed, event is set
		return
	}
	// If we get here, the channel was not closed (event not set)
	// But since we use close() to signal, this branch is theoretically unreachable
	// because once closed, all receives return immediately.
}

// Reset clears the event, causing future Wait calls to block.
func (mre *ManualResetEvent) Reset() {
	if mre.isSet {
		mre.isSet = false
		mre.ch = make(chan struct{}) // Recreate the channel
	}
}
