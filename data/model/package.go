package model

import (
	"errors"
	"fmt"

	"github.com/protobom/protobom/pkg/sbom"
)


type Package struct {
	PkgNevra PackageNevra
	Summary	string
	Description	string
	Provides []RpmCapability
	Requires []RpmCapability
	Files []RpmFile

	//! used by onesystem
	RequirePkgs []string
}


func (pkg *Package) getNEVRA() string {
	return pkg.PkgNevra.GetNEVRA()
}
func (pkg *Package) getName() string {
	return pkg.PkgNevra.Name
}
func (pkg *Package) getVersion() string {
	return pkg.PkgNevra.Version
}
func (pkg *Package) getVersionRelease() string {
	return pkg.PkgNevra.GetVR()
}
//! TODO: This is instant implementation. Need to be fixed.
func (pkg *Package) getPurl() string {
	return fmt.Sprintf("pkg:rpm/generic/%s@%s", pkg.getName(), pkg.getVersionRelease())
}


func (pkg *Package) Convert2ProtobomDocument() (*sbom.Document, error) {
	doc := sbom.NewDocument()

	doc.Metadata.Name = pkg.getNEVRA()
	doc.Metadata.Version = pkg.getVersion()

	// ...and the tool that produced the SBOM:
	// doc.Metadata.Tools = append(
	// 	doc.Metadata.Tools,
	// 	&sbom.Tool{
	// 		Name:    "ACME SBOM Tool",
	// 		Version: "1.0",
	// 		Vendor:  "ACME Corporation"},
	// )

	pkgNode := &sbom.Node{
		Id:               pkg.getPurl(),
		Name:             pkg.getName(),
		Version:          pkg.getVersion(),
	}
	doc.NodeList.AddRootNode(pkgNode)

	for _, req := range pkg.Requires {
		// if req.IsSoLib() {
		reqNode := &sbom.Node{
			Id:          req.GetId(),
			Type:        sbom.Node_PACKAGE,
			Name:        req.Name,
			Version:     req.Version,
		}
		doc.NodeList.AddNode(reqNode)
		doc.NodeList.RelateNodeAtID(reqNode, pkgNode.Id, sbom.Edge_dependsOn)
		// }
	}
	for _, prov := range pkg.Provides {
		// if prov.IsSoLib() {
		reqNode := &sbom.Node{
			Id:          prov.GetId(),
			Type:        sbom.Node_PACKAGE,
			Name:        prov.Name,
			Version:     prov.Version,
		}
		doc.NodeList.AddNode(reqNode)
		doc.NodeList.RelateNodeAtID(reqNode, pkgNode.Id, sbom.Edge_contains)
		// }
	}

	for _, file := range pkg.Files {
		fileNode := &sbom.Node{
			Id:          file.GetId(),
			Type:        sbom.Node_FILE,
			Name:        file.Name,
		}
		doc.NodeList.AddNode(fileNode)
		doc.NodeList.RelateNodeAtID(fileNode, pkgNode.Id, sbom.Edge_contains)
	}

	return doc, nil
}

func (pkg *Package) Convert2MultiProtobomDocument() (*sbom.Document, []*sbom.Document, error) {
	return nil, nil, errors.New("Not Implemented")
}

