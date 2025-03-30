// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// VersionValue version_value
//
// # Holds a version value for an application state
//
// swagger:model version_value
type VersionValue struct {

	// The application state enum index
	ApplicationState int32 `json:"application_state,omitempty"`

	// The version value
	Value string `json:"value,omitempty"`

	// The application state version
	Version int32 `json:"version,omitempty"`
}

// Validate validates this version value
func (m *VersionValue) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this version value based on context it is used
func (m *VersionValue) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *VersionValue) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VersionValue) UnmarshalBinary(b []byte) error {
	var res VersionValue
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
