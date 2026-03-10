package orm

type Observer interface {
	// Created called when the model has been created.
	Created(Event) error
	// Updated called when the model has been updated.
	Updated(Event) error
	// Deleted called when the model has been deleted.
	Deleted(Event) error
	// ForceDeleted called when the model has been force deleted.
	ForceDeleted(Event) error
}

type ObserverWithCreating interface {
	// Creating called when the model is being created.
	Creating(Event) error
}

type ObserverWithDeleting interface {
	// Deleting called when the model is being deleted.
	Deleting(Event) error
}

type ObserverWithForceDeleting interface {
	// ForceDeleting called when the model is being force deleted.
	ForceDeleting(Event) error
}

type ObserverWithRestored interface {
	// Restored called when the model has been restored.
	Restored(Event) error
}

type ObserverWithRestoring interface {
	// Restoring called when the model is being restored.
	Restoring(Event) error
}

type ObserverWithRetrieved interface {
	// Retrieved called when the model is retrieved from the database.
	Retrieved(Event) error
}

type ObserverWithSaved interface {
	// Saved called when the model has been saved.
	Saved(Event) error
}

type ObserverWithSaving interface {
	// Saving called when the model is being saved.
	Saving(Event) error
}

type ObserverWithUpdating interface {
	// Updating called when the model is being updated.
	Updating(Event) error
}

type ModelToObserver struct {
	Model    any
	Observer Observer
}
