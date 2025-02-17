package types

// APIResponse представляет собой структурированный ответ с ошибкой для API,
// с универсальным полем "data", которое может быть любого типа.
type APIStatusResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
