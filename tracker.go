package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Results []struct {
		Name string
		Net  string
		Pad  struct {
			Location struct {
				Id   int
				Name string
			}
		}
		Rocket struct {
			LauncherStages []struct {
				Landing struct {
					Type struct {
						Id int
					}
				}
			} `json:"launcher_stage"`
		}
		VidURLs []struct {
			Url string
		}
	}
}

func main() {
	url := "https://ll.thespacedevs.com/2.2.0/launch/upcoming/?format=json&limit=20&mode=detailed"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	launches, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var foo Response
	if err := json.Unmarshal([]byte(launches), &foo); err != nil {
		panic(err)
	}

	for _, result := range foo.Results {
		description := result.Name

		// Add *'s to launches from Vandenburg
		if 11 == result.Pad.Location.Id {
			description = fmt.Sprintf("***%s***", description)
		}

		// Parse time
		t, err := time.Parse(time.RFC3339, result.Net)
		if err != nil {
			panic(err)
		}

		// Add webcast URL if launch is within 4 hours
		if time.Until(t).Hours() < 4 {
			for _, vid := range result.VidURLs {
				description = fmt.Sprintf("%s %s", description, vid.Url)
			}
		}

		// Detect RTLS landings
		for _, launcherStage := range result.Rocket.LauncherStages {
			if 2 == launcherStage.Landing.Type.Id {
				description = fmt.Sprintf("!!!%s", description)
				break
			}
		}

		fmt.Printf("%s %s\n", t.Local(), description)
	}
}
