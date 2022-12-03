package common

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	NoSub bool   `json:"no_sub"`
}
