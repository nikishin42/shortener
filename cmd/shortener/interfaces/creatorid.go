package interfaces

//go:generate mockgen --build_flags=--mod=mod -package=interfaces -destination=storage_mock.go . Storage
type Storage interface {
	GetID(fullURL string) (string, bool)
	GetFullURL(shortURL string) (string, bool)
	SetPair(shortURL, fullURL string) error
}
