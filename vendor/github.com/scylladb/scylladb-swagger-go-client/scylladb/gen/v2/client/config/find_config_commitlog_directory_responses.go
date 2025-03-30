// Code generated by go-swagger; DO NOT EDIT.

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/scylladb/scylladb-swagger-go-client/scylladb/gen/v2/models"
)

// FindConfigCommitlogDirectoryReader is a Reader for the FindConfigCommitlogDirectory structure.
type FindConfigCommitlogDirectoryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FindConfigCommitlogDirectoryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFindConfigCommitlogDirectoryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFindConfigCommitlogDirectoryDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFindConfigCommitlogDirectoryOK creates a FindConfigCommitlogDirectoryOK with default headers values
func NewFindConfigCommitlogDirectoryOK() *FindConfigCommitlogDirectoryOK {
	return &FindConfigCommitlogDirectoryOK{}
}

/*
FindConfigCommitlogDirectoryOK handles this case with default header values.

Config value
*/
type FindConfigCommitlogDirectoryOK struct {
	Payload string
}

func (o *FindConfigCommitlogDirectoryOK) GetPayload() string {
	return o.Payload
}

func (o *FindConfigCommitlogDirectoryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindConfigCommitlogDirectoryDefault creates a FindConfigCommitlogDirectoryDefault with default headers values
func NewFindConfigCommitlogDirectoryDefault(code int) *FindConfigCommitlogDirectoryDefault {
	return &FindConfigCommitlogDirectoryDefault{
		_statusCode: code,
	}
}

/*
FindConfigCommitlogDirectoryDefault handles this case with default header values.

unexpected error
*/
type FindConfigCommitlogDirectoryDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the find config commitlog directory default response
func (o *FindConfigCommitlogDirectoryDefault) Code() int {
	return o._statusCode
}

func (o *FindConfigCommitlogDirectoryDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *FindConfigCommitlogDirectoryDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *FindConfigCommitlogDirectoryDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
