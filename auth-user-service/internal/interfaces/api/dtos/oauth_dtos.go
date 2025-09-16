package dtos

type OAuthInitiateResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

type OAuthCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

type OAuthUserInfo struct {
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}
