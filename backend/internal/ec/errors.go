package ec

// ProductError informs about possible failures in product requests.
type ProductError string

const (
	// RequestFailedError when service call fails.
	RequestFailedError ProductError = "HTTP_REQUEST_FAILED"

	// ReadDataFailedError when response contains errors.
	ReadDataFailedError ProductError = "READ_DATA_FAILED"

	// UnmarshalDataFailedError when unmarshal data from response fails.
	UnmarshalDataFailedError ProductError = "UNMARSHAL_DATA_FAILED"
)

// Error returns error as a string.
func (e ProductError) Error() string {
	return string(e)
}

// String returns string value.
func (e ProductError) String() string {
	return string(e)
}
