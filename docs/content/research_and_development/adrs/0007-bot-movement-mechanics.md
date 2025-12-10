---
title: "[0007] Bot Movement Mechanics"
description: >
  Defines how bots control their movement in the 2D battle space using thrust-based force application
type: docs
weight: 7
category: "strategic"
status: "accepted"
date: 2025-12-08
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

The BattleBot Universe requires a clearly defined movement system that specifies **how bots control their movement** in the 2D battle space. ADR-0005 established the spatial foundation (2D continuous Euclidean space with rectangular boundaries), and ADR-0006 defined the physics laws that govern movement (surface friction, collisions, Mass property). However, neither ADR has explicitly decided the movement mechanics - specifically, what API do bot developers use to move their bot, and how does the game engine translate those commands into position updates?

This decision defines how bots initiate and sustain movement through the battle space. The choice of movement model has cascading implications for:

- Bot developer experience (complexity of movement control API)
- Tactical depth (ability to manage momentum and positioning)
- Physics consistency (integration with friction and Mass properties)
- Equipment tradeoffs (how loadout weight affects mobility)
- Computational requirements (per-tick physics calculations)
- Determinism (predictability of movement for bot AI)

Without a well-defined movement system, we cannot:

- Provide bot developers with a clear API for movement control
- Ensure consistent integration with friction physics (ADR-0006)
- Create meaningful equipment-based mobility tradeoffs through Mass
- Support deterministic trajectory calculation for bot AI
- Guarantee realistic physics-based gameplay
- Balance different bot builds (heavy vs light loadouts)

## Decision Drivers

* **Physics Integration** - Movement model must seamlessly integrate with friction (F = μ × M × |v|) and Mass (A = F/M) from ADR-0006
* **Spatial Integration** - Must work with continuous 2D Euclidean space and rectangular boundaries from ADR-0005
* **Mass-Based Tradeoffs** - Equipment weight should meaningfully affect movement capability
* **Tactical Depth** - Movement system should enable strategic positioning, momentum management, and skill expression
* **Developer Accessibility** - Movement control API should be understandable and implementable by bot developers
* **Determinism** - Same movement commands should produce identical results across platforms
* **Computational Efficiency** - Per-tick physics calculations must be performant in real-time gameplay (60+ ticks/second target)
* **Realism vs Gameplay** - Balance realistic physics with engaging movement mechanics
* **Future-Proofing** - Should support future mechanics (knockback, explosions, variable thrust equipment)
* **gRPC Integration** - Real-time streaming of position updates (ADR-0004) should work smoothly

## Decision Outcome

### Chosen: Option 2 - Thrust Model

Bots control their movement by applying **thrust force vectors**. The game engine calculates acceleration based on thrust and friction, updates velocity, and applies friction opposing movement. Bots must continuously apply thrust to sustain movement; friction causes natural deceleration when thrust stops.

This model is chosen because:

1. **Perfect physics integration**: Naturally integrates with surface friction (F = μ × M × |v|) and Mass-based acceleration (A = F/M)
2. **Strong equipment tradeoffs**: Heavy equipment significantly reduces mobility through increased Mass
3. **Tactical depth**: Enables momentum management, acceleration/deceleration gameplay, pursuit/kiting dynamics
4. **Deterministic physics**: Standard force-based physics with well-defined formulas for trajectory prediction
5. **Future-proof**: Naturally supports knockback weapons, collision momentum transfer, explosive forces
6. **Consistent with existing ADRs**: ADR-0006 physics laws and ADR-0007→0008 characteristics already assume thrust-based movement

## Movement System Implementation Specification

### Thrust Action API

**Bot Command**: `apply_thrust(force_x: float, force_y: float)`

Alternative representation: `apply_thrust(force_magnitude: float, angle_radians: float)`

**Parameters**:
- `force_x`, `force_y`: Cartesian components of thrust force (in abstract force units)
- Represents the force magnitude and direction a bot applies each game tick
- Thrust is omnidirectional (can be applied in any direction within the 2D plane)

**Thrust Capacity**:
- **Maximum thrust**: Fixed universal constant (TBD value, e.g., 100 force units maximum)
- All bots have equal maximum thrust capacity (not equipment-dependent in initial design)
- Rationale: Mass already creates sufficient mobility differentiation through acceleration formula (A = F/M)
- Future ADRs may introduce equipment-based variable thrust (engines, thrusters) if universal constant becomes limiting

**Action Frequency**: Bots can apply thrust every game tick (continuous control)

### Per-Tick Physics Update

The game engine performs the following physics update **each game tick** to apply thrust-based movement:

**1. Collect Forces**:
- `F_thrust` = thrust force applied by bot command (clamped to maximum capacity)
- `F_friction` = surface friction opposing movement (calculated next)
- `F_collision` = collision response forces (if colliding with other bots or walls)
- `F_external` = knockback, explosions, and other forces (future mechanics)

**2. Calculate Friction Force**:

```
F_friction = μ(position) × M × |v|
```

Where:
- `μ(position)` = position-dependent surface friction coefficient (varies by terrain/arena zone)
- `M` = total bot mass = base mass + sum of equipped item masses
- `|v|` = magnitude of velocity vector
- Direction of friction: opposite to velocity vector (opposes movement)

**3. Calculate Net Force**:

```
F_net = F_thrust - F_friction + F_collision + F_external
```

(Friction opposes thrust; only movement forces contribute to net force)

**4. Calculate Acceleration**:

```
A = F_net / M
```

Where M is total mass from Step 2.

**5. Update Velocity**:

```
v_new = v_current + A × dt
```

Where `dt` is time per game tick (TBD, likely 1/60 second for 60 ticks/second gameplay)

**6. Update Position**:

```
position_new = position_current + v_new × dt
```

**7. Boundary Enforcement**:

- Check if new position exceeds arena boundaries ([-50, 50] × [-50, 50])
- If boundary violation: clamp position to boundary edge
- If bot hits wall: reflect velocity component perpendicular to boundary (elastic collision with infinite-mass wall)
  - Hit x-boundary (x = ±50): reflect x-component: `v_x → -v_x`, keep `v_y` unchanged
  - Hit y-boundary (y = ±50): reflect y-component: `v_y → -v_y`, keep `v_x` unchanged

**8. Collision Detection and Resolution**:

- Detect bot-to-bot collisions (bots touching or overlapping)
- Apply elastic collision momentum transfer (formulas from ADR-0006):

```
v1' = ((m1 - m2) × v1 + 2 × m2 × v2) / (m1 + m2)
v2' = ((m2 - m1) × v2 + 2 × m1 × v1) / (m1 + m2)
```

Where:
- `m1`, `m2` = bot masses
- `v1`, `v2` = velocities before collision (projected onto collision normal)
- `v1'`, `v2'` = velocities after collision

- Update bot velocities based on collision results

### Terminal Velocity

When thrust force balances friction force, bots reach **terminal velocity** (maximum speed):

```
F_thrust = F_friction
F_thrust = μ(position) × M × v_terminal

v_terminal = F_thrust / (μ(position) × M)
```

**Implications**:
- Heavier bots (high M from equipment) reach **lower** terminal velocity for same thrust
- Light bots (minimal equipment) reach **higher** terminal velocity
- **Variable friction zones** affect terminal velocity:
  - Ice terrain (low μ): higher terminal velocity (less friction resistance)
  - Mud terrain (high μ): lower terminal velocity (more friction resistance)
- Creates natural **mobility-firepower tradeoff**: Heavy weapons/armor reduce movement speed

**Example calculations** (TBD values for illustration):
- Heavy bot (M = 50): v_terminal = 100 / (0.5 × 50) = 4 units/tick
- Light bot (M = 10): v_terminal = 100 / (0.5 × 10) = 20 units/tick
- On ice (μ = 0.2): v_terminal = 100 / (0.2 × M) = 500/M (higher)
- On mud (μ = 1.5): v_terminal = 100 / (1.5 × M) = 67/M (lower)

### Movement Characteristics Integration

**Mass** (from ADR-0008, renumbered from ADR-0007) affects movement:

**Force-to-Velocity Relationship**:
```
Acceleration = Applied Thrust Force / Mass
A = F_thrust / M
```

**Implications**:
- Heavier bots require more thrust to achieve same acceleration
- Equipment-derived Mass (weapons, armor) directly affects mobility
- Light bots with minimal equipment are nimble and responsive
- Heavy bots sacrifice speed for firepower and protection
- Creates meaningful equipment loadout decisions (offensive/defensive power vs mobility)

**Thrust Capacity** (not a characteristic; universal constant):
- All bots share same maximum thrust capacity
- Mass determines how efficiently thrust translates to acceleration
- Skill expression: managing thrust within capacity to optimize trajectory

### Boundary and Collision Handling

**Boundary Collisions** (bots hitting arena walls):
- Position clamped to [-50, 50] × [-50, 50]
- Velocity reflected elastically (walls have infinite mass)
- Bot can apply thrust away from boundary to escape
- Prevents indefinite evasion (engagement guaranteed)

**Bot-to-Bot Collisions** (momentum transfer):
- Heavy bots push light bots (mass determines collision outcome)
- Momentum conserved (elastic collision model)
- Light bots can be pushed into unfavorable positions
- Heavy bots harder to displace
- Creates tactical collision-based strategies (ramming, positioning)

**Continuous Movement Requirement**:
- Without thrust, friction decelerates bot to stop
- Bots cannot "coast" indefinitely at high velocity
- Requires active thrust command to sustain movement
- Creates continuous engagement (bots must keep acting)

## Consequences

### Movement Model Advantages

* Good, because thrust model creates realistic physics-based movement with direct Force → Acceleration → Velocity integration
* Good, because Mass (A = F/M) creates strong equipment-based mobility tradeoffs
* Good, because friction (F = μ × M × |v|) creates natural velocity decay without artificial speed caps
* Good, because terminal velocity emerges naturally from physics (not a hard-coded limit)
* Good, because momentum management enables tactical depth and skill expression
* Good, because deterministic force-based physics allows bot developers to predict trajectories
* Good, because consistent with ADR-0005 continuous Euclidean space and ADR-0006 physics laws
* Good, because supports future mechanics (knockback forces, explosions, momentum transfer)
* Good, because elastic collisions with momentum transfer create emergent tactical opportunities
* Good, because variable friction zones (ice, mud) create positional strategy through terrain effects

### Thrust Capacity Decision

* Good, because universal constant simplifies bot design (no thrust capacity stat to allocate)
* Good, because Mass already provides sufficient mobility differentiation
* Neutral, because leaves room for future equipment-based variable thrust engines
* Neutral, because numeric value (max thrust) requires playtesting to balance with Mass

### Implementation Complexity

* Neutral, because force-based control requires learning curve (mitigated by documentation and examples)
* Neutral, because per-tick physics calculations have modest computational overhead
* Neutral, because numeric values (friction coefficients, base mass, max thrust) need playtesting
* Bad, because more complex API than simple velocity-setting model
* Bad, because trajectory prediction requires integration (harder than discrete steps)
* Bad, because requires continuous thrust application (bots must act every tick to sustain movement)

### Integration and Consistency

* Good, because clear separation of concerns: ADR-0006 = physics laws, ADR-0007 = movement mechanics
* Good, because physics laws apply consistently regardless of movement model
* Good, because Mass characteristic directly influences movement through acceleration formula
* Good, because friction creates natural coupling between movement and terrain

## Confirmation

The decision will be confirmed through:

1. **Movement API Implementation**:
   - Bot SDK exposes `apply_thrust(force_x: float, force_y: float)` method
   - Thrust parameters validated and clamped to maximum capacity
   - Movement action successfully processes thrust commands each tick

2. **Physics Update Validation**:
   - Per-tick update correctly calculates: forces → acceleration → velocity → position
   - Friction force correctly calculated: F = μ(pos) × M × |v|
   - Friction opposes velocity direction (not arbitrary direction)
   - Acceleration correctly calculated: A = F_net / M
   - Velocity correctly updated: v_new = v_old + A × dt
   - Position correctly updated: pos_new = pos_old + v_new × dt

3. **Terminal Velocity Verification**:
   - Terminal velocity emerges correctly: v = F_thrust / (μ × M)
   - Heavier bots reach lower terminal velocity than light bots
   - Variable friction zones affect terminal velocity appropriately
   - No artificial speed caps needed (physics alone limits speed)

4. **Mass Integration**:
   - Equipment-derived Mass correctly reduces acceleration
   - Heavy equipment loadouts noticeably reduce bot mobility
   - Light loadouts enable nimble, responsive movement
   - Acceleration formula (A = F/M) validated across different Mass values

5. **Boundary Handling**:
   - Position clamped correctly at arena edges
   - Velocity reflected elastically at boundaries (not lost)
   - Bots can thrust away from boundaries and escape
   - Clamping prevents bots exiting arena

6. **Collision Physics**:
   - Bot-to-bot collisions use elastic formula correctly
   - Momentum conserved in collisions
   - Heavy bots push light bots as expected
   - Velocity components updated correctly
   - Collision detection prevents overlapping bots

7. **Variable Friction Zones**:
   - Different arena zones have different friction coefficients
   - Friction correctly reduces velocity based on terrain
   - Terminal velocity varies by zone (ice faster, mud slower)
   - Smooth transitions between friction zones

8. **Determinism**:
   - Same movement commands produce identical results
   - Floating-point calculations consistent across platforms
   - No chaotic sensitivity to initial conditions
   - Bots can reliably predict trajectories

9. **Developer Accessibility**:
   - Bot SDK documentation explains thrust-based movement clearly
   - Example bots demonstrate movement control (pursuit, kiting, flanking)
   - Trajectory calculation utilities provided
   - Physics formulas exposed for bot AI implementation

10. **Performance Testing**:
    - Physics update loop executes in <16ms for 60 ticks/second (60+ Hz target)
    - Multiple simultaneous bots move smoothly without lag
    - Collision detection performs efficiently
    - No physics-based performance bottlenecks

11. **Gameplay Validation**:
    - Movement feels responsive and controllable
    - Mass-based mobility differences are noticeable and meaningful
    - Friction creates engaging terrain-based tactical decisions
    - No dominant movement strategy (various approaches viable)
    - Skill expression through thrust management evident

12. **Integration with ADRs**:
    - ADR-0006 refactored to remove thrust assumptions, remains physics-agnostic
    - ADR-0008 (renumbered from 0007) updated to reference ADR-0007 for movement mechanics
    - Cross-references in all ADRs accurate and non-circular

## Pros and Cons of the Options

### Option 1: Velocity Model

Bots set a desired velocity vector (vx, vy). Game moves bot based on velocity each tick. Friction continuously reduces velocity. Bot must repeatedly set velocity to maintain movement.

**How it works**:
- Bot command: `set_velocity(vx: float, vy: float)`
- Per-tick: Apply friction to velocity, then update position from velocity
- Friction reduces velocity each tick; bot must reapply velocity to maintain speed
- Velocity can change instantly (no acceleration phase)

**Integration with Friction**:
- Friction decays velocity: `v_new = v_current - (μ × M × |v| × dt)`
- Bot must reapply velocity to overcome friction
- Natural deceleration when bot stops setting velocity
- Terminal velocity emerges naturally

**Integration with Mass**:
- Mass affects friction magnitude (heavier bots decelerate faster)
- Mass does NOT affect initial velocity setting (instant velocity change)
- **Problem**: Violates realistic physics (instant velocity ignores inertia and Mass)

**Pros**:
- Simple API (just set desired velocity)
- Direct control over movement direction
- Easy to calculate trajectories (just vector math)
- Minimal computational overhead
- Natural integration with friction (velocity decay)
- Intuitive for developers (familiar from many games)

**Cons**:
- **Violates realistic physics**: Instant velocity changes ignore Mass and inertia
- **Poor Mass integration**: Mass only affects friction, not acceleration - heavy bots don't feel "heavy"
- **No acceleration phase**: Bots teleport to target velocity (unrealistic)
- **Limited tactical depth**: No momentum management or acceleration gameplay
- **Weak equipment tradeoffs**: Heavy equipment penalty is just faster deceleration (minimal impact)
- **No momentum mechanics**: Collisions would need special handling
- **Inconsistent with ADR-0006**: Physics laws assume forces and acceleration

### Option 2: Thrust Model (CHOSEN)

Bots apply thrust force. Game calculates acceleration (A = F/M), updates velocity, applies friction, updates position. Bots must continuously apply thrust to sustain movement.

**How it works**:
- Bot command: `apply_thrust(force_x: float, force_y: float)`
- Per-tick physics update:
  1. Calculate net force: F_net = F_thrust - F_friction
  2. Calculate acceleration: A = F_net / M
  3. Update velocity: v_new = v_current + A × dt
  4. Update position: pos_new = pos_current + v_new × dt
- Friction opposes thrust force
- Bots reach terminal velocity when thrust equals friction
- Natural deceleration when thrust stops

**Integration with Friction**:
- Friction opposes thrust: F_friction = μ(pos) × M × |v|
- Net force: F_net = F_thrust - F_friction
- Acceleration emerges from net force: A = F_net / M
- Terminal velocity: v_term = F_thrust / (μ × M)
- Variable friction zones affect terminal velocity

**Integration with Mass**:
- **Excellent integration**: Acceleration directly proportional to thrust and inversely proportional to Mass
- Heavy bots (high M) accelerate slower from same thrust
- Light bots (low M) are nimble and responsive
- Creates strong mobility-firepower tradeoff through equipment Mass
- Heavier equipment → higher Mass → lower acceleration → harder to move quickly

**Pros**:
- **Realistic physics**: Force → Acceleration → Velocity → Position (standard physics)
- **Perfect Mass integration**: A = F_thrust / M means heavy bots are truly heavier
- **Strong equipment tradeoffs**: Heavy weapons/armor significantly reduce mobility
- **Tactical depth**: Momentum management, acceleration/deceleration gameplay, pursuit/kiting
- **Consistent with ADR-0006**: Physics laws already assume this model
- **Deterministic and calculable**: Standard force-based physics with predictable trajectories
- **Supports future mechanics**: Naturally supports knockback, collisions, explosions
- **Skill expression**: Managing thrust for optimal movement and positioning
- **Terminal velocity emerges**: No arbitrary speed caps (physics determines max speed)
- **Elastic collisions work naturally**: Momentum transfer based on Mass and velocity

**Cons**:
- **Complex API**: Bot developers must think in forces, not positions
- **Harder trajectory calculation**: Must integrate forces over time (not simple vector math)
- **Requires continuous thrust**: Bots must apply thrust every tick to sustain movement (not passive)
- **Computational overhead**: Force summation, acceleration calculation each tick
- **Learning curve**: Force-based control less intuitive than direct position/velocity control
- **Numeric tuning needed**: Max thrust, base mass, friction coefficients require playtesting

### Option 3: Steps Model

Bot requests to move a distance in a direction. Game moves bot that distance in a single tick (or configurable duration). Bot must issue new step command each tick to continue moving.

**How it works**:
- Bot command: `move_step(distance: float, angle: float)` or `move_step(dx: float, dy: float)`
- Per-tick: Move bot by specified distance (with friction penalty), check boundaries
- Discrete movement: Bot either moves the requested distance or is blocked
- Friction applies as penalty to step distance: `actual_distance = requested_distance × friction_modifier(M, terrain)`

**Integration with Friction**:
- Friction reduces effective step distance per tick
- Discrete penalty (not continuous force opposing movement)
- Different terrains have different distance penalties
- **Awkward mismatch**: Friction is continuous force in physics (ADR-0006), but discrete penalty here
- Uncomfortable dissonance between physics model and movement model

**Integration with Mass**:
- Heavy bots (high M) move less distance per step
- Step distance formula: `d_effective = d_requested × mass_modifier(M)`
- Reasonable mobility reduction, but not physics-based

**Pros**:
- Simple API (bot specifies distance and direction)
- Intuitive (like moving a game piece on a board)
- Easy trajectory planning (discrete steps with fixed distances)
- Bounded per-tick movement (prevents runaway velocity)
- Deterministic (same step = same result)
- No floating-point precision issues
- Minimal computational overhead

**Cons**:
- **Poor friction integration**: Friction is continuous force (ADR-0006), not discrete penalty
- **Discrete feel in continuous space**: Grid-like movement in continuous 2D Euclidean space (mismatch with ADR-0005)
- **No momentum or physics**: Steps are teleports within bounds, not physical movement
- **Collision handling awkward**: No velocity for momentum transfer (bots have no "speed" in physics sense)
- **Limited tactical depth**: No acceleration, deceleration, or momentum management gameplay
- **Inconsistent with ADR-0006**: Physics laws assume continuous forces, not discrete penalties
- **Artificial and game-y**: Doesn't feel like physical robots moving with realistic physics
- **Weak Mass integration**: Mass penalty is multiplicative factor, not force-based
- **Terminal velocity meaningless**: No velocity accumulation (just discrete steps)

## More Information

### Related Documentation

- **[ADR-0005: BattleBot Universe Topological Properties](0005-battlebot-universe-topological-properties.md)**: Defines 2D continuous Euclidean space, rectangular boundaries, and continuous movement requirement that thrust-based movement operates within

- **[ADR-0006: BattleBot Universe Physics Laws](0006-battlebot-universe-physics-laws.md)**: Defines physics laws (surface friction F = μ × M × |v|, air friction, elastic collisions, Mass property) that thrust-based movement relies on. Friction physics and Mass-based acceleration are core to thrust mechanics

- **[ADR-0008: Bot Characteristics System](0008-bot-characteristics-system.md)**: Defines Mass characteristic (base + equipment-derived) that affects thrust-to-acceleration conversion. Movement force-to-velocity relationship depends on Mass

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for real-time streaming of position updates generated by thrust-based physics

### Implementation Notes

**Physics Integration Design**:

The thrust model creates a cohesive physics framework where all systems reinforce each other:

- **Thrust + Friction + Mass = Movement Mechanics**
  - Bot applies thrust force each tick
  - Friction opposes movement (magnitude depends on velocity and Mass)
  - Net force determines acceleration: A = (F_thrust - F_friction) / M
  - Acceleration updates velocity; velocity updates position

- **Terminal Velocity Emerges**
  - When F_thrust = F_friction, net force = 0, acceleration = 0
  - Velocity stabilizes at: v = F_thrust / (μ × M)
  - No artificial speed caps needed (physics alone limits speed)

- **Variable Friction Zones Affect Mobility**
  - Ice terrain (μ = 0.2): higher terminal velocity, easier movement
  - Normal terrain (μ = 0.5): medium terminal velocity
  - Mud terrain (μ = 1.5): lower terminal velocity, harder movement
  - Creates tactical positioning decisions (seek ice for speed, use mud for defensive chokepoints)

- **Equipment Mass Directly Affects Mobility**
  - Heavy weapons increase M, reducing A = F / M
  - Heavy armor increases M, reducing acceleration
  - Heavy equipment creates meaningful mobility penalty
  - Light loadouts sacrifice firepower/protection for mobility
  - Creates strategic equipment loadout decisions

- **Elastic Collisions with Momentum Transfer**
  - Bot-to-bot collisions transfer momentum based on Mass and velocity
  - Heavy bots push light bots (greater momentum)
  - Light bots cannot easily move heavy bots
  - Creates tactical collision-based strategies
  - Momentum conservation ensures physically realistic outcomes

**Numeric Values** (all TBD, require playtesting):

- **Base Mass**: M_base (placeholder value needed)
  - Initial mass of bot without equipment
  - Affects acceleration and friction
  - Every bot starts with same base mass (equipment differentiation added separately)

- **Maximum Thrust Capacity**: F_max (placeholder value, e.g., 100 force units)
  - Universal constant; all bots can apply up to F_max force
  - Higher values enable faster acceleration
  - Lower values make movement more constrained
  - Requires tuning for responsive-but-controlled feel

- **Surface Friction Coefficients**: μ values by terrain type
  - μ_ice ≈ 0.2 (low friction, high terminal velocity)
  - μ_normal ≈ 0.5 (standard friction)
  - μ_mud ≈ 1.5 (high friction, low terminal velocity)
  - Affects how quickly friction opposes thrust

- **Game Tick Duration**: dt (placeholder value, likely 1/60 second)
  - Time interval between physics updates
  - Affects velocity and position change per tick
  - Tied to game server tick rate (ADR-0004)

- **Equipment Mass Contributions**: M_equipment for each item
  - Different weapons have different mass (light rifle, heavy shotgun)
  - Different armor has different mass (light plating, heavy plating)
  - Modules contribute mass based on complexity
  - Affects total bot mass: M_total = M_base + Σ(M_equipment)

**Refinement Process**:

These numeric values will be refined through:

1. **Physics Simulation**: Model stat interactions and balance implications
2. **Playtesting**: Test movement feel with different Mass profiles (light vs heavy bots)
3. **Balance Analysis**: Ensure equipment Mass penalties are meaningful but not punishing
4. **Friction Tuning**: Adjust friction coefficients for appropriate movement speed ranges
5. **Thrust Capacity Tuning**: Ensure max thrust creates responsive-but-controlled movement
6. **Competitive Testing**: Identify dominant strategies and adjust numeric values accordingly

**Key Design Insights**:

- Thrust model emerges naturally from ADR-0006 physics laws (friction, Mass, acceleration)
- Mass creates mobility-firepower tradeoff without needing separate thrust capacity characteristic
- Terminal velocity formula (F / (μ × M)) provides predictable movement limits
- Force-based control enables realistic physics while maintaining complete determinism
- Separation of concerns: ADR-0006 defines laws of physics, ADR-0007 defines how bots move
- Continuous thrust requirement creates natural engagement (bots must keep acting)
- Variable friction zones enable terrain-based tactical decisions
- No arbitrary speed caps (physics alone determines terminal velocity)

**Future Considerations**:

- **Thrust Capacity as Characteristic**: If universal constant feels too limiting, could add thrust capacity as equipment-derived characteristic (engines, thrusters provide thrust)
- **Burst/Boost Mechanic**: Short-term burst of extra thrust beyond capacity (with cooldown or energy cost)
- **Directional Thrust Constraints**: Some bots might have directional movement constraints (tank-like: forward/backward thrust, separate rotation)
- **Energy Cost for Thrust**: Tie thrust to energy resource for additional strategic layer (future ADR)
- **Momentum-Based Knockback**: Weapon knockback could be modeled as impulse force (sudden acceleration)
- **Terrain Effects on Traction**: Different terrains could affect effective thrust (ice slippery, mud sticky)
- **Rotational Movement**: Could add angular velocity and rotation (currently treating as 2D point movement, not rigid body)

### Design Principles

The thrust-based movement system follows these principles:

- **Physics First**: Movement follows realistic force-based physics for predictability and depth
- **Mass Integration**: Equipment weight directly affects mobility through acceleration formula
- **Continuous Forces**: Friction and thrust are continuous, not discrete (consistent with ADR-0006)
- **Deterministic Behavior**: Same forces produce identical results (enabling bot AI prediction)
- **Tactical Depth**: Movement enables strategic positioning, momentum management, and skill expression
- **Developer Accessibility**: Complex physics abstracted into simple thrust API with documentation
- **Extensibility**: Framework supports future mechanics (knockback, explosions, variable engines)
- **Separation of Concerns**: Physics laws (ADR-0006) vs movement mechanics (ADR-0007) clearly separated
- **Engagement Guarantee**: Friction and boundaries combine to prevent indefinite evasion
- **Equipment Tradeoffs**: Heavy equipment creates meaningful, physics-based mobility penalties
