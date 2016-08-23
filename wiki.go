package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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

// viewHandler allow users to view a wiki page.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		handleCommonErrors(&err, &w)
	}
	renderTemplate(&w, "view", p)
}

// editHandler loads the page (or, if it doesn't exist, create an empty Page struct), and displays an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil { // Si no encuentra la página, creará una nueva
		p = &Page{Title: title}
	}
	renderTemplate(&w, "edit", p)
}

func renderTemplate(w *http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("views/" + tmpl + ".html")
	if err != nil {
		handleCommonErrors(&err, w)
	}
	t.Execute(*w, p)
}

// handleCommonErrors handle errors and write the error at web page
func handleCommonErrors(err *error, w *http.ResponseWriter) {
	fmt.Fprintf(*w, "<div class=\"error\">%s</div>", *err)
}

// main executes the program and serve the web server.
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
}
