package utils

import (
	"fmt"
	"strings"
)

const (
	welcomeMessage      = "CON Welcome to Be.Well"
	todoMessage         = "What would you like to do?"
	getCover            = "Get Cover"
	getConsultation     = "Get Consultation"
	getTest             = "Get Test"
	getMedicine         = "Get Medicine"
	getCoaching         = "Get Coaching"
	ussdResponse        = "END Thank you for applying for %s, our team will get back to you shortly."
	defaultResponseText = "CON Invalid choice, please try again"
	textLength          = 1
)

//Menu creates the ussd menu
func Menu() (resp string) {
	resp += "1." + getCover + ".\r\n"
	resp += "2." + getConsultation + ".\r\n"
	resp += "3." + getTest + ".\r\n"
	resp += "4." + getMedicine + ".\r\n"
	resp += "5." + getCoaching + "."
	return resp
}

//ResponseMenu creates USSD response
func ResponseMenu(text string) (resp string) {
	switch text {
	case "1":
		resp = fmt.Sprintf(ussdResponse, getCover)
	case "2":
		resp = fmt.Sprintf(ussdResponse, getConsultation)
	case "3":
		resp = fmt.Sprintf(ussdResponse, getTest)
	case "4":
		resp = fmt.Sprintf(ussdResponse, getMedicine)
	case "5":
		resp = fmt.Sprintf(ussdResponse, getCoaching)
	default:
		resp = DefaultMenu(text)
	}
	return resp

}

//DefaultMenu returns the defaultResponse function and allows a
//user to reenter a text
func DefaultMenu(text string) (resp string) {
	if len(text) == textLength {
		resp = DefaultResponse(text)
	} else {
		userChoice := GetUserChoice(text)
		resp = ResponseMenu(userChoice)
	}
	return resp
}

//DefaultResponse returns a default response and the menu to the user
//incase the user enters an invalid choice
func DefaultResponse(text string) string {
	resp := defaultResponseText + "\r\n"
	resp += Menu()
	return resp
}

//GetUserChoice splits the text sent to callback since all user text is joined by *
//and gets the last entry from the text
func GetUserChoice(text string) string {
	vals := strings.Split(text, "*")
	return vals[len(vals)-textLength]
}

//GetTextValue gets the corresponding value of the text
//to be saved.
func GetTextValue(text string) (resp string) {
	switch text {
	case "1":
		resp = getCover
	case "2":
		resp = getConsultation
	case "3":
		resp = getTest
	case "4":
		resp = getMedicine
	case "5":
		resp = getCoaching
	}
	return resp
}
