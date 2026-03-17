package main

import (
	"regexp"
	"strings"
)

// creating a tuple to hold a string and error function
type ExpectedValueTuple struct {
	value string
	errFn CreateError
}

func newExpectedValueTuple(value string, errFn CreateError) ExpectedValueTuple {
	return ExpectedValueTuple{
		value: value,
		errFn: errFn,
	}
}

func (tuple *ExpectedValueTuple) GetString() string {
	return tuple.value
}

func (tuple *ExpectedValueTuple) GetErrorFunction() CreateError {
	return tuple.errFn
}

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

func (doc *AccessibleDocument) VerifyAccessibility() errList {
	list := make(errList, 0)

	// verify tagging was implemented correctly
	err := doc.VerifyTagging()
	if err != nil {
		list = append(list, err...)
	}

	// verify images have relevant alt text
	err = doc.VerifyGraphics()
	if err != nil {
		list = append(list, err...)
	}

	// verify tables have proper accessibility
	err = doc.VerifyTables()
	if err != nil {
		list = append(list, err...)
	}
	return list
}

func (doc *AccessibleDocument) VerifyTagging() errList {
	list := make(errList, 0)

	// check \DocumentMetadata exists before the document class only once
	prerequisiteMetadataLines := FindLinesWithName(doc.innerDocument.GetPrerequisiteContent(), "DocumentMetadata")
	preambleMetadataLines := FindLinesWithName(doc.innerDocument.GetPreamble(), "DocumentMetadata")
	contentMetadataLines := doc.innerDocument.GetContent().FindAllLines("DocumentMetadata")

	// checking \DocumentMetadat exactly once
	metadataLines := append(prerequisiteMetadataLines, append(preambleMetadataLines, contentMetadataLines...)...)
	if len(metadataLines) > 1 {
		for _, line := range metadataLines {
			list = append(list, DOCUMENT_METADATA_REPEATED(line.startCoordinate, line.endCoordinate))
		}
		return list
	}
	if len(metadataLines) == 0 {
		list = append(list, DummyError(DOCUMENT_METADATA_MISSING))
		return list
	}

	// checking \DocumentMetadata appears before the document class
	metadata := metadataLines[0]
	if len(prerequisiteMetadataLines) == 0 {
		list = append(list, DOCUMENT_METADATA_APPEARED_LATE(metadata.startCoordinate, metadata.endCoordinate))
	}

	// find and parse its only required argument
	if len(metadata.arguments) == 0 {
		return append(list, METADATA_LACKS_CLASS_ARGUMENT(metadata.startCoordinate, metadata.endCoordinate))
	}
	if len(metadata.arguments) > 1 {
		return append(list, METADATA_SEVERAL_ARGUMENTS(metadata.startCoordinate, metadata.endCoordinate))
	}

	// verify it is a required argument
	classArg, ok := metadata.arguments[0].(*ClassArgument)
	if !ok {
		return append(list, METADATA_OTHER_ARGUMENTS_UNSUPPORTED(metadata.startCoordinate, metadata.endCoordinate))
	}

	// parsing as a mapping of key-value pairs
	selectedArg, errFn := newKeyValueArgument(classArg.GetValue().(string))
	if errFn != nil {
		return append(list, errFn(metadata.startCoordinate, metadata.endCoordinate))
	}
	mappedArg := selectedArg.(*KeyValueArgument)

	// checking expected values are used
	// "" is used to indicate that the value is not checked
	checks := map[string]ExpectedValueTuple{
		"tagging":       newExpectedValueTuple("on", METADATA_LACKS_ENABLED_TAGGING_ARGUMENT),
		"tagging-setup": newExpectedValueTuple("", METADATA_LACKS_TAGGING_SETUP_ARGUMENT),
		"pdfstandard":   newExpectedValueTuple("", METADATA_LACKS_PDF_STANDARD_ARGUMENT),
		"lang":          newExpectedValueTuple("", METADATA_LACKS_LANGUAGE_ARGUMENT),
	}

	for key, tuple := range checks {
		val, ok := mappedArg.GetSelectedValue(key)
		if !ok {
			list = append(list, tuple.GetErrorFunction()(metadata.startCoordinate, metadata.endCoordinate))
		} else if tuple.GetString() != "" && val != tuple.GetString() {
			list = append(list, tuple.GetErrorFunction()(metadata.startCoordinate, metadata.endCoordinate))
		}
	}
	return list
}

func (doc *AccessibleDocument) VerifyGraphics() errList {
	list := make(errList, 0)

	// find all graphics
	prerequisiteGraphicsLines := FindLinesWithName(doc.innerDocument.GetPrerequisiteContent(), "includegraphics")
	preambleGraphicsLines := FindLinesWithName(doc.innerDocument.GetPreamble(), "includegraphics")
	contentGraphicsLines := doc.innerDocument.GetContent().FindAllLines("includegraphics")

	// a graphic exists outside the document content
	if len(append(prerequisiteGraphicsLines, preambleGraphicsLines...)) > 0 {
		for _, line := range append(prerequisiteGraphicsLines, preambleGraphicsLines...) {
			list = append(list, GRAPHICS_OUTSIDE_DOCUMENT_CONTENT(line.startCoordinate, line.endCoordinate))
		}
	}

	// check graphicx package is included
	if len(contentGraphicsLines) > 0 {
		list = append(list, doc.VerifyPackageExists("graphicx")...)
	}

	// parsing each argument
	for _, line := range contentGraphicsLines {
		// verify it has exactly 2 arguments
		if len(line.arguments) == 0 {
			list = append(list, GRAPHICS_MISSING_ARGUMENTS(line.startCoordinate, line.endCoordinate))
			continue
		}
		if len(line.arguments) > 2 {
			list = append(list, GRAPHICS_TOO_MANY_ARGUMENTS(line.startCoordinate, line.endCoordinate))
			continue
		}

		// verify, if it only has 1 argument, that it is a required argument
		if len(line.arguments) == 1 {
			_, ok := line.arguments[0].(*ClassArgument)
			if ok {
				list = append(list, GRAPHICS_LACKS_ALT_TEXT(line.startCoordinate, line.endCoordinate))
			} else {
				list = append(list, GRAPHICS_LACKS_SOURCE(line.startCoordinate, line.endCoordinate))
			}
			continue
		}

		// handling of 2 arguments being the correct types
		optionalArg, optionalOk := line.arguments[0].(*OptionArgument)
		if !optionalOk {
			list = append(list, GRAPHICS_FIRST_ARGUMENT_NOT_OPTIONAL(line.startCoordinate, line.endCoordinate))
		}
		_, requiredOk := line.arguments[1].(*ClassArgument)
		if !requiredOk {
			list = append(list, GRAPHICS_SECOND_ARGUMENT_NOT_REQUIRED(line.startCoordinate, line.endCoordinate))
		}
		if !optionalOk || !requiredOk {
			continue
		}

		// ignore the graphic if it is an artifact
		if strings.Contains(strings.Trim(optionalArg.GetValue().(string), WHITESPACE), "artifact") {
			continue
		}

		// verify it has alt text or actualtext
		kvArg, errFn := newKeyValueArgument(optionalArg.GetValue().(string))
		if errFn != nil {
			return append(list, errFn(line.startCoordinate, line.endCoordinate))
		}
		var mappedArg *KeyValueArgument = kvArg.(*KeyValueArgument)
		_, ok := mappedArg.GetSelectedValue("alt")
		if !ok {
			_, ok = mappedArg.GetSelectedValue("actualtext")
			if ok {
				continue
			}
			list = append(list, GRAPHICS_LACKS_ALT_TEXT(line.startCoordinate, line.endCoordinate))
		}
	}

	return list
}

func (doc *AccessibleDocument) VerifyTables() errList {
	list := make(errList, 0)

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
			list = append(list, TABLE_LACKS_REQUIRED_ARGUMENT(group.startCoordinate, group.endCoordinate))
			continue
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
			list = append(list, SERVER_RESPONSIBLE_TABLE_NOT_FOUND(group.startCoordinate, group.endCoordinate))
			continue
		}

		// verify the indexes just before have at least 1 tagpdfsetup
		var tagpdfsetup Component
		for i := index - 1; i >= 0; i-- {
			_, ok := components[i].(*Line)
			// we ignore any groups
			if !ok {
				break
			}

			// ignore non-tagpdfsetup lines
			if components[i].GetName() != "tagpdfsetup" {
				continue
			}

			tagpdfsetup = components[i]
			break
		}
		if tagpdfsetup == nil {
			list = append(list, TABLE_MISSING_TAG_PDF_SETUP(group.startCoordinate, group.endCoordinate))
			continue
		}

		// verify it has exactly one required argument
		if len(tagpdfsetup.GetArguments()) != 1 {
			list = append(list, TAG_PDF_SETUP_REQUIRED_ARGUMENT(group.startCoordinate, group.endCoordinate))
			continue
		}

		// verify the argument is a required one
		arg, ok := tagpdfsetup.GetArguments()[0].(*ClassArgument)
		if !ok {
			list = append(list, TAG_PDF_SETUP_NON_REQUIRED_ARGUMENT(group.startCoordinate, group.endCoordinate))
			continue
		}

		// parse it as a key-value pair
		selectedArg, errFn := newKeyValueArgument(arg.GetValue().(string))
		if errFn != nil {
			list = append(list, errFn(group.startCoordinate, group.endCoordinate))
			continue
		}
		mappedArg := selectedArg.(*KeyValueArgument)

		// verify one of the expected arguments are used
		headerRowCount, ok := mappedArg.GetSelectedValue("table/header-rows")
		if ok {
			// verify header row count
			re := regexp.MustCompile(NUMBERS_ARGUMENT_REGEX)
			if re.FindString(headerRowCount) == "" {
				list = append(list, TAG_PDF_SETUP_HEADER_ROWS_INVALID_VALUE(group.startCoordinate, group.endCoordinate))
			}
		} else {
			presentation, ok := mappedArg.GetSelectedValue("table/tagging")
			if !ok {
				list = append(list, TAG_PDF_SETUP_LACKS_HEADER_ROWS_OR_PRESENTATION(group.startCoordinate, group.endCoordinate))
			}
			if presentation != "presentation" {
				list = append(list, TAG_PDF_SETUP_TABLE_TAGGING_INVALID_VALUE(group.startCoordinate, group.endCoordinate))
			}
		}
	}
	return list
}

func (doc *AccessibleDocument) VerifyPackageExists(pkgName string) errList {
	list := make(errList, 0)

	// get all package lines
	prerequisitePackageLines := FindLinesWithName(doc.innerDocument.GetPrerequisiteContent(), "usepackage")
	preamblePackageLines := FindLinesWithName(doc.innerDocument.GetPreamble(), "usepackage")
	contentPackageLines := doc.innerDocument.GetContent().FindAllLines("usepackage")

	// check it's only within expected spaces
	if len(append(prerequisitePackageLines, contentPackageLines...)) > 0 {
		for _, line := range append(prerequisitePackageLines, contentPackageLines...) {
			list = append(list, USEPACKAGE_OUTSIDE_PREAMBLE(line.startCoordinate, line.endCoordinate))
		}
	}

	// verify the expected pkgName exists
	found := false
	for _, line := range preamblePackageLines {
		// get the required arg if it exists
		var selectedArg Argument
		for _, arg := range line.GetArguments() {
			_, ok := arg.(*ClassArgument)
			if ok {
				selectedArg = arg
				break
			}
		}
		if selectedArg == nil {
			list = append(list, USEPACKAGE_LACKS_CLASS_ARGUMENT(line.startCoordinate, line.endCoordinate))
			continue
		}

		// verify if the package was found
		if strings.Trim(selectedArg.GetValue().(string), WHITESPACE) == pkgName {
			found = true
		}
	}
	if !found {
		list = append(list, DummyError(MISSING_GRAPHICX_PACKAGE))
	}

	return list
}
