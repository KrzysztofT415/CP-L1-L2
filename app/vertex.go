package app

type Vertex struct {
	Id                   int
	ThisVertexChannel    chan *Packet
	NextVerticesChannels []chan<- *Packet

	NextVerticesInfo []*Vertex
	PacketsSeenInfo  []int
}

func NewVertex(id int) *Vertex {
	newVertex := new(Vertex)
	newVertex.Id = id
	newVertex.ThisVertexChannel = make(chan *Packet, 0)
	return newVertex
}
