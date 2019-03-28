package Space_invaders_client

const(
	ENEMY = 0
	PLAYER = 1
	)


type bullet struct {
	x, y, owner int
}

func (*bullet) update_position( p *bullet, x,y int)  {
	p.x = x
	p.y = y
}

func (*bullet) create_bullet(x, y, owner int) bullet{
	b := bullet{x,y,owner}
	return b
}
