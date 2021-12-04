//go:generate mockery --name=Page
package domain

type Page interface {
	GetTitle() string
	GetLinks() []string
}
