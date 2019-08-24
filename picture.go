package main

import (
	"os/exec"
	"log"
)


func TakePicture() {
	err := exec.Command("raspistill", "-o", imgDir + "picture.jpg").Run()
	if err != nil {
		log.Fatal(err)
	}
}
