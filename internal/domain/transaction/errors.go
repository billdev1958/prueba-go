package transaction

import "errors"

var ErrInvalidTransactionAmount = errors.New("el monto de la transacción debe ser mayor a cero")

var ErrMissingComercioID = errors.New("el ID del comercio es obligatorio para la transacción")

var ErrInvalidTransactionID = errors.New("el ID de la transacción es inválido o está vacío")

var ErrTransactionNotFound = errors.New("la transacción no fue encontrada")

// (evitar que te cobren dos veces exactamente la misma transacción por error de red).
var ErrTransactionAlreadyExists = errors.New("esta transacción ya fue procesada anteriormente")

var ErrNegativeNetAmount = errors.New("error de cálculo: el monto neto de la transacción no puede ser negativo")
