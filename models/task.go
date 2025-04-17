package models

type Task struct {
	ID          int     `json:"id"`
	Title       *string `json:"title"`       // pointer to detect null/missing
	Description *string `json:"description"` // pointer type for optional updates
	Completed   *bool   `json:"completed"`
	MainTask	*int 	`json:"maintask"`
}

/*

Using the pointer we can destinguish between

"title": ""
"title": null
missing "title" (means: don’t update it at all)

*/
