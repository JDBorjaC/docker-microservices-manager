package internal

import "errors"

var (
	ErrNotFound           = errors.New("recurso no encontrado")
	ErrDuplicate          = errors.New("recurso ya existe")
	ErrInvalidRequest     = errors.New("parametros de solicitud invalidos")
	ErrUnsupportedLang    = errors.New("lenguaje no soportado")
	ErrBuildFailed        = errors.New("fallo al construir la imagen")
	ErrContainerFailed    = errors.New("fallo al gestionar el contenedor")
)
