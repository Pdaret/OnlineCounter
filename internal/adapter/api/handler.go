package api

import (
	"encoding/json"
	"log"
	"net/http"

	"x-ui-monitor/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/v4/mem"
)

type Server struct {
	userUsecase *usecase.UserUsecase
}

func StartServer(userUsecase *usecase.UserUsecase) {
	server := &Server{userUsecase: userUsecase}

	router := mux.NewRouter()
	router.HandleFunc("/users/total", server.GetUserCount).Methods("GET")
	router.HandleFunc("/users/count/{inbound}", server.GetUserCountByInbound).Methods("GET")
	router.HandleFunc("/users/list/{inbound}", server.GetActiveIPsByInbound).Methods("GET")
	router.HandleFunc("/system/ram-usage", server.CheckHighRAMUsage).Methods("GET")

	log.Println("🌐 API ready on :5000")
	http.ListenAndServe(":5000", router)
}

func (s *Server) GetUserCount(w http.ResponseWriter, r *http.Request) {
	result, err := s.userUsecase.GetTotalUsersCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetUserCountByInbound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inbound := vars["inbound"]

	result, err := s.userUsecase.GetTotalUsersByInbound(inbound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"online_users": result})
}

func (s *Server) GetActiveIPsByInbound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inbound := vars["inbound"]

	ips, err := s.userUsecase.GetActiveIpsByInbound(inbound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string][]string{"active_ips": ips})
}

func (s *Server) CheckHighRAMUsage(w http.ResponseWriter, r *http.Request) {
	usageHigh, err := isRAMUsageHigh()
	if err != nil {
		http.Error(w, "Failed to get RAM usage", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"high_usage": usageHigh})
}

func isRAMUsageHigh() (bool, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return false, err
	}

	return vmStat.UsedPercent > 90.0, nil
}
