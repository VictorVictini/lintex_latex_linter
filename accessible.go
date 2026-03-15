package main

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

// concrete decorator
type AccessibleDocument struct {
	innerDocument IDocument
}

func newAccessibleDocument(doc IDocument) IDocumentBase {
	return &AccessibleDocument{
		innerDocument: doc,
	}
}

func (doc *AccessibleDocument) VerifyAccessibility() CreateError {
	// verify tagging was implemented correctly
	err := doc.VerifyTagging()
	if err != nil {
		return err
	}

	// verify images have relevant alt text
	err = doc.VerifyGraphics()
	if err != nil {
		return err
	}

	// verify tables have proper accessibility
	err = doc.VerifyTables()
	if err != nil {
		return err
	}
	return nil
}

func (doc *AccessibleDocument) VerifyTagging() CreateError {
	// check \DocumentMetadata exists before the document class only once
	prerequisiteMetadataLines := FindLinesWithName(doc.innerDocument.GetPrerequisiteContent(), "DocumentMetadata")
	preambleMetadataLines := FindLinesWithName(doc.innerDocument.GetPreamble(), "DocumentMetadata")
	contentMetadataLines := doc.innerDocument.GetContent().FindAllLines("DocumentMetadata")

	// checking \DocumentMetadat exactly once
	metadataLines := append(prerequisiteMetadataLines, append(preambleMetadataLines, contentMetadataLines...)...)
	if len(metadataLines) > 1 {
		return DOCUMENT_METADATA_REPEATED
	}
	if len(metadataLines) == 0 {
		return DOCUMENT_METADATA_MISSING
	}

	// checking \DocumentMetadata appears before the document class
	if len(prerequisiteMetadataLines) == 0 {
		return DOCUMENT_METADATA_APPEARED_LATE
	}

	// include checks that necessary arguments of \MetaData exist
	// using newKeyValueArgument()

	return nil
}

func (doc *AccessibleDocument) VerifyGraphics() CreateError {
	return nil
}

func (doc *AccessibleDocument) VerifyTables() CreateError {
	return nil
}
