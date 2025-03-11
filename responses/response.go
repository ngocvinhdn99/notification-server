package responses

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	// Data    *echo.Map `json:"data"`

}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Total int `json:"total"`
}
