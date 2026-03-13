package gofig

type Source interface {
	Read(path string) (any, error)
}
