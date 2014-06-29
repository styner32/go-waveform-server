package main

import (
	"fmt"
	"github.com/hoisie/web"
	"github.com/styner32/go-wave-to-json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	DOWNLOAD_FOLDER = "tmp/"
	AUTIO_EXTENSION = "mp3"
	JSON_EXTENSION  = "json"
)

func create(ctx *web.Context, val string) {
	videoPath := videoPath(ctx.Params["stream"])
	audioPath := convertExtension(videoPath, AUTIO_EXTENSION)
	waveformPath := convertExtension(videoPath, JSON_EXTENSION)
	downloadFromUrl(ctx.Params["stream"], videoPath)
	convertToMp3(videoPath, audioPath)
	waveform.Generate(audioPath, waveformPath)
	ctx.WriteString("OK: " + val)
}

func downloadFromUrl(url string, downloadPath string) {
	fmt.Printf("Downloading %s\n", url)

	output, err := os.Create(downloadPath)
	if err != nil {
		fmt.Println("Error while creating", downloadPath, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Printf("The stream has been downloaded: %s bytes\n", n)
}

func videoPath(url string) string {
	tokens := strings.Split(url, "/")
	return fmt.Sprintf("%s%s", DOWNLOAD_FOLDER, tokens[len(tokens)-1])
}

func convertExtension(filePath string, extension string) string {
	tokens := strings.Split(filePath, ".")
	return fmt.Sprintf("%s.%s", tokens[0], extension)
}

func convertToMp3(videoFilePath string, audioPath string) {
	fmt.Printf("Converting %s to %s\n", videoFilePath, audioPath)
	cmd := exec.Command("ffmpeg", "-i", videoFilePath, "-y", "-f", "mp3", "-ab", "192000", "-vn", audioPath)
	_, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	web.Post("/(waveforms.*)", create)
	web.Run("0.0.0.0:9000")
}
