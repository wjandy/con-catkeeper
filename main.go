package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	//"encoding/xml"
	"fmt"
	"strconv"  
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	utils "github.com/converge/catkeeper/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	libvirt "github.com/libvirt/libvirt-go"
	"golang.org/x/net/websocket"
	"io/ioutil"
	//"log"
	"net"
	"net/http"
	_ "net/url"
	"strings"
	"text/template"
	"time"
	//"strconv"
	//"sync/atomic"
	//"time"
)

var User string = "test"

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(func(res http.ResponseWriter, req *http.Request){
			fmt.Println("GetUser")
			Token := req.Header.Get("Token")
			fmt.Println("Token",Token)
			if Token != ""{
				GetUser(Token)
			}		
	})
	//m.Use(martini.Static("web"))
	RegistAPI(m)
	utils.SetupConfig()
	//utils.GetLog().Infoln("server started")
	utils.Logger.Infoln("utils.Config.DataDriver")
	//fmt.Println("ResgisterEvents")
	go ResgisterEvents()
	//fmt.Println("ResgisterEvents over")
	err := http.ListenAndServe(":8081", m)
	utils.Logger.Infoln(err)
	
}

func GetDB() (*sql.DB, error) {
	/* init database */
	db, err := sql.Open(utils.Config.DataDriver, utils.Config.DataSource)
	//db.SetMaxOpenConns(200)
	//db.SetMaxIdleConns(10)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var vncProxy = websocket.Server{Handler: proxyHandler,
	Handshake: func(ws *websocket.Config, req *http.Request) error {
		fmt.Println("vncProxy = websocket.Server")
		ws.Protocol = []string{"base64"}
		return nil
	}}

func RegistAPI(m *martini.ClassicMartini) {
	//TODO
	//listvm
	//
	// m.get（"/listvm",model.go:gasdfasfd）；
	//vnc related
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/consoletest", FetchVNCT)
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/console", FetchVNC)
	m.Get("/websockify", vncProxy.ServeHTTP)
	//m.Get("/websockify", VncProxy)
	//pm related
	m.Get("/catkeeper/v1/pms/(?P<ipaddress>[0-9.]*)", InfoPM)
	m.Get("/catkeeper/v1/pms", ListPM)
	//disk related
	m.Get("/catkeeper/v1/disks/(?P<ipaddress>[0-9.]*)", InfoPMDisk)
	m.Get("/catkeeper/v1/disks", ListPMDisk)
	//net related
	m.Get("/catkeeper/v1/nets/(?P<ipaddress>[0-9.]*)", GetNodeNetInfo)
	m.Get("/catkeeper/v1/net", GetNodeNetInfos)
	//vm related
	m.Get("/catkeeper/v1/servers", ListVM)
	m.Get("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)", InfoVM)
	m.Post("/catkeeper/v1/servers", CreateVM)
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/attach", AttachVolume)
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/detach", DetachVolume)
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/delete", DeleteVM)
	m.Post("/catkeeper/v1/servers/(?P<uuid>[a-zA-Z0-9-]*)/(?P<action>[a-z]*)", ManageVM)
	//disk related
	m.Get("/catkeeper/v1/volumes", ListVolume)
	m.Get("/catkeeper/v1/volumes/(?P<uuid>[a-zA-Z0-9-]*)", InfoVolume)
	m.Post("/catkeeper/v1/volumes", CreateVolume)
	m.Post("/catkeeper/v1/volumes/(?P<uuid>[a-zA-Z0-9-]*)/(?P<action>[a-z]*)", ManageVolume)
	//sold related
	m.Get("/catkeeper/v1/oversell_rate", InfoSold)
    //log related
    	m.Get("/catkeeper/v1/logs", ListLog)
	m.Get("/catkeeper/v1/logs/(?P<uuid>[a-zA-Z0-9-]*)", InfoLog)
	//net speep related
	m.Get("/catkeeper/v1/netspeeds/(?P<ipaddress>[0-9.]*)", InfoNetSpeed)
	m.Get("/catkeeper/v1/netspeeds", ListNetSpeed)
}


func ListNetSpeed(r render.Render){
	utils.Logger.Infoln("get NetSpeed info")
	NetInfos, err :=GetNetSpeedInfos()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	fmt.Println("entry NetSpeedinfo",NetInfos)
	utils.Logger.Infoln(NetInfos)
	r.JSON(202, NetInfos)
	
}

func InfoNetSpeed(r render.Render, params martini.Params){
	utils.Logger.Infoln("get NetSpeed info")
	fmt.Println("entry NetSpeedinfo")
	ipaddress := params["ipaddress"]
	NetInfos, err :=GetNetSpeedInfo(ipaddress)
	fmt.Println("entry NetSpeedinfo")
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	utils.Logger.Infoln(NetInfos)
	r.JSON(200, NetInfos)
	
}


func InfoSold(r render.Render){
	fmt.Println("start InfoSold")
	type SoldCpu struct {
	    PmCpu         int
	    SoldCpu       int
	    VMCPu         int
	    Proportion    int	    
    }
	SCpu := &SoldCpu{}
	utils.Logger.Infoln("get sold  proportion info")
	db, err := GetDB()
	if err != nil {
		utils.Logger.Infoln("get sold  proportion info:",err)
		ReturnError(r, InternalError)
		return
	}
	defer db.Close()
	fmt.Println("get sold  proportion info")
	err = db.QueryRow("select sum(cpu) from physicalmachine").Scan(&SCpu.PmCpu)
	if err != nil {
		fmt.Println("get sold  proportion info",err)
		ReturnError(r, InternalError)
		return
	}
	SCpu.Proportion, _ = strconv.Atoi(utils.Config.Proportion)
	SCpu.SoldCpu = (SCpu.Proportion)*(SCpu.PmCpu)
	err = db.QueryRow("select sum(cpu) from virtualmachine").Scan(&SCpu.VMCPu)
	if err != nil {
		fmt.Println("get sold  proportion info",err)
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, SCpu)
}
func ListLog(r render.Render){
	utils.Logger.Infoln("get vm log info")
	type VMLog struct{
		Uuid          string
		VMUuid        string
		Action        string
		OperateTime   string
		Status        string
		Desc          string
		User          string
		VMName        string
	}
	
	type VMLogs struct{
		VMLogs []VMLog
	}
	db, err := GetDB()
	if err != nil {
		utils.Logger.Infoln("get vm log info:",err)
		ReturnError(r, InternalError)
		return
	}
	defer db.Close()
	Vmlog := VMLog{}
	VmLogs := VMLogs{}
	rows, err := db.Query("select lg.uuid,lg.vmuuid,lg.action,lg.status,lg.operatetime,lg.describ,lg.user,vl.name from virtualmachinelog lg, virtualmachine vl where lg.vmuuid=vl.uuid")
	if err != nil {
		fmt.Println("get vm log info:",err)
		ReturnError(r, InternalError)
		return
	}
	for rows.Next() {
		rows.Scan(
			&Vmlog.Uuid,
			&Vmlog.VMUuid,
			&Vmlog.Action,
			&Vmlog.Status,
			&Vmlog.OperateTime,
			&Vmlog.Desc,
			&Vmlog.User,
			&Vmlog.VMName,	
		)
		VmLogs.VMLogs =append(VmLogs.VMLogs,Vmlog)
	}
	utils.Logger.Infoln(VmLogs)
    r.JSON(200, VmLogs)
	
}

func ManageVMLog(VMUuid string,Action string,Status string,Desc string,User string){
	 fmt.Println("start ManageVMLog")
	 Uuid := uuid.New().String()
	 OperateTime   := NowTime()
	 
	 db, err := GetDB()
	 if err != nil {
	 	fmt.Println("up vm log info:",err)
		utils.Logger.Infoln("up vm log info:",err)
		//ReturnError(r, InternalError)
		return 
	 }
	 defer db.Close()
	 sqltext := fmt.Sprintf("insert into virtualmachinelog (uuid,vmuuid,action,status,operatetime,describ,user) values ('%s','%s','%s','%s','%s','%s','%s')",Uuid,VMUuid,Action,Status,OperateTime,strings.Replace(Desc, "'", "", -1),User)
	 fmt.Println(sqltext)
	 _, err = db.Exec(sqltext)
     if err != nil {
     	fmt.Println("up vm log info:",err)
		utils.Logger.Infoln("up vm log info:",err)
		//ReturnError(r, InternalError)
		return 
	 }
}


func InfoLog(r render.Render,params martini.Params){
	VMuuid := params["uuid"]
	fmt.Println("get vm log info")
	utils.Logger.Infoln("get vm log info")
	type VMLog struct{
		Uuid          string
		VMUuid        string
		VMName        string
		Action        string
		OperateTime   string
		Status        string
		Desc          string
		User          string
	}
	
	type VMLogs struct{
		VMLogs []VMLog
	}
	db, err := GetDB()
	if err != nil {
		utils.Logger.Infoln("get vm log info:",err)
		ReturnError(r, InternalError)
		return 
	}
	defer db.Close()
	Vmlog := VMLog{}
	VmLogs := &VMLogs{}
	sqlText := fmt.Sprintf("select lg.uuid,lg.vmuuid,lg.action,lg.status,lg.operatetime,lg.describ,lg.user,vl.name from virtualmachinelog lg, virtualmachine vl where lg.vmuuid=vl.uuid and  lg.VMUuid  = '%s'", VMuuid)
	fmt.Println(sqlText)
	rows, err := db.Query(sqlText)
	if err != nil {
		fmt.Println("get vm log info:",err)
		ReturnError(r, InternalError)
		return
	}
	for rows.Next() {
		rows.Scan(
			&Vmlog.Uuid,
			&Vmlog.VMUuid,
			&Vmlog.Action,
			&Vmlog.Status,
			&Vmlog.OperateTime,
			&Vmlog.Desc,
			&Vmlog.User,
			&Vmlog.VMName,
		)
		VmLogs.VMLogs =append(VmLogs.VMLogs,Vmlog)
	}
	utils.Logger.Infoln(VmLogs)
    r.JSON(200, VmLogs)
	
}

func ListPMDisk(r render.Render){
	utils.Logger.Infoln("get pm disk info")
	Diskinfos, err :=GetDiskInfos()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	utils.Logger.Infoln(Diskinfos)
	r.JSON(202, Diskinfos)
	
}

func InfoPMDisk(r render.Render, params martini.Params){
	utils.Logger.Infoln("get pm disk info")
	fmt.Println("entry InfoPMDisk")
	ipaddress := params["ipaddress"]
	Diskinfos, err :=GetDiskInfo(ipaddress)
	fmt.Println("entry InfoPMDisk")
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	utils.Logger.Infoln(Diskinfos)
	r.JSON(200, Diskinfos)
	
}

func ListPM(r render.Render) {
	fmt.Println("enter listpm")
	utils.Logger.Infoln("enter listpm")
	//var Name, IpAddress string
	type HostInfo struct {
		/*Name      string
		IpAddress string*/
		IpAddress   string
	    Name        string
	    Description string
	    Cputype     string
        Cpu         int
        Cpufreq     int
        Cpusocket   int
        Cpucore     int
        Cputhread   int
        Numa        int
        Memory      int
	    Active      bool
	}

	type Hosts struct {
		Hosts []HostInfo
	}

	//var hosts []HostInfo
	//get pm info with libvirt before get pm info
	fmt.Println("GetNodesInfos")
	utils.Logger.Infoln("GetNodesInfos")
	
	GetNodesInfos()
	
	h := HostInfo{}
	var hosts Hosts
	db, err := GetDB()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	defer db.Close()
	fmt.Println("select GetNodesInfos")
	rows, err := db.Query("select IpAddress,Description,Name,cputype,cpu,cpufreq,cpusocket,cpukernel,cputhread,numa,memory from physicalmachine")
	if err != nil {
		fmt.Println("select GetNodesInfos",err)
		ReturnError(r, InternalError)
		return
	}
	for rows.Next() {
		rows.Scan(
		  &h.IpAddress, 
		  &h.Description,
		  &h.Name,
		  &h.Cputype,
          &h.Cpu,
          &h.Cpufreq,
          &h.Cpusocket,
          &h.Cpucore,
          &h.Cputhread,
          &h.Numa,
          &h.Memory,
	      //&h.Active,
	   )
		h.Active = true
		fmt.Println("end  ListPM", h)
		hosts.Hosts = append(hosts.Hosts, h)
	}
	utils.Logger.Infoln(hosts)
	r.JSON(200, hosts)
}

func InfoPM(r render.Render, params martini.Params) {
	fmt.Println("entry  InfoPM")
	ipaddress := params["ipaddress"]
	//get  pm info with libvirt before get  pm info
	GetNodeInfo(ipaddress)
	p, err := NewPhysicalMachine(ipaddress)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, p)
}

func ListVM(r render.Render) {
	//var uuid string
	type VMInfo struct {
		/*Uuid string*/
		Uuid           string
        Description    string 
        Hostipaddress  string
        Status         string
        Attachments    string
        Cpu            int 
        Mem            int
        Disk           int
        Imagename      string
        Name           string
        Sysdisk        string
        User           string
        CreateTime     string
        UpdateTime     string
	}
	type VMs struct {
		VMs []VMInfo
	}
	db, err := GetDB()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	defer db.Close()
    v := VMInfo{}
	rows, err := db.Query("select uuid,cpu,mem,disk,hostipaddress,description,status,attachments,imagename,name,sysdisk,user,createtime,updatetime from virtualmachine")
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	var vms VMs
	for rows.Next() {
	    rows.Scan(
		  &v.Uuid,
		  &v.Cpu,
		  &v.Mem,
		  &v.Disk,
		  &v.Hostipaddress,
		  &v.Description,
		  &v.Status,
		  &v.Attachments,
		  &v.Imagename,
		  &v.Name,
		  &v.Sysdisk,
		  &v.User,
		  &v.CreateTime,
		  &v.UpdateTime,
	    )
		vms.VMs = append(vms.VMs, v)
	}
	r.JSON(200, vms)
}

func InfoVM(r render.Render, params martini.Params) {
	uuid := params["uuid"]
	v, err := GetVirtualMachine(uuid)
	if err != nil {
		fmt.Println("InfoVM:",err)
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, v)
}

func ListVolume(r render.Render) {
	//var uuid string
	type VolInfo struct {
		Uuid        string
	    Size        uint64
	    Description string
	    Name        string
	    VolumeType  string
	    Status      string
	    CreateAt    string
	    UpdateAt    string
	    Attachments string
	    HostIp      string
	    DataType    string
	    User        string
	}
	type Vols struct {
		Vols []VolInfo
	}
	db, err := GetDB()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from volume ")
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	 vol  := VolInfo{}
	 vols := Vols{}
	for rows.Next() {
		rows.Scan(
			&vol.Uuid,
		    &vol.Size,
		    &vol.Description,
		    &vol.Name,
		    &vol.VolumeType,
		    &vol.Status,
		    &vol.CreateAt,
		    &vol.UpdateAt,
		    &vol.Attachments,
		    &vol.HostIp,
		    &vol.DataType,
		    &vol.User,
         )
		vols.Vols = append(vols.Vols, vol)
	}
	r.JSON(200, vols)
}
func InfoVolume(r render.Render, params martini.Params) {
	uuid := params["uuid"]
	v, err := GetVolume(uuid)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, v)
}
func CreateVM(r render.Render, req *http.Request) {
	Action := "create"
	Status := "success"
	Desc   := "create successfully"
	
	v := &VirtualMachine{}
	uuuid := uuid.New().String()
	v.Uuid = uuuid
	reqbytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Status ="fail"
		Desc = err.Error()
		//ManageVMLog(v.Uuid,Action, Status, Desc,)
		ReturnError(r, InternalError)
		return
	}
	err = json.Unmarshal(reqbytes, v)
	if err != nil {
		Status ="fail"
		Desc = err.Error()
		fmt.Println(err)
		ManageVMLog(v.Uuid,Action, Status, Desc, v.User)
		ReturnError(r, InternalError)
		return
	}
	fmt.Println("v.User = User", User)
	v.User = User
	v.CreateTime = NowTime()
	v.UpdateTime = NowTime()
	p, err := NewPhysicalMachine(v.HostIpAddress)
	defer p.conn. Close()
	if err != nil {
		Status ="fail"
		Desc = err.Error()
		fmt.Println("err occur where newphyiscalmachine", err)
		ManageVMLog(v.Uuid,Action, Status, Desc, v.User)
		ReturnError(r, InternalError)
		return
	}
	err = p.CreateVM(v)
	if err != nil {
		Status ="fail"
		Desc = err.Error()
		fmt.Println(err)
		ManageVMLog(v.Uuid,Action, Status, Desc, v.User)
		ReturnError(r, InternalError)
		return
	}
	fmt.Println("ManageVMLog")
	ManageVMLog(v.Uuid,Action, Status, Desc, v.User)
	r.JSON(200, v)
}

func ManageVM(r render.Render, params martini.Params, req *http.Request) {
	var err error
	vmuuid := params["uuid"]
	action := params["action"]
	user   := User
	utils.Logger.Infoln("ManageVM")
	//Action := "Manage"
	Status := "success"
	Desc   := fmt.Sprintf("%s successfully", action)
	v, err := GetVirtualMachine(vmuuid)
	if err != nil {
		Status = "fail"
		Desc  = err.Error()
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	d, err := GetDomainFromUuid(vmuuid)
	if err != nil {
		Status = "fail"
		Desc  = err.Error()
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	

	
	switch action {
	case "shutdown":
		//err = d.Destroy()
		//Action = "shutdown"
		utils.Logger.Infoln("Shutdown")
		err = d.Shutdown()
		utils.Logger.Infoln(err)
		v.Status = "closing"
	case "start":
	    //Action = "start"
	    utils.Logger.Infoln("start")
		err = d.Create()
		utils.Logger.Infoln(err)
		v.Status = "starting"
	case "reboot":
	     //Action = "reboot"
	     utils.Logger.Infoln("reboot")
		err = d.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
		utils.Logger.Infoln(err)
		v.Status = "rebooting"
	case "suspend":
	    //Action = "suspend"
	    utils.Logger.Infoln("suspend")
		err = d.Suspend()
		utils.Logger.Infoln(err)
		v.Status = "suspending"
	case "resume":
	    //Action = "resume"
	    utils.Logger.Infoln("resume")
		err = d.Resume()
		utils.Logger.Infoln(err)
		v.Status = "resuming"
	}
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println(err)
		ManageVMLog(v.Uuid,action, Status, Desc, user)
		ReturnError(r, InternalError)
		return
	}
	// 不需要了，回掉函数已经解决问题
//	err = v.SyncDb()
//	if err != nil {
//		Status = "fail"
//		Desc   =  err.Error()
//		ManageVMLog(v.Uuid,action, Status, Desc, user)
//		fmt.Println(err)
//		ReturnError(r, InternalError)
//		return
//	}

	result := Result{}
	result.Result = fmt.Sprintf("%s successfully", action)
	ManageVMLog(v.Uuid,action, Status, Desc, user)
	r.JSON(200, result)
}

func DeleteVM(r render.Render, params martini.Params) {
	var err error
	Action := "Delete"
	Status := "success"
	Desc   := "delete successfully"
	vmuuid := params["uuid"]
	//user   := params["user"]
	
	d, err := GetDomainFromUuid(vmuuid)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println(err)
		ManageVMLog(vmuuid,Action, Status, Desc, User)
		ReturnError(r, InternalError)
		return
	}
	
	err = d.Undefine()
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
	    ManageVMLog(vmuuid,Action, Status, Desc, User)
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	
	sqltext := fmt.Sprintf("delete from virtualmachine where uuid='%s'", vmuuid)
	db, _ := GetDB()
	_, err = db.Exec(sqltext)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
	    ManageVMLog(vmuuid,Action, Status, Desc, User)
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	result := Result{}
	result.Result = "delete successfully"
	ManageVMLog(vmuuid,Action, Status, Desc, User)
	r.JSON(200, result)
}

func AttachVolume(r render.Render, param martini.Params, req *http.Request) {
	
	Action := "AttachVolume"
	Status := "success"
	Desc   := "no err"
	vmuuid := param["uuid"]
	reqbytes, err := ioutil.ReadAll(req.Body)
	a := &AttachVolumeReq{}
	err = json.Unmarshal(reqbytes, a)
	//a.User = User
	d, err := GetDomainFromUuid(vmuuid)
	fmt.Println("GetDomainFromUuid",a)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
	    ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	v, err := GetVirtualMachine(vmuuid)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when GetVirtualMachine", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	vol, err := GetVolume(a.VolUuid)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when GetVolume", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	xmlData, _, err := GetDomainINAXMLFromUuid(vmuuid)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	vol.Attachments = vmuuid
	vol.Status = "in-use"
	fmt.Println("xmlData",xmlData)
	xmlObj, err := ParseDomainXML(xmlData)
	fmt.Println("xmlObj",xmlObj);
	
	//Disks :=  [3]string{"vdc","vdb","vdd"}
	          //xmlObj.Devices.Graphics.VNCPort
	Disks :=  xmlObj.Devices.Disks
	DiskPoint := make([]string, 10, 10)
    var i int
	for _, disk := range Disks {
		fmt.Println("disk.Target.Dev",disk.Target.Dev)
		DiskPoint[i] = disk.Target.Dev
		i = i+1
	}
	//fmt.Println("xmlObj",xmlObj);
	fmt.Println("Disks",DiskPoint)
		//parse disk to get min disk point 
	Points := [10]string{"vdb","vdc","vdd","vde","vdf","vdg","vdh","vdi","vdj","vdk"}	
	var AttachPoint string
	if len(Disks) > 0 {
	    var flag  bool = false 
        for _, point := range Points{
		   AttachPoint = point
		   flag = false  
	       for _, disk := range DiskPoint {
		      if point == disk {
		    	       flag = true
		    	       AttachPoint = ""
		    	       break
		      } 
	      }
	      if !flag {
	    	     break
	      }
	  }
	}else{
		AttachPoint = "vdb"
	}
	fmt.Println("AttachPoint",AttachPoint)
	a.AttachPoint =  AttachPoint
	if v.Attachments == "" {
		v.Attachments = a.VolUuid + ":" + a.AttachPoint
	} else {
		v.Attachments = v.Attachments + "," + a.VolUuid + ":" + a.AttachPoint
	}
	var t *template.Template
	if vol.VolumeType == "rbd" {
		t, err = template.ParseFiles("templates/attachvm.tmpl")
	} else {
		t, err = template.ParseFiles("templates/attachlocal.tmpl")
	}
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when parsefile", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User) 
		ReturnError(r, InternalError)
		return
	}
	fmt.Println(t)
	var xml *bytes.Buffer = new(bytes.Buffer)
	err = t.Execute(xml, a)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when execute template", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		return
	}
	
	xmlStr := xml.String()
	fmt.Println(xmlStr)
	err = d.AttachDeviceFlags(xmlStr, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when AttachDeviceFlags device",err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	err = vol.SyncVolume()
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when SyncVolume", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	err = v.SyncDb()
	if err != nil {
		Status = "fail"
	    Desc   =  err.Error()
		fmt.Println("err occur when SyncDb", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	ManageVMLog(vmuuid,Action, Status, Desc, a.User)
	r.JSON(200, a)
}

func DetachVolume(r render.Render, param martini.Params, req *http.Request) {
	Action := "DetachVolume"
	Status := "success"
	Desc   := "no err"
	vmuuid := param["uuid"]
	reqbytes, err := ioutil.ReadAll(req.Body)
	a := &AttachVolumeReq{}
	err = json.Unmarshal(reqbytes, a)
	a.User = User
	v, err := GetVirtualMachine(vmuuid)
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when GetVirtualMachine", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	vol, err := GetVolume(a.VolUuid)
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when GetVolume", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}

	d, err := GetDomainFromUuid(vmuuid)
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when GetDomainFromUuid", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	vol.Attachments = ""
	vol.Status = "available"
	v.Attachments = strings.Trim(v.Attachments, ","+a.VolUuid+":"+a.AttachPoint)
	//cutstr := ","+a.VolUuid+":"+a.AttachPoint
	//v.Attachments = v.Attachments[:(len(v.Attachments)-len(cutstr))]
	var t *template.Template
	if vol.VolumeType == "rbd" {
		t, err = template.ParseFiles("templates/attachvm.tmpl")
	} else {
		t, err = template.ParseFiles("templates/attachlocal.tmpl")
	}
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when parsefile", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	fmt.Println(t)
	var xml *bytes.Buffer = new(bytes.Buffer)
	err = t.Execute(xml, a)
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when execute template", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		return
	}
	
	xmlStr := xml.String()
	fmt.Println(xmlStr)
	err = d.DetachDeviceFlags(xmlStr, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when DetachDeviceFlags", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	err = vol.SyncVolume()
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when SyncVolume", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	err = v.SyncDb()
	if err != nil {
		Status = "fail"
	    Desc   = err.Error()
		fmt.Println("err occur when SyncDb", err)
		ManageVMLog(vmuuid,Action, Status, Desc, a.User)
		ReturnError(r, InternalError)
		return
	}
	
	ManageVMLog(vmuuid,Action, Status, Desc, a.User)
	r.JSON(200, a)
}

func FetchVNCT(r render.Render, params martini.Params, req *http.Request) {
	fmt.Println("enter fetchvnct")
	var err error
	vmuuid := params["uuid"]
	xmlData, hostip, err := GetDomainXMLFromUuid(vmuuid)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	xmlObj, err := ParseDomainXML(xmlData)
	fmt.Println(*xmlObj, err)
	url := fmt.Sprintf("http://10.72.84.145:8081/vnc_auto.html?path=websockify?ip=%s:%s", hostip, xmlObj.Devices.Graphics.VNCPort)
	fetchRes := FetchVNCRes{RemoteConsole: RemoteConsole{Protocol: "vnc", Url: url}}
	r.JSON(200, fetchRes)
}

func FetchVNC(r render.Render, params martini.Params, req *http.Request) {
	var err error
	vmuuid := params["uuid"]
	xmlData, hostip, err := GetDomainXMLFromUuid(vmuuid)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	fmt.Println("GetDomainXMLFromUuid(vmuuid)   ===",xmlData);
	xmlObj, err := ParseDomainXML(xmlData)
	type VNCRes struct {
		IP   string
		Port string
	}
	fetchRes := VNCRes{hostip, xmlObj.Devices.Graphics.VNCPort}
	r.JSON(200, fetchRes)
}

func CreateVolume(r render.Render, req *http.Request) {
	reqbytes, err := ioutil.ReadAll(req.Body)
	v := &Volume{}
	err = json.Unmarshal(reqbytes, v)
	uuuid := uuid.New().String()
	v.Uuid = uuuid
	v.User = User
	v.DataType = "data"
	if err != nil {
		ReturnError(r, InternalError)
	}
	if v.VolumeType == "local" {
		err = CreateLocalImage(v.Uuid, v.Size, v.HostIp)
	} else {
		v.VolumeType = "rbd"
		_, err = CreateImage(v.Uuid, v.Size)
	}
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	v.Status = "available"
	now := NowTime()
	v.CreateAt = now
	v.UpdateAt = now
	err = v.CreateVolDB()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, v)
}

func ManageVolume(r render.Render, params martini.Params, req *http.Request) {
	voluuid := params["uuid"]
	action := params["action"]
	switch action {
	case "update":
		UpdateVolume(r, voluuid, req)
	case "delete":
		DeleteVolume(r, voluuid)
	case "resize":
		ResizeVolume(r, voluuid, req)
	}
}

func UpdateVolume(r render.Render, voluuid string, req *http.Request) {
	jsonbytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	v, err := GetVolume(voluuid)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	err = json.Unmarshal(jsonbytes, v)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	v.UpdateAt = NowTime()
	err = v.SyncVolume()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, v)
}

func ResizeVolume(r render.Render, voluuid string, req *http.Request) {
	fmt.Println("enter resizevolume")
	var err error
	jsonbytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	v, err := GetVolume(voluuid)
	if err != nil {
		fmt.Println("err when getvolume", err)
		ReturnError(r, InternalError)
		return
	}
	err = json.Unmarshal(jsonbytes, v)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	fmt.Println("v.Volume", v.VolumeType)
	if v.VolumeType == "local" {
		err = ResizeLocalImage(v.Uuid, v.Size, v.HostIp)
	} else {
		err = ResizeImage(v.Uuid, v.Size)
	}
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	err = v.SyncVolume()
	if err != nil {
		fmt.Println("err when syncvolume", err)
		ReturnError(r, InternalError)
		return
	}
	r.JSON(200, v)
}

func DeleteVolume(r render.Render, voluuid string) {
	fmt.Println("enter deletevolume")
	var err error
	v, err := GetVolume(voluuid)
	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	if v.VolumeType == "local" {
		err = RemoveLocalImage(v.Uuid, v.HostIp)
	} else {
		err = RemoveImage(v.Uuid)
	}

	if err != nil {
		fmt.Println(err)
		ReturnError(r, InternalError)
		return
	}
	db, err := GetDB()
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	sqltext := fmt.Sprintf("delete from volume where uuid='%s'", v.Uuid)
	_, err = db.Exec(sqltext)
	if err != nil {
		ReturnError(r, InternalError)
		return
	}
	d := DeleteVolumeRes{"delete volume successfully"}
	r.JSON(200, d)
	return
}

func proxyHandler(ws *websocket.Conn) {
	defer ws.Close()
	r := ws.Request()
	values := r.URL.Query()
	ip, hasIp := values["ip"]
	fmt.Println(ip, hasIp)

	if hasIp == false {
		//log.Println("faile to parse vnc address")
		return
	}

	vc, err := net.Dial("tcp", ip[0])
	if err != nil {
		return
	}
	defer vc.Close()
	done := make(chan bool)

	go func() {
		sbuf := make([]byte, 32*1024)
		dbuf := make([]byte, 32*1024)
		for {
			n, e := ws.Read(sbuf)
			if e != nil {
				done <- true
				return
			}
			n, e = base64.StdEncoding.Decode(dbuf, sbuf[0:n])
			if e != nil {
				done <- true
				return
			}
			n, e = vc.Write(dbuf[0:n])
			if e != nil {
				done <- true
				return
			}
		}
	}()
	go func() {
		sbuf := make([]byte, 32*1024)
		dbuf := make([]byte, 64*1024)
		for {
			n, e := vc.Read(sbuf)
			if e != nil {
				done <- true
				return
			}
			base64.StdEncoding.Encode(dbuf, sbuf[0:n])
			n = ((n + 2) / 3) * 4
			ws.Write(dbuf[0:n])
			if e != nil {
				done <- true
				return
			}
		}
	}()
	select {
	case <-done:
		break
	}
}


func ResgisterEvents() {
	
	fmt.Println("ResgisterEvents start");
	p := &PhysicalMachine{}
	err := libvirt.EventRegisterDefaultImpl()
	if err != nil {
		fmt.Println("update virtualmachine uuid :",err)
		return
	}
	
	callback := func(c *libvirt.Connect, d *libvirt.Domain, event *libvirt.DomainEventLifecycle) {
		//callback := func(c *Connect, d *Domain, event *DomainEventLifecycle) {
		uuid, _ := d.GetUUIDString()
	    status, _, _ := d.GetState()
	    
        utils.Logger.Infoln("callback=========================================================================================================",status)
	    fmt.Println("update virtualmachine uuid :",uuid)
	    //fmt.Println("update virtualmachine uuid string :",string(uuid))
	    //fmt.Println("update virtualmachine uuid string :",string(uuid[:]))
	    fmt.Println("update virtualmachine status :",status)
	    stat :="unknow"
	    switch  status {
	    	    case 0 :
	    	        stat = "no state"
	    	    case 1 :
	    	        stat = "running"
	    	    case 2 :
	    	        stat = "blocked"
	    	    case 3 :
	    	        stat = "paused"
	    	    case 4 :
	    	        stat = "shut down"
	    	    case 5 :
	    	        stat = "shut off"
	    	    case 6 :
	    	        stat = "crashed"     
	    	    case 7 :
	    	        stat = "suspended" 
	    	    case 8 :
	    	        stat = "last"                            
	    	    
	    }
	    sqltext := fmt.Sprintf("update virtualmachine set status='%s' where uuid='%s';",  stat, uuid)
	    fmt.Println("update virtualmachine status :",sqltext)
	    db, _ := GetDB()
	    defer db.Close()
	    db.Exec(sqltext)
	  }
	
	  
	  //libvirt.EventRegisterImpl()
	  db, err := GetDB()
	  defer db.Close()
	  if err != nil {
		 fmt.Println("db",err)
		 return  
  	  }
	
	  rows, err := db.Query("select IpAddress from physicalmachine")
	  if err != nil {
	  	 fmt.Println("physicalmachine",err)
		 return 
	  }
	  
	  for rows.Next() {
		 rows.Scan(&p.IpAddress)
		 conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system",p.IpAddress))
		 if err != nil {
		 	fmt.Println("NewConnect",err)
		 	return 
		 } 
		//evrery node info		
		//every node register event
		//go func(conn *libvirt.Connect) {
		fmt.Println("EventRegisterDefaultImpl",conn);
		fmt.Println("Connect ResgisterEvents");
		callbackId, err := conn.DomainEventLifecycleRegister(nil, callback)
		fmt.Println("DomainEventLifecycleRegister :",callbackId)
		fmt.Println("DomainEventLifecycleRegister :",err)
	    if err != nil {
	       utils.Logger.Infoln(err)
	       fmt.Println("DomainEventLifecycleRegister :",err)
	    }
		//}(conn)
	 }
  
    for {
		err := libvirt.EventRunDefaultImpl()
		if err != nil {
			fmt.Println("Run failed")
			break
		}
    }
}
//获取多个node的节点信息
func GetNodesInfos() error{
	
	fmt.Println("entry GetNodesInfos")
	p := &PhysicalMachine{}
    db, err := GetDB()
	defer db.Close()
	if err != nil {
		return  err
	}
	
	rows, err := db.Query("select IpAddress from physicalmachine")
	if err != nil {
		return err
	}
	
	for rows.Next() {
		rows.Scan(&p.IpAddress);
		GetNodeInfo(p.IpAddress)
	}
	return nil
}
//获取单个节点的信息
func GetNodeInfo(ip string)  error {
	
	fmt.Println("entry GetNodeInfo(ip)")
    db, err := GetDB()
	defer db.Close()
	if err != nil {
		return  err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system", ip))
	if err != nil {
		return  err
	}
	NodeInfo, _ := conn.GetNodeInfo();
	//fmt.Println("get nodeinfo","update physicalmachine set cputype='%s',cpu='%s', cpufreq='%s',cpusocket='%s', cpukernel='%s', cputhread='%s',numa='%s', memory='%s' where ipaddress ='%s'",NodeInfo.Model, NodeInfo.Cpus, NodeInfo.MHz, NodeInfo.Sockets, NodeInfo.Nodes, NodeInfo.Threads, NodeInfo.Nodes, NodeInfo.Memory,ip)
	sqltext := fmt.Sprintf("update physicalmachine set cputype='%s',cpu='%s', cpufreq='%s',cpusocket='%s', cpukernel='%s', cputhread='%s',numa='%s', memory='%s' where ipaddress ='%s'",NodeInfo.Model, strconv.FormatUint(uint64(NodeInfo.Cpus),10), strconv.FormatUint(uint64(NodeInfo.MHz),10), strconv.FormatUint(uint64(NodeInfo.Sockets),10), strconv.FormatUint(uint64(NodeInfo.Cores),10), strconv.FormatUint(uint64(NodeInfo.Threads),10), strconv.FormatUint(uint64(NodeInfo.Nodes),10), strconv.FormatUint(uint64(NodeInfo.Memory),10), ip)
	fmt.Println("get nodeinfo",sqltext)
	db.Exec(sqltext)
	return nil
}
func GetNodeNetInfo(ip string){
	fmt.Println("entry GetNodeNetInfo")
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system",ip))
	if err != nil {
		//return  err
	}
	NodeNetNames, _ := conn.ListNetworks();
	for _, NodeNetName := range NodeNetNames {
		fmt.Println(NodeNetName)
		if NodeNetName != "lo"{
			 Network, _ := conn.LookupNetworkByName(NodeNetName)
			 netxml, _ := Network.GetXMLDesc(1)
			 fmt.Println(netxml)
		   }
			
		}
	
}
func GetNodeNetInfos() error{
	fmt.Println("entry GetNodeNetInfos")
	p := &PhysicalMachine{}
    db, err := GetDB()
	defer db.Close()
	if err != nil {
		return  err
	}
	
	rows, err := db.Query("select IpAddress from physicalmachine")
	if err != nil {
		return err
	}
	
	for rows.Next() {
		rows.Scan(&p.IpAddress);
		GetNodeNetInfo(p.IpAddress)
		
	}
	return nil
}

func GetUser(Token string) (string,error) {
	fmt.Println("GetUser start")
	db, err := sql.Open("mysql", "root:@tcp(10.72.84.145:4000)/iam")
	if err != nil {
		fmt.Println("sql err",err)
		return "", err
	}
	fmt.Println("query user===")
	sql := fmt.Sprintf("select userName from token where token = '%s'", Token)
	err = db.QueryRow(sql).Scan(&User)
	if err != nil {
		fmt.Println("query err",err)
		return "", err
	}
	fmt.Println("user===",User)
	return User, nil
	
}
func NowTime() string {
	now := time.Now().Format("2006-01-02T15:04:05-0700")
	return now
}
