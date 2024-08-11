package login

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // token omitido se não existir
}
