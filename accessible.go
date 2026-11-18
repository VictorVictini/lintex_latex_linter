package main

// consider whether to keep this design pattern for the package or to allow the user to decide the use of the data

// builder object to ensure underlying document meets accessibility requirements
// component interface
type IDocument interface {
	GetPrerequisiteContent() []Line
	GetDocumentClass() Line
	GetPreamble() []Line
	GetContent() Component
}

// concrete component
type Document struct {
	prerequisiteContent []Line // for rare content that goes before the document class
	documentClass       Line
	preamble            []Line
	content             Component
}

func (doc *Document) GetPrerequisiteContent() []Line {
	return doc.prerequisiteContent
}

func (doc *Document) GetDocumentClass() Line {
	return doc.documentClass
}

func (doc *Document) GetPreamble() []Line {
	return doc.preamble
}

func (doc *Document) GetContent() Component {
	return doc.content
}

// base decorator
type IDocumentBase interface{}
