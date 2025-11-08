package payload

type PostRequest struct {
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body"  validate:"required"`
	Status string `json:"status" validate:"required,oneof=PUBLISH DRAFT"`
}
