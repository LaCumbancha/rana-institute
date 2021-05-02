package utils

import (
	"time"
	"unicode"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const HOURS = 'h'
const MINUTES = 'm'
const SECONDS = 's'

type TTL struct {
	Hours			time.Duration
	Minutes			time.Duration
	Seconds			time.Duration
}

func TimeToLive(duration string) TTL {
	durationMap := make(map[rune]int)
	durationAccum := ""
	for _, char := range duration {
		if unicode.IsLetter(char) {
			duration, err := strconv.Atoi(durationAccum)
			if err != nil {
				log.Fatalf("Error converting duration %s to string (type '%c').", durationAccum, char)
			} else {
				durationMap[char] = duration
			}
			durationAccum = ""
		} else {
			durationAccum += string(char)
		}
	}

	return TTL { 
		Hours: 		time.Duration(durationMap[HOURS]), 
		Minutes: 	time.Duration(durationMap[MINUTES]), 
		Seconds:	time.Duration(durationMap[SECONDS]),
	}
}
