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

// NewGossiperGenerationNumberByAddrGetParams creates a new GossiperGenerationNumberByAddrGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGossiperGenerationNumberByAddrGetParams() *GossiperGenerationNumberByAddrGetParams {
	return &GossiperGenerationNumberByAddrGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGossiperGenerationNumberByAddrGetParamsWithTimeout creates a new GossiperGenerationNumberByAddrGetParams object
// with the ability to set a timeout on a request.
func NewGossiperGenerationNumberByAddrGetParamsWithTimeout(timeout time.Duration) *GossiperGenerationNumberByAddrGetParams {
	return &GossiperGenerationNumberByAddrGetParams{
		timeout: timeout,
	}
}

// NewGossiperGenerationNumberByAddrGetParamsWithContext creates a new GossiperGenerationNumberByAddrGetParams object
// with the ability to set a context for a request.
func NewGossiperGenerationNumberByAddrGetParamsWithContext(ctx context.Context) *GossiperGenerationNumberByAddrGetParams {
	return &GossiperGenerationNumberByAddrGetParams{
		Context: ctx,
	}
}

// NewGossiperGenerationNumberByAddrGetParamsWithHTTPClient creates a new GossiperGenerationNumberByAddrGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewGossiperGenerationNumberByAddrGetParamsWithHTTPClient(client *http.Client) *GossiperGenerationNumberByAddrGetParams {
	return &GossiperGenerationNumberByAddrGetParams{
		HTTPClient: client,
	}
}

/*
GossiperGenerationNumberByAddrGetParams contains all the parameters to send to the API endpoint

	for the gossiper generation number by addr get operation.

	Typically these are written to a http.Request.
*/
type GossiperGenerationNumberByAddrGetParams struct {

	/* Addr.

	   The endpoint address
	*/
	Addr string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the gossiper generation number by addr get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GossiperGenerationNumberByAddrGetParams) WithDefaults() *GossiperGenerationNumberByAddrGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the gossiper generation number by addr get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GossiperGenerationNumberByAddrGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) WithTimeout(timeout time.Duration) *GossiperGenerationNumberByAddrGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) WithContext(ctx context.Context) *GossiperGenerationNumberByAddrGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) WithHTTPClient(client *http.Client) *GossiperGenerationNumberByAddrGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAddr adds the addr to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) WithAddr(addr string) *GossiperGenerationNumberByAddrGetParams {
	o.SetAddr(addr)
	return o
}

// SetAddr adds the addr to the gossiper generation number by addr get params
func (o *GossiperGenerationNumberByAddrGetParams) SetAddr(addr string) {
	o.Addr = addr
}

// WriteToRequest writes these params to a swagger request
func (o *GossiperGenerationNumberByAddrGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param addr
	if err := r.SetPathParam("addr", o.Addr); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
