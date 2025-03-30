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

// NewFindConfigEnableShardAwareDriversParams creates a new FindConfigEnableShardAwareDriversParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigEnableShardAwareDriversParams() *FindConfigEnableShardAwareDriversParams {
	return &FindConfigEnableShardAwareDriversParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigEnableShardAwareDriversParamsWithTimeout creates a new FindConfigEnableShardAwareDriversParams object
// with the ability to set a timeout on a request.
func NewFindConfigEnableShardAwareDriversParamsWithTimeout(timeout time.Duration) *FindConfigEnableShardAwareDriversParams {
	return &FindConfigEnableShardAwareDriversParams{
		timeout: timeout,
	}
}

// NewFindConfigEnableShardAwareDriversParamsWithContext creates a new FindConfigEnableShardAwareDriversParams object
// with the ability to set a context for a request.
func NewFindConfigEnableShardAwareDriversParamsWithContext(ctx context.Context) *FindConfigEnableShardAwareDriversParams {
	return &FindConfigEnableShardAwareDriversParams{
		Context: ctx,
	}
}

// NewFindConfigEnableShardAwareDriversParamsWithHTTPClient creates a new FindConfigEnableShardAwareDriversParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigEnableShardAwareDriversParamsWithHTTPClient(client *http.Client) *FindConfigEnableShardAwareDriversParams {
	return &FindConfigEnableShardAwareDriversParams{
		HTTPClient: client,
	}
}

/*
FindConfigEnableShardAwareDriversParams contains all the parameters to send to the API endpoint

	for the find config enable shard aware drivers operation.

	Typically these are written to a http.Request.
*/
type FindConfigEnableShardAwareDriversParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config enable shard aware drivers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigEnableShardAwareDriversParams) WithDefaults() *FindConfigEnableShardAwareDriversParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config enable shard aware drivers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigEnableShardAwareDriversParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) WithTimeout(timeout time.Duration) *FindConfigEnableShardAwareDriversParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) WithContext(ctx context.Context) *FindConfigEnableShardAwareDriversParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) WithHTTPClient(client *http.Client) *FindConfigEnableShardAwareDriversParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config enable shard aware drivers params
func (o *FindConfigEnableShardAwareDriversParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigEnableShardAwareDriversParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
