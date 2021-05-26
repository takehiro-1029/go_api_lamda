package common

import (
	"encoding/json"
)

type Request struct {
	UserName        string `json:"user_name,omitempty"`
	UserEmail       string `json:"user_email,omitempty"`
	UserPass        string `json:"user_pass,omitempty"`
	UserAuthCode    string `json:"user_code,omitempty"`
	UserToken       string `json:"user_token,omitempty"`
	CognitoClientID string `json:"client_id,omitempty"`
}

func ConvertRequestToJSON(req *Request, inputs string) error {
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return err
	}
	return nil
}
