package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings"
	"image"
	"image/png"

	"../../config"
	"../location"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

type Forecast struct {
	Cod     string        `json:"cod"`
	Weather []WeatherData `json:"list"`
	City    CityData      `json:"city"`
}

type WeatherData struct {
	Time   int64       `json:"dt"`
	Main   MainData    `json:"main"`
	Wind   WindData    `json:"wind"`
	Clouds CloudsData  `json:"clouds"`
	WDesc  []WDescData `json:"weather"`
}

func (w WeatherData) TZTime() time.Time {
	return time.Unix(w.Time, 0).UTC().Add(time.Hour * time.Duration(config.General.Timezone))
}

type WDescData struct {
	Id   int64  `json:"id"`
	Main string `json:"main"`
	Desc string `json:"description"`
	Icon string	`json:"icon"`
}

type MainData struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity int     `json:"humidity"`
}

type WindData struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type CloudsData struct {
	All int `json:"all"`
}

type CityData struct {
	Name string `json:"name"`
}
// Super bad code below. Be careful!
func DrawOne(temp, hum, clo int, time, icon string) image.Image {
	dpc := gg.NewContext(300,400)
	dpc.SetRGBA(0,0,0,0)
	dpc.Clear()
	dpc.SetRGB(1, 1, 1)

	res, err := http.Get(fmt.Sprintf("http://openweathermap.org/img/w/%v.png",icon))
	if err != nil || res.StatusCode != 200 {
		fmt.Println(err)
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	dpc.Push()
	dpc.Scale(3,3)
	dpc.DrawImage(m, 25, 12)
	dpc.Pop()

	if err := dpc.LoadFontFace("arial.ttf", 50); err != nil {
		fmt.Printf("Image font: %v",err)
	}
	dpc.DrawStringAnchored(time, 150, 30, 0.5, 0.5)
	dpc.DrawStringAnchored(fmt.Sprintf("H: %v%%",hum), 150, 280, 0.5, 0.5)
	dpc.DrawStringAnchored(fmt.Sprintf("C: %v%%",clo), 150, 330, 0.5, 0.5)

	if err := dpc.LoadFontFace("arial.ttf", 80); err != nil {
		fmt.Printf("Image font: %v",err)
	}
	dpc.DrawStringAnchored(fmt.Sprintf("%v°", temp), 150, 200, 0.5, 0.5)

	return dpc.Image()
}

func GetWeatherImage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		var (
		forecast      Forecast
		city          string = config.Weather.City
	)

	if len(args) > 0 {
		city = strings.Join(args, "+")
	}

	loc, err := location.New(city)
	if err != nil {
		fmt.Printf("Location API: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("location_404"))
		return
	}

	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%v&lon=%v&lang=%v&units=metric&appid=%v",
		newlat, newlng, config.General.Language, config.Weather.WeatherToken))
	if err != nil {
		fmt.Printf("Weather API: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_api_error"))
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		fmt.Printf("Weather Decode: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_parse_error"))
		return
	}
	
	dc := gg.NewContext(1500, 400)
	dc.SetRGBA(0, 0, 0, 0.7)
	dc.Clear()
	for i := 0; i<6; i++ {
		dc.DrawImage(DrawOne(int(forecast.Weather[i].Main.TempMin), 
							 forecast.Weather[i].Main.Humidity, 
							 int(forecast.Weather[i].Clouds.All), 
							 fmt.Sprintf("%.2v:00", forecast.Weather[i].TZTime().Hour()), 
							 forecast.Weather[i].WDesc[0].Icon),300 * i, 0)
	}

	buf := new(bytes.Buffer)
	pngerr := png.Encode(buf, dc.Image())
	if pngerr != nil {
		fmt.Printf("Image: %v",pngerr)
	}
	s.ChannelFileSend(m.ChannelID, "weather.png", buf)
}