package main

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------
type MessageLog struct {
	messages []string
	idx      int
}

func (log *MessageLog) Add(format string, vals ...any) {
	msg := fmt.Sprintf(format, vals...)
	log.messages = append(log.messages, msg)
}

func (log *MessageLog) Clear() {
	log.messages = nil
}

func (log *MessageLog) HasUnread() bool {
	return log.idx < len(log.messages)
}

func (log *MessageLog) Last(n int) []string {
	if n >= len(log.messages) {
		return log.messages
	} else {
		return log.messages[len(log.messages)-n:]
	}
}

func (log *MessageLog) LatestAsStr() string {
	s := ""
	if len(log.messages[log.idx:]) > 0 {
		s = strings.Join(log.messages[log.idx:], " ")
	}
	return s
}

func (log *MessageLog) ClearUnread() {
	log.idx = len(log.messages)
}
