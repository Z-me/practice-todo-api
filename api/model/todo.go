package model

type Todo struct {
  ID        string  `json:"id"`
  Title     string  `json:"title"`
  Status    string  `json:"status"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority"`
}
