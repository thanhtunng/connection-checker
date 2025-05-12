package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type PortCheckRequest struct {
	IPs  []string `json:"ips"`
	Port int      `json:"port"`
}

type PortCheckResponse struct {
	Results []PortResult `json:"results"`
}

type PortResult struct {
	IP     string `json:"ip"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func resolveDomain(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}
	
	var ipStrings []string
	for _, ip := range ips {
		// Only include IPv4 addresses
		if ipv4 := ip.To4(); ipv4 != nil {
			ipStrings = append(ipStrings, ipv4.String())
		}
	}
	return ipStrings, nil
}

func checkPort(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func handlePortCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PortCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := make([]PortResult, 0)
	for _, input := range req.IPs {
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Try to resolve as domain name first
		ips, err := resolveDomain(input)
		if err != nil {
			// If domain resolution fails, treat as IP address
			ip := net.ParseIP(input)
			if ip == nil {
				results = append(results, PortResult{
					IP:     input,
					Status: "error",
					Error:  "Invalid IP address or domain name",
				})
				continue
			}
			ips = []string{input}
		}

		// Check port for each resolved IP
		for _, ip := range ips {
			status := "closed"
			if checkPort(ip, req.Port) {
				status = "open"
			}

			result := PortResult{
				IP:     ip,
				Status: status,
			}
			if input != ip {
				result.IP = fmt.Sprintf("%s (%s)", input, ip)
			}
			results = append(results, result)
		}
	}

	response := PortCheckResponse{
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Wrap the handler with CORS middleware
	http.HandleFunc("/api/check-ports", enableCORS(handlePortCheck))

	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
} 