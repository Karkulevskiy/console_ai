package domain

import (
	"fmt"
)

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
	Code   string `json:"code,omitempty"`
	Lang   string `json:"lang,omitempty"`
}

func (r Response) String() string {
	return fmt.Sprintf("\n\t\"output\": %s\n\t\"code\":\n\n%s\n", r.Output, r.Code)
}
