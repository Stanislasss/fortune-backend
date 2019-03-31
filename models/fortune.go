package models

type FortuneMessage struct {
	Message  string `json:"message"`
	ID       string `json:"id" bson:"id"`
	CheckSum string `bson:"checksum,omitempty"`
}

type FortuneQuery map[string]interface{}

type Json map[string]interface{}
