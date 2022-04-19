package model

type Id int

type Status string

type Todo struct {
  ID        int     `json:"id"`
  Title     string  `json:"title"`
  Status    string  `json:"status"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority"`
}

type Payload struct {
  Title     string  `json:"title"`
  Status    string  `json:"status"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority"`
}

type StatusPayload struct {
  Status    string `json:"status"`
}
