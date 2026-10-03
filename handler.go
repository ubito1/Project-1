package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type PageData struct {
	Result     string
	Input      string
	Banner     string
	Error      string 
	StatusCode int 
}
// 
var templates = template.Must(template.ParseFiles("./Templates/index.html"))

// Read "text" and "banner" from the form values
// Validate: if text is empty, return 400 Bad Request
// Validate: if banner is not one of the three valid names, return 400 Bad Request

func renderError(w http.ResponseWriter, status int, input, banner, msg string) {
	w.WriteHeader(status)
	Data := PageData{
		Input:      input,
		Banner:     banner,
		Error:      msg,
		StatusCode: status,
	}
	if err := templates.ExecuteTemplate(w, "index.html", Data); err != nil {
		log.Printf("Failed to render error template: %v", err)
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	// http.MethodGet litterly means = "GET"
	if req.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.URL.Path != "/" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	var pageData PageData
	if err := templates.ExecuteTemplate(w, "index.html", pageData); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func BannerHandler(w http.ResponseWriter, req *http.Request) {

		// Check that the method is POST
	// If not, return 400 Bad Request
	if req.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// FormValue means go to the HTML word and get me it's value "from the input<>"
	// PostValue("word") is prefered for some reason look it up >>>>
	input := req.FormValue("word")
	banner := req.FormValue("font")
//postvalue & FormValue are used to retrieve the values of form fields submitted via POST requests. In this case, input is assigned the value of the form field named "word", and banner is assigned the value of the form field named "font". These values are then used in the subsequent logic to process the ASCII art rendering based on user input and selected banner type.
	if strings.TrimSpace(input) == "" {
		renderError(w, http.StatusBadRequest, input, banner, "Input text is required")
		return
	}

	if banner == "" {
		renderError(w, http.StatusBadRequest, input, banner, "Banner selection is required")
		return
	}

	validBanners := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	
	
	if !validBanners[banner] {
		renderError(w, http.StatusBadRequest, input, banner, "Invalid banner type selected")
		return
	}

	bannerPath := banner + ".txt"

	bannerLines, err := LoadBanner(bannerPath)
	if err != nil {
		renderError(w, http.StatusInternalServerError, input, banner, "Internal server error: Banner file is missing or corrupted")
		return
	}

	result, err := Render(bannerLines, input)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			renderError(w, http.StatusBadRequest, input, banner, err.Error())
			return
		}
		renderError(w, http.StatusBadRequest, input, banner, "Invalid input characters")
		return
	}

	pageData := PageData{
		Result: result,
		Input:  input,
		Banner: banner,
	}

	if err := templates.ExecuteTemplate(w, "index.html", pageData); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}