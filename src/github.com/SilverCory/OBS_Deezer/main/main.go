package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"

	"strings"

	"github.com/SilverCory/OBS_Deezer/deezer"
)

var possibleFormats string

func init() {
	possibleFormats += "%ALBUM_ID%\n\t\t"
	possibleFormats += "%ALBUM_PICTURE%\n\t\t"
	possibleFormats += "%ALBUM_TITLE%\n\t\t"

	possibleFormats += "%ARTIST_ID%\n\t\t"
	possibleFormats += "%ARTIST_NAME%\n\t\t"

	possibleFormats += "%SONG_ID%\n\t\t"
	possibleFormats += "%SONG_TITLE%\n\t\t"
}

func main() {

	refreshRate := flag.Int("time", 10, "Refresh rate in seconds. Zero or less results in no refresh.")
	id := flag.Int("id", 875499801, "The deezer proile ID.")
	fileName := flag.String("saveName", "deezer_now", "The filename to save to. If empty no save.")
	txtFormat := flag.String("txtFormat", "%SONG_TITLE%\\n\\n%ARTIST_NAME%", "The format of the title. Possible formats are:\n\t\t"+possibleFormats)
	flag.Parse()

	// Create the instance whitch calls fetch first.
	d, err := deezer.CreateDeezer(*id)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Do first output, hash calculation and file writes.
	var currentHash = createHash(d)
	doOutput(d)
	writeFile(d, *fileName, *txtFormat)

	// Return if refreshing is disabled.
	if *refreshRate <= 0 {
		return
	}

	// Refresh using a ticker.
	ticker := time.NewTicker(time.Duration(*refreshRate) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := d.Fetch()
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				newHash := createHash(d)
				if strings.Compare(currentHash, newHash) != 0 {
					currentHash = newHash
					doOutput(d)
					writeFile(d, *fileName, *txtFormat)
				}
			}
		}
	}

}

// Create a hash used for compaison so that we're not writing the same files.
func createHash(d *deezer.Deezer) string {
	var data string

	if !d.Online {
		data = "nol"
	} else {
		data += d.SongData.AlbumID + "_"
		data += d.SongData.SongID + "_"
		data += d.SongData.ArtistID
	}

	return data

}

// Write the files.
func writeFile(d *deezer.Deezer, fileName string, txtFormat string) {

	// Return if filewriting is disabled.
	if len(fileName) <= 0 {
		return
	}

	// Write the JSON file.
	var err error
	var data []byte

	if d.Online {
		data, err = json.Marshal(d.SongData)
		if err != nil {
			fmt.Println(err)
			err = nil
		}
	} else {
		data = []byte("{online: \"false\"}")
	}

	err = ioutil.WriteFile("./"+fileName+".json", data, 0644)
	if err != nil {
		fmt.Println(err)
		err = nil
	}

	// Write the txt file.
	if len(txtFormat) > 0 {

		// Yes this is messy, yes I could have used reflection, no I didn't I know.
		if d.Online {
			txtFormat = strings.Replace(txtFormat, "\\n", "\n", -1)

			txtFormat = strings.Replace(txtFormat, "%ALBUM_ID%", d.SongData.AlbumID, -1)
			txtFormat = strings.Replace(txtFormat, "%ALBUM_PICTURE%", d.SongData.AlbumPicture, -1)
			txtFormat = strings.Replace(txtFormat, "%ALBUM_TITLE%", d.SongData.AlbumTitle, -1)

			txtFormat = strings.Replace(txtFormat, "%ARTIST_ID%", d.SongData.ArtistID, -1)
			txtFormat = strings.Replace(txtFormat, "%ARTIST_NAME%", d.SongData.ArtistName, -1)

			txtFormat = strings.Replace(txtFormat, "%SONG_ID%", d.SongData.SongID, -1)
			txtFormat = strings.Replace(txtFormat, "%SONG_TITLE%", d.SongData.SongTitle, -1)
		} else {
			txtFormat = ""
		}

		err = ioutil.WriteFile("./"+fileName+".txt", []byte(txtFormat), 0644)
		if err != nil {
			fmt.Println(err)
			err = nil
		}

	}

	// Write the image.
	if !d.Online || d.SongData.AlbumImage == nil {
		os.Remove("./" + fileName + ".jpg")
		if err != nil {
			fmt.Println(err)
			err = nil
		}
		return
	}

	out, err := os.Create("./" + fileName + ".jpg")
	if err != nil {
		fmt.Println(err)
		err = nil
	}

	var opt jpeg.Options
	opt.Quality = 80

	err = jpeg.Encode(out, d.SongData.AlbumImage, &opt)
	if err != nil {
		fmt.Println(err)
	}

}

// Log the data to console.
func doOutput(d *deezer.Deezer) {
	if d.Online {
		fmt.Println(d.SongData.SongTitle + " - " + d.SongData.ArtistName)
	} else {
		fmt.Println("User is not online.")
	}
}
