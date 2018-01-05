package duckCurve

import (
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
	bearerToken string
	dayCache    map[string][]sample
	test        map[string]int
)

func init() {
	test = make(map[string]int)
	dayCache = make(map[string][]sample)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	pathParts := strings.Split(r.URL.Path, "/")
	mode := pathParts[2]
	dateString := pathParts[3]
	loc, _ := time.LoadLocation("America/New_York")
	parsedDate, _ := time.ParseInLocation("2006-01-02", dateString, loc)

	switch mode {
	case "day":
		samples, err := getDayDatapoints(w, r, dateString, parsedDate)
		if err != nil {
			fmt.Fprint(w, "Error calling getDayDatapoints")
			return
		}
		output, _ := json.Marshal(samples)
		w.Header().Add("Content-Type", "application/json")
		w.Write(output)
	case "test":
		fmt.Fprint(w, "We will now perform a test\n")

		count := test[dateString] + 1
		test[dateString] = count
		fmt.Fprintf(w, "test map: %#v\n", test)
		fmt.Fprintf(w, "Count: %v", count)
	default:
		fmt.Fprint(w, "Neurio data aggregator API. Example: /api/day/2017-11-01")
	}
}

func getDayDatapoints(w http.ResponseWriter, r *http.Request, dateString string, requestedDate time.Time) ([]sample, error) {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, "sampleContainer", "day_"+dateString, 0, nil)
	v := new(sampleContainer)
	if err := datastore.Get(ctx, key, v); err == nil {
		//fmt.Fprintf(w, "Results: %#v", v)
		return v.Samples, nil
		//http.Error(w, err.Error(), 500)
		//return nil, err
	}

	// if dayCache[dateString] != nil {
	// 	fmt.Fprint(w, "Found in cache")
	// 	return dayCache[dateString], nil
	// }
	//fmt.Fprint(w, "Not in cache")
	startTime := requestedDate.Format(time.RFC3339)
	endTime := requestedDate.Add(time.Hour * 24).Add(time.Minute * -30).Format(time.RFC3339)
	//fmt.Fprintf(w, "%v\n%v", startTime, endTime)

	if bearerToken == "" {
		bearerToken = getBearerToken(w, r)
	}

	neurioURL := "https://api.neur.io/v1/samples?sensorId=" + os.Getenv("NEURIO_SENSOR_ID") + "&start=" + startTime + "&end=" + endTime + "&granularity=minutes&frequency=30&perPage=50"
	//fmt.Fprintf(w, "neurio url: %v<br />", neurioURL)

	//ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", neurioURL, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	res, err := client.Do(req)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, "ioutil.ReadAll() error: %v<br />", err)
		return nil, errors.New("Error")
	}

	var samples []sample
	json.Unmarshal([]byte(string(data)), &samples)

	//dayCache[dateString] = samples
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
