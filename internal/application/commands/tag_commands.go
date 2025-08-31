package commands

type CreateTagCommand struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateTagCommand struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddTagToQuoteCommand struct {
	TagID int `json:"tag_id" binding:"required"`
}

type RemoveTagFromQuoteCommand struct {
	TagID int `json:"tag_id" binding:"required"`
}
