package response

type Author struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type UserActionListResponse struct {
	ActionID uint     `json:"actionId"`
	Author   *Author  `json:"author,omitempty"`
	Article  *Article `json:"article,omitempty"`
}
