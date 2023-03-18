package app

type Packet struct {
	Id              int
	Life            int
	VerticesVisited []int
}

func NewPacket(id int, life int) *Packet {
	newPacket := new(Packet)
	newPacket.Id = id
	newPacket.Life = life
	return newPacket
}
