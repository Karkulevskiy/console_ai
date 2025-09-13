package domain

import (
	"fmt"
	"strings"
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
	var sb strings.Builder
	if strings.TrimSpace(r.Output) != "" {
		sb.WriteString(fmt.Sprintf("\n\t\"output\": %s", r.Output))
	}
	if strings.TrimSpace(r.Lang) != "" {
		sb.WriteString(fmt.Sprintf("\n\t\"lang\": %s", r.Lang))
	}
	if strings.TrimSpace(r.Code) != "" {
		sb.WriteString(fmt.Sprintf("\n\t\"code\":\n%s", r.Code))
	}
	return sb.String()
}
