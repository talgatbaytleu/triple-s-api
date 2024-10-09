package triples

import (
	"log"
	"net/http"
)

var (
	path string
	port string = "8080"
)

// type Page struct {
// 	Title string
//
// 	Body []byte
// }
//
// func (p *Page) save() error {
// 	filename := p.Title + ".txt"
//
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
//
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &Page{Title: title, Body: body}, nil
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, _ := template.ParseFiles(tmpl + ".html")
// 	t.Execute(w, p)
// }

func putBucketHendler(w http.ResponseWriter, r *http.Request) {
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL.Path[1:]

	// HERE MUST BE PATH VALIDATION

	switch r.Method {
	case "PUT":

	case "GET":

	case "DELETE":
	}
	// if err != nil {
	// 	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	// 	return
	// }
	// renderTemplate(w, "view", p)
}

func Run() {
	http.HandleFunc("/", Handler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
