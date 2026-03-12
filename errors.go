package main

type ParseError struct {
	startLine, endLine int
	message            string
}

// Parsing errors : errors related to incorrect parsing of the document
// such as invalid {} [] or begin\end usage

// Structure errors : errors such as missing/mistyping a docclass, having grouping contexts in preamble/prerequisite, having more than 1 component in doc content, missing essential arguments to something (e.g. first class arg to \begin, \end, and \documentclass), doc content not being a begin group containing document
// mostly handled within IDocumentBuilder.buildDocument()

// Accessibility errors : errors related to incorrect accessibility of the document (using outdated methods) or missing accessibility methods (such as not including \Metadata for tagging or relevant code for alt text on figures, images, or tables), may also include some nuances such as table design, text size, colours(?))
