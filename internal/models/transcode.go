package models

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Transcode(filename string) {
	startTime := time.Now()

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", "internal/videos/"+filename)
	// cmd := exec.Command("ffmpeg", "-i", "internal/videos/"+filename, "-v", "quiet", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error getting dimensions:", err)
		return
	}
	fmt.Println("Command output:", out.String())
	dimensions := strings.Split(strings.TrimSpace(out.String()), ",")
	if len(dimensions) != 2 {
		fmt.Println("Failed to parse dimensions")
		return
	}
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	qualities := []string{"1080", "720", "480", "360", "240", "144"}
	
	var wg sync.WaitGroup

	for _, quality := range qualities {
		wg.Add(1)
		go func(quality string) {
			defer wg.Done()

			q, err := strconv.Atoi(quality)
			if err != nil {
				log.Printf("Invalid quality %s: %v", quality, err)
				return
			}
			var scaleFilter string
			if width > height {
				// Landscape
				newHeight := (int(float64(height) / float64(width) * float64(q))/2)*2
				scaleFilter = fmt.Sprintf("scale=%d:%d", q, newHeight)
			} else {
				// Portrait
				newWidth := (int(float64(width) / float64(height) * float64(q))/2)*2
				scaleFilter = fmt.Sprintf("scale=%d:%d", newWidth, q)
			}
			// Here you would adjust the FFmpeg command parameters based on the desired quality
			outputFilename := "output_"+filename+ quality + ".mp4"
			
			cmd := exec.Command("ffmpeg", "-i", "internal/videos/"+filename, "-vf", scaleFilter, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-strict", "experimental", "internal/output/"+outputFilename)
			fmt.Println("Command:", cmd)
			err = cmd.Run()
			if err != nil {
				log.Printf("Error transcoding %s: %v", quality, err)
				return
			}
			log.Printf("Output for %s", quality)
		}(quality)
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("Total processing time: %v\n", endTime.Sub(startTime))
}	