package dtos

type GroupCreateRequest struct {
	SchoolID int `json:"school_id"`
	Name string `json:"name"`
}

type GroupResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
}