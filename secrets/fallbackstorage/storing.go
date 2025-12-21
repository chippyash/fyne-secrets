package fallbackstorage

// Storing interface
type Storing interface {
	// Store writes data
	Store(string, []byte) error
	// Load reads data
	Load(string) ([]byte, error)
	// Exists checks if the data exists
	Exists(string) (bool, error)
	// Delete deletes data
	Delete(string) error
}
