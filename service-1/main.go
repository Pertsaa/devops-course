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
	DiskInfo  string `json:"disk_info"`
	Processes string `json:"process_info"`
}

type Error struct {
	Error string `json:"error"`
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /", infoHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("API listening on port 8080...")
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
		return
	}

	uptime, err := getUptime()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get uptime"})
		return
	}

	diskInfo, err := getDiskInfo()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get disk info"})
		return
	}

	processes, err := getProcessInfo()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get process info"})
		return
	}

	service2Info, err := fetchService2Info()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Error{Error: "failed to get service 2 info"})
		return
	}

	writeJSON(w, http.StatusOK, Info{
		Service1: ServiceInfo{
			Hostname:  hostname,
			Uptime:    uptime,
			DiskInfo:  diskInfo,
			Processes: processes,
		},
		Service2: service2Info,
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

func fetchService2Info() (ServiceInfo, error) {
	info := ServiceInfo{}

	r, err := http.Get("http://service-2:8081")
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		return info, err
	}
	defer r.Body.Close()

	return info, nil
}
