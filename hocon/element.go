package hocon

type MightBeAHoconObject interface {
	// IsObject checks if the current item is object or not.
	// It returns false if the item has cycled reference
	IsObject() bool
	// GetObject returns HoconObject if the item contains one,
	// error - if it cannot be found.
	GetObject() (*HoconObject, error)
}

type HoconElement interface {
	IsString() bool
	GetString() (string, error)
	IsArray() bool
	GetArray() ([]*HoconValue, error)
}
