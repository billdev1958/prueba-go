package comercio

import "errors"

var ErrComercioNotFound = errors.New("el comercio no fue encontrado")

var ErrComercioAlreadyExists = errors.New("el comercio ya existe")

var ErrComercioHasTransactions = errors.New("no se puede eliminar el comercio porque tiene transacciones activas")

var ErrInvalidComercioData = errors.New("los datos del comercio son invalidos")

var ErrInvalidComercioID = errors.New("el id del comercio es invalido")

var ErrInvalidComercioName = errors.New("el nombre del comercio es invalido")

var ErrInvalidComercioComissionRate = errors.New("el porcentaje de comision del comercio es invalido")
