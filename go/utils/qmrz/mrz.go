package qmrz

import (
	"fmt"
	"strings"
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

type MRZ struct {
	DocumentType  string
	DocumentClass string
	Passport      struct {
		Country            string
		Name               string
		DocNumber          string
		HashDocNumber      string
		Nationality        string
		DOB                string
		HashDOB            string
		Sex                string
		ExpiredDate        string
		HashExpiredDate    string
		PersonalNumber     string
		HashPersonalNumber string
		FinalHash          string
		ExpectedHash       struct {
			IsValid            bool
			HashDocNumber      string
			HashDOB            string
			HashExpiredDate    string
			HashPersonalNumber string
			FinalHash          string
		}
	}
	TD1 struct {
		Country         string
		DocNumber       string
		HashDocNumber   string
		AdditionalInfo1 string
		DOB             string
		HashDOB         string
		Sex             string
		ExpiredDate     string
		HashExpiredDate string
		Nationality     string
		AdditionalInfo2 string
		FinalHash       string
		Name            string
		ExpectedHash    struct {
			IsValid         bool
			HashDocNumber   string
			HashDOB         string
			HashExpiredDate string
			FinalHash       string
		}
	}
	TD2 struct {
		Country         string
		Name            string
		DocNumber       string
		HashDocNumber   string
		Nationality     string
		DOB             string
		HashDOB         string
		Sex             string
		ExpiredDate     string
		HashExpiredDate string
		AdditionalInfo  string
		FinalHash       string
		ExpectedHash    struct {
			IsValid         bool
			HashDocNumber   string
			HashDOB         string
			HashExpiredDate string
			FinalHash       string
		}
	}
	VISAA struct {
		Country         string
		Name            string
		DocNumber       string
		HashDocNumber   string
		Nationality     string
		DOB             string
		HashDOB         string
		Sex             string
		ExpiredDate     string
		HashExpiredDate string
		AdditionalInfo  string
		ExpectedHash    struct {
			IsValid         bool
			HashDocNumber   string
			HashDOB         string
			HashExpiredDate string
		}
	}
	VISAB struct {
		Country         string
		Name            string
		DocNumber       string
		HashDocNumber   string
		Nationality     string
		DOB             string
		HashDOB         string
		Sex             string
		ExpiredDate     string
		HashExpiredDate string
		AdditionalInfo  string
		ExpectedHash    struct {
			IsValid         bool
			HashDocNumber   string
			HashDOB         string
			HashExpiredDate string
		}
	}
}

func ParseMRZ(mrz string) (ret MRZ, err error) {
	arr := strings.Split(strings.TrimSpace(mrz), "\n")
	if len(arr) < 2 {
		return ret, fmt.Errorf("Invalid MRZ (Code: 1)")

	}

	docType := rune(arr[0][0])
	charLen := len(strings.TrimSpace(arr[0]))

	switch docType {
	case 'A', 'B', 'C', 'I':
		if charLen == TD1_CHAR_LEN {
			return TD1MRZ(arr, ret)
		} else {
			return TD2MRZ(arr, ret)
		}
	case 'P':
		return PassportMRZ(arr, ret)
	case 'V':
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
	if len(data[0]) < TD1_CHAR_LEN {
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
	if len(data[2]) < TD2_CHAR_LEN {
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
