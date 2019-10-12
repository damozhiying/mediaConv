package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mediaConv/streamconv"
	"net/http"
)

//convStream  Handle media stream convert request
func convStream(w http.ResponseWriter, r *http.Request) {

	var streamConv streamconv.MediaConvReq
	var cntArray [2048]byte
	cntSlice := cntArray[0:]
	n, err := r.Body.Read(cntSlice)
	if err != nil && err != io.EOF {
		fmt.Printf("Read Data error:%v \n", err)
		return
	}

	defer r.Body.Close()

	json.Unmarshal(cntSlice[0:n], &streamConv)
	fmt.Println("the streamConv:", streamConv)
	err = streamconv.HandleMediaConvReq(&streamConv, w, r)
	if err != nil {
		fmt.Println("Handle fail, err:", err)
	}
}

//getAllStatus  Get all media stream transcoding information
func getAllStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get All Media Status ")
	err := streamconv.GetAllMediaConvStatus(w, r)
	if err != nil {
		fmt.Println("getStatus fail, err:", err)
	}
}

func getSingleStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get Single MediaConv Status ")

	var streamConv streamconv.MediaConvReq
	var cntArray [2048]byte
	cntSlice := cntArray[0:]
	n, err := r.Body.Read(cntSlice)
	if err != nil && err != io.EOF {
		fmt.Printf("Read Data error:%v \n", err)
		return
	}

	defer r.Body.Close()
	json.Unmarshal(cntSlice[0:n], &streamConv)
	err = streamconv.GetSingleMediaConvStatus(streamConv.OutputURL, w, r)
	if err != nil {
		fmt.Println("getStatus fail, err:", err)
	}
}

//closeConvIns  Close  media stream transcoding information
func closeConvIns(w http.ResponseWriter, r *http.Request) {
	fmt.Println("closeConvIns ")

	var streamConv streamconv.MediaConvReq
	var cntArray [2048]byte
	cntSlice := cntArray[0:]
	n, err := r.Body.Read(cntSlice)
	if err != nil && err != io.EOF {
		fmt.Printf("Read Data error:%v \n", err)
		return
	}

	defer r.Body.Close()
	json.Unmarshal(cntSlice[0:n], &streamConv)
	err = streamconv.CloseMediaConv(streamConv.OutputURL, w, r)
	if err != nil {
		fmt.Println("CloseMediaConv fail, err:", err)
	}
}

//getSysResource  Close  media stream transcoding information
func getSysResource(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/ConvStream", convStream)
	http.HandleFunc("/GetAllStatus", getAllStatus)
	http.HandleFunc("/CloseConvIns", closeConvIns)
	http.HandleFunc("/GetResourceInfo", getSysResource)
	http.HandleFunc("/GetStatus", getSingleStatus)

	err := http.ListenAndServe(":8092", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
