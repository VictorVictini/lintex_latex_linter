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

	// find and parse its only required argument
	metadata := metadataLines[0]
	var selectedArg Argument
	for i, arg := range metadata.arguments {
		// verify it is the only required argument
		classArg, ok := arg.(*ClassArgument)
		if !ok {
			continue
		}
		if selectedArg != nil {
			return METADATA_SEVERAL_CLASS_ARGUMENTS
		}

		// replace it with a parsed version
		keyValueArg, err := newKeyValueArgument(classArg.GetValue().(string))
		if err != nil {
			return err
		}
		metadata.arguments[i] = keyValueArg
		selectedArg = metadata.arguments[i]
	}

	// verify it was found
	if selectedArg == nil {
		return METADATA_LACKS_CLASS_ARGUMENT
	}
	mappedArg := selectedArg.(*KeyValueArgument)

	// check it contains 'tagging=on'
	val, ok := mappedArg.GetSelectedValue("tagging")
	if !ok {
		return METADATA_LACKS_ENABLED_TAGGING_ARGUMENT
	}
	if val != "on" {
		return METADATA_LACKS_ENABLED_TAGGING_ARGUMENT
	}

	// check it has a tagging setup
	val, ok = mappedArg.GetSelectedValue("tagging-setup")
	if !ok {
		return METADATA_LACKS_TAGGING_SETUP_ARGUMENT
	}

	// check it has a pdf standard
	val, ok = mappedArg.GetSelectedValue("pdfstandard")
	if !ok {
		return METADATA_LACKS_PDF_STANDARD_ARGUMENT
	}

	// check it has a language specified
	val, ok = mappedArg.GetSelectedValue("lang")
	if !ok {
		return METADATA_LACKS_LANGUAGE_ARGUMENT
	}
	return nil
}

func (doc *AccessibleDocument) VerifyGraphics() CreateError {
	return nil
}

func (doc *AccessibleDocument) VerifyTables() CreateError {
	return nil
}
