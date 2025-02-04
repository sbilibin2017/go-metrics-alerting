package apiclient

// ApiClientInterface - интерфейс клиента.
type ApiClientEngineInterface interface {
	Get(path string, query map[string]string, headers map[string]string) (*ApiResponse, error)
	Post(path string, query map[string]string, body any, headers map[string]string) (*ApiResponse, error)
}
