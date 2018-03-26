package main

import (
    "fmt"
    //"database/sql"
	//"github.com/golang/crypto/ssh"
	"golang.org/x/crypto/ssh"
	//utils "github.com/converge/wandycatkeeper/utils"
	//"os"
	"strings"
	"io/ioutil"
	//"bufio"
	"github.com/google/uuid"
)
type DiskInfo struct {
	Uuid      string
	Ip        string
	Name      string
	Capility  string
	Status    string
}

type DiskInfos struct {
		DiskInfos []DiskInfo
}

type sshcmd struct {	
	    user      string
	    password  string
	    ip_port   string
	    cmd       string
}

type Netinfo struct {
	Ip        string
	Netspeed string
}

type Netinfos struct {
	Netinfos []*Netinfo
}
  
const SHPATH = "sh/disk_collector.sh"
const NETPATH = "sh/net_info.sh"


//get all node net speed  info 
func GetNetSpeedInfos() (*Netinfos, error){
	
    var ip string
    
    db, err := GetDB()
    defer db.Close()
    if err != nil {
		return nil, err
	}
    rows, err := db.Query("select IpAddress from physicalmachine")
    if err != nil {
		return nil, err
	}
    Netinfos := &Netinfos{}
	for rows.Next() {
		rows.Scan(&ip)
        Netinfo, _ := GetNetSpeedInfo(ip)
        fmt.Println("end GetNetSpeedInfo",Netinfo)
        Netinfos.Netinfos = append(Netinfos.Netinfos,Netinfo)
        
	}
    //阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
    fmt.Println("end  GetNetSpeedInfo",Netinfos)
    return Netinfos, nil
}

//get one node net speed info
func GetNetSpeedInfo(ip string) (*Netinfo, error){
	
	//cmd := exec.Command(SHPATH, "-c", s)
    //StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
    fmt.Println("entry GetNetSpeedInfo")
    scm := &sshcmd{}
    scm.user = "root"
    scm.password = "1234,qwer"
    scm.ip_port = ip+":22"
    scm.cmd, _ = getFileInfo(NETPATH)
    stdout, err :=  SSH_do(scm.user,scm.password,scm.ip_port,scm.cmd)
    if err != nil {
    	    fmt.Println("SSH_do",err)
        return nil,err
    }
    //创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
    //reader := strings.NewReader(stdout)
    //_, out, _ := bufio.ScanLines(stdout, atEOF)
    //reader := ioutil.
    //实时循环读取输出流中的一行内容
    Netinfo := &Netinfo{}
    //Netinfos := &Netinfos{}
    fmt.Println("end SSH_do",stdout)
    Netinfo.Ip = ip
    Netinfo.Netspeed = strings.Replace(stdout, "\n", "", -1)
    //Netinfos.Netinfos = append(Netinfos.Netinfos, Netinfo)
    fmt.Println("end  GetNetSpeedInfo",Netinfo)
    //阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
    return Netinfo, nil
}


//get one node disk info
func GetDiskInfo(ip string) (*DiskInfos, error){
	
	//cmd := exec.Command(SHPATH, "-c", s)
    //StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
    fmt.Println("entry GetDiskInfo")
    scm := &sshcmd{}
    scm.user = "root"
    scm.password = "1234,qwer"
    scm.ip_port = ip+":22"
    scm.cmd, _ = getFileInfo(SHPATH)
    stdout, err :=  SSH_do(scm.user,scm.password,scm.ip_port,scm.cmd)
    if err != nil {
    	    fmt.Println("SSH_do",err)
        return nil,err
    }
    //创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
    //reader := strings.NewReader(stdout)
    //_, out, _ := bufio.ScanLines(stdout, atEOF)
    //reader := ioutil.
    //实时循环读取输出流中的一行内容
    DiskInfo := DiskInfo{}
    DiskInfos := &DiskInfos{}
    db, err := GetDB()
    defer db.Close()
    if err != nil {
		return nil,err
	}
    fmt.Println("end SSH_do",stdout)
    str := strings.Split(stdout, "\n")  
    for _, line  := range str {
    	     if line!= ""{
    	     	fmt.Println("line",line)
    	        array := strings.Split(line, ":")
            names := strings.Split(array[1], " ")
            name  := names[1]
            fmt.Println("array",array)
            fmt.Println("array[2]",array[2])
            capilitys := strings.Split(array[2]," ")
            fmt.Println("capilitys",capilitys)
            capility  := capilitys[:3]
            fmt.Println("capility",capility)
            status := array[3] 
            DiskInfo.Ip = ip
            DiskInfo.Uuid = uuid.New().String()
            DiskInfo.Name = name
            DiskInfo.Capility = strings.Join(capility, "")
            DiskInfo.Status = status 
             
             fmt.Println("get  physicalmachinediskinfo",DiskInfo.Uuid, DiskInfo.Ip, DiskInfo.Name, DiskInfo.Capility, DiskInfo.Status)
             //_, err :=db.Exec("insert into physicalmachinediskinfo(uuid, Ip, name, capility, status) values ('%s', '%s' ,'%s', '%s', '%s')",DiskInfo.Uuid, DiskInfo.Ip, DiskInfo.name, DiskInfo.capility, DiskInfo.status )
            //if err != nil {
            //	   return nil,err
            //}
            DiskInfos.DiskInfos = append(DiskInfos.DiskInfos, DiskInfo) 
            fmt.Println("end  physicalmachinediskinfo")
    	     }
    } 
    //阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
    return DiskInfos, nil
}


//get all node disk info 
func GetDiskInfos() (*DiskInfos, error){
	
    var ip string
    
    db, err := GetDB()
    defer db.Close()
    if err != nil {
		return nil, err
	}
    rows, err := db.Query("select IpAddress from physicalmachine")
    if err != nil {
		return nil, err
	}
    AllDiskInfos := &DiskInfos{}
	for rows.Next() {
		rows.Scan(&ip)
        DiskInfos, _ := GetDiskInfo(ip)
        fmt.Println("end  physicalmachinediskinfo")
        for _, value := range DiskInfos.DiskInfos {
        	    fmt.Println("range DiskInfos.DiskInfos",value)
            AllDiskInfos.DiskInfos = append(AllDiskInfos.DiskInfos,value)
        }
	}
    //阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
    return AllDiskInfos, nil
}

//读取文件（filePath）里面的内容，返回一个String
func getFileInfo(filePath string) (string, error){
    fd, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "",err
    }
    return string(fd[:]), nil
	 
	
	
	
}
func SSH_do(user, password, ip_port string, cmd string)  (string, error) {

	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{User: user, Auth: PassWd,HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	Client, err  := ssh.Dial("tcp", ip_port, &Conf)
	if err != nil {
		return "",err
	}
	defer Client.Close()
	session, err := Client.NewSession()
	defer session.Close()
	if err != nil {
		return "", err
	}
    Stdout, Stderr := session.Output(cmd)
	return string(Stdout[:]),Stderr
}

