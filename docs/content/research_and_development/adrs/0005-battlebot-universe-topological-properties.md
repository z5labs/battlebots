---
title: "[0005] BattleBot Universe Topological Properties"
description: >
  Mathematical and topological foundation defining the spatial structure of the battle space
type: docs
weight: 5
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

Battle Bots requires a rigorous mathematical foundation for the spatial environment where battles occur. We need to define the topological and geometric properties of the "BattleBot Universe" that govern all spatial interactions, movement mechanics, collision detection, and distance calculations.

This decision defines four fundamental topological properties that characterize the battle space:

1. **Dimensionality**: The number of spatial dimensions (2D vs 3D)
2. **Vector Space**: The mathematical space structure (continuous Euclidean space vs discrete lattice)
3. **Coordinate Chart**: The coordinate system for position representation (Cartesian vs Polar)
4. **Boundary**: The manifold topology (unbounded vs bounded space)

These properties form the mathematical foundation upon which all spatial mechanics, bot characteristics, equipment systems, and action mechanics are built. The choice of topological structure has cascading implications for implementation complexity, computational efficiency, tactical depth, and accessibility.

Without well-defined topological properties, we cannot:
- Implement consistent position and movement mechanics
- Calculate distances for ranged attacks and detection
- Define collision detection algorithms
- Establish arena boundaries and out-of-bounds handling
- Provide predictable physics for bot developers
- Design pathfinding algorithms
- Support visualization and rendering systems
- Create a coherent spatial framework for all game mechanics

## Decision Drivers

* **Mathematical Rigor** - Spatial system should have well-defined mathematical properties
* **Implementation Complexity** - Simpler topologies reduce development and computational cost
* **Tactical Depth** - Spatial structure should enable meaningful strategic positioning
* **Accessibility** - Bot developers should have accessible algorithms (pathfinding, movement)
* **Computational Efficiency** - Spatial calculations should be performant in real-time
* **Predictability** - Physics and movement should be deterministic and understandable
* **Visualization Clarity** - Positions must map cleanly to visual representation
* **Engagement Guarantee** - Topology should ensure bots cannot avoid combat indefinitely
* **Protocol Integration** - Must integrate with gRPC protocol for position updates (ADR-0004)
* **Extensibility** - Should support future enhancements (obstacles, terrain, variable arenas)
* **Standard Tooling** - Should leverage standard mathematical libraries and algorithms

## Decision Outcome

### Property 1: Dimensionality

**Chosen: Option 1.1 - 2D**

Rationale:
- **Simplicity First**: 2D pathfinding (A*, Dijkstra) and collision detection (2D circle intersection) are well-documented and accessible
- **Lower Barrier to Entry**: Bot developers familiar with 2D algorithms from games, robotics simulations, and common CS education
- **Computational Efficiency**: 2D physics orders of magnitude less expensive than 3D (O(n²) vs O(n³) for many operations)
- **Visualization Simplicity**: Direct 2D rendering without camera controls, 3D projection, or depth perception complexity
- **Sufficient Strategic Depth**: 2D space provides adequate complexity for flanking, positioning, range control, and tactical movement
- **Aligns with Other Decisions**: Complements bounded rectangular arena and Cartesian coordinates

**Alternative Considered: Option 1.2 - 3D** would add significant complexity to pathfinding (3D A*), collision detection (3D physics), and visualization (3D rendering, camera controls) without proportional gameplay benefit for 1v1 battles. Can be reconsidered for future game modes if aerial combat or vertical positioning becomes strategically important.

### Property 2: Vector Space

**Chosen: Option 2.1 - R^n (Continuous Euclidean Space)**

Rationale:
- **Smooth Continuous Movement**: Enables fluid, realistic movement that integrates naturally with real-time gameplay model
- **Infinite Precision**: Floating-point coordinates allow sub-unit positioning (no grid snapping artifacts)
- **Standard Physics**: Continuous space supports standard physics (velocity, acceleration, friction) using well-established formulas
- **Tactical Positioning Depth**: Precise positioning enables fine-grained strategy (optimal range, exact angles)
- **Real-time Protocol Integration**: gRPC can stream continuous position updates without discrete grid jumps
- **Extensibility**: Supports future terrain effects, obstacles, dynamic boundaries without discretization constraints

**Alternative Considered: Option 2.2 - n-dimensional Lattice** would simplify collision detection (grid occupancy) and eliminate floating-point precision issues, but would sacrifice movement fluidity and create discrete grid-to-grid jump artifacts that feel unnatural in real-time gameplay. Grid-based pathfinding is simpler, but requires discrete approximation that limits tactical positioning depth.

### Property 3: Coordinate Chart

**Chosen: Option 3.1 - Cartesian**

Rationale:
- **Universal Familiarity**: Cartesian (x, y) coordinates are universally taught and understood
- **Algorithmic Simplicity**: Distance, angle, and vector calculations use standard formulas
- **Library Support**: Every programming language has extensive Cartesian math libraries
- **Rectangular Arena Alignment**: Cartesian coordinates naturally align with rectangular boundaries
- **Grid Visualization**: Maps directly to pixel grids for rendering (x → screen x, y → screen y)
- **Pathfinding Compatibility**: A* and pathfinding algorithms designed for Cartesian grids

**Alternative Considered: Option 3.2 - Polar** (r, θ) coordinates would be more natural for rotational mechanics but require trigonometric conversions for most operations, have less library support, and align poorly with rectangular boundaries. Polar coordinates are better suited for radial-specific mechanics (turrets, circular arenas) which are not part of the current design.

### Property 4: Boundary

**Chosen: Option 4.2 - Rectangular Boundary**

Rationale:
- **Engagement Guarantee**: Fixed boundaries prevent indefinite evasion, forcing interaction
- **Predictable Strategy**: Bots can rely on constant arena size for pathfinding and tactical planning
- **Cartesian Alignment**: Rectangular boundaries align perfectly with Cartesian coordinates
- **Simple Collision Detection**: Boundary checks are trivial comparisons (x < min, x > max)
- **Fair and Balanced**: Equal access to arena space for both bots throughout battle
- **Visualization Clarity**: Rectangular arena displays naturally on screen without distortion
- **Complements Timeout**: Fixed boundaries work with timeout mechanism to ensure conclusion

**Alternative Considered: Option 4.3 - Circular Boundary** would provide uniform distance from center but requires distance calculations for boundary checks (more expensive than rectangular comparisons) and aligns poorly with Cartesian coordinates. Circular boundaries also have no corners, which eliminates certain positional strategies.

**Alternative Rejected: Option 4.1 - Unbounded** would allow indefinite evasion, potentially creating stalemate scenarios even with timeout. Unbounded space also complicates pathfinding (no clear boundaries for navigation) and visualization (camera must follow bots across potentially large distances).

### Spatial System Implementation Specification

The four topological properties define the following concrete spatial system implementation:

#### Coordinate System

The battle space uses a **2D Cartesian coordinate system** with the following properties:

- **Dimensionality**: Two-dimensional space (x, y coordinates only; no z-axis or vertical elevation)
- **Origin (0, 0)**: Located at the center of the arena
- **X-axis**: Horizontal axis, with positive values extending to the right and negative values to the left
- **Y-axis**: Vertical axis, with positive values extending upward and negative values downward
- **Units**: Abstract spatial units (not meters, pixels, or other real-world measurements)
- **Precision**: Floating-point coordinates allow for sub-unit positioning accuracy

This centered origin simplifies calculations for distance, angle, and relative positioning between bots. It also provides symmetry for balanced starting positions in various battle configurations.

#### Boundaries

The battle space is defined by rectangular boundaries:

- **Arena Size**: **100 x 100 units** (TBD - subject to tuning based on playtesting)
- **X-axis Range**: -50 to +50 units
- **Y-axis Range**: -50 to +50 units

**Out-of-Bounds Handling**:

1. **Movement Blocking**: Any movement command that would place a bot outside the boundaries is clamped to the nearest valid position at the boundary edge
2. **No Wrapping**: Coordinates do not wrap around (i.e., exiting the right side does not place a bot on the left side)
3. **Boundary Contact**: Bots may be positioned exactly on the boundary line
4. **Force Effects**: External forces (knockback, explosions, etc.) that would push a bot out-of-bounds will stop at the boundary

#### Bot Positioning and Collision Detection

Each bot occupies a circular area within the battle space:

- **Bot Radius**: **2 units** (TBD - subject to balance tuning)
- **Center Position**: Each bot's coordinates (x, y) represent the center of its circular footprint
- **Minimum Separation**: Bots cannot overlap; their circular areas must not intersect

**Collision Detection** uses simple circle-to-circle distance calculations:

1. **Bot-to-Bot Collision**: Two bots collide when the distance between their centers is less than the sum of their radii (2r = 4 units for identical bots)
2. **Bot-to-Wall Collision**: A bot collides with a wall when its center position plus its radius exceeds the boundary
   - Left wall: `x - radius < -50`
   - Right wall: `x + radius > 50`
   - Bottom wall: `y - radius < -50`
   - Top wall: `y + radius > 50`
3. **Collision Resolution**: When a collision is detected, the movement is adjusted to place the bot in contact with the obstacle without overlapping

**Bot-to-Bot Collision Rules**:
1. **Movement Blocking**: Bots cannot move through each other
2. **Collision Physics**: When two bots attempt to occupy overlapping space:
   - The moving bot stops at the point of contact
   - Both bots remain stationary in their final positions (no pushing or displacement)
3. **Damage**: Bot-to-bot collisions do not inherently cause damage (unless specific collision damage mechanics are introduced)

#### Friction and Movement Physics

The battle space applies friction to all moving bots, which affects their velocity and movement behavior. Friction provides realistic physics that require bots to continuously apply force to maintain movement, and it enables the possibility of variable terrain types with different surface properties.

**Friction Mechanics**:

1. **Friction Force**: Opposes the direction of bot movement, proportional to velocity
2. **Friction Coefficient (μ)**: Determines the strength of friction applied to a bot
   - **Default Coefficient**: **0.1** (TBD - subject to balance tuning)
   - **Range**: 0.0 (frictionless) to 1.0 (maximum friction)
3. **Velocity Decay**: Each update tick, a bot's velocity is reduced by the friction force
4. **Natural Stopping**: Without continuous thrust, a bot will gradually slow to a stop due to friction

**Friction Calculation**:

```
friction_force = -μ × velocity
new_velocity = velocity + friction_force
```

Where:
- `μ` is the friction coefficient at the bot's current position
- `velocity` is the bot's current velocity vector
- The negative sign indicates friction opposes the direction of movement

**Variable Friction Zones**:

1. **Uniform Friction**: By default, the entire battle space has a uniform friction coefficient
2. **Friction Zones**: Specific rectangular or circular areas may define different friction values
   - **Low Friction** (μ < 0.1): "Slippery" surfaces where bots slide more easily
   - **Standard Friction** (μ = 0.1): Normal battle space surface
   - **High Friction** (μ > 0.1): "Rough" surfaces that slow bot movement more quickly
3. **Zone Priority**: When friction zones overlap, the highest friction coefficient applies
4. **Transition Behavior**: Moving between friction zones immediately applies the new coefficient (no gradual transition)

**Tactical Implications**:
- **Movement Planning**: Bots must account for deceleration when planning movements
- **Pursuit and Evasion**: Understanding friction helps predict opponent stopping distances
- **Zone Control**: High-friction zones can limit mobility, while low-friction zones enable faster repositioning
- **Energy Management**: Continuous thrust is required to maintain velocity, affecting energy economy

#### Line of Sight

Line of sight determines whether one bot can "see" or target another bot, which is essential for ranged attacks and targeting systems.

**Line of Sight Rules**:

1. **Direct Path**: A bot has line of sight to another bot if an unobstructed straight line can be drawn between their center positions
2. **Obstacle-Free**: Currently, the only obstacles are other bots. A bot blocks line of sight between two other bots if the line passes through its circular area
3. **Boundary Walls**: Walls do not block line of sight (bots can see through walls but cannot shoot through them - weapon-specific rules apply)

**Line of Sight Calculation**:

To determine if Bot A has line of sight to Bot B:

1. Draw a line segment from A's center to B's center
2. For each other bot C in the arena:
   - Calculate the perpendicular distance from C's center to the line segment
   - If this distance is less than C's radius, line of sight is blocked

This calculation may be optimized by only checking bots that fall within a bounding box around the line segment.

#### Movement Constraints

1. **Speed Limits**: Bots have maximum movement speeds (defined in ADR-0006: Bot Characteristics)
2. **No Teleportation**: Bots cannot instantly move from one position to another; all movement follows continuous paths
3. **Elastic Wall Collisions**: When a bot collides with a wall during movement, it stops at the wall position without bouncing
4. **No Wall Damage**: Wall collisions do not inherently cause damage to bots

### Consequences

#### Dimensionality Decision (2D)

* Good, because simplest spatial implementation (2D algorithms, 2D visualization)
* Good, because lower barrier to entry for bot developers
* Good, because sufficient strategic depth for positioning and tactics
* Good, because reduces computational requirements vs 3D
* Good, because integrates seamlessly with rectangular boundaries and Cartesian coordinates
* Neutral, because limits future aerial or vertical combat mechanics
* Bad, because eliminates vertical positioning as strategic dimension
* Bad, because may feel limiting if users expect 3D movement

#### Vector Space Decision (R^n Continuous)

* Good, because enables smooth, realistic movement
* Good, because supports standard continuous physics formulas
* Good, because infinite precision for tactical positioning
* Good, because integrates naturally with real-time gRPC protocol
* Good, because extensible to terrain effects and obstacles
* Neutral, because requires careful floating-point handling
* Bad, because floating-point edge cases more complex than integer lattice
* Bad, because pathfinding requires discretization step

#### Coordinate Chart Decision (Cartesian)

* Good, because universally familiar coordinate system
* Good, because extensive library and tooling support
* Good, because aligns naturally with rectangular boundaries
* Good, because direct mapping to pixel rendering
* Good, because standard distance and angle formulas
* Neutral, because polar coordinates may be more natural for some rotational mechanics
* Bad, because bot developers must calculate angles for directional actions

#### Boundary Decision (Rectangular)

* Good, because guarantees engagement (no indefinite evasion)
* Good, because simple collision detection (4 comparisons)
* Good, because aligns perfectly with Cartesian coordinates
* Good, because predictable arena for pathfinding
* Good, because fair and balanced (symmetric access)
* Good, because complements timeout mechanism for battle conclusion
* Neutral, because could support variable sizes in future
* Bad, because corners enable defensive camping strategies
* Bad, because lacks additional pressure (no shrinking)

#### Overall Integration

* Good, because all four properties create coherent mathematical foundation
* Good, because choices are mutually reinforcing (Cartesian + Rectangular, 2D + R^n)
* Good, because extensible framework supports future enhancements
* Good, because spatial system implementation follows naturally from topological properties
* Good, because creates foundation for ADR-0006 (characteristics), ADR-0007 (equipment), ADR-0008 (actions)
* Good, because property-based decision structure allows independent tuning and future modifications
* Good, because balances mathematical rigor with practical accessibility

### Confirmation

The decision will be confirmed through:

1. **Topological Consistency**: Verify mathematical properties are correctly implemented
   - 2D coordinate representation in all spatial data structures
   - Continuous floating-point position values (no grid snapping)
   - Cartesian coordinate system throughout codebase
   - Rectangular boundary enforcement

2. **Spatial Mechanics Validation**:
   - Implement collision detection using 2D Euclidean distance
   - Boundary collision detection and clamping working correctly
   - Friction physics applied correctly each tick
   - Line of sight calculations accurate

3. **Developer Accessibility**:
   - Bot SDK exposes Cartesian (x, y) coordinates
   - Movement API uses familiar vector representations
   - Documentation includes standard formulas (distance, angle)
   - Sample bots demonstrate pathfinding in continuous 2D space

4. **Performance Testing**:
   - 2D collision detection meets real-time tick rate requirements
   - Floating-point calculations deterministic across platforms
   - Spatial queries (nearest bot, line of sight) performant

5. **Protocol Integration**:
   - gRPC messages correctly encode 2D positions
   - Position updates stream smoothly in continuous space
   - Boundary violations detected and communicated

6. **Visualization Testing**:
   - 2D Cartesian coordinates map correctly to screen pixels
   - Rectangular arena renders clearly
   - Bot positions and movements display accurately

7. **Extensibility Validation**:
   - System can support future obstacles and terrain
   - Variable rectangular arena sizes possible
   - Spatial queries abstracted for future enhancements

8. **Playtesting**:
   - Arena size (100x100) provides appropriate tactical space
   - Bot radius (2 units) creates reasonable collision frequency
   - Friction coefficient (0.1) balances mobility vs stopping control

9. **Future Consideration**:
   - Document path to 3D extension if needed
   - Evaluate polar coordinates for specific mechanics (turrets, sensors)
   - Consider circular boundaries for alternative game modes

## Pros and Cons of the Options

### Property 1: Dimensionality

#### Option 1.1: 2D (CHOSEN)

Battles occur in two-dimensional space (x, y coordinates).

* Good, because simplest spatial implementation (2D pathfinding, collision detection, physics)
* Good, because lower computational requirements compared to 3D
* Good, because easier visualization (direct 2D display, no camera controls needed)
* Good, because lower barrier to entry for bot developers (2D algorithms more accessible)
* Good, because well-documented algorithms (A*, 2D vector math, 2D physics)
* Good, because sufficient strategic depth for positioning and movement tactics
* Neutral, because appropriate for ground-based combat scenarios
* Neutral, because may be extended to 3D in future if needed
* Bad, because eliminates vertical positioning as strategic dimension
* Bad, because no aerial combat or flying units
* Bad, because may feel limiting if users expect 3D movement

#### Option 1.2: 3D

Battles occur in three-dimensional space (x, y, z coordinates).

* Good, because enables vertical positioning strategy (height advantage)
* Good, because supports aerial combat and flying units
* Good, because additional strategic dimension (above/below positioning)
* Good, because familiar from many modern games
* Neutral, because enables jump mechanics or flight equipment
* Neutral, because may be more engaging for some users
* Bad, because significantly more complex implementation (3D pathfinding, collision, physics)
* Bad, because much higher computational requirements (3D calculations expensive)
* Bad, because complex visualization (3D rendering, camera controls, depth perception)
* Bad, because higher barrier to entry (3D algorithms more complex)
* Bad, because more difficult to debug and visualize battles
* Bad, because may be unnecessary complexity for 1v1 ground combat

### Property 2: Vector Space

#### Option 2.1: R^n - Continuous Euclidean Space (CHOSEN)

Positions represented as continuous floating-point coordinates in standard n-dimensional Euclidean space.

* Good, because enables smooth, continuous movement
* Good, because supports standard continuous physics formulas (velocity, acceleration, friction)
* Good, because infinite precision for tactical positioning (no grid snapping)
* Good, because integrates naturally with real-time gRPC protocol
* Good, because extensible to terrain effects, obstacles, dynamic boundaries
* Good, because familiar from most modern games and simulations
* Neutral, because requires careful floating-point handling for determinism
* Neutral, because pathfinding requires discretization for algorithms like A*
* Bad, because floating-point precision can introduce edge cases
* Bad, because collision detection more complex than grid occupancy
* Bad, because requires careful handling of floating-point comparisons

#### Option 2.2: n-dimensional Lattice of R^n

Positions represented as discrete points on an integer lattice (grid).

* Good, because eliminates floating-point precision issues
* Good, because simpler collision detection (grid cell occupancy)
* Good, because natural fit for grid-based pathfinding (A*, BFS)
* Good, because deterministic integer mathematics
* Good, because easier to reason about and debug
* Neutral, because grid resolution determines precision (finer grids approach continuous)
* Bad, because discrete movement creates grid-to-grid jump artifacts
* Bad, because feels unnatural and less fluid in real-time gameplay
* Bad, because limits tactical positioning precision (can't be "between" grid cells)
* Bad, because requires special handling for diagonal movement distances
* Bad, because poor integration with real-time streaming protocol (discrete jumps)

### Property 3: Coordinate Chart

#### Option 3.1: Cartesian (CHOSEN)

Positions represented using Cartesian coordinates (x, y).

* Good, because universally familiar coordinate system
* Good, because algorithmic simplicity (standard distance and angle formulas)
* Good, because extensive library support in every programming language
* Good, because aligns naturally with rectangular boundaries
* Good, because direct mapping to pixel grids for rendering
* Good, because pathfinding algorithms designed for Cartesian grids
* Good, because simple boundary checks (x < min, x > max)
* Neutral, because appropriate for most spatial scenarios
* Bad, because requires trigonometry for angle calculations
* Bad, because polar coordinates may be more natural for rotational mechanics

#### Option 3.2: Polar

Positions represented using polar coordinates (r, θ) - distance and angle from origin.

* Good, because natural for rotational and radial mechanics
* Good, because distance from center is explicit (r coordinate)
* Good, because angle of position is explicit (θ coordinate)
* Neutral, because appropriate for circular arenas or radial gameplay
* Neutral, because familiar from mathematics and physics
* Bad, because requires trigonometric conversions for most operations
* Bad, because less universal familiarity (more complex than Cartesian)
* Bad, because limited library support (often converted to Cartesian internally)
* Bad, because aligns poorly with rectangular boundaries
* Bad, because distance between two polar points requires conversion
* Bad, because visualization requires conversion to screen coordinates

#### Option 3.3: Discrete Cartesian (for Lattice)

Positions represented using discrete Cartesian coordinates on a lattice (integer x, y).

* Good, because familiar Cartesian coordinate system
* Good, because simple integer arithmetic
* Good, because natural for grid-based pathfinding
* Good, because deterministic (no floating-point issues)
* Neutral, because only viable option for lattice vector space
* Bad, because limited to discrete grid positions
* Bad, because creates grid-to-grid jump artifacts
* Bad, because not chosen due to vector space decision (R^n continuous chosen)

### Property 4: Boundary

#### Option 4.1: Unbounded

No arena boundaries (infinite or very large playable area).

* Good, because enables unlimited strategic space
* Good, because no artificial boundaries constraining movement
* Good, because may feel more "realistic" for some scenarios
* Neutral, because appropriate for exploration-focused gameplay
* Neutral, because eliminates boundary collision checks
* Bad, because allows indefinite evasion (bots can run away forever)
* Bad, because enables stalemate even with timeout (avoid combat entire battle)
* Bad, because complicates pathfinding (no clear navigation boundaries)
* Bad, because difficult visualization (camera must track across large distances)
* Bad, because may create boring gameplay (chasing fleeing opponents)
* Bad, because timeout becomes only victory condition (no engagement guarantee)
* Bad, because eliminates positional strategy (no boundaries to control)
* Bad, because fundamentally incompatible with engagement-focused 1v1 gameplay

#### Option 4.2: Rectangular Boundary (CHOSEN)

Fixed rectangular arena boundaries that remain constant throughout the battle.

* Good, because guarantees engagement (bots cannot escape indefinitely)
* Good, because predictable arena enables reliable pathfinding
* Good, because simplest implementation (fixed boundary checks)
* Good, because fair and balanced (equal arena access)
* Good, because enables positional strategy (corner control, center dominance)
* Good, because complements timeout mechanism for battle conclusion
* Good, because aligns perfectly with Cartesian coordinates
* Good, because clear visualization (fixed arena visible throughout)
* Neutral, because could be enhanced with shrinking in future modes
* Neutral, because requires boundary collision detection
* Bad, because may enable corner camping defensive strategies
* Bad, because lacks additional engagement pressure
* Bad, because fixed size may feel static compared to dynamic boundaries

#### Option 4.3: Circular Boundary

Fixed circular arena boundary centered at origin.

* Good, because provides uniform distance from center for all boundary points
* Good, because guarantees engagement (bots cannot escape)
* Good, because symmetric and aesthetically pleasing
* Good, because no corners to enable camping strategies
* Neutral, because appropriate for radial or rotational gameplay
* Neutral, because requires distance calculation for boundary checks
* Bad, because boundary checks more expensive than rectangular (distance vs comparison)
* Bad, because aligns poorly with Cartesian coordinates (requires distance calculations)
* Bad, because more complex collision detection than rectangular boundaries
* Bad, because visualization may require circular clipping
* Bad, because pathfinding must account for curved boundaries

## More Information

### Related Documentation

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for streaming position updates in continuous 2D space

- **[ADR-0006: Bot Characteristics System](0006-bot-characteristics-system.md)**: Speed and Mass characteristics that govern movement in this spatial system

- **[ADR-0007: Equipment and Loadout System](0007-equipment-loadout-system.md)**: Equipment that modifies movement capabilities within this spatial framework

- **[ADR-0008: Bot Actions and Resource Management](0008-bot-actions-resource-management.md)**: Movement and combat actions that operate within this spatial system

- **[POC User Journey](../user-journeys/0001-poc.md)**: Proof of concept implementation using this spatial foundation

**Note**: This ADR supersedes and integrates the former ADR-0006 (Battle Space Spatial System), which is now deprecated. The spatial implementation details previously in ADR-0006 are now part of this ADR's "Spatial System Implementation Specification" section.

### Implementation Notes

**Mathematical Foundation**:

The BattleBot Universe is mathematically defined as:
- **Topological Space**: 2-dimensional Euclidean space R²
- **Manifold**: Closed rectangular region [−50, 50] × [−50, 50] ⊂ R²
- **Metric**: Standard Euclidean metric d(p,q) = √((x₂−x₁)² + (y₂−y₁)²)
- **Coordinate Chart**: Cartesian coordinates φ: R² → R² where φ(p) = (x, y)
- **Boundary**: ∂M = {(x,y) : x=±50 or y=±50}

**Numeric Value Refinement**:

All numeric values (arena size, bot radius, friction coefficient) are marked TBD and will be refined through:

1. Performance profiling of collision detection at various arena sizes
2. Playtesting to tune arena size for engagement frequency
3. Balance analysis of bot radius for collision frequency
4. Friction coefficient tuning for movement feel
5. Visualization testing for rendering clarity
6. Timeout scenario frequency analysis to tune arena size appropriately
7. Equipment balance testing to ensure stat-based equipment choices remain meaningful

**Key Design Insights**:
- 2D + R^n continuous + Cartesian + Rectangular boundaries create mutually reinforcing topological framework
- Mathematical rigor provides clear foundation for implementation
- Simplicity enables focus on core battle mechanics and accessibility
- Complete spatial framework enables users to implement sophisticated pathfinding and AI solutions
- Extensibility allows future enhancements (3D, fog of war, variable arenas) without disrupting foundation

**Future Considerations**:
- **Variable Arena Sizes**: Different battle modes (quick match vs. tournament) may have different dimensions
- **3D Spatial System**: If aerial combat or vertical positioning becomes strategically important, extend to R³ with (x, y, z) coordinates
- **Circular Boundaries**: Could be added as optional game mode for radial gameplay
- **Dynamic Boundaries**: Shrinking arena boundaries could be added as optional game mode for aggressive pacing
- **Obstacles and Terrain**: Spatial system designed to support static/dynamic obstacles and variable terrain effects
- **Polar Coordinate Option**: For specific mechanics (turret rotation, sensor sweeps) polar coordinates could supplement Cartesian system

### Design Principles

The BattleBot Universe topological properties follow these principles:
- **Mathematical Rigor**: Well-defined topological and geometric properties provide clear foundation
- **Simplicity First**: 2D continuous Euclidean space with Cartesian coordinates and rectangular boundaries reduce complexity
- **Accessibility**: Bot developers have familiar coordinate systems and well-documented algorithms
- **Guaranteed Engagement**: Bounded rectangular arena ensures battles conclude with interaction
- **Predictability**: Deterministic physics and movement with clear boundary rules
- **Fairness**: Symmetric arena and equal access to space for all bots
- **Integration**: Spatial framework seamlessly combines with characteristics, equipment, and actions
- **Extensibility**: Property-based structure allows future enhancements without disrupting core mechanics
