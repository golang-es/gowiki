package main

import (
	"fmt"
	"io/ioutil"
)

// Page is the struct of the page on wiki
type Page struct {
	Title string
	Body  []byte
}

// save saves the page on a text file
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage loads the page from the text file
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// main executes the program and serve the web server.
func main() {
	p1 := &Page{Title: "PruebaPagina", Body: []byte("Esta es una página de ejemplo")}
	p1.save()
	p2, _ := loadPage(p1.Title)
	fmt.Println(string(p2.Body))
}
