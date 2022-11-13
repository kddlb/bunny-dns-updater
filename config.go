package main

type config struct {
	AccessKey string `validate:"required"`
	Zone      string `validate:"required,hostname"`
	Record    string `validate:"required,hostname"`
}
