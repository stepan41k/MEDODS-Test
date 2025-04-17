package models

type User struct {
	GUID []byte `json:"guid"`
	IP   string `json:"ip"`
}
