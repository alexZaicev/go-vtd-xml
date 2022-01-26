package erroring

const (
	CannotBeNil                = "cannot be nil"
	IndexOutOfRange            = "array index out of range"
	InvalidSliceLength         = "invalid slice length"
	InvalidBufferPageSize      = "invalid buffer page size"
	InvalidChar                = "invalid character"
	MaximumDepthExceeded       = "maximum depth exceeded"
	TagPrefixQnameTooLong      = "starting tag prefix or QNAME length too long"
	InvalidCharInText          = "invalid char in text content"
	IllegalBuiltInEntity       = "illegal build-in entity reference"
	XmlIncomplete              = "XML document incomplete"
	AttrNotUnique              = "attribute name not unique"
	AttrNamePrefixQnameTooLong = "attribute name prefix or QNAME length too long"
	AttrNsPrefixQnameTooLong   = "attribute namespace tag prefix or QNAME length too long"
	NonDefaultNsEmpty          = "non-default namespace cannot be empty"
	AttrValueTooLong           = "attribute value is too long"
)
