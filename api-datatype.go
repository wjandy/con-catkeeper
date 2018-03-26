package main

type CreateVMReq struct {
	Uuid        string
	Name        string
	ImageName   string
	Cpu         int
	Mem         int
	Disk        int
	DiskType    string
	Host        string
	Description string
}

type CreateVMRes struct {
	Name string
	Uuid string
}

type Result struct {
	Result string
}

type FetchVNCRes struct {
	RemoteConsole RemoteConsole `json:"remote_console"`
}

type RemoteConsole struct {
	Protocol string
	Url      string
}

type DeleteVolumeRes struct {
	Result string
}

type AttachVolumeReq struct {
	VolUuid     string
	AttachPoint string
	User        string
} 
