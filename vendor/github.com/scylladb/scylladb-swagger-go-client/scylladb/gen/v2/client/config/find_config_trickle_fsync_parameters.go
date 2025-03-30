// Code generated by go-swagger; DO NOT EDIT.

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewFindConfigTrickleFsyncParams creates a new FindConfigTrickleFsyncParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigTrickleFsyncParams() *FindConfigTrickleFsyncParams {
	return &FindConfigTrickleFsyncParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigTrickleFsyncParamsWithTimeout creates a new FindConfigTrickleFsyncParams object
// with the ability to set a timeout on a request.
func NewFindConfigTrickleFsyncParamsWithTimeout(timeout time.Duration) *FindConfigTrickleFsyncParams {
	return &FindConfigTrickleFsyncParams{
		timeout: timeout,
	}
}

// NewFindConfigTrickleFsyncParamsWithContext creates a new FindConfigTrickleFsyncParams object
// with the ability to set a context for a request.
func NewFindConfigTrickleFsyncParamsWithContext(ctx context.Context) *FindConfigTrickleFsyncParams {
	return &FindConfigTrickleFsyncParams{
		Context: ctx,
	}
}

// NewFindConfigTrickleFsyncParamsWithHTTPClient creates a new FindConfigTrickleFsyncParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigTrickleFsyncParamsWithHTTPClient(client *http.Client) *FindConfigTrickleFsyncParams {
	return &FindConfigTrickleFsyncParams{
		HTTPClient: client,
	}
}

/*
FindConfigTrickleFsyncParams contains all the parameters to send to the API endpoint

	for the find config trickle fsync operation.

	Typically these are written to a http.Request.
*/
type FindConfigTrickleFsyncParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config trickle fsync params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigTrickleFsyncParams) WithDefaults() *FindConfigTrickleFsyncParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config trickle fsync params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigTrickleFsyncParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) WithTimeout(timeout time.Duration) *FindConfigTrickleFsyncParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) WithContext(ctx context.Context) *FindConfigTrickleFsyncParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) WithHTTPClient(client *http.Client) *FindConfigTrickleFsyncParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config trickle fsync params
func (o *FindConfigTrickleFsyncParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigTrickleFsyncParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
