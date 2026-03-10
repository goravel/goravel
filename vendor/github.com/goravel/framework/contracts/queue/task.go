package queue

type Task struct {
	ChainJob
	UUID  string     `json:"uuid"`
	Chain []ChainJob `json:"chain"`
}
