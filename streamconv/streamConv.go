package streamconv

import "os/exec"

// MediaConvReq   the struct for stream conv request
type MediaConvReq struct {
	InputURL   string `json:"input_url"`
	OutputURL  string `json:"output_url"`
	ConvParams string `json:"conv_params"`
}

// MediaConvIns   the struct for MediaConvIns
type MediaConvIns struct {
	CmdIns           *exec.Cmd    `json:"-"`
	MediaConvReqInfo MediaConvReq `json:"media-conv-req"`
	BeginTime        string       `json:"begin_time"`
	LastActiveTime   string       `json:"last_active_time"`
	StatusInfo       string       `json:"status"`
}
