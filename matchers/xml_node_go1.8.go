// +build go1.8

package matchers

import "encoding/xml"

type xmlNode struct {
	XMLName xml.Name
	XMLAttr []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []*xmlNode `xml:",any"`
}

func (n *xmlNode) Clean() {
	if len(n.Nodes) == 0 {
		return
	}
	n.Content = nil
	for _, child := range n.Nodes {
		child.Clean()
	}
}
