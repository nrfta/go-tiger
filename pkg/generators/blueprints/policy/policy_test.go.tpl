package policy_test

import (
	am "github.com/nrfta/go-access-management/v4/pkg/access_management"
	"github.com/nrfta/go-platform-security-policy/pkg/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} policies", func() {
	Context("Global", func() {
		Context("USER", func() {
			It("PolicyID: xx-x, xx-x", func() {
				resource := am.NewIdentifiedResource(
					policy.ResourceNamespaces.{{.NamePlural}},
					{{ToLowerCamel .Name}}ID,
				)

				expectAccessOnlyBy(resource, []accessBy{
					domainRoleActions(policy.IdentifiedDomains.SupportLevel1, am.Roles.User, am.Actions.Read),
					domainRoleActions(policy.IdentifiedDomains.SupportLevel2, am.Roles.User, am.Actions.Read),
					domainRoleActions(
						policy.IdentifiedDomains.SupportLevel3,
						am.Roles.User,
						am.Actions.Read,
						am.Actions.Create,
						am.Actions.Update,
						am.Actions.Delete,
					),
				})
			})
		})
	})

	Context("Identified", func() {
		Context("USER", func() {
//			It("PolicyID: xx-x, xx-x, xx-x", func() {
//				resource := am.NewIdentifiedResource(
//					policy.ResourceNamespaces.{{.NamePlural}},
//					{{ToLowerCamel .Name}}ID,
//				)
//
//				ctxSupportLevel3 := authenticateUser(user3, policy.IdentifiedDomains.SupportLevel3, am.Roles.User)
//				err := policy.{{.Name}}Create(ctxSupportLevel3, {{ToLowerCamel .Name}}ID)
//				Expect(err).To(BeNil())
//
//				authenticateUser(user1, policy.IdentifiedDomains.SupportLevel3, am.Roles.User)
//				authenticateUser(user2, policy.NamedDomains.User, am.Roles.User)
//
//				verifyAccess(user1, resource, am.Actions.Read)
//				verifyNoAccess(user2, resource, am.Actions.Any)
//
//				err = policy.{{.Name}}Delete(ctxSupportLevel3, {{ToLowerCamel .Name}}ID)
//				Expect(err).To(BeNil())
//			})
		})
	})
})
