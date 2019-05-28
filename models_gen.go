// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package drophere_go

type Link struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	IsProtected bool    `json:"isProtected"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	Deadline    *string `json:"deadline"`
}

type Message struct {
	Message string `json:"message"`
}

type Token struct {
	LoginToken string `json:"loginToken"`
}

type User struct {
	ID            int     `json:"id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Dropboxauth   *bool   `json:"dropboxauth"`
	Dropboxemail  *string `json:"dropboxemail"`
	Dropboxavatar *string `json:"dropboxavatar"`
}