package app

type Vertex struct {
	Id                        int
	ThisVertexChannel         chan *Packet
	ThisVertexResponseChannel chan bool
	ThisVertexAskChannel      chan chan bool
	HasPacket                 bool
	IsTrapped                 bool
	NextVerticesChannels      []chan<- *Packet
	NextVerticesAskChannels   []chan<- chan bool

	NextVerticesInfo []*Vertex
	PacketsSeenInfo  []int
	Rounds           int
	Tries            int
}

func NewVertex(id int) *Vertex {
	newVertex := new(Vertex)
	newVertex.Id = id
	newVertex.Rounds = 0
	newVertex.Tries = 0
	newVertex.ThisVertexChannel = make(chan *Packet, 0)
	newVertex.ThisVertexResponseChannel = make(chan bool, 0)
	newVertex.ThisVertexAskChannel = make(chan chan bool, 0)
	newVertex.HasPacket = false
	newVertex.IsTrapped = false
	return newVertex
}
