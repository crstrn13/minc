package types

type CreateType struct {
	Name                string
	Provider            string
	UShiftVersion       string
	UShiftImage         string
	UShiftConfig        string
	HTTPSPort           int
	HTTPPort            int
	DisableOverlayCache bool
}

type StatusType struct {
	Container string `json:"container"`
	APIServer string `json:"apiserver"`
	Error     string `json:"error,omitempty"`
}
