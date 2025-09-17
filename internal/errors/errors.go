package deterrs

import (
	"errors"
	"fmt"
)

type detailKey int

const (
	DetOriginalError detailKey = iota
	DetField         detailKey = iota
)

type DetErr struct {
	Type    DetErrType
	Details map[detailKey]interface{}
	Cause   error
}

func (e *DetErr) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Type: %s, Details: %v, Cause: %v}", e.Type, e.Details, e.Cause)
}

func (e *DetErr) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	var t *DetErr
	if errors.As(target, &t) {
		return e.Type == t.Type
	}
	return false
}

func (e *DetErr) As(target interface{}) bool {
	if t, ok := target.(**DetErr); ok {
		*t = e
		return true
	}
	return false
}

func (e *DetErr) Unwrap() error {
	return e.Cause
}

type DetErrType string

const (
	EmptyField   DetErrType = "field must be filled"
	InvalidValue DetErrType = "invalid value"

	OrderActionNotPermittedByStatus DetErrType = "action permitted by order status"

	QueryInsertFailed = "failed to insert"
	QueryUpdateFailed = "failed to update"
	QuerySelectFailed = "failed to select"

	NotImplemented DetErrType = "action not implemented"
	Unknown        DetErrType = "unknow error type"
)

type detErrOption func(*DetErr)

func WithOriginalError(origErr error) detErrOption {
	return func(de *DetErr) {
		de.Details[DetOriginalError] = origErr
	}
}
func WithField(field string) detErrOption {
	return func(de *DetErr) {
		de.Details[DetField] = field
	}
}

func NewDetErr(detErrType DetErrType, opts ...detErrOption) *DetErr {
	detErr := &DetErr{Type: detErrType, Details: make(map[detailKey]interface{})}

	for _, opt := range opts {
		opt(detErr)
	}
	return detErr
}
