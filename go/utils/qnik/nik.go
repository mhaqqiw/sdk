package qnik

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

	"github.com/mhaqqiw/sdk/go/utils/qlog"
)

var (
	stateNIKMap       map[string]string
	cityNIKMap        map[string]string
	districtNIKMap    map[string]string
	initializedNIKMap bool
)

type IDCardData struct {
	ID          string    `json:"id" db:"id"`
	PartnerID   string    `json:"partner_id" db:"partner_id"`
	NIK         string    `json:"nik" db:"nik"`
	Name        string    `json:"name" db:"name"`
	State       string    `json:"state" db:"state"`
	City        string    `json:"city" db:"city"`
	District    string    `json:"district" db:"district"`
	Subdistrict string    `json:"subdistrict" db:"subdistrict"`
	Address     string    `json:"address" db:"address"`
	Gender      string    `json:"gender" db:"gender"`
	DOB         time.Time `json:"dob" db:"dob"`
	ImageID     string    `json:"image_id" db:"image_id"`
	ImageType   int       `json:"image_type" db:"image_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy   string    `json:"updated_by" db:"updated_by"`
	DeletedAt   *string   `db:"deleted_at" json:"deleted_at"`
	DeletedBy   *string   `db:"deleted_by" json:"deleted_by"`
}

type NIKMapItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type option struct {
	sdkPath string
}

type NIKOption func(*option)

func WithPath(path string) NIKOption {
	return func(o *option) {
		o.sdkPath = path
	}
}

func Init(opts ...NIKOption) error {
	opt := &option{
		sdkPath: filepath.Join("files/etc/sdk"),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	stateFile, err := os.ReadFile(opt.sdkPath + "state.json")
	if err != nil {
		return err
	}
	states := make([]NIKMapItem, 0)
	err = json.Unmarshal(stateFile, &states)
	if err != nil {
		return err
	}

	for _, state := range states {
		stateNIKMap[state.Code] = state.Name
	}
	cityFile, err := os.ReadFile(opt.sdkPath + "city.json")
	if err != nil {
		return err
	}

	cities := make([]NIKMapItem, 0)
	err = json.Unmarshal(cityFile, &cities)
	if err != nil {
		return err
	}
	for _, city := range cities {
		cityNIKMap[city.Code] = city.Name
	}

	districtFile, err := os.ReadFile(opt.sdkPath + "district.json")
	if err != nil {
		return err
	}
	districts := make([]NIKMapItem, 0)
	err = json.Unmarshal(districtFile, &districts)
	if err != nil {
		return err
	}
	for _, district := range districts {
		districtNIKMap[district.Code] = district.Name
	}

	initializedNIKMap = true
	return nil
}

func (i *IDCardData) ParseNIK(nik string) error {
	if !initializedNIKMap {
		err := Init()
		if err != nil {
			qlog.Debug(err.Error())
			return errors.New("Failed to initialize NIK Map")
		}
	}

	if len(nik) == 0 {
		return errors.New("Empty NIK")
	}

	if len(nik) < 16 {
		return errors.New("Invalid NIK (Code: 1)")
	}

	for _, code := range nik {
		if !unicode.IsDigit(code) {
			return errors.New("Invalid NIK (Code: 2)")
		}
	}

	state, err := parseNIKState(nik[0:2])
	if err != nil {
		qlog.Debug(err.Error())
		return errors.New("Invalid NIK (Code: 3)")
	}
	i.State = state

	city, err := parseNIKCity(nik[0:4])
	if err != nil {
		qlog.Debug(err.Error())
		return errors.New("Invalid NIK (Code: 4)")
	}
	i.City = city

	district, err := parseNIKDistrict(nik[0:6])
	if err != nil {
		qlog.Debug(err.Error())
		return errors.New("Invalid NIK (Code: 5)")
	}
	i.District = district

	dob, err := parseNIKDOB(nik[6:12])
	if err != nil {
		qlog.Debug(err.Error())
		return errors.New("Invalid NIK (Code: 6)")
	}
	i.DOB = dob

	gender, err := parseNIKGender(nik[6:8])
	if err != nil {
		qlog.Debug(err.Error())
		return errors.New("Invalid NIK (Code: 7)")
	}
	i.Gender = gender

	i.NIK = nik

	return nil
}

func parseNIKState(data string) (string, error) {
	state, ok := stateNIKMap[data]
	if !ok {
		return "", errors.New("state not found")
	}
	return state, nil
}

func parseNIKCity(data string) (string, error) {
	city, ok := cityNIKMap[data]
	if !ok {
		return "", errors.New("city not found")
	}
	return city, nil
}

func parseNIKDistrict(data string) (string, error) {
	district, ok := districtNIKMap[data]
	if !ok {
		return "", errors.New("district not found")
	}
	return district, nil
}

func parseNIKDOB(data string) (time.Time, error) {
	var dob time.Time

	day, err := strconv.Atoi(data[0:2])
	if err != nil {
		return dob, err
	}

	month, err := strconv.Atoi(data[2:4])
	if err != nil {
		return dob, err
	}

	year, err := strconv.Atoi(data[4:6])
	if err != nil {
		return dob, err
	}
	currentYear := time.Now().Year() % 100
	century := 1900
	if year <= currentYear {
		century = 2000
	}

	dob = time.Date(year+century, time.Month(month), day%40, 0, 0, 0, 0, time.Local)

	return dob, nil
}

func parseNIKGender(data string) (string, error) {
	day, err := strconv.Atoi(data[0:2])
	if err != nil {
		return data, err
	}

	if day > 40 {
		return "F", nil
	}
	return "M", nil
}
