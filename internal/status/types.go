package status

type Job struct {
	ID        string `json:"id" dynamodbav:"id"`
	GameName  string `json:"game_name" dynamodbav:"game_name"`
	Status    string `json:"status" dynamodbav:"status"`
	Progress  int    `json:"progress" dynamodbav:"progress"`
	Total     int    `json:"total" dynamodbav:"total"`
	StartedAt int64  `json:"started_at" dynamodbav:"started_at"`
	UpdatedAt int64  `json:"updated_at" dynamodbav:"updated_at"`
}
