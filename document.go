// using the builder design pattern to verify the document is being created correctly and for its construction
package main

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
	reset() CreateError
	addPrerequisiteLine(Component) CreateError
	addDocumentClass(Line) CreateError
	addPreambleLine(Line) CreateError
	addDocumentContent(Component) CreateError
	buildDocument() (Document, CreateError)
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

func (builder *DocumentBuilder) reset() CreateError {
	builder.prerequisiteContent = nil
	builder.documentClass = nil
	builder.preamble = nil
	builder.content = nil

	return nil
}

func (builder *DocumentBuilder) addPrerequisiteLine(component Component) CreateError {
	// ensuring the value is a valid line
	line, ok := component.(*Line)
	if !ok {
		return PREREQUISITE_CONTAINS_GROUP
	}

	// adding it to the document
	builder.prerequisiteContent = append(builder.prerequisiteContent, *line)
	return nil
}

func (builder *DocumentBuilder) addDocumentClass(line Line) CreateError {
	if builder.documentClass != nil {
		return SEVERAL_DOCUMENT_CLASSES_FOUND
	}
	builder.documentClass = &line

	return nil
}

func (builder *DocumentBuilder) addPreambleLine(line Line) CreateError {
	// checking for group construct
	if line.GetName() == "begin" || line.GetName() == "end" {
		return PREAMBLE_CONTAINS_GROUP
	}

	// adding it to the preamble
	builder.preamble = append(builder.preamble, line)
	return nil
}

func (builder *DocumentBuilder) addDocumentContent(content Component) CreateError {
	if builder.content != nil {
		return DOCUMENT_CONTENT_ALREADY_EXISTS
	}
	builder.content = content

	return nil
}

func (builder *DocumentBuilder) buildDocument() (Document, CreateError) {
	// extra verification steps related to the structure of the document

	return Document{
		prerequisiteContent: builder.prerequisiteContent,
		documentClass:       builder.documentClass,
		preamble:            builder.preamble,
		content:             builder.content,
	}, nil
}
