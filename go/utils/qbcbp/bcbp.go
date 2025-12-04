package qbcbp

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type BCBP struct {
	FormatCode   string `json:"format_code"`
	TotalLeg     string `json:"total_leg"`
	Name         string `json:"name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Indicator    string `json:"indicator"`
	PnrCode      string `json:"pnr_code"`
	From         string `json:"from"`
	To           string `json:"to"`
	Airline      string `json:"airline"`
	FlightNumber string `json:"flight_number"`
	Date         string `json:"date"`
	Class        string `json:"class"`
	Seat         string `json:"seat"`
	Sequence     string `json:"sequence"`
	Status       string `json:"status"`
}

func generateData(length int, charset string) string {
	var letters string
	switch charset {
	case "[A-Z]":
		letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "[0-9]":
		letters = "0123456789"
	case "[FCY]":
		letters = "FCY"
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rng.Intn(len(letters))]
	}
	return string(result)
}

func dateToJulian(date string) string {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "000"
	}
	return fmt.Sprintf("%03d", parsedTime.YearDay())
}

func julianDayToGregorian(julianDay string) (string, error) {
	day, err := strconv.Atoi(julianDay)
	if err != nil {
		return "", errors.New("invalid Julian day")
	}
	return time.Date(2025, 1, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02"), nil
}

func GenerateBCBP(lastName, firstName, dateOfFlight, fromAirport, toAirport string) (string, error) {
	if lastName == "" || firstName == "" || dateOfFlight == "" {
		return "", errors.New("missing required parameters")
	}

	pnrCode := generateData(6, "[A-Z]")
	airlineCode := generateData(2, "[A-Z]")
	flightNumber := generateData(4, "[0-9]")
	julianDate := dateToJulian(dateOfFlight)
	seatNumber := generateData(3, "[0-9]") + generateData(1, "[A-Z]")
	checkInSequence := generateData(4, "[0-9]")
	classPassenger := generateData(1, "[FCY]")

	nameField := fmt.Sprintf("%-20s", strings.ToUpper(lastName)+"/"+strings.ToUpper(firstName))
	flightInfo := fmt.Sprintf("%-8s", airlineCode+flightNumber)
	seat := fmt.Sprintf("%-4s", seatNumber)

	bcRawString := fmt.Sprintf("M1%sE%-7s%s%s%-8s%s%s%-4s%-4s1AA", nameField, pnrCode, fromAirport, toAirport, flightInfo, julianDate, classPassenger, seat, checkInSequence)
	return bcRawString, nil
}

func ParseBCBP(data string) (BCBP, error) {
	var result BCBP
	if len(data) < 58 {
		return result, errors.New("invalid BCBP (Code: 1)")
	}

	nameField := strings.TrimSpace(data[2:22])
	nameParts := strings.Split(nameField, "/")
	if len(nameParts) < 1 {
		return result, errors.New("invalid BCBP (Code: 2)")
	}

	firstName := ""
	if len(nameParts) > 1 {
		firstName = nameParts[1]
	}

	date, err := julianDayToGregorian(strings.TrimSpace(data[44:47]))
	if err != nil {
		return result, errors.New("invalid BCBP (Code: 3)")
	}

	airline := strings.TrimSpace(data[36:39])
	if len(airline) < 2 {
		return result, errors.New("invalid BCBP (Code: 4)")
	}
	for _, char := range airline {
		if !(unicode.IsLetter(char) != unicode.IsNumber(char) != unicode.IsLower(char)) {
			return result, errors.New("invalid BCBP (Code: 4)")
		}
	}

	flightNumber, err := sanitizeFlightNumber(data[39:43])
	if err != nil {
		return result, errors.New("invalid BCBP (Code: 5)")
	}

	result = BCBP{
		FormatCode:   strings.TrimSpace(data[0:1]),
		TotalLeg:     strings.TrimSpace(data[1:2]),
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(nameParts[0]),
		Name:         strings.TrimSpace(data[2:22]),
		Indicator:    strings.TrimSpace(data[22:23]),
		PnrCode:      strings.TrimSpace(data[23:30]),
		From:         strings.TrimSpace(data[30:33]),
		To:           strings.TrimSpace(data[33:36]),
		Airline:      airline,
		FlightNumber: flightNumber,
		Date:         date,
		Class:        strings.TrimSpace(data[47:48]),
		Seat:         strings.TrimSpace(data[48:52]),
		Sequence:     strings.TrimSpace(data[52:56]),
		Status:       strings.TrimSpace(data[57:58]),
	}
	return result, nil
}

func sanitizeFlightNumber(flightNumber string) (string, error) {
	flightNumber = strings.TrimSpace(flightNumber)
	if flightNumber != "" && !unicode.IsNumber(rune(flightNumber[0])) {
		return flightNumber, errors.New("invalid Flight Number")
	}
	for i, char := range flightNumber {
		if unicode.IsNumber(char) {
			num, err := strconv.Atoi(string(char))
			if err != nil {
				return flightNumber, err
			}
			if num > 0 {
				return flightNumber[i:], nil
			}
		} else if !unicode.IsLetter(char) {
			return flightNumber, errors.New("invalid Flight Number")
		}
	}
	return flightNumber, errors.New("invalid Flight Number")
}
