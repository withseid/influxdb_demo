package models

type CommonDocument struct {
	CreatedAt    string   `json:"created_at,omitempty"`
	EntityID     string   `json:"kw.entity_id,omitempty"`
	Keywords     []string `json:"ik.keywords,omitempty"`
	PrimaryName  string   `json:"ik.primary_name,omitempty"`
	RankingScore float32  `json:"ranking_score"`
}

type EnterpriseDocument struct {
	CommonDocument
	PE float32 `json:"pe"`
	PB float32 `json:"pb"`
	PS float32 `json:"ps"`
}
