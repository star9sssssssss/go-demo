package gee

import (
	"fmt"
	"strings"
)

//根据trie树构造路由, 每个子节点都是待匹配的路径

type node struct {
	pattern string    //待匹配的路由, 全段路由
	part string       //路由中的一部分
	children []*node  //子节点
	isWild bool       //是否是广泛匹配，含有 * 或者 : 为 true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//根据一小段路由找到第一个匹配的节点， 用于插入新的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//如果部分路径匹配或者属于模糊匹配类型的
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
	
//根据一小段路由找到所有匹配的节点， 用于查询
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

/*
插入一个节点
	pattern 整段路由
	parts 路由以 '/' 分开的数组
	height 当前trie树的高度
*/
func (n *node) insert(pattern string, parts []string, height int) {
	//如果到达整段路由的最后端，将这端路由设置到节点中
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//首先获得当前高度的一小段路径, 从第一层开始
	part := parts[height]
	//通过该段路径，查询是否有节点已经存在
	child := n.matchChild(part)
	if child == nil { //不存在
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//存在
	//以当前节点出发，继续寻找
	child.insert(pattern, parts, height + 1)
}


/*
	根据parts查找是否有该路径的节点
*/

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}


func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}




