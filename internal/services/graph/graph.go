package graph

type Node struct {
	ID   string         `json:"id"`
	Kind string         `json:"kind"`
	Name string         `json:"name"`
	Data map[string]any `json:"data,omitempty"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Builder struct {
	nodesMap map[string]Node
	edgesMap map[string]Edge
}

func NewGraphBuilder() *Builder {
	return &Builder{
		nodesMap: make(map[string]Node),
		edgesMap: make(map[string]Edge),
	}
}

func (b *Builder) AddNode(n Node) {
	if _, exists := b.nodesMap[n.ID]; !exists {
		b.nodesMap[n.ID] = n
	}
}

func (b *Builder) AddEdge(e Edge) {
	key := e.Source + "->" + e.Target
	b.edgesMap[key] = e
}

func (b *Builder) Build() *Graph {
	g := &Graph{
		Nodes: make([]Node, 0, len(b.nodesMap)),
		Edges: make([]Edge, 0, len(b.edgesMap)),
	}

	for _, n := range b.nodesMap {
		g.Nodes = append(g.Nodes, n)
	}
	for _, e := range b.edgesMap {
		g.Edges = append(g.Edges, e)
	}

	return g
}
