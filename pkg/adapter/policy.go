package adapter

import (
	"github.com/bhojpur/web/pkg/adapter/context"
	bcontext "github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
)

// PolicyFunc defines a policy function which is invoked before the controller handler is executed.
type PolicyFunc func(*context.Context)

// FindPolicy Find Router info for URL
func (p *ControllerRegister) FindPolicy(cont *context.Context) []PolicyFunc {
	pf := (*web.ControllerRegister)(p).FindPolicy((*bcontext.Context)(cont))
	npf := newToOldPolicyFunc(pf)
	return npf
}

func newToOldPolicyFunc(pf []web.PolicyFunc) []PolicyFunc {
	npf := make([]PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *context.Context) {
			f((*bcontext.Context)(c))
		})
	}
	return npf
}

func oldToNewPolicyFunc(pf []PolicyFunc) []web.PolicyFunc {
	npf := make([]web.PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *bcontext.Context) {
			f((*context.Context)(c))
		})
	}
	return npf
}

// Policy Register new policy in Bhojpur
func Policy(pattern, method string, policy ...PolicyFunc) {
	pf := oldToNewPolicyFunc(policy)
	web.Policy(pattern, method, pf...)
}
