package domain

import "encoding/json"

type Request struct {
	Input string `json:"input,omitempty"`
	Model string `json:"model,omitempty"`
}

func (r Request) GetModel() string {
	if r.Model == "" {
		return "googleai/gemini-2.5-flash"
	}
	return r.Model
}

type Response struct {
	Output string `json:"output"`
	Code   string `json:"code"`
}

func (r Response) String() string {
	spaces4 := "    "
	bytes, _ := json.MarshalIndent(r, "", spaces4)
	return string(bytes)
}
