package main


const(
	SINGLE_WEAPON = 0
	DOUBLE_WEAPON = 1
	TRIPLE_WEAPON = 2
)


type player struct{
	x, y, life, weapon, index  int
}


func (*player) update_weapon( p *player, weapon int)  {
	p.weapon = weapon
}

func (*player) loose_one_life(p *player){
	p.life --
}

func (*player) create_player(x, y, life, weapon, index int) player{
	p := player{x,y,life, weapon, index}
	return p
}



