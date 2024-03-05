package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var jobIds = make(map[string]int)
var jobIdsMtx sync.Mutex
var JobPubSub = NewPubSub[int]()
var EventSub = NewPubSub[string]()
var UUIDARRAY = []string{}
var UUIDARRAYMtx sync.Mutex
var keys []string
var array []string

const maxWorkers = 2

var wg sync.WaitGroup
var sem = make(chan struct{}, maxWorkers)
var uuidNFile []JobsUUIDFile

type JobsUUIDFile struct {
	UUID     string
	FileName string
}

var totalFramesToEncode = 0.0
var fps = ""
var totalFramesEncoded = 0.0
var currentNumber = 0.0
var currentbytes = 0.0
var percentage = 0.0
var mu sync.Mutex // mutex for progress bar

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Mp4 Converter",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
func (a *App) ConvertFile() {
	for {
		timeNow := time.Now().Format("15:04:05")
		runtime.EventsEmit(a.ctx, "app:tick", timeNow)
		time.Sleep(1 * time.Second)
	}
}
func (a *App) PercentBar() {

	sub := JobPubSub.Subscribe(array)

	for {

		fmt.Println(array[0])
		// i get to here .
		select {
		case msg := <-sub.ch:
			var fileName string

			for _, file := range uuidNFile {

				if file.UUID == msg.Topic {

					fileName = filepath.Base(file.FileName)
					fmt.Println(fileName)
				}

			}

			println("Received:", msg.Topic, msg.Msg, fileName)

			myString := strconv.Itoa(msg.Msg)
			fmt.Println(myString)
			fullString := (msg.Topic + "," + myString + "," + fileName)

			runtime.EventsEmit(a.ctx, "app:took", fullString)

		}
	}
}
func (a *App) doStuff(inputPath string, outPutPath string) {

	folderPath := inputPath        // "C:\\Users\\steve\\OneDrive\\Desktop\\test\\"
	outputFolderPath := outPutPath // "C:\\Users\\steve\\OneDrive\\Desktop\\test\\2"
	fmt.Println("do stuff")
	///////////////
	os.MkdirAll(outputFolderPath, 0755)

	files, _ := filepath.Glob(filepath.Join(folderPath, "*"))

	var UUID string

	jobIdNumber := genrateUUID()
	setUUIDS(jobIdNumber)

	for _, file := range files {
		fmt.Println(file)
		jobIdNumber := genrateUUID()
		setUUIDS(jobIdNumber)
		createFileDetails(jobIdNumber)
		createFileNameAndUUID(jobIdNumber, file)
		fmt.Println(jobIdNumber)
	}
	array = createFileDetails(jobIdNumber)
	fmt.Println(array)
	fmt.Println(array)
	fmt.Println(array)
	fmt.Println(array)
	fmt.Println(array)
	fmt.Println(array)
	fmt.Println(array)

	lensOfArray := len(array)

	for i, ar := range array {

		fmt.Println(i)
		fmt.Println(ar)
		fmt.Println("above is the printed array ")
	}

	for i, file := range files {

		wg.Add(1)
		sem <- struct{}{}
		spawnJob(file, outputFolderPath, &wg, sem, array[i])
		fmt.Println(UUID)

	}
	fmt.Println("returned")

}

func spawnJob(file string, outputFolderPath string, wg *sync.WaitGroup, sem chan struct{}, uuid string) string {

	jobIdNumber := uuid
	setUUIDS(jobIdNumber)

	JobPubSub.Pub(jobIdNumber, 10)
	//EventSub.Pub("newJob", jobIdNumber)
	fmt.Println("11111111111111111111111111111111UNID  below  ")
	fmt.Println(uuid)
	fmt.Println("11111111111111111111111111111111UNID above ")
	outputChan := make(chan string)
	go convertfiles(file, outputFolderPath, outputChan)
	for line := range outputChan {

		mu.Lock() // Lock before modifying shared variables
		fmt.Println(percentage)
		b2 := int(percentage)
		setJobProgress(jobIdNumber, b2)
		mu.Unlock()
		fmt.Println("-----------------------percentage--------------------------------------------------")
		fmt.Println(b2)
		fmt.Println("--------------------------percentage-----------------------------------------------")
		JobPubSub.Pub(jobIdNumber, b2)
		fmt.Println("Received line:", line)

		if strings.Contains(line, "fps") {
			processFPSLine(line)
		} else if strings.Contains(line, "NUMBER_OF_FRAMES:") {
			processFramesLine(line)
		}
	}

	fmt.Println(jobIdNumber)

	JobPubSub.Pub(jobIdNumber, 100)
	setJobProgress(jobIdNumber, 100)
	wg.Done()
	<-sem
	return jobIdNumber
}

func processFPSLine(line string) {
	parts := strings.Split(line, "=")
	if len(parts) > 1 {
		result := strings.TrimSpace(parts[1])
		f, err := strconv.ParseFloat(result, 64)
		if err != nil {
			fmt.Println("Error parsing FPS value:", err)
			return
		}
		mu.Lock() // Lock before modifying shared variables
		totalFramesEncoded += (currentNumber + f)
		mu.Unlock() // Unlock after modification
		workOutPercentage(totalFramesToEncode, totalFramesEncoded)
	} else {
		fmt.Println("No equals sign found in the input string.")
	}
}
func workOutPercentage(totalFramesToEncode, totalFramesEncoded float64) {
	numerator := totalFramesEncoded
	denominator := totalFramesToEncode
	percentage = (numerator / denominator) * 100
	fmt.Printf("%.2f%% of %.2f is %.2f\n", numerator, denominator, percentage)
}

func processFramesLine(line string) {
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		result := strings.TrimSpace(parts[1])
		f, err := strconv.ParseFloat(result, 64)
		if err != nil {
			fmt.Println("Error parsing frames value:", err)
			//	return
		}
		mu.Lock() // Lock before modifying shared variables
		totalFramesToEncode += (currentbytes + f)
		fmt.Println("total frames ")
		fmt.Println(totalFramesToEncode)
		mu.Unlock() // Unlock after modification
	} else {
		fmt.Println("No colon found in the input string.")
	}
}

func convertfiles(file string, outputFolderPath string, outputChan chan<- string) {
	//strings.HasSuffix(file, ".mp4") ||
	if strings.HasSuffix(file, ".mkv") ||
		strings.HasSuffix(file, ".avi") || strings.HasSuffix(file, ".mov") ||
		strings.HasSuffix(file, ".wmv") || strings.HasSuffix(file, ".mpeg") ||
		strings.HasSuffix(file, ".flv") || strings.HasSuffix(file, ".3gp") {
		// Construct the output filename (replace the extension with .mp4)
		outputFile := filepath.Join(outputFolderPath, strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))+".mp4")

		cmd := exec.Command("C:\\Users\\steve\\GolandProjects\\ConvertFilesTomp4\\ffmpeg-6.1.1-essentials_build\\bin\\ffmpeg.exe", "-i", file, "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-y", "-progress", "pipe:1", outputFile)

		// Create pipes for both stdout and stderr
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		// Read from both stdout and stderr
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				outputChan <- line
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				line := scanner.Text()
				outputChan <- line
			}
		}()

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}
	close(outputChan)
}

func createFileNameAndUUID(uuidStr string, filename string) []JobsUUIDFile {
	fmt.Println("create File Details UUID" + uuidStr)
	// Iterate over the map and collect the keys
	set := JobsUUIDFile{uuidStr, filename}

	uuidNFile = append(uuidNFile, set)

	return uuidNFile
}

func createFileDetails(uuidStr string) []string {
	fmt.Println("create File Details UUID" + uuidStr)
	// Iterate over the map and collect the keys

	keys = append(keys, uuidStr)

	return keys
}

func genrateUUID() string {

	uuidWithHyphen := uuid.New()
	fmt.Println(uuidWithHyphen)

	// Convert UUID to a string without hyphens
	uuidString := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuidString
}

func setUUIDS(jobUUID string) {

	UUIDARRAYMtx.Lock()
	UUIDARRAY = append(UUIDARRAY, jobUUID)
	UUIDARRAYMtx.Unlock()

}

func setJobProgress(jobUUID string, progress int) {

	jobIdsMtx.Lock()
	jobIds[jobUUID] = progress
	jobIdsMtx.Unlock()

}
