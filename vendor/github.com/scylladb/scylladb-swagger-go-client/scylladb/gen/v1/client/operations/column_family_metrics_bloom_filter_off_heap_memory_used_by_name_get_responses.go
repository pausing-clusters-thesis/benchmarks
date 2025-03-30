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

// ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetReader is a Reader for the ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGet structure.
type ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK creates a ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK with default headers values
func NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK() *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK {
	return &ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK{}
}

/*
ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK handles this case with default header values.

ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK column family metrics bloom filter off heap memory used by name get o k
*/
type ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault creates a ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault with default headers values
func NewColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault(code int) *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault {
	return &ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics bloom filter off heap memory used by name get default response
func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsBloomFilterOffHeapMemoryUsedByNameGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
