package main


type meteor struct {
	x, y, size int
}

func (*meteor) update_position(m *meteor, x, y int){
	m.x = x
	m.y = y
}

func (*meteor) create_meteor(x, y, size int) meteor{
	m := meteor{x,y,size}
	return m
}
