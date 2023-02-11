package user

type identity struct {
	Provider string `json:"provider,omitempty"`
	UID      string `json:"uid,omitempty"`
}

type userResponse struct {
	Name       string     `json:"name"`
	Role       string     `json:"role,omitempty"`
	Identities []identity `json:"identities,omitempty"`
}
