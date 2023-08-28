package interfaces

//go:generate mockgen --build_flags=--mod=mod -package=interfaces -destination=storage_mock.go . Storage
type Storage interface {
	GetOrSetID
	GetterFullURL
}

type GetOrSetID interface {
	GetID(fullURL string) (string, bool)
	SetPair(shortURL, fullURL string) error
}

type GetterFullURL interface {
	GetFullURL(shortURL string) (string, bool)
}
