package middleware

type SJWTClaims struct {
	Auth    bool   `json:"auth"`
	UserId  int    `json:"userid"`
	Role    int    `json:"role"`
	Service string `json:"service"`
	Hop     int    `json:"hop"`
}
