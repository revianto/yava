package services

type ServiceError struct {
	Code    int
	ErrCode string
	Message string
}

func (e *ServiceError) Error() string { return e.Message }
