package service

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/romik1505/userDetailsService/internal/common"
)

type PersonMessage struct {
	Person      common.CreatePersonRequest `json:"person"`
	ErrorReason string                     `json:"error_reason,omitempty"`
}

func UnmarshalPersonMessage(msg *sarama.ConsumerMessage) (PersonMessage, error) {
	personMessage := PersonMessage{}
	err := json.Unmarshal(msg.Value, &personMessage)
	if err != nil {
		return PersonMessage{}, fmt.Errorf("UnmarshalPersonMessage: %v", err)
	}
	return personMessage, nil
}
