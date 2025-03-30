// Code generated by go-swagger; DO NOT EDIT.

package operations

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

// NewStorageServiceReleaseVersionGetParams creates a new StorageServiceReleaseVersionGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewStorageServiceReleaseVersionGetParams() *StorageServiceReleaseVersionGetParams {
	return &StorageServiceReleaseVersionGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewStorageServiceReleaseVersionGetParamsWithTimeout creates a new StorageServiceReleaseVersionGetParams object
// with the ability to set a timeout on a request.
func NewStorageServiceReleaseVersionGetParamsWithTimeout(timeout time.Duration) *StorageServiceReleaseVersionGetParams {
	return &StorageServiceReleaseVersionGetParams{
		timeout: timeout,
	}
}

// NewStorageServiceReleaseVersionGetParamsWithContext creates a new StorageServiceReleaseVersionGetParams object
// with the ability to set a context for a request.
func NewStorageServiceReleaseVersionGetParamsWithContext(ctx context.Context) *StorageServiceReleaseVersionGetParams {
	return &StorageServiceReleaseVersionGetParams{
		Context: ctx,
	}
}

// NewStorageServiceReleaseVersionGetParamsWithHTTPClient creates a new StorageServiceReleaseVersionGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewStorageServiceReleaseVersionGetParamsWithHTTPClient(client *http.Client) *StorageServiceReleaseVersionGetParams {
	return &StorageServiceReleaseVersionGetParams{
		HTTPClient: client,
	}
}

/*
StorageServiceReleaseVersionGetParams contains all the parameters to send to the API endpoint

	for the storage service release version get operation.

	Typically these are written to a http.Request.
*/
type StorageServiceReleaseVersionGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the storage service release version get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceReleaseVersionGetParams) WithDefaults() *StorageServiceReleaseVersionGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the storage service release version get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceReleaseVersionGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) WithTimeout(timeout time.Duration) *StorageServiceReleaseVersionGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) WithContext(ctx context.Context) *StorageServiceReleaseVersionGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) WithHTTPClient(client *http.Client) *StorageServiceReleaseVersionGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the storage service release version get params
func (o *StorageServiceReleaseVersionGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *StorageServiceReleaseVersionGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
