package handlers

import (
	"encoding/json"
	"net"
	"net/http"
)

// DNSResponse holds the response structure
type DNSResponse struct {
	Domain     string   `json:"domain"`
	RecordType string   `json:"recordType"`
	Records    []string `json:"records,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// DNSHandler processes DNS queries (Updated to GET instead of POST)
func DNSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { // Change from POST to GET
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get domain and recordType from query parameters
	domain := r.URL.Query().Get("domain")
	recordType := r.URL.Query().Get("recordType")

	// Check if both parameters are provided
	if domain == "" || recordType == "" {
		http.Error(w, "Missing domain or recordType parameter", http.StatusBadRequest)
		return
	}

	var records []string
	var err error

	switch recordType {
	case "A":
		records, err = getARecords(domain)
	case "CNAME":
		records, err = getCNAMERecords(domain)
	default:
		http.Error(w, "Invalid record type", http.StatusBadRequest)
		return
	}

	response := DNSResponse{Domain: domain, RecordType: recordType, Records: records}

	if err != nil {
		response.Error = err.Error()
	}

	// Ensure JSON response is set
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getARecords fetches A records for a domain
func getARecords(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	records := make([]string, len(ips))
	for i, ip := range ips {
		records[i] = ip.String()
	}

	return records, nil
}

// getCNAMERecords fetches CNAME records for a domain
func getCNAMERecords(domain string) ([]string, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return nil, err
	}

	return []string{cname}, nil
}
