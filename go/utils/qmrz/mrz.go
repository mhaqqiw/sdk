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
	arr := strings.Split(strings.TrimSpace(mrz), "\n")
	if len(arr) < 1 {
		return ret, fmt.Errorf("Invalid MRZ (Code: 1)")

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
		if charLen == TD1_CHAR_LEN {
			return TD1MRZ(arr, ret)
		} else {
			return TD2MRZ(arr, ret)
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
	data[0] = strings.TrimSpace(data[0])
	if len(data[0]) < TD3_CHAR_LEN {
		return ret, fmt.Errorf("Invalid MRZ in line 1 (Code: 3)")
	}
	ret.DocumentType = clear(data[0][:2])
	ret.DocumentClass = TD3
	ret.Passport.Country = clear(data[0][2:5])
	ret.Passport.Name = clear(data[0][5:])
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
	if len(data[0]) < TD1_CHAR_LEN {
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
	ret.TD2.Name = clear(data[2])
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
