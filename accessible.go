package main

// builder object to ensure underlying document meets accessibility requirements
// component interface
type IDocument interface{}

// concrete component
type Document struct {
	prerequisiteContent []Line // for rare content that goes before the document class
	documentClass       *Line
	preamble            []Line
	content             Component
}

// base decorator
type IDocumentBase interface{}

// concrete decorator
type AccessibleDocument struct {
	innerDocument IDocument
}

func newAccessibleDocument(doc IDocument) IDocumentBase {
	return &AccessibleDocument{
		innerDocument: doc,
	}
}

func (doc *AccessibleDocument) VerifyAccessibility() CustomError {
	// do shid later
	return nil
}
