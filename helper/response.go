package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Status     int         `json:"status"`
	Is_Success bool        `json:"is_success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}
type Paginate struct {
	Limit        int    `json:"limit"`
	Page         int    `json:"page"`
	Sort         string `json:"sort"`
	TotalRows    int    `json:"total_rows"`
	FirstPage    string `json:"first_page"`
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	LastPage     string `json:"last_page"`
	TotalPages   int    `json:"total_pages"`
}

func APIResponse(Status int, Is_success bool, Message string, data interface{}) Response {

	Result := Response{
		Status:     Status,
		Is_Success: Is_success,
		Message:    Message,
		Data:       data,
	}

	return Result
}

func FormatError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
