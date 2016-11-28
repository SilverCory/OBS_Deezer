package deezer

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"strings"

	"image"
	_ "image/jpeg" // Used for JPEG image.

	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// SongData struct for the song data and image.
type SongData struct {
	SongID    string `json:"SNG_ID"`
	SongTitle string `json:"SNG_TITLE"`

	ArtistID   string `json:"ART_ID"`
	ArtistName string `json:"ART_NAME"`

	AlbumID      string      `json:"ALB_ID"`
	AlbumTitle   string      `json:"ALB_TITLE"`
	AlbumPicture string      `json:"ALB_PICTURE"`
	AlbumImage   image.Image `json:"-"`
}

// Deezer data and instance.
type Deezer struct {
	ProfileID int // The profile ID E.G. 875499801
	SongData  SongData
	Online    bool
}

// CreateDeezer initially creates the instance and fetch first.
func CreateDeezer(UserID int) (*Deezer, error) {

	instance := &Deezer{
		ProfileID: UserID,
	}

	err := instance.Fetch()

	return instance, err

}

// Fetch fetches all the data from deezer.
func (d *Deezer) Fetch() error {

	doc, err := goquery.NewDocument("http://www.deezer.com/profile/" + strconv.Itoa(d.ProfileID))
	if err != nil {
		return err
	}

	var data string

	// If this doesn't work add :nth-child(2) to the script.
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "window.__DZR_APP_STATE__ = ") {
			data = strings.TrimPrefix(s.Text(), "window.__DZR_APP_STATE__ = ")
		}
	})

	if len(data) <= 0 {
		return errors.New("No valid data received from deezer!")
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		return err
	}

	dat = dat["TAB"].(map[string]interface{})["home"].(map[string]interface{})

	if val, ok := dat["online"]; ok {

		d.SongData, err = d.processSongData(val.(map[string]interface{}))
		if err != nil {
			d.Online = false
			return err
		}

		d.Online = true
		d.fetchImage()

	} else {
		d.Online = false
		d.SongData = SongData{}
	}

	return nil

}

// fetchImage Fetches the image from deezer cdn.
func (d *Deezer) fetchImage() {

	url := "http://cdn-images.deezer.com/images/cover/"
	url += d.SongData.AlbumPicture
	url += "/400x400-000000-80-0-0.jpg"

	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil || response == nil {
		fmt.Println("There was a http error fetching the album photo!", err)
		d.SongData.AlbumImage = nil
		return
	}

	image, _, err := image.Decode(response.Body)

	if err != nil {
		fmt.Println("There was an error fetching the album photo!", err)
		d.SongData.AlbumImage = nil
		return
	}

	d.SongData.AlbumImage = image

}

// processSongData transforms the data to and from JSON in order to put it in the SongData struct.
func (d *Deezer) processSongData(data map[string]interface{}) (SongData, error) {

	songData := SongData{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return songData, err
	}

	err = json.Unmarshal(jsonData, &songData)
	return songData, err

}
