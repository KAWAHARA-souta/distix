package model

import (
	"fmt"
)


type PackageNevra struct {
	Name	string
	Version	string
	Release	string
	Arch	string
	Epoch	*int
}

func (pkgNevra *PackageNevra) GetNEVRA() string {
	base := fmt.Sprintf("%s-%s-%s.%s", pkgNevra.Name, pkgNevra.Version, pkgNevra.Release, pkgNevra.Arch)
	if pkgNevra.Epoch != nil {
		return fmt.Sprintf("%d:%s", *pkgNevra.Epoch, base)
	}
	return base
}
func (pkgNevra *PackageNevra) GetNVRA() string {
	return fmt.Sprintf("%s-%s-%s.%s",
		pkgNevra.Name, pkgNevra.Version, pkgNevra.Release, pkgNevra.Arch)
}
func (pkgNevra *PackageNevra) GetNVR() string {
	return fmt.Sprintf("%s-%s-%s",
		pkgNevra.Name, pkgNevra.Version, pkgNevra.Release)
}
func (pkgNevra *PackageNevra) GetVR() string {
	return fmt.Sprintf("%s-%s",
		pkgNevra.Version, pkgNevra.Release)
}

