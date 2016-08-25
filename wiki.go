package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	templates = template.Must(template.ParseFiles("./views/edit.html", "./views/view.html", "./views/list.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

// Page is the struct of the page on wiki
type Page struct {
	Title string
	Body  []byte
}

// save saves the page on a text file
func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage loads the page from the text file
func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// viewHandler allow users to view a wiki page.
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil { // Si no encuentra la página, se irá a edición para crear una nueva
		// The http.Redirect function adds an HTTP status code of http.StatusFound (302) and a Location header to the HTTP response.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(&w, "view", p)
}

// editHandler loads the page (or, if it doesn't exist, create an empty Page struct), and displays an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil { // Si no encuentra la página, creará una nueva
		p = &Page{Title: title}
	}
	renderTemplate(&w, "edit", p)
}

// saveHandler saves the information from the edit form
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		handleCommonErrors(err, &w)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// listHandler show a list with all pages names.
func listHandler(w http.ResponseWriter, r *http.Request) {
	pageNames, err := listPages()
	if err != nil {
		handleCommonErrors(err, &w)
	}
	err = templates.ExecuteTemplate(w, "list.html", pageNames)
	if err != nil {
		handleCommonErrors(err, &w)
	}
}

// listPages list all pages on wiki
func listPages() ([]string, error) {
	var names = make([]string, 0)
	files, err := ioutil.ReadDir("./data")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		names = append(names, file.Name()[:len(file.Name())-4])
	}
	return names, nil
}

// renderTemplate refactor to render templates
func renderTemplate(w *http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(*w, tmpl+".html", p)
	if err != nil {
		handleCommonErrors(err, w)
	}
}

// makeHandler returns a function that wrap the edit, view and save functions
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// handleCommonErrors handle errors and write the error at web page
func handleCommonErrors(err error, w *http.ResponseWriter) {
	http.Error(*w, err.Error(), http.StatusInternalServerError)
}

// main executes the program and serve the web server.
func main() {
	fmt.Println("Servidor ejecutandose en: http://localhost:8080")
	fmt.Println("Para ver el contenido digite view/tuarticulo")
	fmt.Println("Para salir presione Ctrl+C")
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/list/", listHandler)
	http.ListenAndServe(":8080", nil)
}
