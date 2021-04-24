package app

type Packet struct {
	Id              int
	VerticesVisited []int
}

func NewPacket(id int) *Packet {
	newPacket := new(Packet)
	newPacket.Id = id
	return newPacket
}
