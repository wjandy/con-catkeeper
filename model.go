package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	//"database/sql"
	libvirt "github.com/libvirt/libvirt-go"
	"text/template"
	utils "github.com/converge/catkeeper/utils"
)

type PhysicalMachine struct {
	/*IpAddress   string
	Name        string
	Description string
	Active      bool*/
	IpAddress   string
	Name        string
	Description string
	Cputype     string
    Cpu         int
    Cpufreq     int
    Cpusocket   int
    Cpukernel   int
    Cputhread   int
    Numa        int
    Memory      int
	Active      bool	
	conn        *libvirt.Connect
}

func NewPhysicalMachine(ipaddress string) (*PhysicalMachine, error) {
	p := &PhysicalMachine{}
	sql := fmt.Sprintf("select IpAddress,Description,Name,cputype,cpu,cpufreq,cpusocket,cpukernel,cputhread,numa,memory from physicalmachine where ipaddress='%s'", ipaddress)
	db, err := GetDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	utils.Logger.Infoln(sql)
	err = db.QueryRow(sql).Scan(
		  &p.IpAddress, 
		  &p.Description,
		  &p.Name,
		  &p.Cputype,
          &p.Cpu,
          &p.Cpufreq,
          &p.Cpusocket,
          &p.Cpukernel,
          &p.Cputhread,
          &p.Numa,
          &p.Memory,
	)
	if err != nil {
		return nil, err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system", p.IpAddress))
	//defer conn.close()
	if err != nil {
		p.Active = false
	}
	p.Active = true
	p.conn = conn
	return p, nil
}

func (p *PhysicalMachine) CreateVM(v *VirtualMachine) error {
	//clone image
	vol := &Volume{}
	vol.Uuid = uuid.New().String()
	vol.Size = v.Disk
	vol.Status = "in-use"
	vol.VolumeType = "rbd"
	vol.CreateAt = NowTime()
	vol.UpdateAt = NowTime()
	vol.Attachments = v.Uuid
	vol.User = v.User
	vol.DataType = "sys"

	v.SysDisk = vol.Uuid
	_, err := Clone(v.ImageName, vol.Uuid)
	if err != nil {
		utils.Logger.Infoln("err occur where clone", err)
		return err
	}
	//resize image
	err = ResizeImage(vol.Uuid, vol.Size)
	if err != nil {
		fmt.Println("err occur where rbd resize", err)
		return err
	}

	t, err := template.ParseFiles("templates/createvm.tmpl")
	if err != nil {
		fmt.Println("err occur when parsefile", err)
		return err
	}
	fmt.Println("template.ParseFiles",t)
	var xml *bytes.Buffer = new(bytes.Buffer)
	err = t.Execute(xml, v)
	if err != nil {
		fmt.Println("err occur when execute template", err)
		return err
	}
	
	fmt.Println("creating kvm insert into dbbase")
	v.Status = "creating"
	err = v.CreateDb()
	if err != nil {
		return err
	}
	err = vol.CreateVolDB()
	if err != nil {
		return err
	}
	
	xmlStr := xml.String()
	fmt.Println(xmlStr)
	d, err := p.conn.DomainDefineXML(xmlStr)
	if err != nil {
		// add roll back function
		v.DeleteVMDb()
		//vol.DeleteVolDb()
		fmt.Println("err occure when define xml", err)
		return err
	}
	err = d.Create()
	if err != nil {
		//update status
		v.Status = "Creating ERROR"
		v.UpdateVMDb()
		// add roll back function
		return err
	}
	
	return err
}
/*
func myrebootcallback(c *libvirt.Connect, d *libvirt.Domain, event *libvirt.DomainEventLifecycle) {
	fmt.Printf("Got event %d\n", event.Event)
	if event.Event == libvirt.DOMAIN_EVENT_STOPPED {
		fmt.Println("rebooting...")
		d.Create()
	}
	name, _ := d.GetName()
	if callbackMap.Check(name) == true {
		callbackId := callbackMap.Get(name).(int)
		c.DomainEventDeregister(callbackId)
		callbackMap.Delete(name)
	}
}
*/

/*
func (p *PhysicalMachine) String() string {
	var result = ""
	result += fmt.Sprintf("%s(%s) running?%t\n", p.Name, p.IpAddress, p.Existing)
	for _, vmPtr := range p.VirtualMachines {
		result += fmt.Sprintf("%s\n", vmPtr)
	}

	return result
}
*/

type VirtualMachine struct {
	Uuid          string
	Name          string
	Cpu           int
	Mem           int
	Disk          uint64
	Description   string
	HostIpAddress string
	Status        string
	User          string
	CreateTime    string
	UpdateTime    string
	Attachments   string
	ImageName     string
	SysDisk       string
}

func (v *VirtualMachine) SyncDb() error {
	sqltext := fmt.Sprintf("update virtualmachine set uuid='%s' , cpu=%d , mem=%d , disk=%d , description='%s' , hostipaddress='%s' , status='%s' , attachments='%s' , imagename='%s' , name='%s' , sysdisk='%s' where uuid='%s'", v.Uuid, v.Cpu, v.Mem, v.Disk, v.Description, v.HostIpAddress, v.Status, v.Attachments, v.ImageName, v.Name, v.SysDisk, v.Uuid)
	fmt.Println(sqltext)
	db, _ := GetDB()
	_, err := db.Exec(sqltext)
	return err
}

func (v *VirtualMachine) UpdateVMDb() error {
	sqltext := fmt.Sprintf("update virtualmachine set uuid='%s' , cpu=%d , mem=%d , disk=%d , description='%s' , hostipaddress='%s' , status='%s' , attachments='%s' , imagename='%s' , name='%s' , sysdisk='%s' where uuid='%s'", v.Uuid, v.Cpu, v.Mem, v.Disk, v.Description, v.HostIpAddress, v.Status, v.Attachments, v.ImageName, v.Name, v.SysDisk, v.Uuid)
	fmt.Println(sqltext)
	db, _ := GetDB()
	_, err := db.Exec(sqltext)
	return err
}


func (v *VirtualMachine) DeleteVMDb() error {
	sqltext := fmt.Sprintf("delete from  virtualmachine where uuid='%s'", v.Uuid)
	db, _ := GetDB()
	defer db.Close();
	fmt.Println("delete  virtualmachine:",sqltext)
	_, err := db.Exec(sqltext)
	return err
}

func (v *VirtualMachine) CreateDb() error {
	sqltext := fmt.Sprintf("insert into virtualmachine(uuid,cpu,mem,disk,description,hostipaddress,status,attachments,imagename,name,sysdisk,user,createtime,updatetime) values('%s',%d,%d,%d,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", v.Uuid, v.Cpu, v.Mem, v.Disk, v.Description, v.HostIpAddress, v.Status, v.Attachments, v.ImageName, v.Name, v.SysDisk,v.User,v.CreateTime,v.UpdateTime)
	db, _ := GetDB()
	defer db.Close();
	fmt.Println("insert virtualmachine:",sqltext)
	_, err := db.Exec(sqltext)
	return err
}

func GetVirtualMachine(uuid string) (*VirtualMachine, error) {
	sqltext := fmt.Sprintf("select uuid,cpu,mem,disk,hostipaddress,description,status,attachments,imagename,name,sysdisk,user,createtime,updatetime from virtualmachine where uuid='%s'", uuid)
	fmt.Println("select VirtualMachine:",sqltext)
	db, _ := GetDB()
	v := &VirtualMachine{}
	err := db.QueryRow(sqltext).Scan(
		&v.Uuid,
		&v.Cpu,
		&v.Mem,
		&v.Disk,
		&v.HostIpAddress,
		&v.Description,
		&v.Status,
		&v.Attachments,
		&v.ImageName,
		&v.Name,
		&v.SysDisk,
		&v.User,
		&v.CreateTime,
		&v.UpdateTime,
	)
	return v, err
}
/*
func ListVirtualMachine() (*[]VirtualMachine, error) {
	type VMs struct {
		VMs []VirtualMachine
	}
	//v := &VirtualMachine{}
	sqltext := fmt.Sprintf("select uuid,cpu,mem,disk,hostipaddress,description,status,attachments,imagename,name,sysdisk from virtualmachine ")
	fmt.Println(sqltext)
	db, _ := GetDB()
	defer db.Close()
	v := &VirtualMachine{}
	rows, err := db.Query(sqltext)
	for rows.Next() {
	    rows.Scan(
		&v.Uuid,
		&v.Cpu,
		&v.Mem,
		&v.Disk,
		&v.HostIpAddress,
		&v.Description,
		&v.Status,
		&v.Attachments,
		&v.ImageName,
		&v.Name,
		&v.SysDisk,
	  )
	  VMs  
	}
	return v, err
}
*/
/*
func (this *VirtualMachine) Start() error {
	err := this.Domain.Create()
	return err
}

func (this *VirtualMachine) Delete(db *sql.DB) error {

	for _, diskpath := range this.Disks {
		log.Printf("deleteing disk %s", diskpath)
		v, err := this.Connect.LookupStorageVolByPath(diskpath)
		if err != nil {
			log.Printf("%s can not be found by libvirt", diskpath)
			continue
		}
		//delete storage
		v.Delete(libvirt.STORAGE_VOL_DELETE_NORMAL)
		v.Free()
	}

	//remove domain
	err := this.Domain.Undefine()
	if err != nil {
		return err
	}

	//remove from database
	_, err = db.Exec("delete from virtualmachine where Id = ?", this.Id)
	if err != nil {
		return err
	}
	//find VM's all mac address, delete all mac<=>ip mappings
	rows, err := db.Query("select MAC from vmmacmapping where VmId = ?", this.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	var mac string
	for rows.Next() {
		rows.Scan(&mac)
		db.Exec("delete from macipmappingcache where MAC = ?", mac)
	}

	_, err = db.Exec("delete from vmmacmapping where VmId = ?", this.Id)
	if err != nil {
		return err
	}
	return nil
}

func (this *VirtualMachine) Stop() error {
	err := this.Domain.Shutdown()
	return err
}

func (this *VirtualMachine) ForceStop() error {
	err := this.Domain.Destroy()
	return err
}

func (this *VirtualMachine) Free() error {
	err := this.Domain.Free()
	return err
}

func (this *VirtualMachine) UpdateDatabase(db *sql.DB, owner string, description string) error {
	if owner != "" && description != "" {
		_, err := db.Exec("update virtualmachine set Owner=?,Description=? where Id=?", owner, description, this.Id)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else { //can not be updated
		return errors.New("owner , description must have values")
	}
}
*/
type Volume struct {
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

func (v *Volume) CreateVolDB() error {
	sqltext := fmt.Sprintf("insert into volume values('%s',%d,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", v.Uuid, v.Size, v.Description, v.Name, v.VolumeType, v.Status, v.CreateAt, v.UpdateAt, v.Attachments, v.HostIp, v.DataType,v.User)
	fmt.Println("insert  volume:",sqltext)
	db, _ := GetDB()
	defer db.Close()
	_, err := db.Exec(sqltext)
	return err
}

func (v *Volume) SyncVolume() (err error) {
	sqltext := fmt.Sprintf("update volume set uuid='%s' , size=%d , description='%s' , name='%s' , volumetype='%s' , status='%s' , createat='%s' , updateat='%s' , attachments='%s' , hostip='%s' , datatype='%s' where uuid='%s';", v.Uuid, v.Size, v.Description, v.Name, v.VolumeType, v.Status, v.CreateAt, v.UpdateAt, v.Attachments, v.HostIp, v.DataType, v.Uuid)
	fmt.Println(sqltext)
	db, _ := GetDB()
	defer db.Close()
	_, err = db.Exec(sqltext)
	return err
}


func (v *Volume) DeleteVolDb() (err error) {
	sqltext := fmt.Sprintf("delete  volume  from  where uuid='%s';", v.Uuid)
	fmt.Println(sqltext)
	db, _ := GetDB()
	defer db.Close()
	_, err = db.Exec(sqltext)
	return err
}

func GetVolume(uuid string) (vol *Volume, err error) {
	sqltext := fmt.Sprintf("select * from volume where uuid='%s'", uuid)
	db, err := GetDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	vol = &Volume{}
	err = db.QueryRow(sqltext).Scan(
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
	return vol, err
}

func GetDomainFromUuid(uuid string) (domain *libvirt.Domain, err error) {
	v, err := GetVirtualMachine(uuid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system", v.HostIpAddress))
	if err != nil {
		return nil, err
	}
	domain, err = conn.LookupDomainByUUIDString(uuid)
	return
}
func GetDomainINAXMLFromUuid(uuid string) (string, string, error) {
	v, err := GetVirtualMachine(uuid)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system", v.HostIpAddress))
	if err != nil {
		return "", "", err
	}
	domain, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return "", "", err
	}
	xmlData, err := domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	return xmlData, v.HostIpAddress, err
}

func GetDomainXMLFromUuid(uuid string) (string, string, error) {
	v, err := GetVirtualMachine(uuid)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+ssh://root@%s/system", v.HostIpAddress))
	if err != nil {
		return "", "", err
	}
	domain, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return "", "", err
	}
	xmlData, err := domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	return xmlData, v.HostIpAddress, err
}
