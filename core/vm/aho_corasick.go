package vm

type ACTrieNode struct {
	children    []*ACTrieNode
	suffix_link *ACTrieNode
	dict_link   *ACTrieNode
	val         byte
	depth       int
	pattern     []byte
	parent      *ACTrieNode
	is_pattern  bool
	sym         uint32
}

// excuslive end [star, end)
type ACResult struct {
	sym   uint32
	start int
	end   int
}

type ACTree struct {
	root *ACTrieNode
}

func (n *ACTrieNode) find_child(c byte) *ACTrieNode {
	for _, child := range n.children {
		if child.val == c {
			return child
		}
	}
	return nil
}

func make_ACTrieNode(val byte) *ACTrieNode {
	return &ACTrieNode{
		children:    make([]*ACTrieNode, 0),
		suffix_link: nil,
		dict_link:   nil,
		val:         val,
		depth:       0,
		pattern:     nil,
		parent:      nil,
		is_pattern:  false,
	}
}

func (t *ACTree) add_pattern(pattern []byte, sym uint32) {
	// construct basic trie
	node := t.root
	for _, c := range pattern {
		child := node.find_child(c)
		if child == nil {
			new_node := make_ACTrieNode(c)
			new_node.parent = node
			new_node.depth = node.depth + 1
			node.children = append(node.children, new_node)
			node = new_node
		} else {
			node = child
		}
	}
	node.sym = sym
	node.pattern = pattern
	node.is_pattern = true
}

func (t *ACTree) add_suffix_links() {
	queue := make([]*ACTrieNode, 0)
	for i := 0; i < len(t.root.children); i++ {
		queue = append(queue, t.root.children[i])
	}
	for len(queue) != 0 {
		node := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		// search from the parent to find suffix link
		cur := node.parent
		for cur != nil {
			cur = cur.suffix_link
			// reaching root
			if cur == nil {
				node.suffix_link = t.root
				break
			}
			// check if the suffix link has the same child
			it := cur.find_child(node.val)
			if it != nil {
				node.suffix_link = it
				break
			}
		}
		// add all children to the queue
		for i := 0; i < len(node.children); i++ {
			queue = append(queue, node.children[i])
		}
	}
}

func (t *ACTree) add_dict_links() {
	queue := make([]*ACTrieNode, 0)
	for i := 0; i < len(t.root.children); i++ {
		queue = append(queue, t.root.children[i])
	}
	for len(queue) != 0 {
		node := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		// find the first pattern node through suffix link
		cur := node.suffix_link
		for cur != nil {
			if cur.is_pattern {
				node.dict_link = cur
				break
			}
			cur = cur.suffix_link
		}
		// add all children to the queue
		for i := 0; i < len(node.children); i++ {
			queue = append(queue, node.children[i])
		}
	}
}

func (t *ACTree) search(input []byte) []ACResult {
	node := t.root
	res := make([]ACResult, 0)
	for i := 0; i < len(input); i++ {
		child := node.find_child(input[i])
		// failed, follow suffix link
		if child == nil {
			for node.suffix_link != nil {
				node = node.suffix_link
				child = node.find_child(input[i])
				if child != nil {
					break
				}
			}
		}
		if child != nil {
            node = child
			dict := node.dict_link
			for dict != nil {
				// found by dictlink
				res = append(res, ACResult{dict.sym, (i - dict.depth + 1), (i + 1)})
				dict = dict.dict_link
			}
            // if child itself is a match
			if node.is_pattern {
				res = append(res, ACResult{node.sym, (i - node.depth + 1), (i + 1)})
			}
		}
	}
	return res
}

func InitSearcherFromDB() *ACTree {
	root := make_ACTrieNode(0)
	tree := ACTree{root: root}
	for i := 256; i < len(si_db); i++ {
		tree.add_pattern(si_db[i], uint32(i))
	}
	tree.add_suffix_links()
	tree.add_dict_links()
	return &tree
}
