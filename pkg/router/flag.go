package router

type FlagSet struct {
	flags []*Flag
}

type Flag struct {
	Name        string
	Description string
	Value       string
}
