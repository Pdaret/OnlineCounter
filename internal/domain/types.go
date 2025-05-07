package domain

type TotalUsersCountResult struct {
	Total    int            `json:"total_users"`
	Inbounds map[string]int `json:"inbounds"`
}
