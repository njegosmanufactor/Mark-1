package Model

type User struct {
	Id           int
	Name         string
	Surename     string
	Username     string
	Password     string
	Company      string
	Email        string
	Applications []string
}
