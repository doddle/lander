package main

import (
	"strings"

	"github.com/withmandala/go-log"
)

func containsString(input, pattern string) bool {
	if strings.Contains(input, pattern) {
		return true
	}
	return false
}

// from an input string.. range through colorschemes and return the first one that matches
func guessColour(logger *log.Logger, schemes []Colorschemes, input string) (result string) {
	logger.Debugf("guessing colours for: %s", input)
	var allVerbs []string

	for _, i := range schemes {
		for _, j := range i.Verbs {
			allVerbs = append(allVerbs, j)
		}
	}

	// see available colours: https://vuetifyjs.com/en/styles/colors/#material-colors

	// no verbs? return a default color
	if len(allVerbs) == 0 {
		return "amber lighten-2"
	}

	// return a default colour if there are no matching verbs
	matches := 0
	for _, i := range allVerbs {
		if containsString(input, i) {
			matches++
		}
	}
	if matches < 1 {
		logger.Warn("no matches")
		return "blue lighten-5"
	}

	// at this point we know there is at least one verb that matches.. lets return the first associated
	for _, i := range schemes {
		for _, j := range i.Verbs {
			if containsString(input, j) {
				return i.Color
			}
		}
	}

	// return a default it all attempts at guessing failed
	return "blue-grey lighten-3"
}

// from an input string.. range through colorschemes and return the first one that matches
func guessHex(logger *log.Logger, schemes []Colorschemes, input string) (result string) {
	logger.Debugf("guessing hex for: %s", input)
	result = "99db9a"
	var allVerbs []string

	for _, i := range schemes {
		for _, j := range i.Verbs {
			allVerbs = append(allVerbs, j)
		}
	}
	// pick some here: https://htmlcolorcodes.com/color-picker/

	// no verbs? return a default color
	if len(allVerbs) == 0 {
		return result
	}

	// return a default colour if there are no matching verbs
	matches := 0
	for _, i := range allVerbs {
		if containsString(input, i) {
			matches++
		}
	}
	if matches < 1 {
		logger.Warn("couldn't match a hex colour")
		return result
	}

	// at this point we know there is at least one verb that matches.. lets return the first associated
	for _, i := range schemes {
		for _, j := range i.Verbs {
			if containsString(input, j) {
				return i.Hex
			}
		}
	}

	// return a default if all attempts at guessing failed
	return result
}
