package models

type API3XUIOnlinesResp struct {
	Success bool     `json:"success"`
	Msg     string   `json:"msg"`
	Emails  []string `json:"obj"`
}

type API3XUIInboundsResp struct {
	Success  bool             `json:"success"`
	Msg      string           `json:"msg"`
	Inbounds []API3XUIInbound `json:"obj"`
}

type API3XUIInbound struct {
	ID             int                  `json:"id"`
	UpTraffic      int                  `json:"up"`
	DownTraffic    int                  `json:"down"`
	TotalTraffic   int                  `json:"total"`
	Remark         string               `json:"remark"`
	Enable         bool                 `json:"enable"`
	ExpiryTime     int                  `json:"expiryTime"`
	ClientStats    []API3XUIClientStats `json:"clientStats"`
	Listen         string               `json:"listen"`
	Port           int                  `json:"port"`
	Protocol       string               `json:"protocol"`
	Settings       string               `json:"settings"`       // ???
	StreamSettings string               `json:"streamSettings"` // ???
	Tag            string               `json:"tag"`
	Sniffing       string               `json:"sniffing"` // ???
	Allocate       string               `json:"allocate"` // ???
}

type API3XUIClientStats struct {
	ID           int    `json:"id"`
	InboundID    int    `json:"inboundId"`
	Enable       bool   `json:"enable"`
	Email        string `json:"email"`
	UpTraffic    int    `json:"up"`
	DownTraffic  int    `json:"down"`
	ExpiryTime   int    `json:"expiryTime"`
	TotalTraffic int    `json:"total"`
	Reset        int    `json:"reset"`
}
