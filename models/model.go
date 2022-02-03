package model

import (
	"fmt"
	"net/http"
)

type AccountData struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
}

type AccountDataList struct {
	Accounts []AccountData `json:"items"`
}

func (a *AccountData) Bind(r *http.Request) error {
	if a.Name == "" || a.Currency == "" || a.Balance == 0 {
		return fmt.Errorf("name is a requred field")
	}

	return nil
}

func (*AccountDataList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*AccountData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
