package response

type Login struct {
	Token string `json:"token"`
}

type RefreshToken struct {
	Token string `json:"token"`
}
