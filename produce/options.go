package produce

// PublishOption customizes CloudEvent attributes for a single publish call.
type PublishOption func(*publishOptions)

type publishOptions struct {
	subject         string
	dataContentType string
	extensions      map[string]any
}

// WithSubject sets the CloudEvent subject field (often used as topic/routing key).
func WithSubject(subject string) PublishOption {
	return func(o *publishOptions) {
		o.subject = subject
	}
}

// WithDataContentType overrides the CloudEvent datacontenttype.
func WithDataContentType(contentType string) PublishOption {
	return func(o *publishOptions) {
		o.dataContentType = contentType
	}
}

// WithExtension adds a CloudEvent extension for a single publish call.
func WithExtension(name string, value any) PublishOption {
	return func(o *publishOptions) {
		if o.extensions == nil {
			o.extensions = make(map[string]any)
		}
		o.extensions[name] = value
	}
}
