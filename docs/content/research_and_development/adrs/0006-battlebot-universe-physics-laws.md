---
title: "[0006] BattleBot Universe Physics Laws"
description: >
    Physics properties defining movement, projectile behavior, and collision mechanics in the battle space
type: docs
weight: 6
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

The BattleBot Universe requires well-defined physics laws to govern gameplay mechanics. ADR-0005 established the mathematical and spatial foundation (2D Euclidean continuous space, rectangular boundaries, Cartesian coordinates), but did not define the physical forces and interactions that govern how objects move, collide, and behave within that space.

This decision defines four fundamental physics properties that characterize the physical laws of the battle space:

1. **Mass** - Universal entity property defining physical weight
2. **Surface Friction** - Resistance force opposing bot movement in the 2D plane
3. **Air Friction** - Resistance force affecting projectile movement through the air
4. **Collisions** - Mechanics governing bot-to-bot and bot-to-wall interactions

These properties form the physics framework upon which all movement mechanics, projectile behavior, collision resolution, and force-based interactions are built. The choice of physics properties has cascading implications for:

- Bot movement characteristics and mobility (integration with Mass)
- Weapon effectiveness and range mechanics (future weapon ADRs)
- Tactical positioning and collision-based strategies
- Computational complexity and determinism
- Developer predictability and accessibility

Without well-defined physics laws, we cannot:

- Implement consistent movement mechanics
- Design weapon systems with predictable projectile behavior
- Resolve collisions deterministically
- Create physics-based tactical depth
- Provide bot developers with predictable physics for strategy implementation
- Ensure deterministic gameplay across all platforms

## Decision Drivers

* **Deterministic Gameplay** - Physics must be consistent and repeatable across all platforms and battles
* **Spatial Integration** - Must work seamlessly with 2D Euclidean continuous space from ADR-0005
* **Tactical Depth** - Physics should enable meaningful strategic positioning and movement decisions
* **Computational Efficiency** - Physics calculations must be performant in real-time tick-based gameplay
* **Developer Predictability** - Bot developers should be able to predict and calculate physics behavior
* **Simplicity First** - Start with simple, constant physics models that can be enhanced in future ADRs
* **Equipment Support** - Physics must integrate with equipment system (ADR-0008) that modifies Mass
* **Future Weapons** - Must support future weapon mechanics without requiring physics redesign
* **Universal Mass Property** - Every entity (bots, projectiles) must have mass
* **Realism vs Gameplay** - Balance realistic physics with engaging gameplay mechanics

## Decision Outcome

We define four fundamental physics properties that create a coherent physics framework for the BattleBot Universe. Each property follows the property-based decision structure from ADR-0005, with multiple options evaluated and chosen options with detailed rationale.

The chosen physics framework consists of:

1. **Mass** - Universal entity property (every bot and projectile has mass)
2. **Surface Friction** - Variable Friction
3. **Air Friction** - Constant Uniform Air Resistance
4. **Collisions** - Elastic Collisions

This physics framework integrates with the spatial system (ADR-0005), equipment system (ADR-0008), and future thrust-based movement to create deterministic, tactically deep gameplay with predictable behavior for bot developers.

## Property 0: Mass

**Every entity in the BattleBot Universe has mass.** This is a fundamental property that affects friction, collisions, projectile behavior, and all force-based interactions.

### Mass Specification

**Bot Mass**:
- **Composition**: Base mass + equipment mass
- **Base Mass**: Every bot has an intrinsic starting mass (TBD: placeholder value)
- **Equipment Mass**: Equipment from ADR-0008 (weapons, armor, modules) contributes to total mass
- **Total Mass Formula**: `M_total = M_base + Σ(M_equipment)`
- **Dynamic Property**: Changes based on equipped items

**Projectile Mass**:
- **Weapon-Specific**: Each weapon type defines projectile mass (TBD)
- **Impacts Behavior**: Mass affects projectile momentum, air friction effects, and collision outcomes
- **Examples**: Rifle projectiles lighter than shotgun pellets

### Mass Impact on Physics

- **Friction**: Higher mass experiences greater friction force (F_friction = μ × M × |v|)
- **Collisions**: Mass determines momentum transfer in collisions (heavier entities push lighter ones)
- **Movement**: Higher mass reduces acceleration from thrust (A = F / M)
- **Projectiles**: Mass affects air resistance and collision damage

### Consequences

* Good, because universal mass property enables consistent physics calculations across all entity types
* Good, because mass provides tuning parameter for equipment balance (heavier weapons, heavier armor)
* Good, because mass creates natural mobility-firepower tradeoffs
* Good, because single mass property simplifies physics implementation
* Neutral, because mass is not directly allocated by developers (derived from equipment choices)
* Bad, because equipment mass must be carefully tuned to avoid dominant builds

---

## Property 1: Surface Friction

**Chosen: Option 1.3 - Variable Friction**

### Options Considered

- **Option 1.1**: None (frictionless surface)
- **Option 1.2**: Constant Uniform Friction
- **Option 1.3**: Variable Friction (CHOSEN)

### Rationale for Chosen Option

- **Terrain-Based Tactical Depth**: Variable friction zones create interesting terrain effects (ice zones, mud zones) enabling tactical positioning decisions
- **Strategic Positioning**: Players can exploit low-friction areas for speed advantages or avoid high-friction areas, adding positional complexity
- **Future Biome Support**: Enables future "biome" mechanics with different terrain types that naturally have different friction characteristics
- **Gameplay Variety**: Different arena regions have distinct movement characteristics, encouraging dynamic tactics and positioning
- **Equipment Balance Opportunities**: Equipment that modifies friction interaction or gives better traction creates additional build diversity
- **Real-World Intuition**: Players intuitively understand that different surfaces have different friction properties
- **Enhanced Strategic Decision-Making**: Movement planning becomes more complex as players must navigate variable friction zones
- **Long-Term Extensibility**: Provides foundation for future weather effects (mud, ice storms) and hazard-based mechanics

### Implementation Formula

```
F_friction = μ(position) × M × |v|
```

Where:
- `μ(position)` = position-dependent surface friction coefficient (varies by arena location/terrain type)
- `M` = bot Mass (calculated from base mass + equipment mass)
- `|v|` = magnitude of bot velocity vector

### Implementation Notes

- Friction coefficient varies based on bot's current position in arena (terrain-dependent)
- Arena map defines friction zones with different coefficient values (e.g., μ_ice = 0.2, μ_mud = 1.5)
- Friction force opposes velocity direction (acts opposite to velocity vector)
- Friction applies continuously to all moving bots every game tick
- Stationary bots (v = 0) have no friction force (no static friction model in initial design)
- Heavier bots (higher Mass from equipment) experience proportionally greater friction force
- Friction force interacts with thrust force to determine acceleration
- Transitions between friction zones should be smooth to prevent unrealistic instant changes

### Alternative Rejected: Option 1.1 - None (Frictionless)

Would create perpetual motion where bots never stop unless they hit boundaries, making movement control extremely difficult and unintuitive. Eliminates integration with Mass characteristic.

### Alternative Rejected: Option 1.2 - Constant Uniform Friction

While simpler to implement, constant friction eliminates opportunities for terrain-based tactical depth. Variable friction provides more engaging gameplay with positional decision-making while still maintaining consistent physics calculations. The framework supports both approaches through the friction coefficient, so variable friction doesn't preclude simpler terrain designs in future arenas.

---

## Property 2: Air Friction

**Chosen: Option 2.2 - Constant Uniform Air Resistance**

### Options Considered

- **Option 2.1**: None (no air resistance)
- **Option 2.2**: Constant Uniform Air Resistance (CHOSEN)

### Rationale for Chosen Option

- **Projectile Range Limiting**: Air friction creates natural maximum range for weapons without arbitrary cutoffs, enabling range-based tactics
- **Tactical Range Decisions**: Weapon effectiveness decreases with distance due to velocity decay, creating optimal engagement ranges
- **Velocity-Based Decay**: Faster projectiles experience more air resistance (quadratic relationship), balancing high-velocity weapons
- **Computational Efficiency**: Constant air resistance coefficient is simple to calculate per projectile per tick
- **Deterministic Behavior**: Same air resistance everywhere ensures predictable projectile paths for weapon aiming
- **Future Weapon Support**: Air friction applies uniformly to all projectile-based weapons (bullets, rockets, thrown objects)
- **No Arbitrary Cutoffs**: Projectiles slow down naturally rather than disappearing at arbitrary range limits

### Implementation Formula

```
F_air = k × v²
```

Where:
- `k` = air resistance coefficient (TBD, requires weapon balance testing)
- `v` = projectile velocity magnitude

### Implementation Notes

- Air friction applies only to projectile objects (bullets, rockets, future thrown weapons)
- Air friction does NOT apply to bots (bot movement governed by surface friction only)
- Projectile velocity decreases over time due to air resistance
- Higher initial velocity = faster deceleration (quadratic relationship)
- Air friction works in combination with gravity to limit projectile lifetime

### Alternative Rejected: Option 2.1 - None

Would create infinite-range projectiles or require arbitrary hard range cutoffs, eliminating range-based tactical decisions and weapon diversity.

---

## Property 3: Collisions

**Chosen: Option 3.2 - Elastic Collisions**

### Options Considered

- **Option 3.1**: Inelastic Collisions (energy absorbed, bots stop or slow significantly)
- **Option 3.2**: Elastic Collisions (CHOSEN)
- **Option 3.3**: Hybrid (context-dependent collision types)

### Rationale for Chosen Option

- **Physics-Based Interactions**: Elastic collisions create momentum transfer based on Mass, rewarding strategic Mass choices
- **Mass-Based Mechanics**: Heavy bots displace light bots in collisions, creating meaningful consequences for equipment loadout Mass
- **Tactical Positioning**: Bots can use collisions to push opponents into unfavorable positions (corners, walls)
- **Knockback Mechanics**: Elastic collisions enable weapon knockback effects for future weapon systems
- **Deterministic Calculations**: Standard elastic collision formulas are well-defined, predictable, and mathematically rigorous
- **Strategic Depth**: Bot developers must account for collision physics in pathfinding, positioning, and tactical decisions
- **Standard Physics**: Elastic collision model is familiar to developers from physics education and game development

### Implementation Formulas

**Bot-to-Bot Collision** (1D collision along collision normal):

```
v1' = ((m1 - m2) × v1 + 2 × m2 × v2) / (m1 + m2)
v2' = ((m2 - m1) × v2 + 2 × m1 × v1) / (m1 + m2)
```

Where:
- `m1`, `m2` = bot masses (base mass + equipment mass)
- `v1`, `v2` = bot velocities before collision (projected onto collision normal)
- `v1'`, `v2'` = bot velocities after collision (along collision normal)

**Bot-to-Wall Collision** (perfect reflection):

```
v_reflected = -v_incident
```

(Wall has infinite mass, perfect elastic reflection)

### Implementation Notes

- **Collision Detection**: Continuous collision detection to prevent bots tunneling through each other or walls
- **Collision Normal**: Collision calculations performed along the line connecting bot centers (1D reduction)
- **Tangent Conservation**: Velocity component perpendicular to collision normal is preserved
- **Coefficient of Restitution**: e = 1.0 (perfectly elastic) - can be tuned if needed
- **Mass Advantage**: Heavy bots (high equipment Mass) push light bots more effectively
- **Boundary Collisions**: Bot-to-wall collisions reflect velocity perfectly (no energy loss)

### Alternative Rejected: Option 3.1 - Inelastic

Would cause bots to lose significant velocity on collision, making movement feel sluggish and frustrating. Eliminates Mass-based tactical advantages in collisions.

### Alternative Rejected: Option 3.3 - Hybrid

Adds complexity without clear benefit for initial implementation. Can be introduced later for specific mechanics.

---

## Physics System Implementation Specification

### Integrated Bot Movement Physics

1. Bots apply thrust force to accelerate
2. Surface friction (μ × M × |v|) opposes movement, creating velocity decay
3. Heavier bots (higher Mass from ADR-0008) require more thrust to overcome friction
4. Bots reach terminal velocity when thrust force equals friction force
5. Collisions transfer momentum based on Mass (elastic collision model)

### Integrated Projectile Physics

1. Projectiles launch with initial velocity (weapon-specific)
2. Air friction (k × v²) causes velocity decay over distance
3. Projectiles disappear when velocity reaches zero or when hitting bots/walls
4. Effective range determined by air friction velocity decay

### Integrated Collision Physics

1. Bot-to-bot collisions use elastic collision formula with momentum conservation
2. Bot-to-wall collisions reflect velocity (wall has infinite mass)
3. Collision detection prevents bots from overlapping or tunneling through walls
4. Mass determines collision outcomes (heavy pushes light)

### Physics Constants (All TBD)

- Surface friction coefficient (μ): TBD
- Air resistance coefficient (k): TBD
- Coefficient of restitution (e): 1.0 (perfectly elastic, may be tuned)

### Deterministic Implementation Requirements

- All physics calculations use deterministic algorithms
- Floating-point calculations must be consistent across platforms
- Physics tick rate: TBD (must align with game server tick rate from ADR-0004)
- Physics update order (per tick):
  1. Calculate all forces (thrust, friction, air resistance, gravity)
  2. Update velocities based on forces
  3. Update positions based on velocities
  4. Resolve collisions and adjust velocities/positions

---

## Consequences

### Overall Integration

* Good, because four physics properties create coherent and predictable gameplay framework
* Good, because physics integrates seamlessly with spatial system (ADR-0005), equipment loadout (ADR-0008), and future thrust actions
* Good, because constant friction and air resistance enable deterministic movement and projectile calculations
* Good, because elastic collisions create tactical depth through Mass-based momentum transfer
* Good, because air friction alone determines projectile range without need for virtual third dimension
* Good, because physics framework is extensible for future enhancements (variable friction zones, terrain effects, weather modifiers)
* Good, because property-based structure allows independent tuning and future modifications
* Good, because universal mass property simplifies physics implementation

### Surface Friction (Variable)

* Good, because creates natural velocity decay and prevents perpetual motion artifacts
* Good, because integrates with Mass to reward/penalize bot weight choices from equipment
* Good, because enables thrust-based movement requiring continuous force application
* Good, because terrain-based friction zones add tactical positioning depth to gameplay
* Good, because naturally supports future "biome" and weather-based mechanics
* Good, because players have intuitive understanding of how different surfaces affect movement
* Neutral, because variable friction zones require arena design and tuning for each map
* Neutral, because friction lookups per position adds modest computational overhead
* Bad, because bot developers must account for variable friction in pathfinding algorithms
* Bad, because slightly more complex to communicate and document physics behavior

### Air Friction (Constant Uniform)

* Good, because creates natural weapon range limitations without arbitrary cutoffs
* Good, because enables tactical range-based gameplay (optimal engagement distances for different weapons)
* Good, because velocity-based decay (v² relationship) rewards careful aim and positioning
* Good, because provides weapon balance tuning parameter
* Good, because computationally efficient (single coefficient calculation per projectile)
* Neutral, because air resistance coefficient k requires weapon balance testing
* Bad, because adds projectile physics calculation complexity

### Collisions (Elastic)

* Good, because elastic collisions create Mass-based tactical interactions
* Good, because heavy bots gain collision advantage, rewarding Mass optimization from equipment
* Good, because enables knockback mechanics for future weapon systems
* Good, because standard physics formulas are well-documented and deterministic
* Good, because creates emergent tactical opportunities (pushing opponents into corners/walls)
* Good, because bot developers can calculate collision outcomes for strategic planning
* Neutral, because elastic collisions may feel "bouncy" if coefficient of restitution not tuned correctly
* Neutral, because collision physics adds computational overhead for detection and resolution
* Bad, because requires careful tuning to prevent exploitative collision-based strategies

---

## Confirmation

The decision will be confirmed through:

1. **Physics Implementation Validation**:
   - Surface friction correctly reduces bot velocity over time based on Mass
   - Air friction correctly reduces projectile velocity over distance
   - Gravity TTL correctly limits projectile lifetime
   - Elastic collisions correctly transfer momentum based on Mass
   - All physics calculations deterministic across platforms and game ticks

2. **Mass System Implementation**:
   - Bot mass correctly calculated as base mass + equipment mass
   - Projectile mass correctly assigned per weapon type
   - Mass is used in all relevant physics calculations (friction, collisions)

3. **Integration Testing**:
   - Physics integrates correctly with spatial system (ADR-0005) coordinate system and boundaries
   - Equipment system (ADR-0008) Mass contributions correctly modify physics behavior
   - Future thrust actions (ADR-0010) correctly overcome friction forces

4. **Gameplay Validation**:
   - Bot movement feels responsive and controllable with surface friction
   - Projectile behavior feels fair and predictable with air friction and gravity TTL
   - Collision mechanics create interesting tactical decisions without exploits
   - Physics constants (μ, k, g, e) tuned for engaging gameplay through playtesting

5. **Performance Testing**:
   - Physics calculations meet real-time tick rate requirements (60+ ticks per second target)
   - Collision detection performs efficiently for multiple bots in arena
   - Floating-point calculations consistent across platforms (Linux, Windows, macOS)
   - No physics-based performance bottlenecks in worst-case scenarios

6. **Developer Accessibility**:
   - Bot SDK exposes physics constants for bot developer calculations
   - Documentation provides physics formulas for movement and collision prediction
   - Sample bots demonstrate physics-aware movement, aiming, and tactical positioning
   - Physics behavior is predictable enough for bot AI implementation

7. **Balance Testing**:
   - Mass-based friction advantages are meaningful but not overwhelming
   - Mass-based collision advantages create tactical depth without dominance
   - Surface friction doesn't make movement feel sluggish or unresponsive
   - Air resistance doesn't make weapons ineffective at intended ranges
   - Collision mechanics don't enable griefing or camping strategies

---

## Pros and Cons of the Options

### Property 0: Mass

**Universal Entity Property** (CHOSEN)

- Good, because every entity having mass creates physically consistent system
- Good, because mass provides tuning parameter for equipment balance
- Good, because mass creates natural mobility-firepower tradeoffs
- Good, because single universal property simplifies implementation
- Good, because mass integrates with all physics calculations
- Neutral, because mass is derived from equipment (not directly allocated)
- Bad, because equipment mass must be carefully tuned

### Property 1: Surface Friction

**Option 1.1: None (Frictionless Surface)**

- Good, because simplest physics model (no friction calculations)
- Good, because eliminates friction calculation overhead
- Good, because maximizes bot responsiveness (instant acceleration/deceleration)
- Bad, because creates perpetual motion (bots never stop without hitting walls)
- Bad, because makes movement control extremely difficult (no natural deceleration)
- Bad, because no integration with Mass characteristic
- Bad, because unrealistic and unintuitive gameplay feel
- Bad, because eliminates thrust-based movement mechanics (no force to overcome)

**Option 1.2: Constant Uniform Friction**

- Good, because simplest friction model (single coefficient everywhere)
- Good, because predictable movement throughout arena
- Good, because integrates with Mass for mobility-weight tradeoffs
- Good, because creates natural velocity decay when thrust not applied
- Good, because enables thrust-based movement requiring continuous force application
- Good, because computationally efficient (single constant coefficient)
- Good, because deterministic behavior across all battles
- Good, because intuitive movement feel (objects slow down without force)
- Neutral, because requires tuning friction coefficient μ for gameplay feel
- Bad, because eliminates opportunities for terrain-based tactical depth
- Bad, because adds physics calculation complexity compared to frictionless

**Option 1.3: Variable Friction** (CHOSEN)

- Good, because creates interesting terrain effects (ice zones, mud zones)
- Good, because adds tactical positioning depth (seek low-friction areas for speed)
- Good, because enables future "biome" mechanics
- Good, because naturally supports weather-based effects (storms, mud, ice)
- Good, because players intuitively understand different surfaces have different friction
- Good, because provides emergent tactical opportunities through map design
- Good, because enriches gameplay with positional decision-making
- Neutral, because provides additional game design opportunities
- Neutral, because variable friction zones require careful arena design
- Neutral, because friction lookups add modest computational overhead per tick
- Bad, because more complex to implement than constant friction
- Bad, because bot developers must account for variable friction in pathfinding
- Bad, because requires additional arena design and content creation (friction maps)
- Bad, because slightly more difficult to visualize and communicate to developers

### Property 2: Air Friction

**Option 2.1: None (No Air Resistance)**

- Good, because simplest projectile physics model
- Good, because eliminates air friction calculations per projectile
- Good, because maximizes weapon range
- Bad, because requires arbitrary hard range cutoffs for weapons
- Bad, because no natural projectile velocity decay over distance
- Bad, because eliminates range-based tactical decisions
- Bad, because projectiles travel indefinitely or require artificial limits
- Bad, because no weapon balance tuning parameter for effective range

**Option 2.2: Constant Uniform Air Resistance** (CHOSEN)

- Good, because natural weapon range limitations based on velocity decay
- Good, because creates optimal engagement distances for tactical gameplay
- Good, because velocity-based decay (v² relationship) balances high-velocity weapons
- Good, because applies uniformly to all projectile weapons (consistency)
- Good, because computationally efficient (single constant coefficient)
- Good, because deterministic behavior across all battles
- Good, because provides weapon balance tuning parameter (k coefficient)
- Good, because eliminates need for virtual third dimension to limit projectile range
- Neutral, because requires tuning air resistance coefficient k for weapon balance
- Bad, because adds projectile physics calculation complexity

### Property 3: Collisions

**Option 3.1: Inelastic Collisions**

- Good, because simpler collision model (bots lose energy and slow down)
- Good, because prevents bouncy collision feel
- Good, because may feel more "realistic" for robot collisions
- Neutral, because could create sticky collision mechanics
- Bad, because bots lose significant velocity on collision (frustrating movement)
- Bad, because Mass has limited impact on collision outcomes (both bots slow)
- Bad, because no momentum transfer or knockback mechanics
- Bad, because eliminates tactical collision-based positioning strategies

**Option 3.2: Elastic Collisions** (CHOSEN)

- Good, because physics-based momentum transfer
- Good, because Mass determines collision outcomes (heavy pushes light)
- Good, because enables knockback mechanics for future weapon systems
- Good, because creates tactical positioning through collision physics
- Good, because standard formulas are well-documented and predictable
- Good, because deterministic calculations for bot AI planning
- Good, because rewards strategic Mass optimization from equipment loadout
- Good, because creates emergent tactical opportunities (push opponents into walls)
- Neutral, because requires tuning coefficient of restitution e for collision feel
- Neutral, because may feel "bouncy" if not tuned correctly
- Bad, because adds computational overhead for momentum transfer calculations
- Bad, because requires careful balancing to prevent collision-based exploits

**Option 3.3: Hybrid (Context-Dependent Collision Types)**

- Good, because flexibility for different collision scenarios (bot vs wall, bot vs bot)
- Good, because could enable specialty mechanics (equipment that changes collision type)
- Neutral, because provides additional game design opportunities
- Bad, because adds significant implementation and calculation complexity
- Bad, because difficult to predict collision outcomes for bot developers
- Bad, because requires context detection overhead per collision
- Bad, because less deterministic behavior (context-dependent outcomes)

---

## More Information

### Related Documentation

- **[ADR-0005: BattleBot Universe Topological Properties](0005-battlebot-universe-topological-properties.md)**: Mathematical and spatial foundation (2D Euclidean space, rectangular boundaries, Cartesian coordinates) that physics laws operate within

- **[ADR-0008: Equipment and Loadout System](0008-equipment-loadout-system.md)**: Equipment system that modifies Mass through equipment weight, affecting friction and collision physics

- **[ADR-0009: Bot Actions and Resource Management](0009-bot-actions-resource-management.md)**: Actions that operate within physics framework (movement actions interact with friction, combat actions create projectiles)

### Implementation Notes

**Physics Integration Design:**

The physics laws create an integrated framework where:

- **Surface friction + Mass + thrust actions = movement mechanics**
  - Bots apply thrust force to overcome friction
  - Heavier bots (more equipment Mass) require more thrust to achieve same velocity
  - Creates natural mobility-firepower tradeoffs through equipment loadout choices

- **Air friction + initial projectile velocity = weapon range mechanics**
  - Projectile velocity decays based on v² relationship
  - Effective weapon range emerges from air resistance coefficient

- **Gravity (simplified TTL) = projectile lifetime mechanics**
  - Projectile TTL calculated from initial velocity, mass, and gravity constant
  - Works in combination with air friction to limit effective range

- **Elastic collisions + Mass = tactical collision mechanics**
  - Momentum transfer based on bot masses
  - Heavy bots push light bots, creating tactical advantages for Mass optimization

**Numeric Value Refinement:**

All physics constants (μ, k, e) are marked TBD (To Be Determined) and will be refined through:

1. Movement feel testing: Tune surface friction coefficient μ for responsive but controlled movement
2. Weapon balance testing: Tune air resistance coefficient k for appropriate weapon effective ranges and projectile lifetime
3. Collision mechanics testing: Tune coefficient of restitution e for satisfying collision feel
4. Mass integration testing: Validate friction-Mass and collision-Mass interactions create meaningful tradeoffs
5. Competitive gameplay: Identify physics exploits or imbalances through player testing
6. Cross-platform testing: Ensure floating-point physics calculations deterministic across platforms

**Key Design Insights:**

- Constant friction and air resistance provide simplicity, predictability, and computational efficiency
- Air resistance alone determines projectile range without introducing a virtual third dimension
- Elastic collisions create Mass-based tactical depth and emergent positioning strategies
- Universal mass property enables consistent physics across all entity types
- Physics framework integrates seamlessly with spatial system (ADR-0005) and equipment loadout (ADR-0008)
- Extensible design allows future enhancements (variable friction zones, weather effects, terrain modifiers) without redesigning core physics

**Future Considerations:**

- Variable Friction Zones: Different terrain types with different friction coefficients (biomes, hazards, power-up zones)
- Weather Effects: Wind affecting air resistance direction/magnitude, rain affecting surface friction
- Collision Damage: High-velocity collision-based damage for ramming strategies
- Equipment-Modified Collision Type: Equipment that changes collision behavior (shock absorbers for inelastic, spikes for damage)
- Momentum-Based Knockback Weapons: Weapons that apply impulse forces for knockback effects
- Projectile Drag Coefficient Tuning: Per-projectile type air resistance for different weapon characteristics

### Design Principles

The physics laws follow these principles:

- **Simplicity First**: Constant coefficients over position-dependent complexity for initial implementation
- **Deterministic Behavior**: Consistent physics calculations across all platforms and battles
- **Mass Integration**: Physics interacts with Mass to create strategic equipment tradeoffs
- **Tactical Depth**: Physics enables strategic positioning, range management, and collision-based tactics
- **Extensibility**: Framework supports future enhancements without breaking core mechanics
- **Developer Predictability**: Bot developers can calculate and predict physics behavior for AI implementation
- **Realism Balanced with Gameplay**: Physics feels consistent and intuitive but prioritizes engaging gameplay over simulation accuracy
- **Computational Efficiency**: Physics calculations performant in real-time tick-based gameplay model
