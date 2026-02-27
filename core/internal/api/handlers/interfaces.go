package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
)

func InterfacesHandler(w http.ResponseWriter, r *http.Request) {

	ifaces, err := system.GetInterfaces()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ifaces)
}

type InterfaceEditRequest struct {
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

func EditInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	return
    }
	//debe existir candidate
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	println("Candidate created, interfaces:", len(fw.Interfaces))

	//obtener nombre interfaz desde URL
     parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

     if len(parts) < 3 {
     	http.Error(w, "invalid path", http.StatusBadRequest)
	  return
      }

     ifaceName := parts[len(parts)-1]


     println("fw nil?", fw == nil)
     println("interfaces nil?", fw.Interfaces == nil)
     println("ifaceName:", ifaceName)

	 iface := fw.Interfaces[ifaceName]
	iface, ok := fw.Interfaces[ifaceName]
	if !ok {
		http.Error(w, "interface not found", 404)
		return
	}
      
	var req InterfaceEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	//aplicar cambios solo en candidate
	if req.IP != "" {
		iface.IP = req.IP
	}
	if req.Zone != "" {
		iface.Zone = req.Zone
	}
	if req.State != "" {
		iface.State = req.State
	}
	if req.Management != nil {
		iface.Management = req.Management
	}

	w.Write([]byte("interface updated in candidate"))
}