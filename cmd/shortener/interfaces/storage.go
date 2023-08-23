package interfaces

//go:generate mockgen --build_flags=--mod=mod -package=interfaces -destination=creatorid_mock.go . CreatorID
type CreatorID interface {
	CreateID(data []byte, base string) (string, error)
}
