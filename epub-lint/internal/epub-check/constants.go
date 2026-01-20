package epubcheck

const (
	invalidIdPrefix         = "Error while parsing file: value of attribute \""
	invalidAttribute        = "Error while parsing file: attribute \""
	missingUniqueIdentifier = "The unique-identifier \""
	EmptyMetadataProperty   = "Error while parsing file: character content of element \""
	invalidPlayOrder        = "Error while parsing file: identical playOrder values for navPoint/navTarget/pageTarget that do not refer to same target"
	duplicateIdPrefix       = "Error while parsing file: Duplicate ID \""
	invalidBlockquote       = "Error while parsing file: element \"blockquote\" incomplete;"
	missingImgAlt           = "Error while parsing file: element \"img\" missing required attribute \"alt\""
	unexpectedSectionEl     = "Error while parsing file: element \"section\" not allowed here"
)
