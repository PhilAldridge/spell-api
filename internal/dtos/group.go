package dtos

type GroupCreateRequest struct {
	SchoolID int `json:"school_id"`
	Name string `json:"name"`
}