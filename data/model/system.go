package model

import (
	// "errors"
	"fmt"
	"os"

	"github.com/protobom/protobom/pkg/sbom"

	// "github.com/distix-pj/distix/data/model"
)


type System struct {
	HostName string
	Packages []Package
}


func (sys *System) BaseProtobomDocument() (*sbom.Document, error) {
	doc := sbom.NewDocument()

	doc.Metadata.Name = sys.HostName
	doc.Metadata.Version = "v9.9.9"

	systemNode := &sbom.Node{
		Id:               sys.HostName,
		Name:             sys.HostName,
		// Version:          pkg.getVersion(),
	}
	doc.NodeList.AddRootNode(systemNode)

	rpmTopNode := &sbom.Node{
		Id:          "RPM-Packages",
		Type:        sbom.Node_PACKAGE,
		Name:        "RPM-Packages",
		// Version:     req.Version,
	}
	doc.NodeList.AddNode(rpmTopNode)
	doc.NodeList.RelateNodeAtID(rpmTopNode, systemNode.Id, sbom.Edge_contains)

	return doc, nil
}

func (sys *System) Convert2ProtobomDocument() (*sbom.Document, error) {
	doc, err := sys.BaseProtobomDocument()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	rpmTopNode := doc.NodeList.GetNodeByID("RPM-Packages")
	for _, pkg := range sys.Packages {
		// 	// if req.IsSoLib() {
		pkgNode := &sbom.Node{
			Id:          pkg.getPurl(),
			Type:        sbom.Node_PACKAGE,
			Name:        pkg.getName(),
			Version:     pkg.getVersion(),
		}
		doc.NodeList.AddNode(pkgNode)
		doc.NodeList.RelateNodeAtID(pkgNode, rpmTopNode.Id, sbom.Edge_contains)
		// 	// }
	}

	for _, pkg := range sys.Packages {
		// pkgNode := doc.NodeList.GetNodeByID(pkg.getPurl())
		// if pkgNode == nil {
		// 	err := errors.New("Unexpected Internal Error. Can't find pkgNode in OneSystem.Convert2ProtobomDocument()")
		// 	fmt.Fprintln(os.Stderr, err)
		// 	return nil, err
		// }
		for _, req := range pkg.Requires {
pkgReqLoop:
			for _, extPkg := range sys.Packages {
				extPkgNode := doc.NodeList.GetNodeByID(extPkg.getPurl())
				for _, prov := range extPkg.Provides {
					if req.Name == prov.Name {
						if !isEdgeExists(doc.NodeList, pkg.getPurl(), extPkg.getPurl(), sbom.Edge_dependsOn) && !isSelfDepend(pkg.getPurl(), extPkg.getPurl()) {
							doc.NodeList.RelateNodeAtID(extPkgNode, pkg.getPurl(), sbom.Edge_dependsOn)
						}
						break pkgReqLoop
					}
				}
				for _, file := range extPkg.Files {
					if req.Name == file.Name {
						if !isEdgeExists(doc.NodeList, pkg.getPurl(), extPkg.getPurl(), sbom.Edge_dependsOn) && !isSelfDepend(pkg.getPurl(), extPkg.getPurl()) {
							doc.NodeList.RelateNodeAtID(extPkgNode, pkg.getPurl(), sbom.Edge_dependsOn)
						}
						break pkgReqLoop
					}
				}
			}
		}
	}

	return doc, nil
}

func (sys *System) Convert2MultiProtobomDocument() (*sbom.Document, []*sbom.Document, error) {
	doc, err := sys.BaseProtobomDocument()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, nil, err
	}
	pkgDocs := make([]*sbom.Document, len(sys.Packages))
	for i, pkg := range sys.Packages {
		pkgDoc, err := pkg.Convert2ProtobomDocument()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}
		pkgDocs[i] = pkgDoc
	}
	return doc, pkgDocs, nil
}



func isEdgeExists(nl *sbom.NodeList, fromID, toID string, edgeType sbom.Edge_Type) bool {
	for _, edge := range nl.GetEdges() {
		if edge.GetFrom() == fromID && edge.GetType() == edgeType {
			for _, to := range edge.GetTo() {
				if to == toID {
					return true
				}
			}
		}
	}
	return false
}

func isSelfDepend(fromID, toID string) bool {
	if fromID == toID {
		return true
	}
	return false
}

