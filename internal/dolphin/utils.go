package dolphin

import (
	pcast "github.com/pingcap/parser/ast"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type nodeSearch struct {
	list  []pcast.Node
	check func(pcast.Node) bool
}

func (s *nodeSearch) Enter(n pcast.Node) (pcast.Node, bool) {
	if s.check(n) {
		s.list = append(s.list, n)
	}
	return n, false // skipChildren
}

func (s *nodeSearch) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func collect(root pcast.Node, f func(pcast.Node) bool) []pcast.Node {
	if root == nil {
		return nil
	}
	ns := &nodeSearch{check: f}
	root.Accept(ns)
	return ns.list
}

type nodeVisit struct {
	fn func(pcast.Node)
}

func (s *nodeVisit) Enter(n pcast.Node) (pcast.Node, bool) {
	s.fn(n)
	return n, false // skipChildren
}

func (s *nodeVisit) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func visit(root pcast.Node, f func(pcast.Node)) {
	if root == nil {
		return
	}
	ns := &nodeVisit{fn: f}
	root.Accept(ns)
}

// Maybe not useful?
func text(nodes []pcast.Node) []string {
	str := make([]string, len(nodes))
	for i := range nodes {
		if nodes[i] == nil {
			continue
		}
		str[i] = nodes[i].Text()
	}
	return str
}

func parseTableName(n *pcast.TableName) *ast.TableName {
	return &ast.TableName{
		Schema: n.Schema.String(),
		Name:   n.Name.String(),
	}
}
