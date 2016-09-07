package types

type BitBucketRequest struct {
	Actor       *Actor      `json:"actor"`
	Push        *Push       `json:"push"`
	Repository  *Repository `json:"repository"`
	AsanaTaskID string      `json:"task_id"`
}
type Actor struct {
	Username string `json:"username"`
	Links    *Links `json:"links"`
}
type Links struct {
	HTML *HTML `json:"html"`
}
type HTML struct {
	Href string `json:"href"`
}
type Push struct {
	Changes []*Change `json:"changes"`
}
type Change struct {
	New   *New   `json:"new"`
	Links *Links `json:"links"`
}
type New struct {
	Target *Target `json:"target"`
}
type Target struct {
	Date    string  `json:"date"`
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Author  *Author `json:"author"`
}
type Author struct {
	User *User `json:"user"`
}
type User struct {
	Username string `json:"username"`
	Type     string `json:"type"`
}
type Repository struct {
	Name string `json:"name"`
}
