package app_err

type BusinessError struct {
	code    string
	message string
}

func (b BusinessError) Error() string {
	return b.message
}

func (b BusinessError) Code() string {
	return b.code
}

func NewBusinessError(message string) error {
	return BusinessError{
		code:    "BusinessError",
		message: message,
	}
}

func NewAuthorizationError(message string) error {
	return BusinessError{
		code:    "AuthorizationError",
		message: message,
	}
}

func NewUnauthorizedError() error {
	return BusinessError{
		code:    "UnauthorizedError",
		message: "Вы не авторизованы. Пожалуйста, авторизуйтесь",
	}
}

func NewConflictError(message string) error {
	return BusinessError{
		code:    "ConflictError",
		message: message,
	}
}
