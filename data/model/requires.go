package model

import (
	"fmt"
	"regexp"
)


type RpmCapability struct {
	Name    string
	Version string
	Flags   uint32
}

func (rCap *RpmCapability) IsSoLib() bool {
	rpmSoPattern := regexp.MustCompile(`^.*\.so(\.\d+)*(\([^)]*\))+$`)
	return rpmSoPattern.MatchString(rCap.Name)
}

func (rCap *RpmCapability) GetId() string {
	return fmt.Sprintf("RPM-CAP-%s", rCap.Name)
}


type RpmFile struct {
	Name string
	// fileInfo rpmtuils.FileInfo
}

func (rFile *RpmFile) GetId() string {
	return fmt.Sprintf("RPM-FILE-%s", rFile.Name)
}

