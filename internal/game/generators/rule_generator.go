package generators

import (
	"GolangGame251130/internal/game"
)

// RuleBasedGenerator uses multiple strategies to generate levels
type RuleBasedGenerator struct {
	nextSpawnX   float32
	lastFoodX   float32
	
	currentRule GenRule
	itemsLeft   int
}

type GenRule func(g *game.Game, startX float32) (items []SpawnDef, nextOffset float32)

type SpawnDef struct {
	Type    int
	LineIdx int
	Lane    int
	OffsetX float32
}

func NewRuleBasedGenerator() *RuleBasedGenerator {
	return &RuleBasedGenerator{
		nextSpawnX: 300,
		lastFoodX:  0,
		itemsLeft:  0,
	}
}

func (r *RuleBasedGenerator) ShouldSpawn(g *game.Game) bool {
	cameraX := g.GetCameraX()
	spawnThreshold := cameraX + 300 
	return r.nextSpawnX < spawnThreshold
}

func (r *RuleBasedGenerator) SpawnItem(g *game.Game) {
	// Rule selection logic
	// We generate a "chunk" of items based on a selected rule
	
	// 1. Check Food
	level := g.GetLevel()
	foodInterval := float32(125 + (level-1)*25)
	
	if r.nextSpawnX - r.lastFoodX > foodInterval {
		r.spawnOne(g, 2, game.Intn(2), game.Intn(2), 0)
		r.lastFoodX = r.nextSpawnX
		r.nextSpawnX += 60 + float32(game.Intn(40))
		return
	}

	// 2. Select Rule randomly
	ruleIdx := game.Intn(4)
	
	var items []SpawnDef
	var width float32

	switch ruleIdx {
	case 0: // Static Pattern (Reuse simple ones)
		items, width = ruleStatic()
	case 1: // Zipper
		items, width = ruleZipper()
	case 2: // Tunnel
		items, width = ruleTunnel()
	case 3: // Random Field
		items, width = ruleRandomField()
	}

	// 3. Spawn Items
	lines := g.GetLines()
	for _, itemDef := range items {
		spawnX := r.nextSpawnX + itemDef.OffsetX
		
		targetLine := lines[itemDef.LineIdx]
		var newItem game.Item
		
		switch itemDef.Type {
		case 1: 
			newItem = game.NewRock(targetLine, spawnX, itemDef.Lane)
		case 3:
			newItem = game.NewGoldRock(targetLine, spawnX, itemDef.Lane)
		case 4:
			newItem = game.NewHardRock(targetLine, spawnX, itemDef.Lane)
		}
		
		if newItem != nil {
			targetLine.AddItem(newItem)
		}
	}

	r.nextSpawnX += width
}

func (r *RuleBasedGenerator) spawnOne(g *game.Game, typeID, lineIdx, lane int, offset float32) {
	lines := g.GetLines()
	targetLine := lines[lineIdx]
	spawnX := r.nextSpawnX + offset
	
	var newItem game.Item
	switch typeID {
	case 1: newItem = game.NewRock(targetLine, spawnX, lane)
	case 2: newItem = game.NewFood(targetLine, spawnX, lane)
	case 3: newItem = game.NewGoldRock(targetLine, spawnX, lane)
	case 4: newItem = game.NewHardRock(targetLine, spawnX, lane)
	}
	
	if newItem != nil {
		targetLine.AddItem(newItem)
	}
}

func (r *RuleBasedGenerator) OnCoordinateReset(offset float32) {
	r.nextSpawnX -= offset
	r.lastFoodX -= offset
}

// Rules Implementation

func ruleStatic() ([]SpawnDef, float32) {
	// Simple static pattern
	// 50% Horizontal or Vertical double
	if game.Intn(2) == 0 {
		// Double Vertical
		return []SpawnDef{
			{1, 0, 0, 0},
			{1, 1, 0, 0},
		}, 60
	} else {
		// Side by Side
		return []SpawnDef{
			{1, 0, 0, 0},
			{1, 0, 1, 0},
		}, 60
	}
}

func ruleZipper() ([]SpawnDef, float32) {
	// Up Down Up Down
	count := 3 + game.Intn(3) // 3 to 5 items
	items := make([]SpawnDef, count)
	
	startLane := game.Intn(2)
	
	for i := 0; i < count; i++ {
		lineIdx := i % 2
		items[i] = SpawnDef{
			Type: 1, 
			LineIdx: lineIdx, 
			Lane: startLane, 
			OffsetX: float32(i * 30),
		}
	}
	
	return items, float32(count * 30 + 30)
}

func ruleTunnel() ([]SpawnDef, float32) {
	// Block one lane completely for a while
	length := 3 + game.Intn(3) // 3 to 5 segments
	blockedLane := game.Intn(2) // 0 or 1
	
	items := make([]SpawnDef, 0, length*2)
	
	for i := 0; i < length; i++ {
		// Place rocks on both lines in the blocked lane
		items = append(items, SpawnDef{1, 0, blockedLane, float32(i * 32)})
		items = append(items, SpawnDef{1, 1, blockedLane, float32(i * 32)})
	}
	
	return items, float32(length * 32 + 50)
}

func ruleRandomField() ([]SpawnDef, float32) {
	// Place rocks randomly but ensure at least one path is open
	// We do this column by column
	cols := 3 + game.Intn(4) 
	items := make([]SpawnDef, 0, cols*2)
	
	for i := 0; i < cols; i++ {
		// 4 possible slots: (0,0), (0,1), (1,0), (1,1)
		// We pick 1 to 3 slots to fill
		
		// Simple approach: Pick one slot to remain OPEN
		openLine := game.Intn(2)
		openLane := game.Intn(2)
		
		for line := 0; line < 2; line++ {
			for lane := 0; lane < 2; lane++ {
				if line == openLine && lane == openLane {
					continue // Safe spot
				}
				
				// 50% chance to spawn rock in other spots
				if game.Intn(2) == 0 {
					// 10% GoldRock
					typeID := 1
					if game.Intn(10) == 0 { typeID = 3 }
					
					items = append(items, SpawnDef{typeID, line, lane, float32(i * 40)})
				}
			}
		}
	}
	
	return items, float32(cols * 40 + 40)
}
