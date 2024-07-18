package main

import (
	"fmt"
	"sort"
	"time"
)

var (
	HULLS   = []string{"Wasp", "Hornet", "Hopper", "Viking", "Hunter", "Crusader", "Paladin", "Dictator", "Ares", "Titan", "Mammoth"}
	TURRETS = []string{"Firebird", "Freeze", "Isida", "Tesla", "Hammer", "Twins", "Ricochet", "Smoky", "Striker", "Vulcan", "Thunder", "Scorpion", "Railgun", "Magnum", "Gauss", "Shaft"}
	//DRONES  = []string{"Crisis", "Brutus", "Saboteur", "Trickster", "Mechanic", "Booster", "Defender", "Hyperion"}
)

type Datastamp struct {
	Timestamp      time.Time
	Name           string
	Rank           int
	Kills          int
	Deaths         int
	EarnedCrystals int
	GearScore      int
	Hulls          map[string]Thing
	Turrets        map[string]Thing
	Drones         map[string]Thing
	SuppliesUsed   map[string]int
}

type Thing struct {
	ScoreEarned int
	TimePlayed  int
}

type kv struct {
	Key   string
	Value Thing
}

func sortedSliceByScore(m map[string]Thing) []kv {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value.ScoreEarned > ss[j].Value.ScoreEarned
	})
	return ss
}

func msToHours(microseconds int) int {
	return microseconds / (1000 * 60 * 60)
}

func (d *Datastamp) Print() {
	line := "------------------------------------------------\n"
	fmt.Print("Time:", d.Timestamp)
	fmt.Print("\nName: ", d.Name, "\nGS: ", d.GearScore)
	fmt.Print("\nKills: ", d.Kills, "\nDeaths: ", d.Deaths, "\nK/D:", float32(d.Kills)/float32(d.Deaths))

	fmt.Print("\nTurret\t\tScore\t\tTime played, h\n", line)
	turrets := sortedSliceByScore(d.Turrets)
	for _, a := range turrets[:5] {
		fmt.Print(a.Key, "\t\t", a.Value.ScoreEarned, "\t\t", msToHours(a.Value.TimePlayed), "\n")
	}

	hulls := sortedSliceByScore(d.Hulls)
	fmt.Print("\nHull\t\tScore\t\tTime played, h\n", line)
	for _, a := range hulls[:5] {
		fmt.Print(a.Key, "\t\t", a.Value.ScoreEarned, "\t\t", msToHours(a.Value.TimePlayed), "\n")
	}

}

func (d *Datastamp) Store(data ResponseWrapper) {
	r := data.Response
	d.Name, d.Rank, d.Kills, d.Deaths, d.EarnedCrystals, d.GearScore =
		r.Name, r.Rank, r.Kills, r.Deaths, r.EarnedCrystals, r.GearScore

	d.Hulls = make(map[string]Thing)
	d.Turrets = make(map[string]Thing)
	d.Drones = make(map[string]Thing)
	d.SuppliesUsed = make(map[string]int)

	for _, hull := range HULLS {
		d.Hulls[hull] = Thing{0, 0}
	}

	for _, a := range r.HullsPlayed {
		hull := d.Hulls[a.Name]
		hull.TimePlayed += a.TimePlayed
		hull.ScoreEarned += a.ScoreEarned
		d.Hulls[a.Name] = hull
	}

	for _, turret := range TURRETS {
		d.Turrets[turret] = Thing{0, 0}
	}
	for _, a := range r.TurretsPlayed {
		turret := d.Turrets[a.Name]
		turret.TimePlayed += a.TimePlayed
		turret.ScoreEarned += a.ScoreEarned
		d.Turrets[a.Name] = turret
	}

	for _, a := range r.DronesPlayed {
		d.Drones[a.Name] = Thing{a.ScoreEarned, a.TimePlayed}
	}

	for _, a := range r.SuppliesUsage {
		d.SuppliesUsed[a.Name] = a.Usages
	}

	d.Timestamp = time.Now().Truncate(time.Hour)

}
