package payload

type CommentRequest struct {
	Comment string `json:"comment" validate:"required"`
	PostId  string `json:"post_id" validate:"required"`
}
