package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

const DiskPath string = "/home/"

func RemoveLocalImage(name, hostip string) error {
	cmd := exec.Command("ssh", hostip, "rm", "-rf", DiskPath+name)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Wait()
	return err
}

func ResizeLocalImage(name string, size uint64, hostip string) error {
	s := strconv.FormatUint(size, 10) + "G"
	cmd := exec.Command("ssh", hostip, "qemu-img", "resize", DiskPath+name, s)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func CreateLocalImage(name string, size uint64, hostip string) error {
	s := strconv.FormatUint(size, 10) + "G"
	cmd := exec.Command("ssh", hostip, "qemu-img", "create", "-f", "qcow2", DiskPath+name, s)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

/*
func main() {
	e := ResizeLocalImage("6dabc9ca-49a1-4f87-bb76-a43c070a38c3", 4, "10.72.84.146")
	fmt.Println(e)
}
*/
