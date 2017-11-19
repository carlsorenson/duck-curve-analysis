package duckCurve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	switch pathParts[3] {
	case "day":
		loc, _ := time.LoadLocation("America/New_York")
		parsedDate, _ := time.ParseInLocation("2006-01-02", pathParts[4], loc)
		getDayDatapoints(w, r, parsedDate)
	default:
		fmt.Fprint(w, "Neurio data aggregator API. Example: /api/v1/day/2017-11-01")
	}
}

func getDayDatapoints(w http.ResponseWriter, r *http.Request, requestedDate time.Time) {
	startTime := requestedDate.Format(time.RFC3339)
	endTime := requestedDate.Add(time.Hour * 24).Format(time.RFC3339)
	//fmt.Fprintf(w, "%v\n%v", startTime, endTime)

	token := getBearerToken(w, r)
	neurioURL := "https://api.neur.io/v1/samples?sensorId=" + os.Getenv("NEURIO_SENSOR_ID") + "&start=" + startTime + "&end=" + endTime + "&granularity=minutes&frequency=30&perPage=50"
	//fmt.Fprintf(w, "neurio url: %v<br />", neurioURL)

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", neurioURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll() error: %v<br />", err)
		return
	}

	var samples []sample
	json.Unmarshal([]byte(string(data)), &samples)
	output, _ := json.Marshal(samples)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(output)

	//fmt.Fprintf(w, "%v %v", string(data), samples)
}

type sample struct {
	ConsumptionPower int    `json:"consumptionPower"`
	Timestamp        string `json:"timestamp"`
}

func getBearerToken(w http.ResponseWriter, r *http.Request) string {
	neurioURL := "https://api.neur.io/v1/oauth2/token"

	v := url.Values{}
	v.Set("grant_type", "client_credentials")
	v.Set("client_id", os.Getenv("NEURIO_CLIENT_ID"))
	v.Set("client_secret", os.Getenv("NEURIO_CLIENT_SECRET"))

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.PostForm(neurioURL, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	//fmt.Fprintf(w, "HTTP call returned %v", resp)

	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// c := &http.Client{}
	// resp, err := c.Do(req)
	// if err != nil {
	// 	fmt.Fprintf(w, "http.Do() error: %v<br />", err)
	// 	return
	// }
	// defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll() error: %v<br />", err)
		return ""
	}

	var oar oauthResponse
	json.Unmarshal([]byte(string(data)), &oar)
	//fmt.Fprintf(w, "access_token? %v", oar.AccessToken)
	//fmt.Fprintf(w, "read resp.Body successfully:<br />%v<br />", string(data))
	return oar.AccessToken
}

type oauthResponse struct {
	AccessToken string `json:"access_token"`
}
