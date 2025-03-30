// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// RepairAsyncStatusResponse repair_async_statusResponse
// Example: RUNNING
//
// swagger:model repair_async_statusResponse
type RepairAsyncStatusResponse string

func NewRepairAsyncStatusResponse(value RepairAsyncStatusResponse) *RepairAsyncStatusResponse {
	return &value
}

// Pointer returns a pointer to a freshly-allocated RepairAsyncStatusResponse.
func (m RepairAsyncStatusResponse) Pointer() *RepairAsyncStatusResponse {
	return &m
}

const (

	// RepairAsyncStatusResponseRUNNING captures enum value "RUNNING"
	RepairAsyncStatusResponseRUNNING RepairAsyncStatusResponse = "RUNNING"

	// RepairAsyncStatusResponseSUCCESSFUL captures enum value "SUCCESSFUL"
	RepairAsyncStatusResponseSUCCESSFUL RepairAsyncStatusResponse = "SUCCESSFUL"

	// RepairAsyncStatusResponseFAILED captures enum value "FAILED"
	RepairAsyncStatusResponseFAILED RepairAsyncStatusResponse = "FAILED"
)

// for schema
var repairAsyncStatusResponseEnum []interface{}

func init() {
	var res []RepairAsyncStatusResponse
	if err := json.Unmarshal([]byte(`["RUNNING","SUCCESSFUL","FAILED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		repairAsyncStatusResponseEnum = append(repairAsyncStatusResponseEnum, v)
	}
}

func (m RepairAsyncStatusResponse) validateRepairAsyncStatusResponseEnum(path, location string, value RepairAsyncStatusResponse) error {
	if err := validate.EnumCase(path, location, value, repairAsyncStatusResponseEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this repair async status response
func (m RepairAsyncStatusResponse) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateRepairAsyncStatusResponseEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this repair async status response based on context it is used
func (m RepairAsyncStatusResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
