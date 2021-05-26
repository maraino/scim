package patch

import (
	"github.com/elimity-com/scim/errors"
	f "github.com/elimity-com/scim/internal/filter"
	"github.com/elimity-com/scim/schema"
	"net/http"
)

// validateRemove validates the remove operation contained within the validator based on on Section 3.5.2.2 in RFC 7644.
// More info: https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.2
func (v OperationValidator) validateRemove() error {
	// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
	if v.path == nil {
		return &errors.ScimError{
			ScimType: errors.ScimTypeNoTarget,
			Status:   http.StatusBadRequest,
		}
	}

	refAttr, err := v.getRefAttribute(v.path.AttributePath)
	if err != nil {
		return err
	}
	if v.path.ValueExpression != nil {
		if err := f.NewFilterValidator(v.path.ValueExpression, schema.Schema{
			Attributes: f.MultiValuedFilterAttributes(*refAttr),
		}).Validate(); err != nil {
			return err
		}
	}
	if subAttrName := v.path.SubAttributeName(); subAttrName != "" {
		if _, err := v.getRefSubAttribute(refAttr, subAttrName); err != nil {
			return err
		}
	}
	return nil
}