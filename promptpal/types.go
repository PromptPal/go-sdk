package promptpal

type Configuration struct {
	Input struct {
		HTTP *struct {
			URL   string `json:"url"`
			Token string `json:"token"`
		} `json:"http,omitempty"`
	} `json:"input"`
	Output struct {
		Schema  string `json:"schema"`
		GoTypes *struct {
			Prefix      string `json:"prefix"`
			PackageName string `json:"package_name"`
			Output      string `json:"output"`
		} `json:"go_types,omitempty"`
		TypeScriptTypes *struct {
			Prefix string `json:"prefix"`
			Output string `json:"output"`
		} `json:"typescript_types,omitempty"`
	} `json:"output"`
}
