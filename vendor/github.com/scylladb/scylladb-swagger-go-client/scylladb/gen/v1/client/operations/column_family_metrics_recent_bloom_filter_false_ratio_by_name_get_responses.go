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

// ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetReader is a Reader for the ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGet structure.
type ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK creates a ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK with default headers values
func NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK() *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK {
	return &ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK{}
}

/*
ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK handles this case with default header values.

ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK column family metrics recent bloom filter false ratio by name get o k
*/
type ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault creates a ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault with default headers values
func NewColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault(code int) *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault {
	return &ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics recent bloom filter false ratio by name get default response
func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalseRatioByNameGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
