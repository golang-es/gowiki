package main

import (
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
	if err != nil { // Si no encuentra la página, se irá a edición para crear una nueva
		// The http.Redirect function adds an HTTP status code of http.StatusFound (302) and a Location header to the HTTP response.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
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

// saveHandler saves the information from the edit form
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		handleCommonErrors(err, &w)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// renderTemplate refactor to render templates
func renderTemplate(w *http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("views/" + tmpl + ".html")
	if err != nil {
		handleCommonErrors(err, w)
		return
	}
	err = t.Execute(*w, p)
	if err != nil {
		handleCommonErrors(err, w)
	}
}

// handleCommonErrors handle errors and write the error at web page
func handleCommonErrors(err error, w *http.ResponseWriter) {
	http.Error(*w, err.Error(), http.StatusInternalServerError)
}

// main executes the program and serve the web server.
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
