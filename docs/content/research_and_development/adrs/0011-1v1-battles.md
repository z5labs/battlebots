---
title: "[0011] 1v1 Battles"
description: >
  Arena concept definition and battle properties for 1v1 game mode
type: docs
weight: 11
category: "strategic"
status: "accepted"
date: 2025-12-09
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions (frameworks, auth strategies, cross-cutting patterns)
- user-journey: Solutions for specific user journey problems (feature implementation approaches)
- api-design: API endpoint design decisions (pagination, filtering, bulk operations)
-->

## Context and Problem Statement

BattleBot 1v1 battles require a formal definition of what constitutes an arena and how battle instances are configured. Currently, the BattleBot Universe (ADR-0005) defines the universal topological space (R² Euclidean coordinates with Cartesian system) and fixed boundary dimensions ([-50, 50] × [-50, 50]). However, this establishes a single static arena rather than a configurable system that enables tactical variety through property selection.

This ADR introduces the **Arena** concept as the bridge between the universal BattleBot Universe and concrete 1v1 battle instances.

### Arena Definition

An **Arena** is a configured instance of the BattleBot Universe, defined by a combination of selectable and fixed properties:

- **Arena = Bounded Region (from ADR-0005) + Biome (from ADR-0010) + Battle Configuration**

The arena concept enables:
- Multiple different battle configurations using the same underlying spatial and physics systems
- Pre-battle tactical selection through property choices
- Gameplay variety without changing universal topology or physics laws
- Clear separation: ADR-0005 defines universal spatial foundation; ADR-0011 defines how battle instances use that foundation

### Boundary Extraction and Architectural Shift

ADR-0005 "BattleBot Universe Topological Properties" defined boundary as Property 4 with specific dimensions:
- **Rectangular boundary**: [-50, 50] × [-50, 50] (100×100 units)
- **Mathematical definition**: Closed manifold M = {(x,y) ∈ R² : -50 ≤ x ≤ 50, -50 ≤ y ≤ 50}
- **Boundary**: ∂M = {(x,y) : x=±50 or y=±50}

This universal boundary definition was appropriate for establishing the spatial system's mathematical foundation. However, hardcoding a single boundary dimension:
1. **Prevents arena variety** - All battles use identical dimensions
2. **Couples topology to game mode** - Changing boundaries requires changing universal topology
3. **Limits future game modes** - Different game modes need different arena sizes
4. **Reduces pre-battle strategy** - Removes player choice in arena configuration

**Architectural Solution**: Transfer boundary ownership from ADR-0005 (universal topology) to ADR-0011 (battle instance configuration).

- **Stays in ADR-0005**: R² Euclidean topological space, Cartesian coordinate system, spatial properties
- **Moves to ADR-0011**: Boundary dimensions, shapes, and selection mechanism
- **Benefit**: Enables selectable boundaries while preserving topological rigor

### Problem Statement

We must formalize the 1v1 battle type by defining five core properties:

1. **Biome Selection** - How is arena geography determined?
2. **Boundary Configuration** - How are arena dimensions determined?
3. **Visibility System** - What information do bots receive?
4. **Start Positioning** - Where do bots begin the battle?
5. **Win Conditions** - What determines battle outcome?

Each property has multiple options with distinct gameplay and implementation implications. This ADR evaluates options for each property and establishes the foundation for 1v1 battles as the game's initial battle type.

### Relationship to Other ADRs

- **ADR-0005** (BattleBot Universe Topological Properties): Defines universal R² space, Cartesian coordinates, spatial metrics. ADR-0011 applies this space to concrete battle instances.
- **ADR-0006** (BattleBot Universe Physics Laws): Defines variable friction F_friction = μ(position) × M × |v|. ADR-0011 applies physics laws within configured arenas.
- **ADR-0010** (Biome Definition Framework): Establishes Geography as the biome property, where μ(position) comes from terrain. ADR-0011 enables biome selection for arenas.
- **ADR-0004** (Bot to Battle Server Interface): gRPC protocol for state exchange. ADR-0011 defines configuration transmitted via this protocol.

## Decision Drivers

* **Biome Integration** - Arena properties must enable ADR-0010 biome framework to create tactical variety through geography
* **Boundary Flexibility** - Different battle contexts should support different arena sizes without changing universal topology
* **Gameplay Variety** - Property-based configuration should enable diverse tactical environments and strategic choices
* **Implementation Feasibility** - Initial properties should be implementable in MVP, with complexity deferred to future ADRs
* **Extensibility** - Five-property framework should accommodate future additions (dynamic boundaries, new visibility modes) without breaking existing system
* **Balanced Gameplay** - Arena properties must create fair, competitive experiences that reward skill and adaptability
* **Player Agency** - Pre-battle property selection should provide meaningful strategic choices that affect gameplay
* **Determinism** - Battle outcomes must be reproducible given identical arena configuration and initial random seed
* **Physics Integration** - Arena configuration must work seamlessly with variable friction from ADR-0006 and movement mechanics from ADR-0007
* **Spatial Consistency** - Arena must respect 2D Euclidean space and Cartesian coordinates from ADR-0005
* **Topology Independence** - Selecting arena properties should not require changing universal spatial system properties
* **Engagement Guarantee** - Arena configuration should ensure meaningful interaction and prevent indefinite evasion

## Considered Options

This ADR evaluates five distinct battle properties:

* **Property 1: Biome Selection** - How is terrain/friction determined for the arena?
* **Property 2: Boundary Configuration** - What are the arena's shape and dimensions?
* **Property 3: Visibility System** - What information do bots receive about opponents?
* **Property 4: Start Positioning** - Where are bots placed when battle begins?
* **Property 5: Win Conditions** - What determines how battle concludes and who wins?

Each property has multiple options with distinct advantages and limitations. The outcome section details chosen options and integration rationale.

## Decision Outcome

### Chosen Arena Configuration

1. **Biome**: Selectable (user chooses from available options)
2. **Boundary**: Selectable (user chooses from available configurations)
3. **Visibility**: Full (all bots see complete information)
4. **Start Positioning**: Random (bots placed randomly within arena)
5. **Win Conditions**: Health Depletion + Timeout + Disconnect

This five-property framework creates the **Arena** as a configured instance of the BattleBot Universe. Different property selections produce different arenas with distinct tactical characteristics.

### Arena Concept: Mathematical Foundation

An **Arena** (A) is formally defined as:

**A = (M, B, μ(·), Σ, V, C)**

Where:
- **M** = Arena manifold (boundary configuration) - a bounded region of R²
- **B** = Boundary ∂M - the frontier of the arena region
- **μ(·)** = Position-dependent friction coefficient function (from selected biome)
- **Σ** = Starting positions (initial bot placements)
- **V** = Visibility configuration (information available to bots)
- **C** = Win conditions (rules determining battle conclusion)

**Key Properties**:
- Inherits 2D Cartesian coordinate system (x, y) from ADR-0005
- Inherits physics laws and movement mechanics from ADR-0006, ADR-0007
- Respects variable friction through μ(position) provided by biome selection
- Enables deterministic reproducibility through seeded randomness

**Instance Example**:
```
Arena_Desert_100x100_Full_Random = (
  M = {(x,y) ∈ R² : -50 ≤ x ≤ 50, -50 ≤ y ≤ 50},        // Rectangular, 100×100 units
  B = {(x,y) : x=±50 or y=±50},                          // Rectangular boundary
  μ(·) = Desert_Biome_Friction_Map,                       // Desert friction topology
  Σ = RandomPlacement(arena=A, seed=s, separation=20),   // Random start positions
  V = FullVisibility,                                     // All bots see everything
  C = [HealthDepletion, Timeout(300s), Disconnect(30s)]  // Multi-path win conditions
)
```

### Property 1: Biome

**Chosen Option**: Option 1.2 - Selectable

#### Definition

User selects a biome from available options before battle initialization. The selected biome determines the arena's geography property, including:
- Friction coefficient map μ(position): Arena → R⁺
- Visual texture representation
- Tactical positioning differences

#### Biome Integration

Biomes operate per ADR-0010 "Biome Definition Framework":

- **Geography Property**: Biomes define arena topology through position-dependent friction
- **Friction Function**: Each biome specifies μ(position) values across arena regions
- **Tactical Differentiation**: Different friction topologies create distinct movement characteristics:
  - **Low friction zones** (ice): Increased speed, reduced maneuverability
  - **High friction zones** (mud): Reduced speed, increased control
  - **Mixed zones**: Complex tactical positioning opportunities
- **Visual Clarity**: Texture representation communicates friction properties intuitively

#### Initial Biome Options

At ADR-0011 acceptance, no concrete biomes are implemented. Future ADRs will define specific biomes following ADR-0010 framework:

- **Example biomes** (from ADR-0010):
  - **Desert**: Sand regions (μ = 0.6), rocky outcrops (μ = 0.4), gravel paths (μ = 1.0)
  - **Arctic**: Ice plains (μ = 0.2), snow banks (μ = 0.5), frozen rock (μ = 0.4)
  - **Forest**: Grass clearings (μ = 0.5), muddy swamps (μ = 1.2), dirt paths (μ = 0.6)

**Implementation Note**: MVP may initially support single biome (uniform friction μ(x,y) = constant) until concrete biome ADRs define geography specifics.

#### Physics Integration

Friction from selected biome integrates with ADR-0006 physics:

**Friction Formula**: F_friction = μ(position) × M × |v|

Where:
- μ(position) = position-dependent coefficient from selected biome
- M = bot mass (from ADR-0008 characteristics)
- |v| = velocity magnitude
- Affects terminal velocity: v_terminal = F_thrust / (μ(position) × M)

**Gameplay Impact**: Bot movement characteristics vary by terrain location:
- Same bot in ice zone (low μ) accelerates faster, travels farther
- Same bot in mud zone (high μ) accelerates slower, stops quicker
- Forces positioning decisions: faster zones for retreat, control zones for combat

#### Rationale for Selectable Biomes

* **Framework Validation** - Operationalizes ADR-0010 biome framework in actual gameplay; validates biome concept through practice
* **Tactical Variety** - Different friction topologies create distinct tactical environments; forces bot adaptability
* **Strategic Selection** - Biome choice becomes pre-battle strategic decision based on bot capabilities and opponent analysis
* **Adaptability Testing** - Forces bot developers to handle variable friction across arena rather than assuming constant
* **Gameplay Depth** - Geography-based positioning creates additional layer of tactical complexity beyond direct combat
* **Future Extensibility** - Framework enables adding new biomes without changing core arena system

#### Future Considerations

- **Concrete Biome ADRs** (future): Separate ADRs will define Desert, Arctic, Forest, and additional biomes with specific friction topologies
- **Biome Selection UI**: Game implementation may randomize, present curated choices, or enable full selection
- **Biome Rotation**: Tournament modes might rotate biomes for fairness or focus on single biome for consistency
- **Dynamic Biomes** (future): Climate property (deferred in ADR-0010) could add weather effects modifying friction

### Property 2: Boundary

**Chosen Option**: Option 2.2 - Selectable

#### Definition

User selects boundary configuration from available options before battle initialization. Boundary defines:
- **Shape** - Rectangular, circular, hexagonal, or other geometric form
- **Dimensions** - Size of bounded region (e.g., 100×100 units)
- **Edge Behavior** - Collision mechanics at boundaries

#### Boundary Extraction from ADR-0005

ADR-0005 Property 4 established the universal boundary as part of topological properties. This ADR transfers boundary ownership:

**Transition Rationale**:
1. **ADR-0005 correctly defined** universal R² space with fixed boundary for establishing spatial foundation
2. **Hardcoding prevented** arena variety - single boundary couples topology to game mode
3. **Future flexibility needed** - Different game modes require different arena sizes and shapes
4. **Clear separation** - Universal topology (ADR-0005) separate from battle instance configuration (ADR-0011)

**What Remains in ADR-0005**:
- 2D Euclidean topological space R²
- Cartesian coordinate system and origin definition
- Spatial metrics and mathematical properties
- Coordinate system transformation rules

**What Moves to ADR-0011**:
- Specific boundary dimensions (e.g., 100×100)
- Boundary shapes (rectangular, circular, etc.)
- Boundary selection mechanism
- Per-boundary configuration options

**Preservation of Mathematical Rigor**:
- Mathematical definitions from ADR-0005 are preserved and reused
- Manifold notation (M, ∂M) maintained in ADR-0011 context
- Boundary mechanics and collision resolution detailed here

#### Initial Boundary Option: Rectangular 100×100 Units (Default)

**Shape**: Rectangular (aligned with x and y axes)

**Mathematical Definition**:
- **Arena Manifold**: M = {(x,y) ∈ R² : -50 ≤ x ≤ 50, -50 ≤ y ≤ 50}
- **Boundary**: ∂M = {(x,y) : x = -50 or x = 50 or y = -50 or y = 50}
- **Boundary Orientation**: Closed manifold with boundary
- **Interior**: M° = {(x,y) ∈ R² : -50 < x < 50, -50 < y < 50}
- **Metric**: Standard Euclidean d(p₁, p₂) = √((x₂−x₁)² + (y₂−y₁)²)

**Geometric Properties**:
- **Arena Center**: (0, 0)
- **Arena Dimensions**: 100 units × 100 units
- **Corners**: (-50, -50), (50, -50), (-50, 50), (50, 50)
- **X-axis Range**: -50 to +50 units
- **Y-axis Range**: -50 to +50 units
- **Perimeter**: 400 units
- **Area**: 10,000 square units

**Symmetry**: Rectangular boundary with origin at center provides symmetric access to all four quadrants for both bots (fair and balanced).

#### Movement Clamping Mechanics

Any bot movement command attempting to place bot outside arena boundary is clamped to valid position:

**Clamping Algorithm**:
```
x_final = clamp(x_attempted, -50 + radius, 50 - radius)
y_final = clamp(y_attempted, -50 + radius, 50 - radius)
```

Where:
- `x_attempted` = requested x position from movement command
- `radius` = bot collision radius (typically 2 units)
- `clamp(v, min, max)` = constrains v to [min, max] range

**Behavior**:
- Bot position clamped to boundary edge if movement exceeds limits
- No wrapping (exiting right side doesn't place bot on left side)
- Bot can position exactly on boundary line
- Ensures bot never overlaps arena boundary

#### Collision Detection and Resolution

**Rectangular Boundary Collision Detection**:

Simple comparison checks determine wall collision:

```
Left wall collision:   x - radius < -50
Right wall collision:  x + radius > 50
Bottom wall collision: y - radius < -50
Top wall collision:    y + radius > 50
```

**Elastic Collision Resolution** (per ADR-0006):

Wall treated as infinite-mass object. Velocity components are transformed:

**Collision with vertical wall** (left or right):
- x-component velocity: v_x_final = -v_x_initial (reversed)
- y-component velocity: v_y_final = v_y_initial (preserved)
- Example: Bot moving right hits right wall → bounces left with same rightward energy redirected leftward

**Collision with horizontal wall** (top or bottom):
- x-component velocity: v_x_final = v_x_initial (preserved)
- y-component velocity: v_y_final = -v_y_initial (reversed)
- Example: Bot moving up hits top wall → bounces down with same upward energy redirected downward

**Corner collision** (both x and y exceed):
- Both velocity components reversed
- x-component velocity: v_x_final = -v_x_initial
- y-component velocity: v_y_final = -v_y_initial

**No Damage**: Boundary collisions do not reduce bot Health (colliding with walls is free).

**Friction Interaction**: Variable friction μ(position) applies to deceleration during boundary contact (friction reduces velocity magnitude).

#### Future Boundary Options

**Option 2.3: Rectangular 150×150 units** (larger arena)
- **Use case**: Extended battles, emphasis on maneuvering over direct combat
- **Manifold**: M = {(x,y) ∈ R² : -75 ≤ x ≤ 75, -75 ≤ y ≤ 75}
- **Engagement dynamics**: Lower encounter frequency, more room for positioning strategies

**Option 2.4: Rectangular 75×75 units** (smaller arena)
- **Use case**: Fast-paced combat, emphasis on direct engagement
- **Manifold**: M = {(x,y) ∈ R² : -37.5 ≤ x ≤ 37.5, -37.5 ≤ y ≤ 37.5}
- **Engagement dynamics**: Higher encounter frequency, less room for evasion

**Option 2.5: Circular boundary** (future ADR)
- **Shape**: Circle centered at origin
- **Advantage**: Uniform distance from center, no corner camping strategies
- **Challenge**: More expensive boundary checks (distance calculations vs comparisons)
- **Implementation**: Requires O(1) distance check: d = √(x² + y²); collide if d > radius
- **Status**: Deferred - rectangular implementation first, then circular as variant

**Option 2.6: Hexagonal boundary** (future ADR)
- **Shape**: Regular hexagon centered at origin
- **Advantage**: Balanced symmetry with six edges and corners
- **Challenge**: More complex collision detection (six edge checks)
- **Gameplay**: Creates different positional dynamics than rectangular (curved transition)
- **Status**: Deferred - more complex than rectangular and circular

#### Rationale for Selectable Boundaries

* **Arena Variety** - Different sizes create different engagement dynamics: small arenas force frequent interaction, large arenas reward maneuvering
* **Game Mode Flexibility** - Tournament mode might use 100×100, quick play 75×75, tactical 150×150 without changing universal topology
* **No Hardcoding** - Avoids coupling single boundary dimension to spatial system; boundary becomes data not code
* **Balance Tuning** - Arena size can be adjusted empirically through playtesting without modifying ADRs
* **Engagement Dynamics** - Smaller boundaries increase combat frequency and timeout pressure; larger boundaries reward avoidance and positioning
* **Future Game Modes** - Different game modes (team battles, tournaments, special events) can use different boundaries

#### Integration with Other Systems

**Physics (ADR-0006)**:
- Elastic collisions apply at boundary
- Friction from biome applies during wall contact
- No additional forces required for boundary interaction

**Movement (ADR-0007)**:
- Thrust-based movement respected but clamped at boundaries
- Movement commands execute within arena bounds
- Velocity reversal on collision is part of physics resolution

**Biome (ADR-0010)**:
- Friction map μ(position) defined over boundary region
- Different boundaries may have different friction topology layouts
- Boundary size doesn't affect friction coefficient values

### Property 3: Visibility

**Chosen Option**: Option 3.1 - Full Visibility

#### Definition

All bots have complete visibility of all other bots at all times. No fog of war, visibility radius limitations, or occlusion mechanics.

**Information Available**:
- Complete position (x, y) of all opponents
- Velocity and direction of all opponents
- Health status of all opponents
- Equipment loadout of all opponents
- State of ongoing actions and cooldowns

#### Implementation

**State Protocol** (ADR-0004):
- Each game tick, battle server broadcasts complete state to all bots
- State includes positions of all bots (no filtering or occlusion)
- Bots receive complete information simultaneously

**Visibility Calculations**: None required
- No visibility radius checks
- No line-of-sight occlusion
- No state filtering or partial updates
- No information asymmetry between bots

**Determinism**:
- Complete information enables pure strategy without luck
- Same game state + same bot logic → same decisions
- Enables perfect bot logic and strategic play

#### Rationale for Full Visibility

* **Simplest Implementation** - MVP requires no visibility system; complete state updates sufficient
* **Deterministic Gameplay** - All information available enables pure strategic decisions without information asymmetry
* **Accessibility** - Bot developers need not implement visibility queries or partial information handling
* **Baseline Validation** - Core combat mechanics can be tested without visibility complexity
* **Reproducibility** - Complete information ensures identical replays with same seed

#### Alternative: Constant Fog of War (Deferred)

**Definition**: Bots see only within fixed radius around their position (e.g., 30 units).

**Gameplay Implications**:
- Creates reconnaissance gameplay and scouting strategy
- Enables ambush tactics and positioning for information advantage
- Requires visibility calculation: can bot at (x₁, y₁) see opponent at (x₂, y₂)?
- Distance check: d = √((x₂−x₁)² + (y₂−y₁)²); visible if d ≤ radius

**Implementation Complexity**:
- Requires visibility radius system architecture
- State filtering: only broadcast visible bots
- Sensor equipment (ADR-0009) deferred; enables future sensor mechanics
- Computational cost increases with bot count

**Status**: Deferred to future ADR after full visibility baseline validated through playtesting

**Future ADR**: "Constant Fog of War Visibility System" will define visibility mechanics comprehensively

#### Alternative: Vanishing Fog of War (Deferred)

**Definition**: Bots see within radius, but explored areas remain visible indefinitely (fog "vanishes" once explored).

**User Note**: Specified as providing "very nice" transition from full visibility to constant fog of war.

**Gameplay Implications**:
- Creates map control mechanics and information persistence
- Rewards positional control and exploration
- Bots can "remember" areas they've previously explored
- Information advantage grows with exploration during battle
- Strategic pacing: early scouting provides mid-game information advantage

**Implementation Complexity** (Highest):
- Requires per-bot visibility memory (which areas explored by which bot)
- Fog state tracking: for each grid cell, is it visible now or was visible?
- State synchronization: must communicate explored areas to bot
- Memory management: growing visibility state over battle duration
- Visualization: show explored areas differently from unseen areas

**Status**: Deferred pending full visibility and constant fog of war validation

**Future ADR**: "Vanishing Fog of War Visibility System" will define advanced visibility mechanics after simpler alternatives validated

#### Future Visibility Roadmap

1. **ADR-0011** (Current): Full visibility
2. **Future ADR-A**: Constant fog of war with configurable radius
3. **Future ADR-B**: Vanishing fog of war with explored area persistence
4. **Future Game Modes**: Different visibility modes for different game types

### Property 4: Start Positioning

**Chosen Option**: Option 4.2 - Random

#### Definition

At battle initialization, each bot is assigned a random position within the arena boundaries. Position generation ensures valid placement:
- Randomized coordinates (x, y) within arena boundary
- Minimum separation between bots (prevent initial collision)
- Deterministic through seeded randomness (enables replay)

#### Positioning Rules

**Random Coordinate Selection**:
- x_coordinate ∈ [-50, 50] (uniform random distribution within x-range)
- y_coordinate ∈ [-50, 50] (uniform random distribution within y-range)
- Both coordinates selected independently

**Collision Radius Accommodation**:
- Bot collision radius ≈ 2 units
- Bot center must maintain 2-unit clearance from boundary
- Valid position range: x_center ∈ [-48, 48], y_center ∈ [-48, 48]

**Minimum Separation Requirement**:
- Default separation distance: 20 units (TBD based on playtesting)
- Distance between bot centers must be ≥ 20 units
- Prevents spawning overlapping or immediately touching
- Provides reaction time before first engagement

**Collision-Free Generation**:
- Position generator repeats until valid non-overlapping position found
- Ensures no bots spawn touching or overlapping
- Ensures all bots have minimum separation

**Deterministic Randomness**:
- Arena configuration includes random seed for reproducibility
- Same seed → same starting positions (enables battle replay and debugging)
- Seed can be derived from match ID, round number, or explicitly specified

**Orientation** (Future):
- Bot facing direction may also be randomized (if directional weapons exist)
- Prevents orientation-based spawn advantage

#### Position Generation Algorithm

```
function generate_start_positions(arena, num_bots, seed):
  rng = SeededRandomGenerator(seed)
  positions = []

  for each bot:
    valid = false
    while not valid:
      x = rng.uniform(-48, 48)
      y = rng.uniform(-48, 48)
      position = (x, y)

      # Check separation from all existing positions
      valid = true
      for each existing_position in positions:
        distance = sqrt((x - existing_position.x)² + (y - existing_position.y)²)
        if distance < 20:
          valid = false
          break

    positions.append(position)

  return positions
```

#### Rationale for Random Positioning

* **Forces Adaptability** - Bots cannot rely on fixed spawn points; must adapt to variable initial positions
* **Prevents Hardcoded Openings** - Eliminates "spawn camping" or fixed-position opening tactics that exploit knowledge of spawns
* **Simplifies Biome Design** - Biomes don't require predefined spawn point layouts per biome
* **Simplifies Boundary Design** - Boundary configuration doesn't constrain spawn logic; works with any valid boundary
* **Fair and Balanced** - Uniform probability distribution ensures no positional advantage at battle start
* **Encourages General AI** - Bot logic must work from any starting position and orientation, not optimized for specific spawns
* **Reproducible** - Seed-based generation enables exact replay of same starting positions for debugging

#### Rejected Alternative: Fixed Positioning

**Definition**: Predefined spawn points (e.g., opposite corners, center edges).

**Example Configuration**:
- Bot A always spawns at (-40, -40)
- Bot B always spawns at (40, 40)

**Advantages**:
- Predictable and consistent
- Simple for bot developers to optimize opening sequences
- Easier battle replay (no seed needed)
- Simpler implementation (hardcoded positions)

**Disadvantages**:
- **Hardcoded opening sequences**: Bots can exploit known spawns → meta-strategy instead of adaptability
- **Per-boundary definitions required**: Each boundary size needs spawn point definition (100×100 different from 150×150)
- **Positional advantages**: Different spawn locations create unequal advantages (corner spawns vs center spawns)
- **Biome interaction**: Friction topology may favor certain spawns (ice near spawn = speed advantage)
- **Reduces AI quality**: Bot optimization focuses on specific spawns rather than general play
- **Less emergent gameplay**: Predictability reduces strategic depth and surprise tactics

**Rejection Rationale**: Sacrifices adaptability and bot quality for minor implementation convenience. Fixed spawns create meta-strategies that reward hardcoding over general intelligence.

#### Rejected Alternative: User-Selected Positioning

**Definition**: Players choose spawn points for their bots before battle.

**Advantages**:
- Strategic pre-battle positioning choice
- Player agency in arena configuration
- Interesting metagame of spawn point selection

**Disadvantages**:
- **UI Complexity**: Requires spawn point selection interface
- **Information asymmetry**: Who picks first? Creates advantage for second player
- **Metagaming**: Spawn point selection becomes meta-strategy (always pick corners vs always pick center)
- **Startup delay**: Spawn selection phase delays battle start
- **Automation complexity**: Unclear how automated/AI bots choose spawn points
- **Fairness issues**: Different spawn points create unequal opportunities
- **Not applicable to all modes**: Random automated tournaments don't fit player selection model

**Rejection Rationale**: Adds implementation and balance complexity without proportional strategic depth gain. Metagaming around spawn points is less interesting than adaptive gameplay.

#### Implementation Notes

**Numeric Value Tuning**:
- **Minimum separation: 20 units** (TBD - adjust based on playtesting)
  - Too small: Bots spawn near collision, immediate panic
  - Too large: Bots spawn far apart, artificial delay to combat
  - Optimal: Enough reaction time for first action (~1-2 ticks)
- **Arena position range: [-48, 48]** (derived from 100×100 - 2 unit radius)
  - Ensures bot collision radius doesn't exceed boundary
  - Automatic for different boundary sizes

**Seed Management**:
- Arena configuration must include explicit seed or seed derivation rule
- Seed should be tournament-round-stable (same seed for replay analysis)
- Can use hash(battle_id, round_number) as seed source

### Property 5: Win Conditions

**Chosen Option**: Option 5.3 - Multi-Path Victory (Health Depletion + Timeout + Disconnect)

#### Win Condition 1: Health Depletion (Primary)

**Trigger**: Bot's Health characteristic reaches 0

**Health Definition** (from ADR-0008):
- Each bot has Health value (initial health from characteristics + equipment)
- Health reduced by opponent attacks
- Health depleted when reduced to 0

**Victory Resolution**:
- **Immediate**: Battle ends immediately when any bot's Health reaches 0
- **Winner**: Bot with Health > 0
- **Loser**: Bot with Health = 0 (eliminated)
- **Decisive**: Cleanest victory condition; primary path for engagement-focused battles

**Gameplay Impact**:
- Rewards dealing damage to opponent
- Punishes taking damage without reciprocation
- Creates pressure to engage (can't win without opponent damage)
- Eliminates stalling strategies (must take damage eventually)

#### Win Condition 2: Timeout

**Trigger**: Battle duration reaches maximum time limit

**Time Limit**: **5 minutes (300 seconds)** [TBD - Subject to Playtesting]

**Tuning Rationale**:
- Initial estimate based on expected engagement pacing
- 5 minutes allows ~60 combat exchanges at 2-second resolution
- Sufficient for positioning and maneuvering phases
- Short enough to prevent tournament bottleneck
- Will be empirically tuned through playtesting

**Expected Tuning Ranges**:
- **Quick Play**: 3 minutes (120 seconds) - faster pacing
- **Standard Match**: 5 minutes (300 seconds) - balanced default
- **Tournament**: 7-10 minutes - allows complex strategies
- **Training**: Unlimited or 30 minutes - practice mode

**Winner Determination**:

When timeout expires, battle concludes and victory awarded:

**Case 1: Health Difference**
- If bots have different Health values:
  - **Winner**: Bot with higher Health remaining
  - **Loser**: Bot with lower Health remaining
  - **Rationale**: Rewards both damage dealing AND damage avoidance

**Case 2: Equal Health** (Tiebreaker)

Multiple resolution options [TBD - requires playtesting to determine player preference]:

- **Option A: Draw** - Timeout with equal health = draw (both qualify for tournament advancement or both eliminated)
- **Option B: Last Damage Wins** - Whoever dealt the most recent damage wins (rewards aggression and engagement)
- **Option C: Total Damage Dealt** - Whoever dealt most cumulative damage wins (measures sustained aggression)
- **Option D: Damage Efficiency** - Whoever dealt most damage relative to health remaining wins
- **Recommendation**: Option B (Last Damage Wins) - simplest, rewards engaging gameplay

**Default Tiebreaker (Recommended)**: Last bot to deal damage wins
- **Rationale**: Encourages aggressive engagement rather than passive waiting
- **Implementation**: Track timestamp of last damage dealt by each bot; winner is bot with more recent timestamp
- **Gameplay Effect**: Losing bot has incentive to take action before timeout (risk reward decision)

#### Win Condition 3: Disconnect

**Trigger**: Bot disconnects from battle server

**Disconnect Causes**:
- Network connection loss
- Bot process crash or shutdown
- Protocol violation or invalid message
- Server-side connection termination
- Client timeout

**Reconnection Grace Period**: **30 seconds** [TBD - Subject to Testing]

**Graceful Recovery Procedure**:
1. **Initial Disconnect**: Bot loses connection (0-5 seconds: initial connection failure)
2. **Grace Period Active**: Bot has 30 seconds to reconnect (5-35 seconds from disconnect)
   - Battle continues with disconnected bot in last known state
   - Disconnected bot cannot send commands
   - Disconnected bot may be affected by opponent actions (unclear state)
3. **Grace Period Expires** (35 seconds): If not reconnected, proceed to forfeit
4. **Forfeit**: Disconnected bot loses battle immediately

**Reconnection Success**:
- **Within Grace Period**: Bot resumes control from last synchronized state
- **Recovery State**: Bot state (position, health, cooldowns) preserved from moment of disconnect
- **Battle Continues**: Battle proceeds normally with reconnected bot
- **No Penalties**: Reconnection carries no Health or action penalties (encourages recovery)

**Rationale for Grace Period**:
- **Fairness for Transient Issues**: Brief network hiccup shouldn't forfeit match immediately
- **Recovery Time**: Allows reconnection within reasonable timeframe (~30 seconds typical for reconnection)
- **Progress Guarantee**: If grace period expires, battle progresses rather than hanging indefinitely
- **Competitive Clarity**: Clear rules for tournament play and ranking systems
- **Connection Stability**: Encourages implementing reliable connection logic; stability becomes competitive skill

**Numeric Value Tuning** [TBD]:
- **Casual Play**: 60 seconds (forgiving, tolerates temporary network issues)
- **Standard Match**: 30 seconds (balanced default)
- **Competitive**: 15 seconds (strict, values consistent connection)
- **Tournaments**: 30 seconds (standard fairness)

#### Win Condition Priority

When multiple win conditions could trigger, priority order:

1. **Disconnect Grace Period Expiration** (highest priority)
   - If grace period expires, disconnected bot forfeits immediately
   - Prevents hung battles waiting for reconnection indefinitely
   - Example: Timeout clock at 1 second, disconnect at 2 seconds → disconnect wins (forfeit overrides timeout)

2. **Health Depletion**
   - If bot reaches Health = 0, battle ends immediately
   - Higher priority than timeout (decisive victory)
   - Example: Bot defeated at 4:50 (10 seconds before timeout) → immediate victory

3. **Timeout**
   - If neither disconnect nor health depletion triggered, timeout determines winner
   - Lowest priority (fallback condition)
   - Example: Both bots at positive health, 5:00 elapsed → higher health bot wins

**Priority Ensures**:
- Technical failures don't hang battles (disconnect highest)
- Decisive victories end battles immediately (health > timeout)
- Battles always conclude within maximum duration (timeout guarantees conclusion)

#### Rationale for Multi-Path Approach

* **Prevents Single Strategy Dominance** - Timeout prevents pure evasion (must deal damage or maintain health), health prevents pure attrition (timeout forces engagement)
* **Guaranteed Battle Conclusion** - Every battle concludes within 5 minutes maximum; no indefinite stalemates
* **Handles Technical Failures** - Disconnect condition prevents hung battles from network/server issues
* **Competitive Clarity** - Unambiguous win determination enables tournaments, rankings, and competitive play
* **Encourages Engagement** - Timeout + health comparison rewards aggressive tactics over pure evasion
* **Risk-Reward Decisions** - Multiple win paths create strategic choices (play safe for timeout, or risk damage for health depletion)
* **Fairness** - All conditions must be satisfied or triggered to end battle; no arbitrary early stopping

#### Future Win Condition Extensions

**Additional Conditions for Future Game Modes**:
- **Capture Point**: Control arena zone for duration (King of the Hill mode)
- **Objective Completion**: Destroy/protect specific targets (objective-based mode)
- **Elimination**: All opponents eliminated (free-for-all or team modes)
- **Score Threshold**: First bot to reach score target (best-of-N matches)

**Different Timeout Values by Mode**:
- **Quick Play**: 3 minutes (faster matches)
- **Tournament**: 10 minutes (strategic depth)
- **Training**: No timeout (practice mode)

**Dynamic Boundaries** (Future):
- Shrinking arena boundaries could add pressure (arena gets smaller over time)
- Ring of fire approaching boundary (environmental hazard)
- Would require additional ADR defining shrinking mechanics

### Consequences

#### Overall Arena System

* Good, because property-based framework enables diverse battle configurations through composition without changing core systems
* Good, because Arena concept cleanly separates universal topology (ADR-0005) from battle instances (ADR-0011)
* Good, because selectable properties (biome, boundary) create pre-battle strategic choices and player agency
* Good, because framework is extensible: add new properties, boundary shapes, biomes, visibility modes without breaking core system
* Good, because integrates seamlessly with existing systems (biome framework ADR-0010, physics ADR-0006, spatial system ADR-0005)
* Good, because provides foundation for future game modes (team battles, tournaments) with different arena properties
* Neutral, because property values (timeout duration, reconnection period, minimum separation) require playtesting for tuning
* Neutral, because property framework adds system complexity compared to single fixed arena
* Bad, because provides many configuration options that may overwhelm casual players without guidance
* Bad, because requires UI/UX for property selection and tournament configuration

#### Property 1: Biome Selection

* Good, because operationalizes ADR-0010 biome framework in gameplay; validates biome concept through practice
* Good, because different friction topologies create tactical variety and environment differentiation
* Good, because pre-battle biome choice becomes strategic decision based on bot capabilities
* Good, because forces bot adaptability to handle variable friction across arena
* Good, because geography-based positioning creates additional tactical complexity
* Neutral, because initial implementation may support only single biome until concrete biome ADRs defined
* Bad, because biome balance tuning complex (different friction topologies may favor different bot types)

#### Property 2: Boundary Configuration

* Good, because selectable boundaries enable arena variety without changing universal topology
* Good, because different sizes create different engagement dynamics (small = frequent combat, large = maneuvering)
* Good, because enables different game modes with appropriate arena sizes
* Good, because allows balance tuning through boundary size adjustment based on playtesting
* Good, because rectangular boundaries have simple collision detection (comparisons vs complex geometry)
* Neutral, because initial implementation supports only single rectangular size until additional boundary ADRs defined
* Neutral, because boundary size affects engagement pacing and requires tuning per game mode
* Bad, because corners enable defensive camping strategies (bots can position in corners safely)
* Bad, because rectangular boundaries provide less engagement pressure than shrinking or circular alternatives

#### Property 3: Visibility System

* Good, because full visibility is simplest implementation for MVP (no visibility system required)
* Good, because complete information enables pure strategic gameplay without information luck
* Good, because bot developers don't need visibility query implementation for initial release
* Good, because baseline validation of combat mechanics before adding visibility complexity
* Neutral, because full visibility may feel unrealistic or reduce immersion for some players
* Neutral, because complete information is advantage to experienced players (can execute perfect knowledge strategies)
* Bad, because eliminates reconnaissance and information advantage gameplay
* Bad, because defers interesting fog of war mechanics (ambush tactics, scouting, information control)
* Bad, because may not support future sensor equipment or tactical reconnaissance gameplay

#### Property 4: Start Positioning

* Good, because random positioning forces bot adaptability to any starting location
* Good, because prevents hardcoded opening sequences that exploit fixed spawns
* Good, because simplifies biome and boundary design (no spawn point definitions needed)
* Good, because fair and balanced (uniform distribution ensures no positional advantage)
* Good, because encourages general-purpose bot AI rather than spawn-optimized logic
* Good, because seed-based randomness enables reproducible battles for debugging
* Neutral, because minimum separation distance requires playtesting tuning
* Neutral, because adds minor implementation complexity for position generation
* Bad, because unpredictable start positions may feel chaotic for casual players
* Bad, because eliminates strategic control over starting positions (some players prefer control)

#### Property 5: Win Conditions

* Good, because multi-path approach prevents single strategy dominance (no pure evasion, no infinite attrition)
* Good, because guaranteed battle conclusion within 5 minutes ensures tournaments and rankings work
* Good, because health depletion provides decisive victory condition
* Good, because timeout ensures battles progress even with evasive opponents
* Good, because disconnect handling prevents hung battles from technical failures
* Good, because clear priority order eliminates ambiguity in edge cases
* Neutral, because timeout tiebreaker (last damage) requires playtesting validation
* Neutral, because 5-minute duration and 30-second reconnection grace period require empirical tuning
* Bad, because multiple win conditions add system complexity and player learning curve
* Bad, because timeout mechanics may reward passive play (wait for opponent to make mistakes)
* Bad, because reconnection grace period could allow unfair advantage (lag exploits)

## Confirmation

The Arena system and 1v1 battle properties will be confirmed through:

### 1. Arena System Implementation

- [ ] Arena configuration correctly specifies biome, boundary, visibility, positioning, win conditions
- [ ] Battle server instantiates arena from configuration
- [ ] Arena properties correctly constrain battle gameplay (bots respect boundaries, see correct information, spawn at designated positions)
- [ ] Multiple arena instances can be created and configured independently

### 2. Biome Integration

- [ ] Selected biome provides μ(position) function for friction physics
- [ ] Friction topology from biome correctly affects bot movement per ADR-0006 formula: F_friction = μ(position) × M × |v|
- [ ] Different biomes create observable tactical differences (ice zones faster, mud zones slower)
- [ ] Concrete biome ADRs successfully instantiate biome framework (Desert, Arctic, Forest)
- [ ] Friction map queries meet real-time tick rate requirements

### 3. Boundary Mechanics

- [ ] Bot movement clamped at boundaries (no escape outside arena)
- [ ] Elastic collisions apply correctly at walls
- [ ] Rectangular boundary collision detection uses simple comparisons (efficient)
- [ ] Bots can position exactly on boundary line without clipping
- [ ] Different boundary sizes work correctly (100×100, 150×150, etc.)

### 4. Visibility System

- [ ] State updates include complete position information for all bots
- [ ] No occlusion calculations or filtering performed
- [ ] Bot implementations receive full information each tick
- [ ] Deterministic gameplay with complete information

### 5. Start Positioning

- [ ] Random position generation produces valid placements
- [ ] Minimum 20-unit separation between bots enforced
- [ ] Seed-based generation enables reproducible start positions
- [ ] No bots spawn overlapping or outside boundaries
- [ ] Same seed produces identical starting positions (for replay)

### 6. Win Condition Detection

- [ ] Health Depletion: Battle ends immediately when Health = 0
- [ ] Timeout: Battle concludes at 5 minutes with higher-health winner
- [ ] Disconnect: Grace period tracked, forfeit triggered at expiration
- [ ] Win condition priority correct (disconnect > health > timeout)

### 7. Integration Testing

- [ ] ADR-0005 topological properties respected (R² space, Cartesian coordinates)
- [ ] ADR-0006 physics laws applied correctly (friction, collisions, movement)
- [ ] ADR-0007 movement mechanics constrained by boundaries
- [ ] ADR-0008 health characteristic used for health depletion win condition
- [ ] ADR-0009 equipment mass affects movement within friction system
- [ ] ADR-0010 biome framework provides geography for selected biomes
- [ ] ADR-0004 gRPC protocol transmits arena configuration and state

### 8. Playtesting Validation

- [ ] Arena size (100×100) provides appropriate tactical space (not too cramped, not too sprawling)
- [ ] Engagement frequency appropriate for 5-minute timeout (not excessive running, not immediate combat)
- [ ] Biome friction variations create meaningful positioning decisions
- [ ] Random spawning produces varied and interesting battles
- [ ] Win conditions appropriately distributed (reasonable health vs timeout comparison frequency)
- [ ] Disconnect grace period (30 seconds) appropriately balances fairness and progress

## Pros and Cons of the Options

### Property 1: Biome

#### Option 1.1: Fixed (Single Biome)

All 1v1 battles use identical biome with uniform friction across entire arena.

* Good, because simplest implementation (no biome selection system required)
* Good, because consistent gameplay experience (no biome learning curve)
* Good, because avoids biome-specific UI and configuration complexity
* Good, because all bots compete on same terrain (fair from environment perspective)
* Neutral, because appropriate if only one biome is ever defined
* Bad, because eliminates tactical variety from geography differentiation
* Bad, because fails to validate ADR-0010 biome framework in practice
* Bad, because removes pre-battle strategic choice related to terrain
* Bad, because limits extensibility (makes future biome addition more complex)
* **Status**: REJECTED - Reduces gameplay variety and fails to validate biome framework

#### Option 1.2: Selectable (CHOSEN)

User selects biome from available options when starting battle. Biome determines friction topology.

* Good, because operationalizes ADR-0010 biome framework in gameplay
* Good, because validates biome concept through competitive practice
* Good, because creates tactical variety through different friction topologies
* Good, because enables pre-battle strategic selection based on bot capabilities
* Good, because forces bot adaptability to variable friction (tests AI quality)
* Good, because increases gameplay depth through terrain-based positioning
* Good, because extensible: adding new biomes doesn't change core system
* Neutral, because requires future ADRs to define concrete biomes before meaningful choices exist
* Neutral, because biome balance requires careful tuning per biome
* Bad, because adds biome selection UI/UX complexity
* Bad, because increases testing scope (must balance multiple biomes)
* Bad, because players need education on biome differences and friction effects

### Property 2: Boundary

#### Option 2.1: Fixed (Single Boundary)

All 1v1 battles use identical 100×100 rectangular boundary.

* Good, because simplest implementation (no boundary selection system)
* Good, because consistent gameplay experience (predictable arena size)
* Good, because standard arena size familiar to all players
* Good, because matches existing documentation and expectation
* Neutral, because appropriate if single arena size sufficient
* Bad, because eliminates arena variety through size variation
* Bad, because prevents different game modes with different arena needs
* Bad, because hardcodes single boundary to topological system (tight coupling)
* Bad, because limits future expansion to different arena types
* Bad, because removes pre-battle strategic choice of arena scale
* **Status**: REJECTED - Reduces flexibility and couples topology to game mode

#### Option 2.2: Selectable (CHOSEN)

User selects boundary configuration from available options. Boundary determines arena shape and dimensions.

* Good, because enables arena variety without changing universal topology
* Good, because different sizes create different engagement dynamics
* Good, because allows different game modes with appropriate arena sizes
* Good, because supports balance tuning through boundary size adjustment
* Good, because enables pre-battle strategic selection (prefer maneuvering space or engagement frequency)
* Good, because rectangular boundaries have simple collision detection
* Good, because extracting boundary from ADR-0005 cleanly separates topology from game mode
* Good, because extensible: new boundary shapes and sizes don't change core system
* Neutral, because initial implementation supports only single rectangular size until additional boundary ADRs defined
* Neutral, because rectangular boundaries may feel static compared to dynamic or circular alternatives
* Bad, because corners enable defensive camping strategies
* Bad, because rectangular boundaries provide less engagement pressure than shrinking/circular
* Bad, because adds boundary selection UI/UX complexity

#### Option 2.3: Dynamic (Shrinking)

Arena boundaries shrink over battle duration, forcing increasing interaction.

* Good, because creates increasing pressure for engagement over time
* Good, because prevents indefinite evasion (shrinking walls force approach)
* Good, because provides dynamic pacing (slow start, fast ending)
* Good, because creates comeback mechanics (losing bot has pressure incentive)
* Bad, because significantly more complex implementation (shrinking geometry, collision updates)
* Bad, because requires complex math for shrinking boundary collision detection
* Bad, because unfamiliar mechanic to most players (needs learning)
* Bad, because changes physics at runtime (friction, collision points shift)
* Bad, because may cause unexpected or frustrating boundary behavior
* **Status**: DEFERRED - High implementation complexity; defer until core arena system validated

### Property 3: Visibility

#### Option 3.1: Full (CHOSEN)

All bots see complete information about all opponents at all times.

* Good, because simplest implementation (no visibility system required)
* Good, because complete information enables pure strategic gameplay
* Good, because deterministic (same state + same logic = same decisions)
* Good, because bot developers don't need visibility queries
* Good, because validates core combat mechanics without visibility complexity
* Good, because accessible to new bot developers (no advanced visibility concepts)
* Good, because baseline for future visibility system additions
* Neutral, because may feel unrealistic (unlimited visibility unusual)
* Neutral, because may provide advantage to experienced players (perfect information strategies)
* Bad, because eliminates reconnaissance and scouting gameplay
* Bad, because removes information advantage mechanics
* Bad, because doesn't support future sensor equipment or visibility items

#### Option 3.2: Constant Fog of War

Bots see only within fixed radius around their position (e.g., 30 units).

* Good, because creates reconnaissance gameplay and scouting strategy
* Good, because enables information advantage through positioning
* Good, because supports ambush tactics and positional surprise
* Good, because enables sensor equipment mechanics (ADR-0009 deferred)
* Good, because more realistic visibility model
* Neutral, because appropriate after full visibility baseline validated
* Bad, because requires visibility radius system architecture
* Bad, because requires state filtering and partial updates
* Bad, because increases computational cost (visibility checks per tick)
* Bad, because adds bot developer complexity (visibility query implementation)
* Bad, because requires UI updates (fog of war visualization)
* **Status**: DEFERRED to future ADR - Defer until full visibility validated

#### Option 3.3: Vanishing Fog of War

Bots see within radius, but explored areas remain visible indefinitely (fog "vanishes" once explored).

* Good, because creates map control mechanics and exploration strategy
* Good, because provides persistent information advantage for early exploration
* Good, because creates nice transition from full visibility to constant fog (user-noted)
* Good, because enables scouting-focused bot strategies
* Good, because supports positional reconnaissance tactics
* Neutral, because more complex but not most complex
* Bad, because requires per-bot visibility memory (explored area state)
* Bad, because requires fog state tracking and synchronization
* Bad, because state grows over battle duration (memory management)
* Bad, because highest implementation complexity of visibility options
* Bad, because unfamiliar mechanic requiring player education
* **Status**: DEFERRED to future ADR - Defer until simpler fog of war options validated

### Property 4: Start Positioning

#### Option 4.1: Fixed

Predefined spawn points (e.g., opposite corners: Bot A at (-40, -40), Bot B at (40, 40)).

* Good, because predictable and consistent for all battles
* Good, because simplest implementation (hardcoded positions)
* Good, because enables bot optimization for specific spawns
* Good, because simple for replay and debugging
* Bad, because enables hardcoded opening sequences exploiting known spawns
* Bad, because requires per-boundary spawn definitions (100×100 different from 150×150)
* Bad, because different spawns create unequal positional advantages
* Bad, because biome interaction may favor certain spawns (ice near spawn = unfair speed advantage)
* Bad, because discourages general-purpose bot AI (rewards spawn-specific optimization)
* Bad, because reduces emergent gameplay and strategic variety
* **Status**: REJECTED - Hardcoded spawns enable exploit strategies over adaptive play

#### Option 4.2: Random (CHOSEN)

Bots placed at random positions within arena, seeded for reproducibility.

* Good, because forces bot adaptability to any starting position
* Good, because prevents hardcoded opening sequences that exploit fixed spawns
* Good, because simplifies biome and boundary design (no per-configuration spawn points needed)
* Good, because fair and balanced (uniform distribution, no positional advantage)
* Good, because encourages general-purpose bot AI rather than spawn optimization
* Good, because seed-based reproducibility enables debugging and replay
* Good, because enables emergent gameplay and tactical variety
* Neutral, because minimum separation distance (20 units) requires empirical tuning
* Neutral, because adds position generation implementation complexity
* Bad, because unpredictable spawning may feel chaotic to casual players
* Bad, because removes player control over starting positions
* Bad, because may spawn bots far apart (longer initial engagement delays)

#### Option 4.3: User-Selected

Players choose spawn points for their bots before battle.

* Good, because provides strategic pre-battle positioning choice
* Good, because gives player agency in arena configuration
* Good, because creates interesting metagame of spawn selection
* Bad, because requires spawn point selection UI
* Bad, because creates information asymmetry (who picks first advantage)
* Bad, because enables spawn metagaming (always pick corners vs always pick center)
* Bad, because delays battle start for spawn selection phase
* Bad, because unclear for automated bots (how do AI bots select spawns?)
* Bad, because fairness issues (different spawns = unequal opportunities)
* **Status**: REJECTED - Adds complexity without proportional strategic depth

### Property 5: Win Conditions

#### Option 5.1: Health Depletion Only

Battle ends immediately when any bot's Health reaches 0.

* Good, because simplest win condition (single trigger)
* Good, because decisive victory condition
* Good, because rewards aggressive damage dealing
* Bad, because allows indefinite evasion (bot with health lead can avoid opponent)
* Bad, because creates stalemate scenarios (both bots evading, no damage dealt)
* Bad, because provides advantage to defensive bot (can outlast aggressive opponent)
* Bad, because no time pressure (battles could theoretically last indefinitely)
* **Status**: REJECTED - Allows evasion strategies that prevent engagement

#### Option 5.2: Health Depletion + Timeout

Health depletion is primary; timeout is fallback conclusion if neither bot depleted after 5 minutes.

* Good, because health depletion is decisive victory
* Good, because timeout prevents indefinite evasion
* Good, because simpler than three-path system (no disconnect handling)
* Neutral, because two conditions simpler than three
* Bad, because doesn't handle disconnect scenarios (hanging battles on technical failure)
* Bad, because doesn't address reconnection fairness
* Bad, because technical failures (network loss, bot crash) would hang battles
* **Status**: CONSIDERED - Works but incomplete without disconnect handling

#### Option 5.3: Health Depletion + Timeout + Disconnect (CHOSEN)

Multi-path victory with three independent win conditions each with priority.

* Good, because health depletion provides decisive victory (first path)
* Good, because timeout guarantees battle conclusion within 5 minutes (second path)
* Good, because disconnect handling prevents hung battles from technical failures (third path)
* Good, because prevents single strategy dominance (timeout prevents evasion, health prevents attrition)
* Good, because clear priority order eliminates ambiguity
* Good, because handles all edge cases (technical failures, prolonged evasion, sudden elimination)
* Good, because competitive clarity for tournaments and rankings
* Good, because encourages engagement (multiple paths reward different strategies)
* Good, because graceful recovery from transient connection issues (grace period)
* Neutral, because three conditions more complex than one or two
* Neutral, because timeout values and tiebreakers require empirical tuning
* Neutral, because disconnection grace period requires adjustment per use case (casual vs competitive)
* Bad, because more complex system than single win condition
* Bad, because requires timeout tiebreaker (last damage, health, score?)
* Bad, because timeout mechanics may reward passive play (wait for opponent to make mistakes)
* Bad, because reconnection grace period could enable lag exploitation

## More Information

### Related Documentation

**Core Spatial and Physics Foundation**:
- **ADR-0005: BattleBot Universe Topological Properties** - Defines R² Euclidean space, Cartesian coordinates, universal spatial properties
- **ADR-0006: BattleBot Universe Physics Laws** - Defines variable friction, elastic collisions, movement physics
- **ADR-0007: Bot Movement Mechanics** - Defines thrust-based movement, terminal velocity, movement constraints

**Bot Configuration and Characteristics**:
- **ADR-0008: Bot Characteristics System** - Defines Health, Defense, Mass properties that affect gameplay
- **ADR-0009: Equipment and Loadout System** - Defines weapons and armor affecting characteristics and mass

**Biome and Arena Framework**:
- **ADR-0010: Biome Definition Framework** - Defines Geography property and biome meta-framework
- **ADR-0004: Bot to Battle Server Interface** - gRPC protocol for arena configuration and state transmission

**User Documentation**:
- **1v1 Battles Documentation** - User-facing guide to 1v1 game mechanics and arena system

### Boundary Extraction from ADR-0005

**What ADR-0005 Previously Stated**:
```
Property 4: Boundary - Rectangular Boundary (CHOSEN)

Arena Size: 100 x 100 units (TBD - subject to tuning based on playtesting)
X-axis Range: -50 to +50 units
Y-axis Range: -50 to +50 units

Mathematical Definition:
- Topological Space: 2-dimensional Euclidean space R²
- Manifold: Closed rectangular region [−50, 50] × [−50, 50] ⊂ R²
- Metric: Standard Euclidean metric d(p,q) = √((x₂−x₁)² + (y₂−y₁)²)
- Coordinate Chart: Cartesian coordinates φ: R² → R² where φ(p) = (x, y)
- Boundary: ∂M = {(x,y) : x=±50 or y=±50}
```

**Why Extraction Was Necessary**:

1. **Coupling Issue** - Hardcoding boundary in universal topology couples spatial system to specific game mode (1v1)
2. **Inflexibility** - Cannot modify boundary size without changing topological definition
3. **Future Scalability** - Other game modes need different boundaries; would require topology changes for each
4. **Architectural Clarity** - Separates universal properties (R²) from game-mode-specific properties (arena dimensions)

**What Transfers to ADR-0011**:
- Rectangular boundary dimensions ([-50, 50] × [-50, 50])
- Boundary shapes and size options
- Boundary selection mechanism
- Movement clamping and collision mechanics
- Mathematical definitions (preserved from ADR-0005)

**What Remains in ADR-0005**:
- 2D Euclidean topological space R²
- Cartesian coordinate system and origin definition
- Spatial metrics and mathematical properties
- Universal coordinate system rules

**Integration Approach**:
- ADR-0011 preserves and reuses mathematical notation from ADR-0005
- ADR-0011 applies boundary within specific game mode (1v1)
- Future game modes can apply different boundaries without topology changes
- Enables selective application of boundaries per game type

### Future Arena Extensions

**Visibility Systems** (future ADRs):
- **ADR-TBD: Constant Fog of War** - Visibility radius mechanics for reconnaissance gameplay
- **ADR-TBD: Vanishing Fog of War** - Explored area persistence for map control gameplay

**Boundary Variants** (future ADRs):
- **ADR-TBD: Circular Boundaries** - Circular arena for radial gameplay
- **ADR-TBD: Hexagonal Boundaries** - Hexagonal arena for balanced symmetry
- **ADR-TBD: Dynamic Boundaries** - Shrinking arena for increasing engagement pressure

**Biome Implementations** (future ADRs):
- **ADR-TBD: Desert Biome** - Sand, rock, gravel terrain with defined friction topology
- **ADR-TBD: Arctic Biome** - Ice, snow, frozen rock terrain with low-friction zones
- **ADR-TBD: Forest Biome** - Grass, swamp, dirt terrain with varied friction

**Additional Win Conditions** (future ADRs):
- **ADR-TBD: King of the Hill** - Control arena center zone for duration
- **ADR-TBD: Capture Point** - Objective-based victory conditions
- **ADR-TBD: Free-for-All** - Elimination with multiple bots

**Game Modes** (future specifications):
- **Quick Play**: 3-minute timeout, 75×75 arena, random biome
- **Tournament**: 10-minute timeout, 100×100 arena, predefined biome selection
- **Training**: No timeout, 100×100 arena, selectable biome, practice mode
- **Team Battles**: Multiple 1v1 arenas in coordination (requires future ADR for team mechanics)

### Design Principles

**Property-Based Framework**: Arena configured through composition of independent properties rather than predefined templates. Enables combinations and future extensibility.

**Extensibility First**: Property system designed to add new options without breaking existing configurations. Each property can be extended independently.

**Simplicity First**: Initial implementation uses simplest option for each property (Full visibility, Random positioning, Health + Timeout + Disconnect). More complex options deferred until basics validated through gameplay.

**Leverage Existing Systems**: Builds on ADR-0005 topology, ADR-0006 physics, ADR-0010 biome framework. Does not invent new physics or spatial concepts.

**Validate Before Expanding**: Biome framework must be proven through concrete implementations before properties like Climate, Vegetation added. Visibility must be tested with full visibility before introducing fog of war.

**Clear Separation of Concerns**:
- ADR-0005: Universal topological space (R²)
- ADR-0006: Universal physics laws (friction, collisions)
- ADR-0010: Biome framework (what defines biomes)
- ADR-0011: Battle configuration (how to configure specific battles)

**Player Agency**: Selectable properties (biome, boundary) enable pre-battle strategic choices without adding excessive complexity.

### Implementation Notes

**Numeric Values Marked TBD**:

All numeric constants should be marked as TBD and tuned through playtesting:

- **Timeout duration: 5 minutes** [TBD]
  - Adjust based on average battle engagement frequency
  - If battles timeout too frequently: increase timeout
  - If battles rarely timeout: decrease timeout
  - Target: ~10-20% of battles end via timeout, 70-80% via health, <5% via disconnect

- **Reconnection grace period: 30 seconds** [TBD]
  - Adjust based on typical reconnection times observed
  - If players can't reconnect in time: increase period
  - If period too long: battles stall waiting for reconnection
  - Target: 80%+ successful reconnections, <1% timeout waiting for reconnect

- **Minimum spawn separation: 20 units** [TBD]
  - Adjust based on first-engagement timing
  - If spawns too close: immediate panic and poor opening gameplay
  - If spawns too far: excessive running before engagement
  - Target: 1-2 ticks reaction time before first possible engagement

- **Arena size 100×100** [TBD]
  - Adjust based on engagement dynamics
  - If battles feel cramped: increase to 150×150
  - If battles feel sparse: decrease to 75×75
  - Target: Balanced mix of positioning and direct combat

**Playtesting Priorities**:

1. **Balance Testing**: Do win condition paths trigger with appropriate frequency? (Health 70-80%, Timeout 10-20%, Disconnect <5%)
2. **Engagement Testing**: Is arena size appropriate for engagement frequency?
3. **Biome Testing**: Do different biomes create observably different gameplay?
4. **Spawn Testing**: Is minimum separation appropriate? Do bots have reaction time?
5. **Timeout Testing**: Is 5-minute duration appropriate for expected match pacing?
6. **Disconnect Testing**: Is 30-second grace period fair for connection recovery?

### Integration Matrix

| ADR | Property | How Used |
|-----|----------|----------|
| ADR-0004 | Boundary, Biome, Visibility, Positioning, Win Conditions | gRPC transmits arena configuration, state includes all properties |
| ADR-0005 | Boundary (Mathematical), Coordinates | Arena uses R² topology, Cartesian coordinates, inherits metric |
| ADR-0006 | Boundary (Collisions), Biome (Friction) | Physics applies within boundaries; friction from biome |
| ADR-0007 | Boundary (Movement Constraint), Positioning | Movement constrained by boundaries; random positioning generated at start |
| ADR-0008 | Win Conditions (Health) | Health characteristic used for health depletion win condition |
| ADR-0009 | Boundary (Mass affects movement), Biome (Friction interaction) | Equipment mass interacts with biome friction |
| ADR-0010 | Biome (Geography) | Selected biome provides friction topology μ(position) |

### Design Patterns Used

**Property-Based Composition**: Arena defined as collection of independent properties:
- Each property selectable independently
- Properties compose to form complete arena
- Similar pattern used in ADR-0010 (Biome properties: Geography, Climate, etc.)

**Deferred Alternatives**: Complex alternatives documented but not implemented:
- Fog of war modes documented but deferred
- Dynamic boundaries documented but deferred
- Alternative biomes documented but deferred
- Enables clear roadmap without initial bloat

**Seed-Based Determinism**: Random positioning reproducible through seed:
- Same seed → same starting positions
- Enables battle replay and analysis
- Maintains randomness for gameplay variety

**Graceful Degradation**: System supports both connected and disconnected states:
- Grace period allows recovery from transient failures
- Clear timeout prevents indefinite hangs
- No requirement for perfect connectivity

---

## Summary

ADR-0011 introduces the **Arena** concept as the bridge between the universal BattleBot Universe (ADR-0005) and concrete 1v1 battle instances. Five key properties define arena configuration:

1. **Biome** (Selectable) - Terrain type determining friction topology
2. **Boundary** (Selectable) - Arena shape and dimensions
3. **Visibility** (Full) - Information available to bots
4. **Start Positioning** (Random) - Initial bot placement
5. **Win Conditions** (Multi-Path) - Health + Timeout + Disconnect

This property-based framework enables gameplay variety, player agency, and future extensibility while maintaining consistent spatial and physics foundations. The arena system formally separates universal topology (ADR-0005) from battle instance configuration (ADR-0011), enabling different game modes with different arena properties without changing core systems.
