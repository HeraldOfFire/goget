package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// The data struct decoded config.json
type Config struct {
	BasePath string `json:"basePath"`
	Groups   []struct {
		Path         string   `json:"path"`
		Format       string   `json:"format"`
		URLTemplate  string   `json:"urlTemplate"`
		URLVariables []string `json:"urlVariables"`
	} `json:"groups"`
}

//Download file from URL
func downloadFile(filepath string, url string) (err error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMovedPermanently {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("STARTED: " + filepath)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//Read and decode config.json
func getConfig() (Config, error) {

	data, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var conf Config

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil

}

func main() {

	conf, err := getConfig()
	if err != nil {
		fmt.Println("ERROR: config.json invalid or not found")
		return
	}

	//init waitGroup for the goroutines
	var wg sync.WaitGroup

	//loop through config groups
	for i := 0; i < len(conf.Groups); i++ {

		//utility variables
		dirPath := conf.BasePath + conf.Groups[i].Path
		fileFormat := conf.Groups[i].Format
		urlTemplate := conf.Groups[i].URLTemplate
		urlVariables := conf.Groups[i].URLVariables

		//make directory if not present
		os.MkdirAll(dirPath, os.ModePerm)

		fmt.Println("DOWNLOADING GROUP: " + dirPath)

		//add slots in the waitgroup
		//wg.Add(len(urlVariables))

		for j := 0; j < len(urlVariables); j++ {

			//manage variables "from-to"
			if strings.Contains(urlVariables[j], "->") {

				fromTo := strings.Split(urlVariables[j], "->")

				from, _ := strconv.Atoi(fromTo[0])
				to, _ := strconv.Atoi(fromTo[1])

				wg.Add(to-from)

				for k := from; k <= to; k++ {

					//start a goroutine for each urlVariable
					go func(filepath string, url string) {
	
						//when all the goroutines are done, end the waitgroup
						defer wg.Done()
	
						//download requested files
						err = downloadFile(filepath, url)
						if err != nil {
							fmt.Println("FAILED: " + url + " ERROR: " + err.Error())
							return
						}
	
						fmt.Println("DONE: " + filepath)
	
					}(dirPath+"/"+fmt.Sprint(k)+fileFormat, strings.Replace(urlTemplate, "<<variable>>", fmt.Sprint(k), 1))
				}

			} else {

				wg.Add(1)

				//start a goroutine for each urlVariable
				go func(filepath string, url string) {

					//when all the goroutines are done, end the waitgroup
					defer wg.Done()

					//download requested files
					err = downloadFile(filepath, url)
					if err != nil {
						fmt.Println("FAILED: " + filepath + " ERROR: " + err.Error())
						return
					}

					fmt.Println("DONE: " + filepath)

				}(dirPath+"/"+urlVariables[j]+fileFormat, strings.Replace(urlTemplate, "<<variable>>", urlVariables[j], 1))
			}

		}
		//wait for the goroutines in the waitgroup to end
		wg.Wait()
	}
}
