package triples

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

func JustHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test1")
}

func JustHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test2")
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("put handler")

	// HERE MUST BE PATH VALIDATION
	// if err != nil {
	// 	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	// 	return
	// }
	// renderTemplate(w, "view", p)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get handler")
	fullPath := r.URL.Path[1:]
	fmt.Println(fullPath)
	fmt.Println(strings.Split(fullPath, "/"))
	fmt.Println(len(strings.Split(fullPath, "/")))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete handler")
}

func Run() {
	http.HandleFunc("/", JustHandler2)
	http.HandleFunc("PUT /", PutHandler)
	http.HandleFunc("GET /", GetHandler)
	http.HandleFunc("DELETE /", DeleteHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
