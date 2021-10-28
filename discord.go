package main

import (
  "net/http"
  "fmt"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "log"
  "time"
  "io/ioutil"
  "encoding/json"
  "os"
  "strconv"
)

var (
  gaugeUsersByChannel = prometheus.NewGaugeVec(
     prometheus.GaugeOpts{
        Namespace: "discord",
        Name:      "channel_users",
        Help:      "Represents the number of discord users by channel",
     },
		[]string{
			"channel",
		})
)
var (
  gaugeUsers = prometheus.NewGaugeVec(
     prometheus.GaugeOpts{
        Namespace: "discord",
        Name:      "users",
        Help:      "Represents the online users",
     },
		[]string{
			"user",
		})
)
var (
  gaugeTotalUsers = prometheus.NewGauge(
     prometheus.GaugeOpts{
        Namespace: "discord",
        Name:      "total_users",
        Help:      "Represents the number of total discord users",
      })
)

type discord struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	InstantInvite interface{} `json:"instant_invite"`
	Channels      []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Position int    `json:"position"`
	} `json:"channels"`
	Members []struct {
		ID            string      `json:"id"`
		Username      string      `json:"username"`
		Discriminator string      `json:"discriminator"`
		Avatar        interface{} `json:"avatar"`
		Status        string      `json:"status"`
		Game          struct {
			Name string `json:"name"`
		} `json:"game,omitempty"`
		AvatarURL string `json:"avatar_url"`
		Deaf      bool   `json:"deaf,omitempty"`
		Mute      bool   `json:"mute,omitempty"`
		SelfDeaf  bool   `json:"self_deaf,omitempty"`
		SelfMute  bool   `json:"self_mute,omitempty"`
		Suppress  bool   `json:"suppress,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	} `json:"members"`
	PresenceCount int `json:"presence_count"`
}

func getChannelNameById(id string, people1 discord) string {
  for _, element := range people1.Channels {
    if (element.ID == id) {
      return element.Name
    }
  }

  return "Unknown"
}

func main() {
  if _, ok := os.LookupEnv("SERVERID"); !ok {
    log.Fatal("Missing environment variable 'SERVERID'")
  }
  if _, ok := os.LookupEnv("REFRESH_INTERVAL"); !ok {
    os.Setenv("REFRESH_INTERVAL", "60")
  }

  http.Handle("/metrics", promhttp.Handler())

  prometheus.MustRegister(gaugeUsersByChannel)
  prometheus.MustRegister(gaugeUsers)
  prometheus.MustRegister(gaugeTotalUsers)

  var numberOfTotalUser float64 = 0
  var numberOfUsersByChannel = make(map[string]float64)

  go func() {
     for {
        //fmt.Println("Fetching new data from discords API")

        url := fmt.Sprintf("https://discordapp.com/api/guilds/%s/widget.json", os.Getenv("SERVERID"))
        spaceClient := http.Client{
          Timeout: time.Second * 2,
	      }
        req, err := http.NewRequest(http.MethodGet, url, nil)

        if err != nil {
          log.Fatal(err)
	      }

        res, getErr := spaceClient.Do(req)

        if getErr != nil {
          log.Fatal(getErr)
	      }

        //fmt.Printf("HTTP: %s\n", res.Status)

        body, readErr := ioutil.ReadAll(res.Body)

        if readErr != nil {
          log.Fatal(readErr)
	      }

        widgetData := discord{}
	      jsonErr := json.Unmarshal(body, &widgetData)

        if jsonErr != nil {
          log.Fatal(jsonErr)
	      }

        numberOfTotalUser = 0
        numberOfUsersByChannel = map[string]float64{}
        gaugeUsers.Reset()
        for _, element := range widgetData.Members {
          if (len(element.ChannelID) > 0) {
            numberOfTotalUser++
            numberOfUsersByChannel[element.ChannelID]++
            gaugeUsers.WithLabelValues(element.Username).Set(1)
          }
        }

        gaugeUsersByChannel.Reset()
        for _, channel := range widgetData.Channels {
          if _, ok := numberOfUsersByChannel[channel.ID]; ok {
            gaugeUsersByChannel.WithLabelValues(getChannelNameById(channel.ID, widgetData)).Set(numberOfUsersByChannel[channel.ID])
          } else {
            gaugeUsersByChannel.WithLabelValues(getChannelNameById(channel.ID, widgetData)).Set(0)
          }
        }

        gaugeTotalUsers.Set(numberOfTotalUser)

        refreshInterval, _ := strconv.Atoi(os.Getenv("REFRESH_INTERVAL"))
        time.Sleep(time.Second * time.Duration(refreshInterval))
     }
  }()

  log.Fatal(http.ListenAndServe(":8080", nil))
}
