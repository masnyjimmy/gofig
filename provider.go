package gofig

type Provider interface {
	Source() (Source, error)
}
