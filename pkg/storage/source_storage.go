package storage

import "github.com/antonchaban/news-aggregator/pkg/model"

// SourceStorage is an interface that defines the methods for interacting with the source storage.
//
// Key Responsibilities:
//
// 1. GetAll retrieve all stored sources.
//
// 2. Save a single source to the storage.
//
// 3. Save multiple sources to the storage.
//
// 4. Delete a source from the storage by its ID.
//
// 5. GetByID retrieve a source by its ID.
//
// Expected Behaviors or Guarantees:
// - The GetAll method should return all sources available in the storage or an error if the operation fails.
//
// - The Save method should store the provided source and return the saved source (with updated fields such as ID) or an error if the operation fails.
//
// - The SaveAll method should store all provided sources and return an error if the operation fails for any reason.
//
// - The Delete method should remove the source with the specified ID from the storage and return an error if the operation fails.
//
// - The GetByID method should return the source with the specified ID or an error if the source does not exist or the operation fails.
//
// Common Errors or Exceptions and Handling:
// - `error`: This general error can occur in any of the methods. It should be handled by logging the error and returning an appropriate message to the user or retrying the operation if possible.
//
// - Data validation errors: Validate input data before attempting to save or retrieve sources. For example, check if the `id` in Delete or GetByID methods is a valid positive integer.
//
// Known Limitations or Restrictions:
// - The Delete method does not specify what happens if the source does not exist. It should be clarified whether it returns an error or silently succeeds.
// - The methods do not define any constraints on the size or format of the sources being saved.
//
// Usage Guidelines or Best Practices:
// - Use transactions where necessary to ensure data consistency, for example, when saving multiple sources with SaveAll.
// - Ensure proper error handling and logging to facilitate debugging and monitoring of storage operations.
// - Validate input data thoroughly before performing any operations to prevent injection attacks or corrupt data.
type SourceStorage interface {
	GetAll() ([]model.Source, error)
	Save(src model.Source) (model.Source, error)
	SaveAll(sources []model.Source) error
	Delete(id int) error
	GetByID(id int) (model.Source, error)
}
