package main

import (
	"encoding/json"
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	//	"io"
	//	"os"
	//	"os/exec"
)

type Ceph struct {
	Conn   *rados.Conn
	Rbdctx *rados.IOContext
	Ckctx  *rados.IOContext
}

func NewCephConnect() (c *Ceph, err error) {
	c = &Ceph{}
	conn, err := rados.NewConn()
	if err != nil {
		return
	}
	conn.ReadDefaultConfigFile()
	conn.Connect()
	c.Conn = conn
	rbdctx, err := conn.OpenIOContext("rbd")
	if err != nil {
		return
	}
	c.Rbdctx = rbdctx
	ckctx, err := conn.OpenIOContext("catkeeper")
	if err != nil {
		return
	}
	c.Ckctx = ckctx
	return
}

func (c *Ceph) GetImage(name string) *rbd.Image {
	image := rbd.GetImage(c.Rbdctx, name)
	all, err := rbd.GetImageNames(c.Rbdctx)
	fmt.Println(all, err)
	return image
}

func Clone(imageName, destName string) (dest *rbd.Image, err error) {
	fmt.Println("enter clone:", imageName, destName)
	c, err := NewCephConnect()
	if err != nil {
		fmt.Println("connection err", err)
		return
	}
	image := c.GetImage(imageName)
	i, _ := json.Marshal(image)
	fmt.Println(string(i))
	//	dest, err = image.Clone("rbd/"+imageName+"@first", c.Ckctx, destName, 3, 22)
	dest, err = image.Clone("snapshot", c.Ckctx, destName, 3, 22)
	if err != nil {
		fmt.Println("clone error", err)
		return
	}
	return
}

func GetImage(name string) (image *rbd.Image, err error) {
	c, err := NewCephConnect()
	if err != nil {
		fmt.Println("newceph error", err)
		return
	}
	image = rbd.GetImage(c.Ckctx, name)
	return
}

func CreateImage(name string, size uint64) (image *rbd.Image, err error) {
	c, err := NewCephConnect()
	if err != nil {
		fmt.Println("newceph error", err)
		return
	}
	size = size * 1 << 30
	image, err = rbd.Create(c.Ckctx, name, size, 22)
	return
}

func ResizeImage(name string, size uint64) error {
	i, err := GetImage(name)
	if err != nil {
		return err
	}
	err = i.Open()
	defer i.Close()
	if err != nil {
		return err
	}
	size = size << 30
	err = i.Resize(size)
	return err
}

func RemoveImage(name string) error {
	i, err := GetImage(name)
	if err != nil {
		return err
	}
	err = i.Remove()
	return err
}

/*
func Export(snapName, desName string) (err error) {
	c, err := NewCephConnect()
	if err != nil {
		fmt.Println("connection err", err)
		return
	}
	image := c.GetImage(snapName)
	err = image.Open()
	if err != nil {
		return
	}
	dst, err := os.OpenFile("/tmp/"+desName, os.O_WRONLY|os.O_CREATE, 0644)
	_, err = io.Copy(dst, image)
	return
}

func Clone(imageName, destName string) (err error) {
	cmd := exec.Command("rbd", "clone", "rbd/vm222@first", "catkeeper/vm2")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return
}
*/

/*
func main() {
	d, e := Clone("vm222", "test128")
	fmt.Println(d, e)
}
*/
/*
func main() {
	_, e := CreateImage("wao", 2)
	fmt.Println(e)
}
*/
