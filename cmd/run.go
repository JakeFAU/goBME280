/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	aio "github.com/jakefau/goAdafruit"
	"github.com/jakefau/rpi-devices/dev"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/io/i2c"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ioUser := viper.GetString("IOUSER")
		ioKey := viper.GetString("IOKEY")

		//get the applicationn ready
		app := aio.NewClient(ioKey, ioUser)
		app.BaseURL, _ = url.Parse(baseURL)

		//get the feeds
		tempFeed := getFeed(fmt.Sprintf("%v.temperature", feedPrefix), *app)
		humidFeed := getFeed(fmt.Sprintf("%v.humidy", feedPrefix), *app)
		pressureFeed := getFeed(fmt.Sprintf("%v.pressure", feedPrefix), *app)

		//get the data
		d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, address)
		if err != nil {
			log.Fatal(err)
		}
		b := dev.New(d)
		err = b.Init()
		if err != nil {
			log.Fatal(err)
		}
		t, h, p, _ := b.EnvData()
		t = convert64(t)
		h = convert64(h)
		p = convert64(p)
		if units == "english" {
			t = toFahrenheit(t)
			p = toMercury(p)
		}

		//now set the feeds
		app.SetFeed(tempFeed)
		app.Data.Create(&aio.Data{Value: t})
		app.SetFeed(humidFeed)
		app.Data.Create(&aio.Data{Value: h})
		app.SetFeed(pressureFeed)
		app.Data.Create(&aio.Data{Value: p})
		log.Printf("Temp: %fF, Press: %f, Hum: %f%%\n", t, p, h)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func getFeed(feedKey string, client aio.Client) *aio.Feed {
	feed, _, err := client.Feed.Get(feedKey)
	if err != nil {
		log.Fatal(err)
	}
	return feed
}

func toFahrenheit(c float64) float64 {
	return (c * 9 / 5.0) + 32
}

func toMercury(m float64) float64 {
	return m * 0.02953
}

func convert64(value float64) string {
	return strconv.FormatFloat(float64(value), 'e', 2, 32)
}
