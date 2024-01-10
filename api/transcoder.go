package api

import (
	"log"
	"os/exec"
	"sync"
)

func Transcoder() {
	qualities := []string{"720", "480", "360", "240", "144"}
	
	var wg sync.WaitGroup

	for _, quality := range qualities {
		wg.Add(1)
		go func(quality string) {
			defer wg.Done()
			// Here you would adjust the FFmpeg command parameters based on the desired quality
			outputFilename := "output_" + quality + ".mp4"
			cmd := exec.Command("ffmpeg", "-i", "input/input.webm", "-vf", "scale='trunc(oh*a/2)*2:"+quality, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-strict", "experimental", "output/"+outputFilename)

			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error transcoding %s: %v", quality, err)
				return
			}
			log.Printf("Output for %s: %s", quality, output)
		}(quality)
	}
	wg.Wait()
}