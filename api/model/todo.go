package model

type Status string

type Todo struct {
  ID        int     `json:"id"`
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

type Payload struct {
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

type StatusPayload struct {
  Status    string `json:"status" binding:"required"`
}
