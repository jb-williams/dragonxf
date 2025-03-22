package main

// GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dragonxf_intel
// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dragonxf_linux
// GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dragonxf_arm

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

type Heroxf struct {
	XP     int
	ST     int
	AC     int
	MX     int
	HP     int
	Attk   int
	Dth    int
	PrevXP int // Track previous XP to determine level-up conditions
}

type Dragonxf struct {
	PR   int
	AC   int
	MX   int
	HP   int
	Attk int
	Dth  int
}

var hero Heroxf
var dragon Dragonxf

func getSaveFilePathxf() string {
	usr, _ := user.Current()
	saveDir := filepath.Join(usr.HomeDir, ".local", "share")
	os.MkdirAll(saveDir, os.ModePerm)
	return filepath.Join(saveDir, "dragonxf_game.save")
}

func initGamexf() {
	saveFile := getSaveFilePathxf()
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		hero = Heroxf{0, 1, 14, 50, 50, 4, 0, 0}
		dragon = Dragonxf{15, 18, 100, 100, 6, 0}
		saveGamexf()
	} else {
		loadGamexf()
	}
}

func saveGamexf() {
	file, _ := os.Create(getSaveFilePathxf())
	defer file.Close()
	fmt.Fprintf(file, "%d %d %d %d %d %d %d %d\n", hero.XP, hero.ST, hero.AC, hero.MX, hero.HP, hero.Attk, hero.Dth, hero.PrevXP)
	fmt.Fprintf(file, "%d %d %d %d %d %d\n", dragon.PR, dragon.AC, dragon.MX, dragon.HP, dragon.Attk, dragon.Dth)
}
func loadGamexf() {
	file, err := os.Open(getSaveFilePathxf())
	if err != nil {
		fmt.Println("⚠ Error loading save file, starting a new game.")
		initGamexf()
		return
	}
	defer file.Close()

	_, err = fmt.Fscanf(file, "%d %d %d %d %d %d %d %d", &hero.XP, &hero.ST, &hero.AC, &hero.MX, &hero.HP, &hero.Attk, &hero.Dth, &hero.PrevXP)
	if err != nil {
		fmt.Println("⚠ Error reading Chadwick's stats, resetting game.")
		initGamexf()
		return
	}

	_, err = fmt.Fscanf(file, "%d %d %d %d %d %d", &dragon.PR, &dragon.AC, &dragon.MX, &dragon.HP, &dragon.Attk, &dragon.Dth)
	if err != nil {
		fmt.Println("⚠ Error reading Dragon's stats, resetting game.")
		initGamexf()
	}
}

func heroDamagexf() int {
	damage := rand.Intn(8) + 1 + hero.Attk
	return damage
}

func heroCriticalHitsxf() int {
	damage := 2 * heroDamagexf()
	// 	damage := 2 * (8 + 1 + hero.Attk)
	return damage
}

func dragonDamagexf() int {
	damage := rand.Intn(8) + 1 + dragon.Attk
	return damage
}

func dragonCriticalxf() int {
	damage := 2 * dragonDamagexf()
	// 	damage := 2 * (rand.Intn(8) + 1 + dragon.Attk)
	return damage
}

func attackxf() {
	var wgxf sync.WaitGroup
	wgxf.Add(2) // We have two concurrent tasks

	heroChan := make(chan string, 1)
	dragonChan := make(chan string, 1)

	// Hero's attack in a goroutine
	go func() {
		defer wgxf.Done()
		var result string
		if hero.ST == 0 {
			heroRoll1, heroRoll2 := rand.Intn(20)+1+hero.Attk, rand.Intn(20)+1+hero.Attk
			if heroRoll1 > heroRoll2 {
				heroRoll := heroRoll1
				if heroRoll-hero.Attk == 20 {
					damage := heroCriticalHitsxf()
					dragon.HP -= damage
					result = fmt.Sprintf("⚔ Chadwick Attacks with advantage!\n🎲 Chadwick rolled: %d\n✅ Chadwick CRITICAL HITS for %d damage, doing double  damage!!!!\n", heroRoll, damage)
				} else if heroRoll >= dragon.AC {
					damage := heroDamagexf()
					dragon.HP -= damage
					result = fmt.Sprintf("⚔ Chadwick Attacks with advantage!\n🎲 Chadwick rolled: %d\n✅ Chadwick hits for %d damage!\n", heroRoll, damage)
				} else {
					result = fmt.Sprintf("⚔ Chadwick Attacks with advantage!\n🎲 Chadwick rolled: %d\n❌ Chadwick misses!\n", heroRoll)
				}
			} else {
				heroRoll := heroRoll2
				if heroRoll >= dragon.AC {
					damage := heroDamagexf()
					dragon.HP -= damage
					result = fmt.Sprintf("⚔ Chadwick Attacks with advantage!\n🎲 Chadwick rolled: %d\n✅ Chadwick hits for %d damage!\n", heroRoll, damage)
				} else {
					result = fmt.Sprintf("⚔ Chadwick Attacks with advantage!\n🎲 Chadwick rolled: %d\n❌ Chadwick misses!\n", heroRoll)
				}
			}
		} else {
			heroRoll := rand.Intn(20) + 1 + hero.Attk
			if heroRoll-hero.Attk == 20 {
				damage := heroCriticalHitsxf()
				dragon.HP -= damage
				result = fmt.Sprintf("⚔ Chadwick Attacks!\n🎲 Chadwick rolled: %d\n✅ Chadwick CRITICAL HITS for %d damage, doing double damage!!!!\n", heroRoll, damage)
			} else if heroRoll >= dragon.AC {
				damage := heroDamagexf()
				dragon.HP -= damage
				result = fmt.Sprintf("⚔ Chadwick Attacks!\n🎲 Chadwick rolled: %d\n✅ Chadwick hits for %d damage!\n", heroRoll, damage)
			} else {
				result = fmt.Sprintf("⚔ Chadwick Attacks!\n🎲 Chadwick rolled: %d\n❌ Chadwick misses!\n", heroRoll)
			}
		}
		heroChan <- result
	}()

	// Dragon's counterattack in a goroutine
	go func() {
		defer wgxf.Done()
		var result string
		if hero.ST == 0 {
			dragonRoll1, dragonRoll2 := rand.Intn(20)+1+dragon.Attk, rand.Intn(20)+1+dragon.Attk
			if dragonRoll1 > dragonRoll2 {
				dragonRoll := dragonRoll2
				if dragonRoll-dragon.Attk == 20 {
					damage := dragonCriticalxf()
					dragon.HP -= damage
					result = fmt.Sprintf("\n🐉 Dragon counterattacks with disadvantage!\n🎲 Dragon rolled: %d\n✅ Dragon CRITICAL HITS for %d damage, doing double damage!!!!\n", dragonRoll, damage)
				} else if dragonRoll >= hero.AC {
					damage := dragonDamagexf()
					hero.HP -= damage
					fmt.Println()
					result = fmt.Sprintf("\n🐉 Dragon counterattacks with disadvantage!\n🎲 Dragon rolled: %d\n✅ Dragon hits for %d damage!\n", dragonRoll, damage)
				} else {
					fmt.Println()
					result = fmt.Sprintf("\n🐉 Dragon counterattacks with disadvantage!\n🎲 Dragon rolled: %d\n❌ Dragon misses!\n", dragonRoll)
				}
			} else {
				dragonRoll := dragonRoll1
				if dragonRoll >= hero.AC {
					damage := dragonDamagexf()
					hero.HP -= damage
					fmt.Println()
					result = fmt.Sprintf("\n🐉 Dragon counterattacks with disadvantage!\n🎲 Dragon rolled: %d\n✅ Dragon hits for %d damage!\n", dragonRoll, damage)
				} else {
					fmt.Println()
					result = fmt.Sprintf("\n🐉 Dragon counterattacks with disadvantage!\n🎲 Dragon rolled: %d\n❌ Dragon misses!\n", dragonRoll)
				}
			}
		} else {
			dragonRoll := rand.Intn(20) + 1 + dragon.Attk
			if dragonRoll-dragon.Attk == 20 {
				damage := dragonCriticalxf()
				dragon.HP -= damage
				result = fmt.Sprintf("\n🐉 Dragon counterattacks!\n🎲 Dragon rolled: %d\n✅ Dragon CRITICAL HITS for %d damage, doing double damage!!!!\n", dragonRoll, damage)
			} else if dragonRoll >= hero.AC {
				damage := dragonDamagexf()
				hero.HP -= damage
				fmt.Println()
				result = fmt.Sprintf("\n🐉 Dragon counterattacks!\n🎲 Dragon rolled: %d\n✅ Dragon hits for %d damage!\n", dragonRoll, damage)
			} else {
				fmt.Println()
				result = fmt.Sprintf("\n🐉 Dragon counterattacks!\n🎲 Dragon rolled: %d\n❌ Dragon misses!\n", dragonRoll)
			}
		}
		dragonChan <- result
	}()

	// Wait for both attacks to finish
	wgxf.Wait()
	close(heroChan)
	close(dragonChan)

	// Print results in order
	fmt.Print(<-heroChan)
	fmt.Print(<-dragonChan)

	resolvexf()
}

func resolvexf() {
	if dragon.HP <= 0 {
		fmt.Printf("\n\033[1mChadwick Strongpants has slain a dragon!\033[0m\n")
		dragon.HP = dragon.MX
		dragon.Dth++
		hero.XP += 100
	}

	if hero.HP <= 0 {
		hero.HP = hero.MX
		hero.Dth++
	}

	// Only upgrade stats if the previous XP was an increment of 300
	if hero.PrevXP%300 == 0 && hero.XP%300 != 0 {
		hero.MX += 5
		hero.HP = hero.MX
		if hero.AC < 20 {
			hero.AC++
		}
		if hero.Attk < 7 {
			hero.Attk++
		}
	}

	hero.PrevXP = hero.XP // Update previous XP to track next level-up condition

	if hero.ST == 0 {
		hero.ST = 1 // Reset stealth only after an attack
	}

	saveGamexf()
}

func displayxf() {
	var isStealth bool
	if hero.ST == 0 {
		isStealth = true
	} else if hero.ST == 1 {
		isStealth = false
	}

	fmt.Printf("\n📜 Current Stats:\n")
	fmt.Printf("🦸 Chadwick - HP: %d/%d, XP: %d, Stealth?: %t, Attack: %d, Defense: %d, Deaths: %d\n", hero.HP, hero.MX, hero.XP, isStealth, hero.Attk, hero.AC, hero.Dth)
	fmt.Printf("🐉 Dragon - HP: %d/%d, Perception: %d, Attack: %d, Defense: %d, Deaths: %d\n\n", dragon.HP, dragon.MX, dragon.PR, dragon.Attk, dragon.AC, dragon.Dth)
}

func main() {
	listFlag := flag.Bool("l", false, "List current stats")
	flag.BoolVar(listFlag, "list", false, "List current stats")
	stealthFlag := flag.Bool("s", false, "Attempt to stealth")
	flag.BoolVar(stealthFlag, "stealth", false, "Attempt to stealth")
	attackFlag := flag.Bool("a", false, "Attack the dragon")
	flag.BoolVar(attackFlag, "attack", false, "Attack the dragon")
	resetFlag := flag.Bool("r", false, "Reset the game")
	flag.BoolVar(resetFlag, "reset", false, "Reset the game")
	helpFlag := flag.Bool("h", false, "Show available flags and their descriptions")
	flag.BoolVar(helpFlag, "help", false, "Show available flags and their descriptions")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("+-------------------------------------------------+")
		fmt.Println("\t\033[1mChadwick Strongpants, Dragon Slayer\033[0m")
		fmt.Println("+-------------------------------------------------+")
		fmt.Println("\tHow many Deaths to slay the Dragon!")
		fmt.Println("      Advancement at 100xp, 400xp, 700xp, etc...")
		fmt.Println("+-------------------------------------------------+")
		fmt.Println()
		fmt.Println("Usage: dragonxf [flags]")
		fmt.Println("Flags:")
		fmt.Println("  -l, --list\tList current stats")
		fmt.Println("  -s, --stealth\tAttempt to stealth")
		fmt.Println("  -a, --attack\tAttack the dragon")
		fmt.Println("  -r, --reset\tReset the game")
		fmt.Println("  -h, --help\tShow this help message")
		fmt.Println()
	}

	flag.Parse()

	if !*listFlag && !*attackFlag && !*resetFlag && !*stealthFlag && !*helpFlag {
		flag.Usage()
		return
	}

	if *helpFlag {
		flag.Usage()
		return
	}

	initGamexf()

	if *resetFlag {
		os.Remove(getSaveFilePathxf())
		fmt.Println()
		fmt.Println("Game reset!")
		initGamexf()
		displayxf()
		return
	}

	if *listFlag {
		displayxf()
	}

	if *stealthFlag {
		stealthRoll := rand.Intn(20) + 1 + hero.Attk
		if stealthRoll >= dragon.PR {
			hero.ST = 0
			fmt.Println()
			fmt.Printf("🎲 Stealth roll: %d vs Dragon Perception: %d\n", stealthRoll, dragon.PR)
			fmt.Println()
			fmt.Println("Chadwick is stealthed")
		} else {
			hero.ST = 1
			fmt.Println()
			fmt.Printf("🎲 Stealth roll: %d vs Dragon Perception: %d\n", stealthRoll, dragon.PR)
			fmt.Println()
			fmt.Println("Chadwick failed to stealth")
		}
		saveGamexf()
		displayxf()
		// 		return
	}

	if *attackFlag {
		attackxf()
		displayxf()
	}
}
