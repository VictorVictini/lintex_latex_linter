// using the builder design pattern to verify the document is being created correctly and for its construction
package main

import (
	"errors"
	"fmt"
)

// defining the document to be worked on
// i.e. the product
type Document struct {
	prerequisiteContent []Line // for rare content that goes before the document class
	documentClass       *Line
	preamble            []Line
	content             Component
}

// ensures all content have necessary group requirements
func (doc *Document) CheckContentGroupings() bool {
	return doc.content.CheckGroupings(make(map[string]bool))
}

// builder interface
type IDocumentBuilder interface {
	reset() error
	addPrerequisiteLine(Component) error
	addDocumentClass(Line) error
	addPreambleLine(Line) error
	addDocumentContent(Component) error
	buildDocument() (Document, error)
}

// concrete builder
type DocumentBuilder struct {
	prerequisiteContent []Line // for rare content that goes before the document class
	documentClass       *Line
	preamble            []Line
	content             Component
}

func newDocumentBuilder() *DocumentBuilder {
	return &DocumentBuilder{}
}

func (builder *DocumentBuilder) reset() error {
	builder.prerequisiteContent = nil
	builder.documentClass = nil
	builder.preamble = nil
	builder.content = nil

	return nil
}

func (builder *DocumentBuilder) addPrerequisiteLine(component Component) error {
	// ensuring the value is a valid line
	line, ok := component.(*Line)
	if !ok {
		return fmt.Errorf("Should not have grouping construct before document class")
	}

	// adding it to the document
	builder.prerequisiteContent = append(builder.prerequisiteContent, *line)
	return nil
}

func (builder *DocumentBuilder) addDocumentClass(line Line) error {
	if builder.documentClass != nil {
		return errors.New("Document cannot have several classes")
	}
	builder.documentClass = &line

	return nil
}

func (builder *DocumentBuilder) addPreambleLine(line Line) error {
	// checking for group construct
	if line.GetName() == "begin" || line.GetName() == "end" {
		return fmt.Errorf("Preamble cannot contain a grouping construct (\\%s)", line.GetName())
	}

	// adding it to the preamble
	builder.preamble = append(builder.preamble, line)
	return nil
}

func (builder *DocumentBuilder) addDocumentContent(content Component) error {
	if builder.content != nil {
		return errors.New("Document content should be a single group under 'document'")
	}
	builder.content = content

	return nil
}

func (builder *DocumentBuilder) buildDocument() (Document, error) {
	// extra verification steps related to the structure of the document

	return Document{
		prerequisiteContent: builder.prerequisiteContent,
		documentClass:       builder.documentClass,
		preamble:            builder.preamble,
		content:             builder.content,
	}, nil
}
