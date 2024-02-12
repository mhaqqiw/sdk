package qmodule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mhaqqiw/sdk/go/qentity"
)

func CountElapsed(start time.Time) time.Duration {
	return time.Since(start)
}

func GetDataFromUUID(uuidStr string) (qentity.UUIDData, error) {
	nodeID := "00:00:00:00:00:00"
	res := qentity.UUIDData{}
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return res, err
	}

	t := u.Time()
	sec, nsec := t.UnixTime()
	timeStamp := time.Unix(sec, nsec)

	byteData, err := u.MarshalText()
	if err != nil {
		return res, err
	}

	val, err := u.Value()
	if err != nil {
		return res, err
	}

	n := u.NodeID()

	if len(n) == 6 {
		nodeID = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", n[0], n[1], n[2], n[3], n[4], n[5])
	}

	res.UUID = u
	res.Timestamp = timeStamp
	res.ClockSequence = u.ClockSequence()
	res.Domain = u.Domain().String()
	res.NodeID = nodeID
	res.ID = u.ID()
	res.MarshalText = string(byteData)
	res.String = u.String()
	res.Version = u.Version().String()
	res.URN = u.URN()
	res.Value = val.(string)
	res.Variant = u.Variant().String()

	if res.Version != "VERSION_1" {
		res.NodeID = "00:00:00:00:00:00"
		nodeID = "00:00:00:00:00:00"
	}
	return res, nil
}

func CheckRecaptcha(secret, token, validateURL string) (qentity.SiteVerifyResponse, error) {
	var body qentity.SiteVerifyResponse
	req, err := http.NewRequest(http.MethodPost, validateURL, nil)
	if err != nil {
		return body, err
	}
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", token)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

func GenerateUUIDV1() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
