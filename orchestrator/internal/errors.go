package internal

import "errors"

var (
	ErrNotFound       = errors.New("recurso no encontrado")
	ErrDuplicate      = errors.New("recurso ya existe")
	ErrInvalidRequest = errors.New("parametros de solicitud invalidos")
)
