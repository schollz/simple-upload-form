package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type JsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := JsonResponse{
			Success: false,
			Message: "Method not allowed",
		}
		sendJsonResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	// Parse the form data, including the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		response := JsonResponse{
			Success: false,
			Message: "Unable to parse form",
		}
		sendJsonResponse(w, http.StatusBadRequest, response)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response := JsonResponse{
			Success: false,
			Message: "Unable to get file from form",
		}
		sendJsonResponse(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// Create a new file on the server
	dst, err := os.Create(handler.Filename)
	if err != nil {
		response := JsonResponse{
			Success: false,
			Message: "Unable to create file on server",
		}
		sendJsonResponse(w, http.StatusInternalServerError, response)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file on the server
	_, err = io.Copy(dst, file)
	if err != nil {
		response := JsonResponse{
			Success: false,
			Message: "Unable to save file on server",
		}
		sendJsonResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := JsonResponse{
		Success: true,
		Message: fmt.Sprintf("File uploaded successfully: %s", handler.Filename),
	}
	sendJsonResponse(w, http.StatusOK, response)
}

func sendJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	port := 8080
	fmt.Printf("Server listening on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
