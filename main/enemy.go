package main

type enemy struct{
	x, y, life, weapon  int
}

func (*enemy) update_position( p *enemy, x,y int)  {
	p.x = x
	p.y = y
}

func (*enemy) update_weapon( p *enemy, weapon int)  {
	p.weapon = weapon
}

func (*enemy) loose_one_life(p *enemy){
	p.life --
}

func (*enemy) create_enemy(x, y, life, weapon, index int) enemy{
	p := enemy{x,y,life, weapon}
	return p
}
