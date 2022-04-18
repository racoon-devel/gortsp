package rtsp

// StatusCode represents RTSP response status code
type StatusCode int

const (
	Continue                       StatusCode = 100
	Ok                             StatusCode = 200
	Created                        StatusCode = 201
	LowOnStorageSpace              StatusCode = 250
	MultipleChoices                StatusCode = 300
	MovedPermanently               StatusCode = 301
	MovedTemporarily               StatusCode = 302
	SeeOther                       StatusCode = 303
	NotModified                    StatusCode = 304
	UseProxy                       StatusCode = 305
	BadRequest                     StatusCode = 400
	Unauthorized                   StatusCode = 401
	PaymentRequired                StatusCode = 402
	Forbidden                      StatusCode = 403
	NotFound                       StatusCode = 404
	MethodNotAllowed               StatusCode = 405
	NotAcceptable                  StatusCode = 406
	ProxyAuthenticationRequired    StatusCode = 407
	RequestTimeout                 StatusCode = 408
	Gone                           StatusCode = 410
	LengthRequired                 StatusCode = 411
	PreconditionFailed             StatusCode = 412
	RequestEntityTooLarge          StatusCode = 413
	RequestURITooLarge             StatusCode = 414
	UnsupportedMediaType           StatusCode = 415
	ParameterNotUnderstood         StatusCode = 451
	ConferenceNotFound             StatusCode = 452
	NotEnoughBandwidth             StatusCode = 453
	SessionNotFound                StatusCode = 454
	MethodNotValidInThisState      StatusCode = 455
	HeaderFieldNotValidForResource StatusCode = 456
	InvalidRange                   StatusCode = 457
	ParameterIsReadOnly            StatusCode = 458
	AggregateOperationNotAllowed   StatusCode = 459
	OnlyAggregateOperationAllowed  StatusCode = 460
	UnsupportedTransport           StatusCode = 461
	DestinationUnreachable         StatusCode = 462
	InternalServerError            StatusCode = 500
	NotImplemented                 StatusCode = 501
	BadGateway                     StatusCode = 502
	ServiceUnavailable             StatusCode = 503
	GatewayTimeout                 StatusCode = 504
	VersionNotSupported            StatusCode = 505
	OptionNotSupported             StatusCode = 551
)

// String returns text description of specified StatusCode
func (c StatusCode) String() string {
	statusCodeStrings := map[StatusCode]string{
		Continue:                       "Continue",
		Ok:                             "Ok",
		Created:                        "Created",
		LowOnStorageSpace:              "Low On Storage Space",
		MultipleChoices:                "Multiple Choices",
		MovedPermanently:               "Moved Permanently",
		MovedTemporarily:               "Moved Temporarily",
		SeeOther:                       "See Other",
		NotModified:                    "Not Modified",
		UseProxy:                       "Use Proxy",
		BadRequest:                     "Bad Request",
		Unauthorized:                   "Unauthorized",
		PaymentRequired:                "Payment Required",
		Forbidden:                      "Forbidden",
		NotFound:                       "Not Found",
		MethodNotAllowed:               "Method Not Allowed",
		NotAcceptable:                  "Not Acceptable",
		ProxyAuthenticationRequired:    "Proxy Authentication Required",
		RequestTimeout:                 "Request Time-out",
		Gone:                           "Gone",
		LengthRequired:                 "Length Required",
		PreconditionFailed:             "Precondition Failed",
		RequestEntityTooLarge:          "Request Entity Too Large",
		RequestURITooLarge:             "Request-URI Too Large",
		UnsupportedMediaType:           "Unsupported Media Type",
		ParameterNotUnderstood:         "Parameter Not Understood",
		ConferenceNotFound:             "Conference Not Found",
		NotEnoughBandwidth:             "Not Enough Bandwidth",
		SessionNotFound:                "Session Not Found",
		MethodNotValidInThisState:      "Method Not Valid in This State",
		HeaderFieldNotValidForResource: "Header Field Not Valid for Resource",
		InvalidRange:                   "Invalid Range",
		ParameterIsReadOnly:            "Parameter Is Read-Only",
		AggregateOperationNotAllowed:   "Aggregate operation not allowed",
		OnlyAggregateOperationAllowed:  "Only aggregate operation allowed",
		UnsupportedTransport:           "Unsupported transport",
		DestinationUnreachable:         "Destination unreachable",
		InternalServerError:            "Internal Server Error",
		NotImplemented:                 "Not Implemented",
		BadGateway:                     "Bad Gateway",
		ServiceUnavailable:             "Service Unavailable",
		GatewayTimeout:                 "Gateway Time-out",
		VersionNotSupported:            "RTSP Version not supported",
		OptionNotSupported:             "Option not supported",
	}

	s, ok := statusCodeStrings[c]
	if !ok {
		return statusCodeStrings[InternalServerError]
	}

	return s
}
