package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"Nova/handlers"
)

// ServeStatic serves the index.html file
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// LocationHandler handles the IP location request
func LocationHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "IP address not provided", http.StatusBadRequest)
		return
	}

	location := getIPLocation(ip)

	// Respond with the location info in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"location": location,
	})
}

// getIPLocation gets the location of the IP address using external or internal methods.
func getIPLocation(ip string) map[string]interface{} {
	// Using the ipinfo.io API to get location information
	resp, err := http.Get(fmt.Sprintf("http://ipinfo.io/%s/geo", ip))
	if err != nil {
		return map[string]interface{}{"error": "Error fetching location"}
	}
	defer resp.Body.Close()

	// Replace ioutil.ReadAll with io.ReadAll as ioutil is deprecated
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"error": "Error reading response body"}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]interface{}{"error": "Error parsing location data"}
	}

	// Return the location information
	return result
}

func main() {
	log.Println("📢 Starting Security Scanner...")

	// Serve the main page correctly
	http.HandleFunc("/", ServeStatic)

	// Handlers for different functionalities
	http.HandleFunc("/scan", handlers.ScanHandler)
	http.HandleFunc("/reports", handlers.ReportsHandler)
	http.HandleFunc("/logs", handlers.LogsHandler)
	//http.HandleFunc("/api/reports", handlers.FetchReportsHandler)

	// DNS Checker API (Ensure no redirect issue)
	http.HandleFunc("/dns-check", handlers.DNSHandler)

	// New location route for fetching IP location
	http.HandleFunc("/location", LocationHandler)

	// DNS check page route
	http.HandleFunc("/dnscheck", ServeStatic)

	// Serve static files (JavaScript, CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	log.Println("🚀 Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
