package errors

type DuplicateError struct{}

func (e *DuplicateError) Error() string {
	return "duplicate entry"
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "entry not found"
}
