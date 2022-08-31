package entity

type User struct {
	Phone          string   `json:"phone" form:"phone"`
	Password       string   `json:"password" form:"password"`
	Uid            string   `json:"uid" form:"uid"`
	Tid            string   `json:"tid" form:"tid"`
	Name           string   `json:"name" form:"name"`
	AuthorizedPids []string `json:"authorizedPids" form:"authorizedPids"`
	Resources      []string `json:"resources" form:"resources"`
}
