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

// ColumnFamilyMetricsBloomFilterFalsePositivesGetReader is a Reader for the ColumnFamilyMetricsBloomFilterFalsePositivesGet structure.
type ColumnFamilyMetricsBloomFilterFalsePositivesGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsBloomFilterFalsePositivesGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsBloomFilterFalsePositivesGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsBloomFilterFalsePositivesGetOK creates a ColumnFamilyMetricsBloomFilterFalsePositivesGetOK with default headers values
func NewColumnFamilyMetricsBloomFilterFalsePositivesGetOK() *ColumnFamilyMetricsBloomFilterFalsePositivesGetOK {
	return &ColumnFamilyMetricsBloomFilterFalsePositivesGetOK{}
}

/*
ColumnFamilyMetricsBloomFilterFalsePositivesGetOK handles this case with default header values.

ColumnFamilyMetricsBloomFilterFalsePositivesGetOK column family metrics bloom filter false positives get o k
*/
type ColumnFamilyMetricsBloomFilterFalsePositivesGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsBloomFilterFalsePositivesGetDefault creates a ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault with default headers values
func NewColumnFamilyMetricsBloomFilterFalsePositivesGetDefault(code int) *ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault {
	return &ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics bloom filter false positives get default response
func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsBloomFilterFalsePositivesGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
