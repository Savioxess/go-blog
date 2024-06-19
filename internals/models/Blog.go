package models

type Blog struct {
	ID        []byte `json:"id"`
	AuthorID  []byte `json:"author_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	PostedOn string `json:"posted_on"`
}
