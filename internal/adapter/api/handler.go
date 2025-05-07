package api

import (
	"encoding/json"
	"net/http"

	"x-ui-monitor/internal/usecase"

	"github.com/gorilla/mux"
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
