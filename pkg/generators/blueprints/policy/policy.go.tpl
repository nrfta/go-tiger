package policy

import (
	"context"

	am "github.com/nrfta/go-access-management/v4/pkg/access_management"
	"github.com/nrfta/go-jwt-middleware/v4/pkg/resolver"
)

// {{.Name}}Update checks context user permissions to create the identified {{ToLowerCamel .Name}}.
func {{.Name}}Create(
	ctx context.Context,
	{{ToLowerCamel .Name}}ID string,
) error {
	return {{ToLowerCamel .Name}}Create_v0_2_25(ctx, {{ToLowerCamel .Name}}ID)
}

func {{ToLowerCamel .Name}}Create_v0_2_25(
	ctx context.Context,
	{{ToLowerCamel .Name}}ID string,
) error {
	// Check PolicyID: xx-x
	if err := resolver.IsAuthorized(
		ctx,
		GlobalResources.{{.NamePlural}},
		am.Actions.Create,
	); err != nil {
		return err
	}

	// TODO: Add identified policies if any

	return nil
}

// {{.Name}}Update checks context user permissions to udpate the identified {{ToLowerCamel .Name}}.
func {{.Name}}Update(ctx context.Context, {{ToLowerCamel .Name}}ID string) error {
	res := am.NewIdentifiedResource(
		ResourceNamespaces.{{.NamePlural}},
		{{ToLowerCamel .Name}}ID,
	)

	// Check PolicyID: xx-x
	if err := resolver.IsAuthorized(ctx, res, am.Actions.Update); err != nil {
		return err
	}
	return nil
}

// {{.Name}}Delete checks context user permissions to delete the identified {{ToLowerCamel .Name}}.
func {{.Name}}Delete(ctx context.Context, {{ToLowerCamel .Name}}ID string) error {
	res := am.NewIdentifiedResource(ResourceNamespaces.{{.NamePlural}}, {{ToLowerCamel .Name}}ID)

	// Check PolicyID: xx-x
	if err := resolver.IsAuthorized(ctx, res, am.Actions.Delete); err != nil {
		return err
	}

	// TODO: remove identified policies if any
}

func applyGlobalPoliciesFor{{.NamePlural}}() error {
	if err := applyGlobalPoliciesFor{{.NamePlural}}Users(); err != nil {
		return err
	}

	if err := applyGlobalPoliciesFor{{.NamePlural}}Admins(); err != nil {
		return err
	}

	return nil
}

func applyGlobalPoliciesFor{{.NamePlural}}Users() error {
	// PolicyID `xx-x`: All support users can read {{ToLowerCamel .NamePlural}}
	if err := am.AddPolicyForRole(
		am.Roles.User,
		IdentifiedDomains.AllSupportLevels,
		GlobalResources.{{.NamePlural}},
		am.Actions.Read,
	); err != nil {
		return err
	}

	// PolicyID `xx-x`: Support level 3 users can create, update and delete {{ToLowerCamel .NamePlural}}
	if err := am.AddPolicyForRole(
		am.Roles.User,
		IdentifiedDomains.SupportLevel3,
		GlobalResources.{{.NamePlural}},
		am.Actions.Create,
		am.Actions.Update,
		am.Actions.Delete,
	); err != nil {
		return err
	}

	return nil
}

func applyGlobalPoliciesFor{{.NamePlural}}Admins() error {
	return nil
}
