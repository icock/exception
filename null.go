package exception

// NullReference is thrown when an application attempts to use null
// in a case where an a non-null value is required.
// NullReference implements Throwable.
type NullReference struct {
	*Exception
}

// NewNullReference constructs a NullReference exception.
func NewNullReference(message string, cause Throwable) *NullReference {
	return &NullReference{New(message, cause)}
}


