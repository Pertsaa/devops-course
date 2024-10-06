package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type Info struct {
	Service1 ServiceInfo `json:"service_1"`
	Service2 ServiceInfo `json:"service_2"`
}

type ServiceInfo struct {
	Hostname  string `json:"hostname"`
	Uptime    string `json:"uptime"`
	DiskInfo  string `json:"diskInfo"`
	Processes string `json:"processes"`
}

type Error struct {
	Error string `json:"error"`
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /", infoHandler)

	server := http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	fmt.Println("API listening on port 8081...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not Found")
		return
	}

	hostname, err := getHostname()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get hostname"})
	}

	uptime, err := getUptime()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get uptime"})
	}

	diskInfo, err := getDiskInfo()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get disk info"})
	}

	processes, err := getProcessInfo()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get process info"})
	}

	writeJSON(w, http.StatusOK, Info{
		Service1: ServiceInfo{
			Hostname:  hostname,
			Uptime:    uptime,
			DiskInfo:  diskInfo,
			Processes: processes,
		},
		Service2: ServiceInfo{},
	})
}

func writeJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func getHostname() (string, error) {
	cmd := exec.Command("hostname")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", "")
	return output, nil
}

func getUptime() (string, error) {
	cmd := exec.Command("uptime")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", "")
	return output, nil
}

func getDiskInfo() (string, error) {
	cmd := exec.Command("df", "-h")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", "")
	return output, nil
}

func getProcessInfo() (string, error) {
	cmd := exec.Command("ps", "-ax")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", "")
	return output, nil
}
