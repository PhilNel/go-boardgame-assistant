package logger

import (
	"encoding/json"
	"log"
	"time"
)

type IncomingRequest struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	GameName  string    `json:"game_name"`
	Question  string    `json:"question"`
}

type SuccessfulQAPair struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	GameName  string    `json:"game_name"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
}

func LogIncomingRequest(gameName, question string) {
	logEntry := IncomingRequest{
		Type:      "INCOMING_REQUEST",
		Timestamp: time.Now().UTC(),
		GameName:  gameName,
		Question:  question,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Failed to marshal incoming request log: %v", err)
		log.Printf("INCOMING REQUEST - Game: %s, Question: %s", gameName, question)
		return
	}

	log.Printf("STRUCTURED_LOG: %s", string(jsonData))
}

func LogSuccessfulQAPair(gameName, question, answer string) {
	logEntry := SuccessfulQAPair{
		Type:      "SUCCESSFUL_QA_PAIR",
		Timestamp: time.Now().UTC(),
		GameName:  gameName,
		Question:  question,
		Answer:    answer,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Failed to marshal QA pair log: %v", err)
		log.Printf("SUCCESSFUL_QA_PAIR - Game: %s | Question: %s | Answer: %s", gameName, question, answer)
		return
	}

	log.Printf("STRUCTURED_LOG: %s", string(jsonData))
}
