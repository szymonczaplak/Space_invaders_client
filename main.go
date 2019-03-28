package Space_invaders_client

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	INITIAL_PLAYER_X = 150
	INITIAL_PLAYER_Y = 200
	INITIAL_WEAPON = 0
	INITIAL_LIFE = 2
	MOVEMENT_SPEED = 1

	ENEMY_SIZE = 10
	INITIAL_ENEMY_Y = -20

	BULLET_SPEED = 5

	ENEMY_SPEED = 1

	MIN_BULLET_Y = -1

	MAX_BULLET_Y = 400

	ENEMY_BULLET_SPEED = 1.3
	)

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func recFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  // parse arguments, you have to call this by yourself
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	re:= regexp.MustCompile("(?P<func>.*)=(?P<value>.*)")
	match := re.FindStringSubmatch(r.URL.Path[1:])
	//result := make(map[string]string)
	function := ""
	value := ""

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			switch name {
			case "func":
				function = match[i]
			case "value":
				value = match[i]
			}
		}
	}


	switch function {
	case "register_player":
		register_player(w)
	case "bye":
		sayBye(w,value)
	case "player_move":
		move_player(value[1:], value[:1])
	case "update_status":
		update_status(w)
	case "player_bullet_create":
		player_bullet_create(value)
	}

	 // send data to client side
}

func move_player(dir, index string)  {
	ind, _ := strconv.Atoi(index)
	switch dir {
	case "up":
		if players[ind].y <=0 {
			return
		}
		players[ind].y -= MOVEMENT_SPEED
	case "down":
		if players[ind].y >=290 {
			return
		}
		players[ind].y += MOVEMENT_SPEED
	case "left":
		if players[ind].x <=0 {
			return
		}
		players[ind].x -= MOVEMENT_SPEED
	case "right":
		if players[ind].x >=290 {
			return
		}
		players[ind].x += MOVEMENT_SPEED
		}

}

func player_bullet_create(index string){
	ind, _ := strconv.Atoi(index)
	parent := players[ind]
	b := bullet{parent.x, parent.y, PLAYER }
	bullets = append(bullets, &b)
}

func register_player( w http.ResponseWriter ){
	index := len(players)
	p := player{INITIAL_PLAYER_X, INITIAL_PLAYER_Y, INITIAL_LIFE, INITIAL_WEAPON, index}
	players = append(players, &p)

	fmt.Fprintf(w,strconv.Itoa(len(players)-1))
}

func update_status(w http.ResponseWriter){
	status :=  "OK\n" + update_players() + update_enemies() + update_bullets()
	fmt.Fprintf(w, status)
}

func update_players() string{
	x := ""
	y := ""
	life := ""
	weapon := ""
	index := ""
	for _, p := range players {
		x += strconv.Itoa(p.x) + " "
		y += strconv.Itoa(p.y) + " "
		life += strconv.Itoa(p.life) + " "
		weapon += strconv.Itoa(p.weapon) + " "
		index += strconv.Itoa(p.index) + " "
	}
	return x + "\n" + y + "\n" + life +"\n" + weapon + "\n" + index + "\n"
}

func update_enemies() string{
	x := ""
	y := ""
	life := ""
	weapon := ""

	for _, p := range enemies {
		x += strconv.Itoa(p.x) + " "
		y += strconv.Itoa(p.y) + " "
		life += strconv.Itoa(p.life) + " "
		weapon += strconv.Itoa(p.weapon) + " "
	}
	return x + "\n" + y + "\n" + life +"\n" + weapon + "\n"
}

func update_bullets() string{
	x := ""
	y := ""
	owner := ""

	for ind, b := range bullets {

		if b.y >= MAX_BULLET_Y || b.y < MIN_BULLET_Y {
			bullets = append(bullets[:ind], bullets[ind+1:]...)
			continue
		}
		x += strconv.Itoa(b.x) + " "
		y += strconv.Itoa(b.y) + " "
		owner += strconv.Itoa(b.owner) + " "
	}

	return x + "\n" + y + "\n" + owner + "\n"
}

func sayBye (w http.ResponseWriter, name string){
	fmt.Fprintf(w, "Bye " + name) // send data to client side
}

func listen(){
	err := http.ListenAndServe(":1212", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
}

func create_wave(){

	wave_array := []int {50, 80, 110, 140, 170, 200, 230}

	for _, x := range wave_array{
		for i := 0; i< epoch; i++ {
			if i%2 == 0{
				e := enemy{x + ENEMY_SIZE +5, INITIAL_ENEMY_Y + i * ENEMY_SIZE, INITIAL_LIFE, INITIAL_WEAPON}
				enemies = append(enemies, &e)
			}else {
				e := enemy{x, INITIAL_ENEMY_Y + i*ENEMY_SIZE, INITIAL_LIFE, INITIAL_WEAPON}
				enemies = append(enemies, &e)
			}
		}
	}
	epoch++
}

func add_bullets(){
	if time.Since(enemy_last_bullet_create).Seconds() < ENEMY_BULLET_SPEED{
		return
	}
	enemy_last_bullet_create = time.Now()
	for _, e := range enemies{
		if e.y >0 {
			bullets = append(bullets, &bullet{e.x, e.y, ENEMY})
		}
	}

}

func create_enemy_bullets()  {
	min_y := 300
	for _, e := range bullets{
		if e.y < min_y{
			min_y = e.y
		}
	}

	min_enemy_y := 300
	for _, e := range enemies{
		if e.y < min_enemy_y{
			min_enemy_y = e.y
		}
	}
		add_bullets()
	}

func move_enemies(){
	for _, e := range enemies{
		if time.Since(enemy_last_move_switch_time).Seconds() > 2.5{
			enemy_move_switch *= -1
			enemy_last_move_switch_time = time.Now()
		}
			e.x += enemy_move_switch * ENEMY_SPEED
			e.y += ENEMY_SPEED
	}
}

func bullet_move()  {
	for _, b := range bullets {
		if b.owner == PLAYER {
			b.y -= BULLET_SPEED
		}else {
			b.y += BULLET_SPEED
		}
	}
}


func remove_bullet(i int){
	if len(bullets) == 1 {
		bullets = [] *bullet{}
		return
	}
	bullets[i] = bullets[len(bullets)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	bullets = bullets[:len(bullets)-1]
}
func remove_enemy(i int){
	fmt.Print("Removing" + strconv.Itoa(i) + "from len:" + strconv.Itoa(len(enemies)) + "\n")
	//mu.Lock()
	if len(enemies) == 1 {
		enemies = []*enemy{}
		return
	}
	enemies[i] = enemies[len(enemies)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	enemies = enemies[:len(enemies)-1]
	//mu.Unlock()
}


func detect_collision(){
	to_remove := []int{}
	for ind, b := range bullets{
		switch b.owner {
		case ENEMY:
			for _, p := range players{
				if b.x <= p.x+10 && b.x >= p.x {
					if b.y <= p.y+10 && b.y >=p.y{
						p.loose_one_life()
						to_remove =append(to_remove, ind)
					}
				}
			}
		case PLAYER:
			for _, e := range enemies{
				if b.x <= e.x+10 && b.x >= e.x {
					if b.y <= e.y+10 && b.y >=e.y{
						e.loose_one_life()
						to_remove =append(to_remove, ind)
					}
				}
			}
		}
	}
	removed := 0
	for i := 0; i<len(to_remove); i++{
		remove_bullet(to_remove[i] - removed)
		removed++
	}
}


func check_enemies_life(){
	to_remove := []int{}
	for ind, e := range enemies{
		if e.life <= 0 {
			to_remove = append(to_remove, ind)
		}
	}
	removed := 0
	for i := 0; i<len(to_remove); i++{
		remove_enemy(to_remove[i] - removed)
		removed++
	}
}

func main_loop(){
	for {
		if len(players) == 0{
			fmt.Print("No players in game\n")
			time.Sleep(60 * time.Millisecond)
			continue
		} else{
			fmt.Print(strconv.Itoa(len(players)) + "players in game\n")
		}

		if len(enemies) == 0 {
			create_wave()
		}

		bullet_move()
		move_enemies()
		detect_collision()
		create_enemy_bullets()
		check_enemies_life()

		time.Sleep(60 * time.Millisecond)
	}
}

var players = [] *player{}
var bullets = [] *bullet{}

var bonuses = []*bonus{}

var meteors = []*meteor{}

var enemies = []*enemy{}

var enemy_move_switch = 1

var enemy_last_move_switch_time = time.Now()

var enemy_last_bullet_create = time.Now()

var epoch = 4

var mu = sync.Mutex{}

var gameState = "stop"
func main() {
	go main_loop()

	flag.Parse()

	http.HandleFunc("/", recFunc)            // set router
	listen()

	for{}
}