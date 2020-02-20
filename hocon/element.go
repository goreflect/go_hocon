package hocon

type MightBeAHoconObject interface {
	IsObject() bool
	GetObject() (*HoconObject, error)
}

type HoconElement interface {
	IsString() bool
	GetString() (string, error)
	IsArray() bool
	GetArray() ([]*HoconValue, error)
}
