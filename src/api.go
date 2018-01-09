package duckCurve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

var (
	bearerToken   string
	dayCache      map[string][]sample
	location      *time.Location
	dateFormatYMD string
)

func init() {
	dayCache = make(map[string][]sample)
	location, _ = time.LoadLocation("America/New_York")
	dateFormatYMD = "2006-01-02"
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	pathParts := strings.Split(r.URL.Path, "/")
	mode := pathParts[2]

	switch mode {
	case "day":
		force := (len(pathParts) > 4) && pathParts[4] == "force"
		dateString := pathParts[3]

		parsedDate, _ := time.ParseInLocation(dateFormatYMD, dateString, location)
		samples, err := getDayDatapoints(w, r, dateString, parsedDate, force)
		if err != nil {
			fmt.Fprint(w, "Error calling getDayDatapoints")
			return
		}
		output, _ := json.Marshal(samples)
		w.Header().Add("Content-Type", "application/json")
		w.Write(output)
	case "average":
		w.Header().Add("Content-Type", "text/plain")
		dayTypes := pathParts[3]
		dateString := pathParts[4]
		parsedDate, _ := time.ParseInLocation(dateFormatYMD, dateString, location)
		workingSamples, _ := getDayDatapoints(w, r, dateString, parsedDate, false)
		parsedDate = truncateAtLocation(parsedDate.Add(time.Hour * 25))
		daysIncluded := 1
		for currentDate := parsedDate; currentDate.Month() == parsedDate.Month(); currentDate = truncateAtLocation(currentDate.Add(time.Hour * 25)) {
			day := currentDate.Weekday()
			include := false
			switch dayTypes {
			case "all":
				include = true
			case "weekends":
				include = day == time.Sunday || day == time.Saturday
			case "weekdays":
				include = day != time.Sunday && day != time.Saturday
			}
			if include {
				dateString := currentDate.Format(dateFormatYMD)
				temp, _ := getDayDatapoints(w, r, dateString, currentDate, false)
				if len(temp) > 0 {
					daysIncluded++
					for i, v := range temp {
						if i < len(workingSamples)-1 {
							workingSamples[i].ConsumptionPower += v.ConsumptionPower
						}
					}
				}
			}
		}
		for i := range workingSamples {
			workingSamples[i].ConsumptionPower /= daysIncluded
		}
		output, _ := json.Marshal(workingSamples)
		w.Header().Add("Content-Type", "application/json")
		w.Write(output)
	case "warmup":
		command := pathParts[3]
		checkWarmup(w, r, command)
	default:
		fmt.Fprint(w, "Neurio data aggregator API. Example: /api/day/2017-11-01")
	}
}

type warmupStatus struct {
	Value string
}

func checkWarmup(w http.ResponseWriter, r *http.Request, command string) {
	fmt.Fprint(w, "About to run warmup case\n")

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "warmupStatus", "warmup_status", 0, nil)

	v := new(warmupStatus)
	v.Value = "Test"

	switch command {
	case "run":
		kickOffWarmup(ctx, w, r, key)
	case "status":
		v := new(warmupStatus)
		if err := datastore.Get(ctx, key, v); err == nil {
			fmt.Fprintf(w, "%v", v.Value)
		} else {
			fmt.Fprintf(w, "No warmup has been run. %v", err)
		}
	}
}

func kickOffWarmup(ctx context.Context, w http.ResponseWriter, r *http.Request, key *datastore.Key) {
	v := new(warmupStatus)
	v.Value = "Status: Starting warmup\n"
	fmt.Fprint(w, v)
	datastore.Put(ctx, key, v)
	startDate := time.Date(2017, time.Month(9), 15, 0, 0, 0, 0, location)
	today := truncateAtLocation(time.Now())
	for currentDate := today; currentDate.After(startDate); currentDate = truncateAtLocation(currentDate.Add(time.Hour * -23)) {
		dateString := currentDate.Format(time.RFC3339)
		v.Value = fmt.Sprintf("Currently retrieving %v at %v\n", dateString, time.Now())
		fmt.Fprintf(w, "%v", v.Value)
		datastore.Put(ctx, key, v)
		getDayDatapoints(w, r, dateString, currentDate, false)
	}
	v.Value = "Warmup complete\n"
	datastore.Put(ctx, key, v)
}

func getDayDatapoints(w http.ResponseWriter, r *http.Request, dateString string, requestedDate time.Time, force bool) ([]sample, error) {
	// Do not try to get data from today or later
	if time.Now().Sub(requestedDate) < time.Hour*24 {
		return []sample{}, nil
	}

	// Can we return a cached copy of the data?
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "sampleContainer", "day_"+dateString, 0, nil)
	v := new(sampleContainer)
	if !force {
		if err := datastore.Get(ctx, key, v); err == nil {
			return v.Samples, nil
		}
	}

	// Get the data from Neurio
	startTime := requestedDate.Format(time.RFC3339)
	endTime := requestedDate.Add(time.Hour * 24).Add(time.Minute * -30).Format(time.RFC3339)
	if bearerToken == "" {
		bearerToken = getBearerToken(w, r)
	}
	neurioURL := "https://api.neur.io/v1/samples?sensorId=" + os.Getenv("NEURIO_SENSOR_ID") + "&start=" + startTime + "&end=" + endTime + "&granularity=minutes&frequency=30&perPage=50"
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", neurioURL, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	res, err := client.Do(req)

	// Process the data
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll() error: %v<br />", err)
		return nil, errors.New("Error")
	}
	var samples []sample
	json.Unmarshal([]byte(string(data)), &samples)
	v.Samples = samples
	if _, err := datastore.Put(ctx, key, v); err != nil {
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	return samples, nil
}

type sampleContainer struct {
	Samples []sample
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll() error: %v<br />", err)
		return ""
	}

	var oar oauthResponse
	json.Unmarshal([]byte(string(data)), &oar)
	return oar.AccessToken
}

type oauthResponse struct {
	AccessToken string `json:"access_token"`
}

func truncateAtLocation(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
