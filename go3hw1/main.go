package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

/*
В контексте приглашений на мероприятия, RSVP — это запрос подтверждения от приглашённого человека или людей.
RSVP — это акроним французской фразы Répondez s’il vous plaît,
означающей буквально «Будьте добры ответить» или «Пожалуйста, ответьте».
*/

type Rsvp struct {
	Name, Email, Phone, Company string
	MailOrCall                  bool
}

var responces = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 5)

func LoadTemplates() {
	// TODO - load template here
	// There are 5 templates: welcome.html, form.html, list.html, thanks.html, sorry.html + base template(layout.html)
	templateNames := [5]string{"welcome", "form", "list", "resume", "thanks"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

// welcomeHandler handles root URL
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	templates["welcome"].Execute(w, nil)
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	templates["resume"].Execute(w, nil)
}

// listHandler handles /list URL
func listHandler(w http.ResponseWriter, r *http.Request) {
	templates["list"].Execute(w, responces)
}

type formData struct {
	*Rsvp
	Errors []string
}

// formHandler handles /form URL
func formHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Результат работы - пустая форма
	// GET localhost/form
	// Взять данные из запроса(request), проверить, что данные не пустые и добавить очередное приглашение в список
	// POST localhost/form

	if r.Method == http.MethodGet {
		templates["form"].Execute(w, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})

	} else if r.Method == http.MethodPost {
		r.ParseForm() // Парсим данные из request и записываем их в request.Form
		responceData := Rsvp{
			Name:       r.FormValue("name"),
			Email:      r.Form["email"][0], // the same as above
			Phone:      r.FormValue("phone"),
			Company:    r.FormValue("company"),
			MailOrCall: r.FormValue("mailorcall") == "true",
		}

		errors := []string{}
		// Проверка значений полей формы. Пустые поля недопустимы.
		if responceData.Name == "" {
			errors = append(errors, "Please, enter your name!")
		}
		if responceData.Email == "" {
			errors = append(errors, "Please, enter your email!")
		}
		if responceData.Phone == "" {
			errors = append(errors, "Please, enter your phone!")
		}
		if responceData.Company == "" {
			errors = append(errors, "Please, enter your company name!")
		}
		if len(errors) > 0 {
			templates["form"].Execute(w, formData{
				Rsvp: &responceData, Errors: errors,
			})
		} else {
			responces = append(responces, &responceData)
			templates["thanks"].Execute(w, responceData.Name)
		}

	}

}

func main() {
	LoadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/resume", resumeHandler)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
