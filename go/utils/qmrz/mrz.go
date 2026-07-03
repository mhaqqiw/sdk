package qmrz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const TD1_CHAR_LEN = 30
const TD2_CHAR_LEN = 36
const TD3_CHAR_LEN = 44
const VISA_A_CHAR_LEN = 44
const VISA_B_CHAR_LEN = 36
const TD1 = "TD1"
const TD2 = "TD2"
const TD3 = "TD3"
const VISA_A = "MRV-A"
const VISA_B = "MRV-B"

type Passport struct {
	Country            string `json:"country"`
	Name               string `json:"name"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	DocNumber          string `json:"doc_number"`
	HashDocNumber      string `json:"hash_doc_number"`
	Nationality        string `json:"nationality"`
	DOB                string `json:"dob"`
	HashDOB            string `json:"hash_dob"`
	Sex                string `json:"sex"`
	ExpiredDate        string `json:"expired_date"`
	HashExpiredDate    string `json:"hash_expired_date"`
	PersonalNumber     string `json:"personal_number"`
	HashPersonalNumber string `json:"hash_personal_number"`
	FinalHash          string `json:"final_hash"`
	ExpectedHash       struct {
		IsValid            bool   `json:"is_valid"`
		HashDocNumber      string `json:"hash_doc_number"`
		HashDOB            string `json:"hash_dob"`
		HashExpiredDate    string `json:"hash_expired_date"`
		HashPersonalNumber string `json:"hash_personal_number"`
		FinalHash          string `json:"final_hash"`
	} `json:"expected_hash"`
}

type MRZ struct {
	DocumentType  string   `json:"document_type"`
	DocumentClass string   `json:"document_class"`
	Passport      Passport `json:"passport"`
	TD1           struct {
		Country         string `json:"country"`
		DocNumber       string `json:"doc_number"`
		HashDocNumber   string `json:"hash_doc_number"`
		AdditionalInfo1 string `json:"additional_info_1"`
		DOB             string `json:"dob"`
		HashDOB         string `json:"hash_dob"`
		Sex             string `json:"sex"`
		ExpiredDate     string `json:"expired_date"`
		HashExpiredDate string `json:"hash_expired_date"`
		Nationality     string `json:"nationality"`
		AdditionalInfo2 string `json:"additional_info_2"`
		FinalHash       string `json:"final_hash"`
		Name            string `json:"name"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ExpectedHash    struct {
			IsValid         bool   `json:"is_valid"`
			HashDocNumber   string `json:"hash_doc_number"`
			HashDOB         string `json:"hash_dob"`
			HashExpiredDate string `json:"hash_expired_date"`
			FinalHash       string `json:"final_hash"`
		} `json:"expected_hash"`
	} `json:"td1"`
	TD2 struct {
		Country         string `json:"country"`
		Name            string `json:"name"`
		DocNumber       string `json:"doc_number"`
		HashDocNumber   string `json:"hash_doc_number"`
		Nationality     string `json:"nationality"`
		DOB             string `json:"dob"`
		HashDOB         string `json:"hash_dob"`
		Sex             string `json:"sex"`
		ExpiredDate     string `json:"expired_date"`
		HashExpiredDate string `json:"hash_expired_date"`
		AdditionalInfo  string `json:"additional_info"`
		FinalHash       string `json:"final_hash"`
		ExpectedHash    struct {
			IsValid         bool   `json:"is_valid"`
			HashDocNumber   string `json:"hash_doc_number"`
			HashDOB         string `json:"hash_dob"`
			HashExpiredDate string `json:"hash_expired_date"`
			FinalHash       string `json:"final_hash"`
		} `json:"expected_hash"`
	} `json:"td2"`
	VISAA struct {
		Country         string `json:"country"`
		Name            string `json:"name"`
		DocNumber       string `json:"doc_number"`
		HashDocNumber   string `json:"hash_doc_number"`
		Nationality     string `json:"nationality"`
		DOB             string `json:"dob"`
		HashDOB         string `json:"hash_dob"`
		Sex             string `json:"sex"`
		ExpiredDate     string `json:"expired_date"`
		HashExpiredDate string `json:"hash_expired_date"`
		AdditionalInfo  string `json:"additional_info"`
		ExpectedHash    struct {
			IsValid         bool   `json:"is_valid"`
			HashDocNumber   string `json:"hash_doc_number"`
			HashDOB         string `json:"hash_dob"`
			HashExpiredDate string `json:"hash_expired_date"`
		} `json:"expected_hash"`
	} `json:"visa_a"`
	VISAB struct {
		Country         string `json:"country"`
		Name            string `json:"name"`
		DocNumber       string `json:"doc_number"`
		HashDocNumber   string `json:"hash_doc_number"`
		Nationality     string `json:"nationality"`
		DOB             string `json:"dob"`
		HashDOB         string `json:"hash_dob"`
		Sex             string `json:"sex"`
		ExpiredDate     string `json:"expired_date"`
		HashExpiredDate string `json:"hash_expired_date"`
		AdditionalInfo  string `json:"additional_info"`
		ExpectedHash    struct {
			IsValid         bool   `json:"is_valid"`
			HashDocNumber   string `json:"hash_doc_number"`
			HashDOB         string `json:"hash_dob"`
			HashExpiredDate string `json:"hash_expired_date"`
		} `json:"expected_hash"`
	} `json:"visa_b"`
}

func pad(text string, length int, char rune) string {
	for len(text) < length {
		text += string(char)
	}
	if len(text) > length {
		return text[:length]
	}
	return text
}

func formatName(name string) string {
	parts := strings.Fields(strings.ToUpper(name))
	if len(parts) == 0 {
		return ""
	}
	surname := parts[0]
	givenNames := strings.Join(parts[1:], "<")
	formatted := surname + "<<" + givenNames
	return strings.ReplaceAll(formatted, " ", "<")
}

func formatDate(dateStr string) string {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "000000"
	}
	return parsed.Format("060102")
}

func charValue(char rune) int {
	switch {
	case unicode.IsDigit(char):
		return int(char - '0')
	case char >= 'A' && char <= 'Z':
		return int(char-'A') + 10
	case char == '<':
		return 0
	default:
		return 0
	}
}

func computeCheckDigit(input string) int {
	weights := []int{7, 3, 1}
	sum := 0
	for i, r := range input {
		sum += charValue(r) * weights[i%3]
	}
	return sum % 10
}

func GenerateMRZ(mrzType string, mrz MRZ) (string, error) {
	if mrzType != TD3 {
		return "", errors.New("Unsupported MRZ Type")
	}
	data := mrz.Passport
	line1 := "P<" + pad(strings.ToUpper(data.Country), 3, '<') + pad(formatName(data.Name), 39, '<')

	passportNumber := pad(strings.ToUpper(data.DocNumber), 9, '<')
	passportNumberCheck := computeCheckDigit(passportNumber)

	dob := formatDate(data.DOB)
	dobCheck := computeCheckDigit(dob)

	expiry := formatDate(data.ExpiredDate)
	expiryCheck := computeCheckDigit(expiry)

	personalNumber := pad(data.PersonalNumber, 14, '<')
	personalNumberCheck := computeCheckDigit(personalNumber)

	sex := "<" // default placeholder if empty
	if len(data.Sex) > 0 {
		sex = strings.ToUpper(string(data.Sex[0]))
	}

	line2Body := passportNumber +
		fmt.Sprintf("%d", passportNumberCheck) +
		pad(strings.ToUpper(data.Nationality), 3, '<') +
		dob +
		fmt.Sprintf("%d", dobCheck) +
		sex + // convert 'female' -> 'F'
		expiry +
		fmt.Sprintf("%d", expiryCheck) +
		personalNumber +
		fmt.Sprintf("%d", personalNumberCheck)

	compositeCheck := computeCheckDigit(line2Body)
	line2 := line2Body + fmt.Sprintf("%d", compositeCheck)

	return line1 + "\n" + line2, nil
}

func ParseMRZ(mrz string) (ret MRZ, err error) {
	if len(mrz) == 0 {
		return ret, errors.New("Empty MRZ")
	}
	arr := strings.Split(strings.TrimSpace(mrz), "\n")
	if len(arr) < 1 {
		return ret, fmt.Errorf("Invalid MRZ (Code: 1)")

	}

	if len(arr[0]) == 0 {
		return ret, errors.New("Empty MRZ")
	}

	docType := rune(arr[0][0])
	charLen := len(strings.TrimSpace(arr[0]))

	switch docType {
	case 'A', 'B', 'C', 'I':
		if len(arr[0]) > TD2_CHAR_LEN {
			if charLen == (3 * TD1_CHAR_LEN) {
				arr = splitByN(arr[0], TD1_CHAR_LEN)
			} else {
				arr = splitByN(arr[0], TD2_CHAR_LEN)
			}
		}
		switch len(arr) {
		case 3:
			return TD1MRZ(arr, ret)
		case 2:
			return TD2MRZ(arr, ret)
		default:
			return ret, fmt.Errorf("Invalid MRZ line count")
		}
	case 'P':
		if len(arr[0]) > TD3_CHAR_LEN {
			arr = splitByN(arr[0], TD3_CHAR_LEN)
		}
		return PassportMRZ(arr, ret)
	case 'V':
		if len(arr[0]) > VISA_A_CHAR_LEN {
			if charLen == (2 * VISA_A_CHAR_LEN) {
				arr = splitByN(arr[0], VISA_A_CHAR_LEN)
			} else {
				arr = splitByN(arr[0], VISA_B_CHAR_LEN)
			}
		}
		if charLen == VISA_A_CHAR_LEN {
			return VISAAMRZ(arr, ret)
		} else {
			return VISABMRZ(arr, ret)
		}
	default:
		return ret, fmt.Errorf("MRZ not supported (Code: 2)")
	}
}

func PassportMRZ(data []string, ret MRZ) (MRZ, error) {
	if len(data) < 2 {
		return ret, fmt.Errorf("Invalid MRZ (Code: 3)")
	}
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < TD3_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = TD3
	ret.Passport.Country = clear(data[0][2:5])
	parts := strings.SplitN(data[0][5:], "<<", 2)
	if len(parts) > 0 {
		ret.Passport.LastName = clear(parts[0])
	}
	if len(parts) > 1 {
		ret.Passport.FirstName = clear(parts[1])
	}
	ret.Passport.Name = strings.TrimSpace(ret.Passport.FirstName + " " + ret.Passport.LastName)
	data[1] = strings.TrimSpace(data[1])
	if len(data[1]) < TD3_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 2 (Code: 3)")
	}
	ret.Passport.DocNumber = clear(data[1][:9])
	ret.Passport.HashDocNumber = clear(data[1][9:10])
	ret.Passport.Nationality = clear(data[1][10:13])
	ret.Passport.DOB = clear(data[1][13:19])
	ret.Passport.HashDOB = clear(data[1][19:20])
	ret.Passport.Sex = clear(data[1][20:21])
	ret.Passport.ExpiredDate = clear(data[1][21:27])
	ret.Passport.HashExpiredDate = clear(data[1][27:28])
	ret.Passport.PersonalNumber = clear(data[1][28:42])
	ret.Passport.HashPersonalNumber = clear(data[1][42:43])
	ret.Passport.FinalHash = clear(data[1][43:])
	return ret, nil
}

func TD1MRZ(data []string, ret MRZ) (MRZ, error) {
	if len(data) < 3 {
		return ret, fmt.Errorf("Invalid MRZ (Code: 3)")
	}
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < TD1_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = TD1
	ret.TD1.Country = clear(data[0][2:5])
	ret.TD1.DocNumber = clear(data[0][5:14])
	ret.TD1.HashDocNumber = clear(data[0][14:15])
	ret.TD1.AdditionalInfo1 = clear(data[0][15:])
	data[1] = strings.TrimSpace(data[1])
	if len(data[1]) < TD1_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 2 (Code: 3)")
	}
	ret.TD1.DOB = clear(data[1][:6])
	ret.TD1.HashDOB = clear(data[1][6:7])
	ret.TD1.Sex = clear(data[1][7:8])
	ret.TD1.ExpiredDate = clear(data[1][8:14])
	ret.TD1.HashExpiredDate = clear(data[1][14:15])
	ret.TD1.Nationality = clear(data[1][15:18])
	ret.TD1.AdditionalInfo2 = clear(data[1][18:29])
	ret.TD1.FinalHash = clear(data[1][29:])
	data[2] = strings.TrimSpace(data[2])
	if len(data[2]) < TD1_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 3 (Code: 3)")
	}
	parts := strings.SplitN(data[2], "<<", 2)
	if len(parts) > 0 {
		ret.TD1.LastName = clear(parts[0])
	}
	if len(parts) > 1 {
		ret.TD1.FirstName = clear(parts[1])
	}
	ret.TD1.Name = strings.TrimSpace(ret.TD1.FirstName + " " + ret.TD1.LastName)
	return ret, nil
}
func TD2MRZ(data []string, ret MRZ) (MRZ, error) {
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < TD2_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = TD2
	ret.TD2.Country = clear(data[0][2:5])
	ret.TD2.Name = clear(data[0][5:])
	data[1] = strings.TrimSpace(data[1])
	if len(data[1]) < TD2_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 2 (Code: 3)")
	}
	ret.TD2.DocNumber = clear(data[1][:9])
	ret.TD2.HashDocNumber = clear(data[1][9:10])
	ret.TD2.Nationality = clear(data[1][10:13])
	ret.TD2.DOB = clear(data[1][13:19])
	ret.TD2.HashDOB = clear(data[1][19:20])
	ret.TD2.Sex = clear(data[1][20:21])
	ret.TD2.ExpiredDate = clear(data[1][21:27])
	ret.TD2.HashExpiredDate = clear(data[1][27:28])
	ret.TD2.AdditionalInfo = clear(data[1][28:35])
	ret.TD2.FinalHash = clear(data[1][35:])
	return ret, nil
}
func VISAAMRZ(data []string, ret MRZ) (MRZ, error) {
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < VISA_A_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = VISA_A
	ret.VISAA.Country = clear(data[0][2:5])
	ret.VISAA.Name = clear(data[0][5:])
	data[1] = strings.TrimSpace(data[1])
	if len(data[1]) < VISA_A_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 2 (Code: 3)")
	}
	ret.VISAA.DocNumber = clear(data[1][:9])
	ret.VISAA.HashDocNumber = clear(data[1][9:10])
	ret.VISAA.Nationality = clear(data[1][10:13])
	ret.VISAA.DOB = clear(data[1][13:19])
	ret.VISAA.HashDOB = clear(data[1][19:20])
	ret.VISAA.Sex = clear(data[1][20:21])
	ret.VISAA.ExpiredDate = clear(data[1][21:27])
	ret.VISAA.HashExpiredDate = clear(data[1][27:28])
	ret.VISAA.AdditionalInfo = clear(data[1][28:])
	return ret, nil
}
func VISABMRZ(data []string, ret MRZ) (MRZ, error) {
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < VISA_B_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = VISA_B
	ret.VISAB.Country = clear(data[0][2:5])
	ret.VISAB.Name = clear(data[0][5:])
	data[1] = strings.TrimSpace(data[1])
	if len(data[1]) < VISA_B_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 2 (Code: 3)")
	}
	ret.VISAB.DocNumber = clear(data[1][:9])
	ret.VISAB.HashDocNumber = clear(data[1][9:10])
	ret.VISAB.Nationality = clear(data[1][10:13])
	ret.VISAB.DOB = clear(data[1][13:19])
	ret.VISAB.HashDOB = clear(data[1][19:20])
	ret.VISAB.Sex = clear(data[1][20:21])
	ret.VISAB.ExpiredDate = clear(data[1][21:27])
	ret.VISAB.HashExpiredDate = clear(data[1][27:28])
	ret.VISAB.AdditionalInfo = clear(data[1][28:])
	return ret, nil
}
func clear(str string) string {
	arr := strings.Split(str, "<")
	ret := []string{}
	for _, i := range arr {
		if len(i) > 0 {
			ret = append(ret, i)
		}
	}
	return strings.Join(ret, " ")
}

func splitByN(s string, n int) []string {
	var result []string
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

func ParseMRZExpiry(expiry string) (time.Time, error) {
	var ret time.Time
	if len(expiry) != 6 {
		return ret, errors.New("invalid expiry")
	}

	yearPart := expiry[:2]
	monthPart := expiry[2:4]
	dayPart := expiry[4:6]
	baseCentury := (time.Now().Year() / 100) * 100

	yy, err := strconv.Atoi(yearPart)
	if err != nil {
		return ret, err
	}
	fullYear := baseCentury + yy

	date := fmt.Sprintf("%d-%s-%s", fullYear, monthPart, dayPart)

	ret, err = time.Parse("2006-01-02", date)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func ParseMRZDOB(dob string) (time.Time, error) {
	var ret time.Time
	if len(dob) != 6 {
		return ret, errors.New("invalid dob")
	}

	yearPart := dob[:2]
	monthPart := dob[2:4]
	dayPart := dob[4:6]

	yy, err := strconv.Atoi(yearPart)
	if err != nil {
		return ret, err
	}

	mm, err := strconv.Atoi(monthPart)
	if err != nil {
		return ret, err
	}

	dd, err := strconv.Atoi(dayPart)
	if err != nil {
		return ret, err
	}

	// Guess century: MRZ DOB is usually 1900s or 2000s
	// If year > current year (2-digit), assume 1900s, else 2000s
	currentYear := time.Now().Year() % 100
	century := 1900
	if yy <= currentYear {
		century = 2000
	}
	fullYear := century + yy

	date := fmt.Sprintf("%04d-%02d-%02d", fullYear, mm, dd)
	ret, err = time.Parse("2006-01-02", date)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func GenerateMRZPassport(p Passport) (string, string, error) {
	// Parse name
	surname, given := parseName(p.Name)

	// Line 1
	line1 := fmt.Sprintf("P<%s%s<<%s",
		formatField(p.Country, 3),
		formatNameGen(surname),
		formatNameGen(given),
	)

	line1 = padRight(line1, 44, '<')

	// Date conversion
	dob := formatDate(p.DOB)

	exp := formatDate(p.ExpiredDate)

	docNumber := formatField(p.DocNumber, 9)
	nationality := formatField(p.Nationality, 3)
	sex := strings.ToUpper(string(p.Sex[0]))

	// Checksums
	docCheck := checkDigit(docNumber)
	dobCheck := checkDigit(dob)
	expCheck := checkDigit(exp)

	optional := strings.Repeat("<", 14)

	line2 := fmt.Sprintf("%s%d%s%s%d%s%s%d%s",
		docNumber,
		docCheck,
		nationality,
		dob,
		dobCheck,
		sex,
		exp,
		expCheck,
		optional,
	)

	// Final check digit
	finalCheck := checkDigit(line2)
	line2 = fmt.Sprintf("%s%d", line2, finalCheck)

	line2 = padRight(line2, 44, '<')

	return line1, line2, nil
}

func parseName(name string) (string, string) {
	parts := strings.Fields(strings.ToUpper(name))
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], "<")
}

func formatNameGen(s string) string {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "<")
	return filterMRZChars(s)
}

func formatField(s string, max int) string {
	s = strings.ToUpper(s)
	s = filterMRZChars(s)
	if len(s) > max {
		return s[:max]
	}
	return padRight(s, max, '<')
}

func padRight(s string, length int, pad rune) string {
	for len(s) < length {
		s += string(pad)
	}
	return s
}

func filterMRZChars(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('<')
		}
	}
	return result.String()
}

func checkDigit(input string) int {
	weights := []int{7, 3, 1}
	total := 0

	for i, c := range input {
		value := charValue(c)
		total += value * weights[i%3]
	}

	return total % 10
}
