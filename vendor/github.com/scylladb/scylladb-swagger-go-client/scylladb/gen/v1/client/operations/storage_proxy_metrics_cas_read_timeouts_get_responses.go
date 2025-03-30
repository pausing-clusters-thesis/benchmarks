// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/scylladb/scylladb-swagger-go-client/scylladb/gen/v1/models"
)

// StorageProxyMetricsCasReadTimeoutsGetReader is a Reader for the StorageProxyMetricsCasReadTimeoutsGet structure.
type StorageProxyMetricsCasReadTimeoutsGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageProxyMetricsCasReadTimeoutsGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageProxyMetricsCasReadTimeoutsGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageProxyMetricsCasReadTimeoutsGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageProxyMetricsCasReadTimeoutsGetOK creates a StorageProxyMetricsCasReadTimeoutsGetOK with default headers values
func NewStorageProxyMetricsCasReadTimeoutsGetOK() *StorageProxyMetricsCasReadTimeoutsGetOK {
	return &StorageProxyMetricsCasReadTimeoutsGetOK{}
}

/*
StorageProxyMetricsCasReadTimeoutsGetOK handles this case with default header values.

StorageProxyMetricsCasReadTimeoutsGetOK storage proxy metrics cas read timeouts get o k
*/
type StorageProxyMetricsCasReadTimeoutsGetOK struct {
	Payload interface{}
}

func (o *StorageProxyMetricsCasReadTimeoutsGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *StorageProxyMetricsCasReadTimeoutsGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStorageProxyMetricsCasReadTimeoutsGetDefault creates a StorageProxyMetricsCasReadTimeoutsGetDefault with default headers values
func NewStorageProxyMetricsCasReadTimeoutsGetDefault(code int) *StorageProxyMetricsCasReadTimeoutsGetDefault {
	return &StorageProxyMetricsCasReadTimeoutsGetDefault{
		_statusCode: code,
	}
}

/*
StorageProxyMetricsCasReadTimeoutsGetDefault handles this case with default header values.

internal server error
*/
type StorageProxyMetricsCasReadTimeoutsGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage proxy metrics cas read timeouts get default response
func (o *StorageProxyMetricsCasReadTimeoutsGetDefault) Code() int {
	return o._statusCode
}

func (o *StorageProxyMetricsCasReadTimeoutsGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageProxyMetricsCasReadTimeoutsGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageProxyMetricsCasReadTimeoutsGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
