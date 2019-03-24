package main


type bonus struct {
	x, y int
}

func (b *bonus) update_position(x, y int){
	b.x = x
	b.y = y
}

func (*bonus) create_bonus(x, y int) bonus{
	b := bonus{x,y}
	return b
}
