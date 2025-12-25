package generators

import (
	"GolangGame251130/internal/game"
)

// ChunkParams holds the parameters for the current level chunk.
type ChunkParams struct {
	LaneSwitchChance int // 0-20%
	LineSwitchChance int // 0-10% (Probability to switch target pickaxe owner)
	RockSpawnRate    int // 30-80% (Probability of rock on target path)
	FoodSpawnRate    int // 0-30% (Probability of food on non-target path)
	ObstacleDensity  int // 20-60% (Probability of obstacle off-path)
}

// PathGenerator handles level generation with a specific path logic.
// Grid size: 24px
type PathGenerator struct {
	nextSpawnX         float32
	pathLanes          []int // Current safe lane for each line (0 or 1)
	switchSafety       []int // Counter for safety duration after switch
	targetPickaxeOwner int   // Which player *should* have the pickaxe (0 or 1)

	// Chunk management
	chunkRemaining int         // Number of grids remaining in current chunk
	currentChunk   ChunkParams // Current chunk parameters
}

func NewPathGenerator() *PathGenerator {
	gen := &PathGenerator{
		nextSpawnX:         400,
		pathLanes:          []int{0, 1}, // Initial lanes
		switchSafety:       []int{0, 0},
		targetPickaxeOwner: 0,
		chunkRemaining:     0, // Will trigger new chunk immediately
	}
	return gen
}

func (g *PathGenerator) ShouldSpawn(gameInst *game.Game) bool {
	// Spawn ahead of camera
	spawnThreshold := gameInst.GetCameraX() + 320
	return g.nextSpawnX < spawnThreshold
}

func (g *PathGenerator) SpawnItem(gameInst *game.Game) {
	lines := gameInst.GetLines()
	gridSize := float32(24.0)

	// --- 0. Update Chunk State ---
	if g.chunkRemaining <= 0 {
		// Start new chunk
		g.chunkRemaining = game.RandomIntn(16) + 15 // 15 to 30 grids

		level := gameInst.GetLevel()

		// Randomize parameters scaling with level
		// Level 1 -> 10 Scaling

		g.currentChunk.LaneSwitchChance = getScaledValue(level, 0, 10, 5, 30)
		g.currentChunk.LineSwitchChance = getScaledValue(level, 0, 5, 5, 20)
		g.currentChunk.RockSpawnRate = getScaledValue(level, 10, 20, 30, 60)
		g.currentChunk.FoodSpawnRate = getScaledValue(level, 0, 20, 0, 10)
		g.currentChunk.ObstacleDensity = getScaledValue(level, 20, 40, 30, 80)
	}
	g.chunkRemaining--
	params := g.currentChunk

	// --- 1. Update Generator State (Lane switches / Pickaxe Target switch) ---

	// Chance to switch Target Pickaxe Owner
	// Use LineSwitchChance from chunk params
	if game.RandomIntn(100) < params.LineSwitchChance {
		g.targetPickaxeOwner = 1 - g.targetPickaxeOwner
	}

	for i := range g.pathLanes {
		// Update safety counter
		if g.switchSafety[i] > 0 {
			g.switchSafety[i]--
		}

		// Probability to switch lane for this line
		// Use LaneSwitchChance from chunk params
		// Only switch if not currently in safety period
		if g.switchSafety[i] == 0 {
			if game.RandomIntn(100) < params.LaneSwitchChance {
				g.pathLanes[i] = 1 - g.pathLanes[i]
				g.switchSafety[i] = 2 // "Treat 2 grids as path" -> Safety for 2 grids
			}
		}
	}

	// --- 2. Spawn Items ---

	for lineIdx := range lines {
		line := lines[lineIdx]
		pathLane := g.pathLanes[lineIdx]
		isSafety := g.switchSafety[lineIdx] > 0
		isTargetOwner := (lineIdx == g.targetPickaxeOwner)

		// Generate for both lanes in this line (0 and 1)
		for lane := 0; lane < 2; lane++ {
			// Calculate Spawn X with Variance: 0 ~ 7
			variance := float32(game.RandomIntn(8)) // 0 to 7
			spawnX := g.nextSpawnX + variance

			isPath := (lane == pathLane)

			// --- Safety / Transition Zone ---
			if isSafety {
				// During safety switch, keep area clear.
				continue
			}

			// --- Path Logic ---
			if isPath {
				if isTargetOwner {
					// Target Path: Spawn Rock/GoldRock based on RockSpawnRate
					r := game.RandomIntn(100)
					if r < params.RockSpawnRate {
						if game.RandomIntn(100) < 10 {
							line.AddItem(game.NewGoldRock(line, spawnX, lane))
						} else {
							line.AddItem(game.NewRock(line, spawnX, lane))
						}
					} else {
					}
				} else {
					if game.RandomIntn(100) < params.FoodSpawnRate {
						line.AddItem(game.NewFood(line, spawnX, lane))
					}
				}
				continue
			}

			if game.RandomIntn(100) < params.ObstacleDensity {
				r := game.RandomIntn(100)
				if r < 40 {
					line.AddItem(game.NewRock(line, spawnX, lane))
				} else if r < 70 {
					line.AddItem(game.NewHardRock(line, spawnX, lane))
				} else if r < 85 {
					line.AddItem(game.NewGoldRock(line, spawnX, lane))
				} else {
					line.AddItem(game.NewFood(line, spawnX, lane))
				}
			}
		}
	}

	g.nextSpawnX += gridSize
}

func (g *PathGenerator) OnCoordinateReset(offset float32) {
	g.nextSpawnX -= offset
}

func getScaledValue(level, minV, maxV, minTarget, maxTarget int) int {
	if level < 1 {
		level = 1
	}
	if level > 10 {
		level = 10
	}

	t := float32(level-1) / 9.0 // 0.0 at Lv1, 1.0 at Lv10

	// Lerp for min and max of the random range
	currentMin := float32(minV) + (float32(minTarget)-float32(minV))*t
	currentMax := float32(maxV) + (float32(maxTarget)-float32(maxV))*t

	iMin := int(currentMin)
	iMax := int(currentMax)

	if iMin > iMax {
		iMin, iMax = iMax, iMin
	}

	// Ensure range is valid for RandomIntn
	rangeSize := iMax - iMin + 1
	if rangeSize <= 0 {
		return iMin
	}

	return iMin + game.RandomIntn(rangeSize)
}
