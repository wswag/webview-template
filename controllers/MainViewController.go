package controllers

import (
	"fmt"
)

// MainViewController represents the controller interface for the main view
type MainViewController struct {
	SynchronizableModel

	Title string `json:"title"`
	Answer int `json:"answer"`
}

// Hi is a Smoke Test
func (m *MainViewController) Hi() {
	m.Answer++
	m.Title += " Hi!"
	fmt.Println("Invoked MainViewController.Hi()")
}