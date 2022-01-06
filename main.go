package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
)

type Course struct {
	ID     string  `json:"_id"`
	Name   string  `json:"courseName"`
	Price  string  `json:"price"`
	Author *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullName"`
	Website  string `json:"website"`
}

//mock db
var courseDB []Course

//middleware
func (c *Course) isEmpty() bool {
	return c.Name == ""
}

func main() {
	fmt.Println("Server running on localhost port 4000")
	//seeed mock database
	courseDB = append(courseDB, Course{
		ID:    "1",
		Name:  "Learn Go",
		Price: "30$",
		Author: &Author{
			FullName: "John Doe",
			Website:  "JohnDoe.dev.com",
		},
	})
	r := mux.NewRouter()
	r.HandleFunc("/getAllCourse", getAllCourses).Methods("GET")
	r.HandleFunc("/getCourseByID/{_id}", getCourseById).Methods("GET")
	r.HandleFunc("/createCourse", createCourse).Methods("POST")
	r.HandleFunc("/updateCourse/{_id}", updateCourse).Methods("POST")
	r.HandleFunc("/deleteCourse/{_id}", deleteCourse).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", r))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")
	json.NewEncoder(w).Encode(courseDB)
}

func getCourseById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")
	params := mux.Vars(r)

	for _, course := range courseDB {
		if course.ID == params["_id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No Course with this id")

}

func createCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Error creating course")
		return
	}

	var course Course

	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.isEmpty() {
		json.NewEncoder(w).Encode("There is no json data in request")
		return
	}

	//generate uuid
	id, _ := uuid.NewV4()
	course.ID = id.String()
	courseDB = append(courseDB, course)
	json.NewEncoder(w).Encode("course created")

}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")

	params := mux.Vars(r)

	for index, course := range courseDB {
		if course.ID == params["_id"] {
			courseDB = append(courseDB[:index], courseDB[index+1:]...)
			var course Course
			json.NewDecoder(r.Body).Decode(&course)
			id, _ := uuid.NewV4()
			course.ID = id.String()
			courseDB = append(courseDB, course)
			json.NewEncoder(w).Encode("Course updated")
			return
		}
	}
	json.NewEncoder(w).Encode("course ID does not exist")

}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")
	courseID := mux.Vars(r)["_id"]

	for index, course := range courseDB {
		if course.ID == courseID {
			courseDB = append(courseDB[:index], courseDB[index+1:]...)
			json.NewEncoder(w).Encode("Course deleted")
			return
		}
	}

	json.NewEncoder(w).Encode("ID does not exist")
}
