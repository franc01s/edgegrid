package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/go-co-op/gocron/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

type Application struct {
	client  *http.Client
	results *simplejson.Json
}

func newApplication() *Application {

	app := Application{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	// init http client 
	app.client = &http.Client{
		Timeout: 10 * time.Second, // Set timeout
		Jar:     jar,
	}

	return &app
}

var count int

func (a *Application) pullEdge() {

	uri := fmt.Sprintf("https://monitoringapi.solaredge.com/site/%s/currentPowerFlow?api_key=%s", os.Getenv("EDGEGRID_SITE"), os.Getenv("EDGEGRID_API_KEY"))
	req, _ := http.NewRequest("GET", uri, nil)
	resp, err := a.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 429 {
		slog.Error("failed to pull edgeGrid tags: %s", string(body))
		return
		// recreate http client, prevent daily rate limit
		// I don't understand why, but it seems to work
	} else if resp.StatusCode == 429 {
		jar, _ := cookiejar.New(nil)
		a.client = &http.Client{
			Timeout: 10 * time.Second, // Set timeout
			Jar:     jar,
		}
		return
	}

	count++
	slog.Info(fmt.Sprintf("updating EgeGrid values: %d", count))

	currentPowerFlow, err := simplejson.NewJson(body)
	if err != nil {
		slog.Error("error making simple json", err.Error())
	}

	a.results = currentPowerFlow

}

func (a *Application) getResult(w http.ResponseWriter, req *http.Request) {
	power := struct {
		Load float64 `json:"load"`
		Grid float64 `json:"grid"`
	}{Load: a.results.Get("siteCurrentPowerFlow").Get("LOAD").Get("currentPower").MustFloat64(),
		Grid: a.results.Get("siteCurrentPowerFlow").Get("GRID").Get("currentPower").MustFloat64()}

	powerJson, _ := json.Marshal(power)
	w.Write(powerJson)
}

func main() {

	app := newApplication()

	s, err := gocron.NewScheduler(
		gocron.WithGlobalJobOptions(
			gocron.WithSingletonMode(
				gocron.LimitModeReschedule)),
	)
	if err != nil {
		panic(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			30*time.Second,
		),
		gocron.NewTask(
			app.pullEdge,
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(j.ID())
	s.Start()
	http.HandleFunc("/", app.getResult)

	http.ListenAndServe(":8080", nil)

}