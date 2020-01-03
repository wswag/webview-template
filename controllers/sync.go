package controllers

// Synchronizer interface serves functions for setting a js-sync function
// and executing it. It can be implemented if displayed data changes without
// interaction of the webserver. Call Synchronize() to update js model values then.
type Synchronizer interface {
	SetSyncFunction(sync func())
	Synchronize()
}

// SynchronizableModel is an embeddable type implementing the Synchronizer Interface
type SynchronizableModel struct {
	sync func()
}

// Synchronize will call the sync() method set by SetSyncFunction()
func (s *SynchronizableModel) Synchronize() {
	if (s.sync != nil) {
		s.sync()
	}
}

// SetSyncFunction sets the function called by Synchronize()
func (s *SynchronizableModel) SetSyncFunction(sync func()) {
	s.sync = sync
}