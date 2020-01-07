package verify

import (
	"github.com/trrtly/esign/context"
)

// Verify struct
type Verify struct {
	Individual   *Individual
	Organization *Organization
}

// NewVerify init
func NewVerify(ctx *context.Context) *Verify {
	return &Verify{
		Individual:   NewIndividual(ctx),
		Organization: NewOrganization(ctx),
	}
}
