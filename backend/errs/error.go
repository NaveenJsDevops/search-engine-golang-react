package errs

// Error is to capture error
type Error struct {
	Code    int    `json:"errCode"`
	Message string `json:"message"`
	Err     error  `json:"error"`
	Module  string `json:"-"`
	IsDbErr bool   `json:"-"`
}
