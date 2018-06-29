package exceptions

//Exception is a base type to manage error
type Exception struct {
	StatusCode int          `json:"status_code,omitempty"`
	ErrorTag   string       `json:"error,omitempty"`
	Message    string       `json:"message,omitempty"`
	Causes     []*Exception `json:"causes,omitempty"`
}

func (e *Exception) Error() string {
	return e.Message
}

//Status returns http status from exception
func (e *Exception) Status() int {
	return e.StatusCode
}

//AddCause to an exception
func (e *Exception) AddCause(err error) {
	switch t := err.(type) {
	case *Exception:
		e.Causes = append(e.Causes, t)
	default:
		e.Causes = append(e.Causes, &Exception{Message: err.Error()})
	}
}

//NewComponentException returns an exception with component error tag
func NewComponentException(err error) *Exception {
	return newException("internal_error", 500, err)
}

//NewIntegrationException returns an exception when has error in other platform component
func NewIntegrationException(err error) *Exception {
	return newException("integration_error", 500, err)
}

//NewInvalidArgumentException returns an exception with component error tag
func NewInvalidArgumentException(err error) *Exception {
	return newException("invalid_arguments", 400, err)
}

func newException(tag string, httpStatus int, err error) *Exception {
	if err == nil {
		return nil
	}
	ex := new(Exception)
	ex.ErrorTag = tag
	ex.StatusCode = httpStatus
	ex.Message = err.Error()
	ex.Causes = make([]*Exception, 0, 0)
	return ex
}
