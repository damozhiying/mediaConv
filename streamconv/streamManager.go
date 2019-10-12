package streamconv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	mediaSteamTypeInvalid = iota
	mediaSteamTypeRTSP    = iota
	mediaSteamTypeRTMP    = iota
	mediaSteamTypeHTTP    = iota
)

var (
	// ErrMediaConvExist   error code for MediaStream exist
	ErrMediaConvExist = errors.New("MediaStream Exist")

	// ErrMediaURLInvalid   error code for MediaStream exist
	ErrMediaURLInvalid = errors.New("MediaStream URL Invalid")

	// ErrMediaOutputURLNoExist   error code for MediaStream output url not exist
	ErrMediaOutputURLNoExist = errors.New("MediaStream OutputURL Not Exist")

	//SucMediaOutputInOpeningInfo   indicate the MediaStream Transcode in Opening
	SucMediaOutputInOpeningInfo = "MediaStream Trancoding in Opening"

	//SucMediaOutputInClosingInfo   indicate the MediaStream Transcode in Closeing
	SucMediaOutputInClosingInfo = "MediaStream Trancoding in Closing"
)

// MediaConvManager   the struct for MediaConv Controller
type MediaConvManager struct {
	mediaConvPool map[string]*MediaConvIns
	iMaxPoolSize  int
	mutex         sync.Mutex
}

// responseInfo       the struct for the response info
type responseInfo struct {
	Code int    `json:"code"`
	Info string `jsong:"info"`
}

// media conv status Infos       the struct for the media conv status info
type mediaConvStatus struct {
	Code           int            `json:"conv_count"`
	MediaConvArray []MediaConvIns `json:"media_conv_array"`
}

var mediaConvControll MediaConvManager

func init() {
	cmd := exec.Command("ffmpeg", "-version")
	errOut := bytes.NewBuffer(nil)
	cmd.Stderr = errOut
	err := cmd.Run()
	errLines := strings.Split(errOut.String(), "\n")
	for _, errLine := range errLines {
		fmt.Printf("ffmpeg output:%s", errLine)
	}

	if err != nil {
		panic(fmt.Sprintf("Run ffmpeg error:%v", err))
	}

	mediaConvControll.iMaxPoolSize = 1000
	mediaConvControll.mediaConvPool = make(map[string]*MediaConvIns)
	fmt.Println("init the Max pool size to 1000")
}

func runMediaConvCmd(cmdParams []string, mediaConvIns *MediaConvIns) {

	fmt.Printf("cmd:%v \n", cmdParams)
	cmd := exec.Command(cmdParams[0], cmdParams[1:]...)
	mediaConvIns.CmdIns = cmd

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stderr.Close()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
		return
	}

	wait := make(chan error, 1)
	go func() {
		wait <- cmd.Wait()
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	rErr := bufio.NewReader(stderr)
	iLoop := 1
	for iLoop == 1 {
		select {

		case err = <-wait:
			{
				fmt.Printf("media conv end, err: %v, \n", err)
				mediaConvIns.LastActiveTime = time.Now().String()
				mediaConvIns.StatusInfo = "Media Conv End"
				iLoop = 0
			}

		case <-ticker.C:
			{
				iInfoSize := rErr.Buffered()
				StatusInfo := ""

				line, err := rErr.ReadBytes('\r')
				for iInfoSize > 1 {
					if err == nil {
						StatusInfo = string(line)
						iInfoSize = iInfoSize - len(line)
					}
					line, err = rErr.ReadBytes('\r')
				}

				mediaConvIns.StatusInfo = StatusInfo
				mediaConvIns.LastActiveTime = time.Now().String()

				if err != nil {
					fmt.Printf("ReadLine err:%s\n", err.Error())
				}
			}
		}
	}

	mediaConvControll.mutex.Lock()
	delete(mediaConvControll.mediaConvPool, mediaConvIns.MediaConvReqInfo.OutputURL)
	mediaConvControll.mutex.Unlock()
}

func (mediaConvManager *MediaConvManager) getMediaConvIns(outputURL string) *MediaConvIns {
	mediaConvManager.mutex.Lock()
	defer mediaConvManager.mutex.Unlock()

	mediaIns, ok := mediaConvManager.mediaConvPool[outputURL]
	if ok {
		return mediaIns
	}

	return nil
}

func (mediaConvManager *MediaConvManager) buildMediaConvParams(mediaConvReq *MediaConvReq) ([]string, error) {
	inputStreamType := mediaSteamTypeInvalid
	outputStreamType := mediaSteamTypeInvalid

	switch {
	case strings.HasPrefix(mediaConvReq.InputURL, "rtsp://"):
		{
			inputStreamType = mediaSteamTypeRTSP
		}

	case strings.HasPrefix(mediaConvReq.InputURL, "rtmp://"):
		{
			inputStreamType = mediaSteamTypeRTMP
		}
	}

	switch {
	case strings.HasPrefix(mediaConvReq.OutputURL, "rtsp://"):
		{
			outputStreamType = mediaSteamTypeRTSP
		}

	case strings.HasPrefix(mediaConvReq.OutputURL, "rtmp://"):
		{
			outputStreamType = mediaSteamTypeRTMP
		}
	}

	if inputStreamType == mediaSteamTypeInvalid || outputStreamType == mediaSteamTypeInvalid {
		return nil, ErrMediaURLInvalid
	}

	//pre input
	cmdParams := []string{}
	cmdParams = append(cmdParams, "ffmpeg")
	cmdParams = append(cmdParams, "-fflags")
	cmdParams = append(cmdParams, "nobuffer")
	if inputStreamType == mediaSteamTypeRTSP {
		cmdParams = append(cmdParams, "-rtsp_flags")
		cmdParams = append(cmdParams, "prefer_tcp")
	}

	//input
	cmdParams = append(cmdParams, "-i")
	cmdParams = append(cmdParams, mediaConvReq.InputURL)

	//transcode params
	cmdParams = append(cmdParams, "-an")
	cmdParams = append(cmdParams, "-s")
	cmdParams = append(cmdParams, "1280x720")
	cmdParams = append(cmdParams, "-vcodec")
	cmdParams = append(cmdParams, "libx264")
	cmdParams = append(cmdParams, "-vb")
	cmdParams = append(cmdParams, "1000k")

	//output
	switch {
	case outputStreamType == mediaSteamTypeRTSP:
		{
			cmdParams = append(cmdParams, "-f")
			cmdParams = append(cmdParams, "rtsp")
			cmdParams = append(cmdParams, "-rtsp_flags")
			cmdParams = append(cmdParams, "prefer_tcp")
			cmdParams = append(cmdParams, mediaConvReq.OutputURL)
		}

	case outputStreamType == mediaSteamTypeRTMP:
		{
			cmdParams = append(cmdParams, "-f")
			cmdParams = append(cmdParams, "flv")
			cmdParams = append(cmdParams, mediaConvReq.OutputURL)
		}
	}

	return cmdParams, nil

}

func (mediaConvManager *MediaConvManager) doMediaConv(mediaConvReq *MediaConvReq) error {
	cmdParams, err := mediaConvManager.buildMediaConvParams(mediaConvReq)
	if err != nil {
		fmt.Printf("build media conv fail, %v \n", err)
		return err
	}

	//Add the ConvIns
	mediaConvManager.mutex.Lock()
	mediaConvIns := &MediaConvIns{nil, (*mediaConvReq), time.Now().String(), time.Now().String(), "Prepareing"}
	mediaConvManager.mediaConvPool[mediaConvReq.OutputURL] = mediaConvIns
	mediaConvManager.mutex.Unlock()

	//Run Media Conv
	go runMediaConvCmd(cmdParams, mediaConvIns)
	return nil
}

// HandleMediaConvReq  handle MediaConvReq
func HandleMediaConvReq(mediaConvReq *MediaConvReq, w http.ResponseWriter, r *http.Request) error {
	mediaConvIns := mediaConvControll.getMediaConvIns(mediaConvReq.OutputURL)
	if mediaConvIns != nil {
		return ErrMediaConvExist
	}

	err := mediaConvControll.doMediaConv(mediaConvReq)
	if err != nil {
		errResp := responseInfo{http.StatusBadRequest, mediaConvReq.InputURL + "  " + ErrMediaConvExist.Error()}
		jsonCnt, _ := json.Marshal(&errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonCnt)
	} else {
		okResp := responseInfo{http.StatusOK, mediaConvReq.InputURL + "  " + SucMediaOutputInOpeningInfo}
		w.WriteHeader(http.StatusOK)
		jsonCnt, _ := json.Marshal(&okResp)
		w.Write(jsonCnt)
	}
	return err
}

// CloseMediaConv  handle Close MediaConv
func CloseMediaConv(outputURL string, w http.ResponseWriter, r *http.Request) error {
	mediaConvIns := mediaConvControll.getMediaConvIns(outputURL)
	if mediaConvIns == nil {
		errResp := responseInfo{http.StatusBadRequest, outputURL + "  " + ErrMediaOutputURLNoExist.Error()}
		jsonCnt, _ := json.Marshal(&errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonCnt)
		return ErrMediaOutputURLNoExist
	}

	if mediaConvIns.CmdIns != nil {
		mediaConvIns.CmdIns.Process.Kill()
		resp := responseInfo{http.StatusOK, outputURL + "  " + SucMediaOutputInClosingInfo}
		jsonCnt, _ := json.Marshal(&resp)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonCnt)
	}

	return nil
}

// GetAllMediaConvStatus  handle All MediaConv Status
func GetAllMediaConvStatus(w http.ResponseWriter, r *http.Request) error {

	mediaconvArray := [256]MediaConvIns{}
	arrayIdx := 0

	mediaConvControll.mutex.Lock()
	for _, mediaIns := range mediaConvControll.mediaConvPool {
		if arrayIdx < 256 {
			mediaconvArray[arrayIdx] = *mediaIns
			arrayIdx++
		} else {
			break
		}
	}
	mediaConvControll.mutex.Unlock()

	respInfo := mediaConvStatus{arrayIdx, mediaconvArray[0:arrayIdx]}
	w.WriteHeader(http.StatusOK)
	jsonCnt, _ := json.Marshal(&respInfo)
	w.Write(jsonCnt)
	return nil
}

// GetSingleMediaConvStatus  handle single media conv status
func GetSingleMediaConvStatus(outputURL string, w http.ResponseWriter, r *http.Request) error {
	mediaconvArray := [1]MediaConvIns{}

	mediaConvIns := mediaConvControll.getMediaConvIns(outputURL)
	if mediaConvIns == nil {
		errResp := responseInfo{http.StatusBadRequest, outputURL + "  " + ErrMediaOutputURLNoExist.Error()}
		jsonCnt, _ := json.Marshal(&errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonCnt)
		return ErrMediaOutputURLNoExist
	}
	mediaconvArray[0] = *mediaConvIns
	respInfo := mediaConvStatus{1, mediaconvArray[0:]}
	w.WriteHeader(http.StatusOK)
	jsonCnt, _ := json.Marshal(&respInfo)
	w.Write(jsonCnt)
	return nil
}
