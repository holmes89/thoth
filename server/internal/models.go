package internal

type Game struct {
	SK    string  `json:"-"`
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Path  *string `json:"path,omitempty"`
}

type Room struct {
	SK          string `json:"-"`
	ID          string `json:"id"`
	Host        string `json:"host"`
	ViewerCount int    `json:"viewers"`
}

type User struct {
	SK          string `json:"-"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type ErrorObj struct {
	Description *string `json:"description,omitempty"`
	Id          int     `json:"id"`
	Name        string  `json:"name"`
}
