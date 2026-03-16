package main

import (
	"regexp"
	"strings"
)

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
	if len(metadata.arguments) == 0 {
		return METADATA_LACKS_CLASS_ARGUMENT
	}
	if len(metadata.arguments) > 1 {
		return METADATA_SEVERAL_ARGUMENTS
	}

	// verify it is a required argument
	classArg, ok := metadata.arguments[0].(*ClassArgument)
	if !ok {
		return METADATA_OTHER_ARGUMENTS_UNSUPPORTED
	}

	// parsing as a mapping of key-value pairs
	selectedArg, err := newKeyValueArgument(classArg.GetValue().(string))
	if err != nil {
		return err
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
	// find all graphics
	prerequisiteGraphicsLines := FindLinesWithName(doc.innerDocument.GetPrerequisiteContent(), "includegraphics")
	preambleGraphicsLines := FindLinesWithName(doc.innerDocument.GetPreamble(), "includegraphics")
	contentGraphicsLines := doc.innerDocument.GetContent().FindAllLines("includegraphics")

	// a graphic exists outside the document content
	if len(append(prerequisiteGraphicsLines, preambleGraphicsLines...)) > 0 {
		return GRAPHICS_OUTSIDE_DOCUMENT_CONTENT
	}

	// parsing each argument
	for _, line := range contentGraphicsLines {
		// verify it has exactly 2 arguments
		if len(line.arguments) == 0 {
			return GRAPHICS_MISSING_ARGUMENTS
		}
		if len(line.arguments) > 2 {
			return GRAPHICS_TOO_MANY_ARGUMENTS
		}

		// verify, if it only has 1 argument, that it is a required argument
		if len(line.arguments) == 1 {
			_, ok := line.arguments[0].(*ClassArgument)
			if ok {
				return GRAPHICS_LACKS_ALT_TEXT
			} else {
				return GRAPHICS_LACKS_SOURCE
			}
		}

		// handling of 2 arguments being the correct types
		optionalArg, ok := line.arguments[0].(*OptionArgument)
		if !ok {
			return GRAPHICS_FIRST_ARGUMENT_NOT_OPTIONAL
		}
		_, ok = line.arguments[1].(*ClassArgument)
		if !ok {
			return GRAPHICS_SECOND_ARGUMENT_NOT_REQUIRED
		}

		// ignore the graphic if it is an artifact
		if strings.Contains(optionalArg.GetValue().(string), "artifact") {
			continue
		}

		// verify it has alt text or actualtext
		kvArg, err := newKeyValueArgument(optionalArg.GetValue().(string))
		if err != nil {
			return err
		}
		var mappedArg *KeyValueArgument = kvArg.(*KeyValueArgument)
		_, ok = mappedArg.GetSelectedValue("alt")
		if !ok {
			_, ok = mappedArg.GetSelectedValue("actualtext")
			if ok {
				continue
			}
			return GRAPHICS_LACKS_ALT_TEXT
		}
	}

	return nil
}

func (doc *AccessibleDocument) VerifyTables() CreateError {
	// verify all tables' formats
	tableGroups := doc.innerDocument.GetContent().FindAllGroups("tabular")
	for _, group := range tableGroups {
		// verify it has at least 1 required argument
		var selectedArg Argument
		for _, arg := range group.arguments {
			_, ok := arg.(*ClassArgument)
			if ok {
				selectedArg = arg
				break
			}
		}
		if selectedArg == nil {
			return TABLE_LACKS_REQUIRED_ARGUMENT
		}

		// find what component index the group occurs at
		components := group.outerGroup.components
		index := -1
		for i, component := range components {
			currGroup, ok := component.(*Group)
			if !ok {
				continue
			}
			if currGroup.GetStartCoordinate().line == group.GetStartCoordinate().line && currGroup.GetStartCoordinate().charPos == group.GetStartCoordinate().charPos {
				index = i
				break
			}
		}
		if index == -1 {
			return SERVER_RESPONSIBLE_TABLE_NOT_FOUND
		}

		// verify the index just before has tagpdfsetup
		if index == 0 {
			return TABLE_MISSING_TAG_PDF_SETUP
		}
		tagpdfsetup, ok := components[index-1].(*Line)
		if !ok {
			return TABLE_MISSING_TAG_PDF_SETUP
		}
		if tagpdfsetup.GetName() != "tagpdfsetup" {
			return TABLE_MISSING_TAG_PDF_SETUP
		}

		// verify it has exactly one required argument
		if len(tagpdfsetup.arguments) != 1 {
			return TAG_PDF_SETUP_REQUIRED_ARGUMENT
		}

		// verify the argument is a required one
		arg, ok := tagpdfsetup.arguments[0].(*ClassArgument)
		if !ok {
			return TAG_PDF_SETUP_NON_REQUIRED_ARGUMENT
		}

		// parse it as a key-value pair
		selectedArg, err := newKeyValueArgument(arg.GetValue().(string))
		if err != nil {
			return err
		}
		mappedArg := selectedArg.(*KeyValueArgument)

		// verify one of the expected arguments are used
		headerRowCount, ok := mappedArg.GetSelectedValue("table/header-rows")
		if ok {
			// verify header row count
			re := regexp.MustCompile(NUMBERS_ARGUMENT_REGEX)
			if re.FindString(headerRowCount) == "" {
				return TAG_PDF_SETUP_HEADER_ROWS_INVALID_VALUE
			}
		} else {
			presentation, ok := mappedArg.GetSelectedValue("table/tagging")
			if !ok {
				return TAG_PDF_SETUP_LACKS_HEADER_ROWS_OR_PRESENTATION
			}
			if presentation != "presentation" {
				return TAG_PDF_SETUP_TABLE_TAGGING_INVALID_VALUE
			}
		}
	}
	return nil
}
