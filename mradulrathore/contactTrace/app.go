package main

import (
	"fmt"
	"html/template"
	"log"
	"mradulrathore/models/contacts"
	"mradulrathore/models/user"
	"net/http"
	"path"

	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func users(w http.ResponseWriter, r *http.Request) {
	book := string("")

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if r.Method == "POST" {
		r.ParseForm()       // parse arguments, you have to call this by yourself
		fmt.Println(r.Form) // print form information in server side
		fmt.Println("path", r.URL.Path)
		fmt.Println("scheme", r.URL.Scheme)
		fmt.Println(r.Form["url_long"])

		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
		layout := "2006-01-02"
		str := r.FormValue("birthDate")
		dateOfBirth, err := time.Parse(layout, str)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("dfsd")

		newUser := user.User{
			ID:          primitive.NewObjectID(),
			Name:        r.FormValue("Name"),
			BirthDate:   dateOfBirth,
			PhoneNumber: r.FormValue("phoneNumber"),
			Email:       r.FormValue("email"),
			CreatedOn:   time.Now(),
		}
		user.CreateUser(newUser)
	} else {
		http.Redirect(w, r, "/users", http.StatusFound)
	}
}

func success(w http.ResponseWriter, r *http.Request) {

	book := string("")

	fp := path.Join("templates", "success.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

}

func contact(w http.ResponseWriter, r *http.Request) {
	book := string("")

	fp := path.Join("templates", "contact.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if r.Method == "POST" {

		r.ParseForm() // parse arguments, you have to call this by yourself

		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
		newContact := contacts.Contact{
			UserIDOne: r.FormValue("userIdOne"),
			UserIDTwo: r.FormValue("userIdTwo"),
			Timestamp: time.Now(),
		}

		contacts.CreateContact(newContact)

	} else {
		http.Redirect(w, r, "/users", http.StatusFound)
	}

}

func getPrimaryContacts(w http.ResponseWriter, r *http.Request) {
	book := string("")

	fp := path.Join("templates", "getUserDetails.html")
	tmpl, err := template.ParseFiles(fp)
	affectedContacts := []contacts.Contact{}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if r.Method == "GET" {

		r.ParseForm() // parse arguments, you have to call this by yourself
		id := r.URL.Path[len("/contact/?user="):]
		i := strings.Index(r.URL.Path, "&")
		infectionTimestamp := r.URL.Path[i+len("infection_timestamp="):]
		layout := "2006-01-02"
		endRange, err := time.Parse(layout, infectionTimestamp)
		startRange := endRange.Add(14 * time.Hour)
		if err != nil {
			log.Fatal(err)
			return
		}

		res, err := contacts.GetAllContacts(id)

		for i, v := range res {
			if v.Timestamp.After(startRange) && v.Timestamp.Before(endRange) {
				fmt.Println(i)
				fmt.Println(v.UserIDTwo)
				affectedContacts = append(affectedContacts, v)
			}
		}

	} else {
		http.Redirect(w, r, "/users", http.StatusFound)
	}

}

func getUser(w http.ResponseWriter, r *http.Request) {
	book := string("")
	id := r.URL.Path[len("/users/"):]
	fp := path.Join("templates", "getUserDetails.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if r.Method == "GET" {
		r.ParseForm() // parse arguments, you have to call this by yourself

		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}

		res, err := user.GetUsersByID(id)
		v := res
		if err == nil {
			fmt.Fprintf(w, "%s %s %s %s %s\n", v.ID.String(), v.Name, v.BirthDate, v.Email, v.PhoneNumber)
		}
	} else {
		http.Redirect(w, r, "/users/", http.StatusFound)
	}

}

func handleRequests() {
	http.HandleFunc("/users", users) // set router
	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/success", success)              // set router
	http.HandleFunc("/contacts", contact)             // set router
	http.HandleFunc("/contacts/", getPrimaryContacts) // set router
	err := http.ListenAndServe(":9090", nil)          // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {

	// newUser1 := user.User{
	// 	ID:          primitive.NewObjectID(),
	// 	Name:        "Gaurav Rathore",
	// 	BirthDate:   getDOB(2000, 12, 16),
	// 	PhoneNumber: "8933810311",
	// 	Email:       "er.mradulrathore@gmail.com",
	// 	CreatedOn:   time.Now(),
	// }
	// newUser2 := user.User{
	// 	ID:          primitive.NewObjectID(),
	// 	Name:        "Pinkesh Rathore",
	// 	BirthDate:   getDOB(2000, 2, 8),
	// 	PhoneNumber: "9827049691",
	// 	Email:       "ajayrathore@gmail.com",
	// 	CreatedOn:   time.Now(),
	// }
	// newContact := contacts.Contact{
	// 	UserIDOne: newUser1.ID,
	// 	UserIDTwo: newUser2.ID,
	// 	Timestamp: time.Now(),
	// }
	// user.DeleteAllUsers()
	//contacts.DeleteAllContacts()
	// user.CreateUser(newUser1)
	// user.CreateUser(newUser2)
	//contacts.CreateContact(newContact)
	//users, _ := user.GetAllUsers()
	// PrintList(users)

	// users = []user.User{
	// 	user.User{ID: primitive.NewObjectID(), Name: "Ajay Rathore", BirthDate: getDOB(1998, 2, 8), PhoneNumber: "32410311", Email: "ajayrathore@gmail.com", CreatedOn: time.Now()},
	// 	user.User{ID: primitive.NewObjectID(), Name: "Mitali Rathore", BirthDate: getDOB(1972, 2, 8), PhoneNumber: "842310311", Email: "mitalirathore@gmail.com", CreatedOn: time.Now()},
	// }
	// user.CreateMany(users)
	// user.GetAllUsers()
	handleRequests()
}
