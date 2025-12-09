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
- Establish arena boundaries and out-of-bounds handling
- Provide a predictable spatial framework for bot developers
- Design algorithms that work with the spatial structure
- Support visualization and rendering systems
- Create a coherent foundation for all spatial game mechanics

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

### Property 4: Boundary Configuration

**Note**: Boundary configuration has been transferred to **ADR-0011 (1v1 Battles)** where it is treated as a selectable arena property rather than a fixed universal topological property.

**Rationale for Transfer**: While a default rectangular boundary ([-50, 50] × [-50, 50]) was appropriate for establishing the spatial system's mathematical foundation, enabling different game modes to use different arena sizes requires decoupling boundary dimensions from universal topology. ADR-0011 provides comprehensive boundary definition and selection mechanism for 1v1 battle instances.

**Reference**: See **[ADR-0011: 1v1 Battles](0011-1v1-battles.md)** for complete boundary specification, alternative boundary options, and battle instance configuration.

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

Boundary configuration for specific battle types is defined in type-specific ADRs. See **[ADR-0011: 1v1 Battles](0011-1v1-battles.md)** for:

- **Boundary Dimensions**: Arena size and shape configuration (initially 100 x 100 units rectangular)
- **Boundary Selection**: How arena boundaries are chosen for specific battles
- **Out-of-Bounds Handling**: Movement clamping, wrapping rules, boundary contact, force effects
- **Alternative Boundaries**: Future boundary shapes and sizes (circular, hexagonal, dynamic)

**Principle**: The BattleBot Universe defines the universal 2D Euclidean topological space and Cartesian coordinate system. Individual game modes (1v1 battles, team battles, etc.) configure specific bounded regions within this space through arena selection.

#### Movement Constraints

The following basic movement constraints apply to the battle space:

1. **Continuous Movement**: Bots cannot instantly teleport from one position to another; all movement follows continuous paths through the 2D Euclidean space
2. **Boundary Enforcement**: Any movement that would place a bot outside the rectangular boundaries is prevented (specific collision mechanics will be defined in a separate ADR)
3. **No Coordinate Wrapping**: The space does not wrap around (i.e., exiting one side does not place a bot on the opposite side)
4. **Deterministic Physics**: All spatial calculations use deterministic algorithms to ensure consistent behavior across platforms

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

#### Overall Integration

* Good, because all four properties create coherent mathematical foundation
* Good, because choices are mutually reinforcing (Cartesian + Rectangular, 2D + R^n)
* Good, because extensible framework supports future enhancements
* Good, because spatial system implementation follows naturally from topological properties
* Good, because creates foundation for ADR-0006 (physics), ADR-0007 (movement), ADR-0008 (characteristics), ADR-0009 (equipment), ADR-0010 (actions)
* Good, because property-based decision structure allows independent tuning and future modifications
* Good, because balances mathematical rigor with practical accessibility

### Confirmation

The decision will be confirmed through:

1. **Topological Consistency**: Verify mathematical properties are correctly implemented
   - 2D coordinate representation in all spatial data structures
   - Continuous floating-point position values (no grid snapping)
   - Cartesian coordinate system throughout codebase
   - Metric calculations accurate (distance formulas)

2. **Spatial Mechanics Validation**:
   - Continuous movement paths validated
   - Coordinate calculations accurate and deterministic
   - Coordinate transformations correct

3. **Developer Accessibility**:
   - Bot SDK exposes Cartesian (x, y) coordinates
   - Movement API uses familiar vector representations
   - Documentation includes standard formulas (distance, angle)
   - Sample bots demonstrate pathfinding in continuous 2D space

4. **Performance Testing**:
   - 2D coordinate calculations meet real-time tick rate requirements
   - Floating-point calculations deterministic across platforms
   - Spatial queries (distance, angle) performant

5. **Protocol Integration**:
   - gRPC messages correctly encode 2D positions
   - Position updates stream smoothly in continuous space
   - Coordinate system correctly transmitted and received

6. **Visualization Testing**:
   - 2D Cartesian coordinates map correctly to screen pixels
   - Bot positions and movements display accurately
   - Coordinate system visualization clear

7. **Extensibility Validation**:
   - System can support future battle types and game modes
   - Coordinate system independent of specific arena configurations
   - Spatial queries abstracted for future enhancements

8. **Boundary Validation**: See ADR-0011 for boundary-specific confirmation criteria
   - Arena boundaries implemented per ADR-0011
   - Boundary enforcement working correctly for selected arenas
   - Out-of-bounds handling per game mode specification

9. **Future Consideration**:
   - Document path to 3D extension if needed (would add z coordinate)
   - Evaluate polar coordinates for specific game mechanics (turrets, sensors)
   - Consider extending to support non-Euclidean spaces if needed

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

### Property 4: Boundary Configuration

**Note**: Boundary options and pros/cons have been moved to **[ADR-0011: 1v1 Battles](0011-1v1-battles.md)** where boundary selection is defined as a configurable arena property.

For comprehensive evaluation of boundary options (rectangular, circular, hexagonal, dynamic) and their trade-offs, see ADR-0011 Property 2: Boundary.

## More Information

### Related Documentation

**Spatial and Physics Foundation**:
- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for streaming position updates in continuous 2D space

- **[ADR-0006: BattleBot Universe Physics Laws](0006-battlebot-universe-physics-laws.md)**: Physics properties (friction, collisions, gravity) that govern movement mechanics in this spatial system

- **[ADR-0007: Bot Movement Mechanics](0007-bot-movement-mechanics.md)**: Defines how bots apply thrust forces to control their movement within this spatial system

- **[ADR-0008: Bot Characteristics System](0008-bot-characteristics-system.md)**: Mass and other characteristics that govern movement and combat in this spatial system

- **[ADR-0009: Equipment and Loadout System](0009-equipment-loadout-system.md)**: Equipment that modifies characteristics and movement capabilities within this spatial framework

- **[ADR-0010: Bot Actions and Resource Management](0010-bot-actions-resource-management.md)**: Movement and combat actions that operate within this spatial system

**Battle Configuration and Arena System**:
- **[ADR-0011: 1v1 Battles](0011-1v1-battles.md)**: Defines Arena concept and configurable battle properties (boundary, biome, visibility, positioning, win conditions) for 1v1 battles. Boundary ownership transferred from ADR-0005 to ADR-0011.

**Implementation and Testing**:
- **[POC User Journey](../user-journeys/0001-poc.md)**: Proof of concept implementation using this spatial foundation

**Note**: This ADR supersedes and integrates the former ADR-0006 (Battle Space Spatial System), which is now deprecated.

**Future ADRs**: The following topics were removed from this ADR and will be addressed in separate architectural decisions:
- Collision Detection and Bot Positioning (bot size, collision mechanics, collision resolution)
- Friction and Movement Physics (friction coefficients, velocity decay, variable friction zones)
- Line of Sight (visibility calculations, obstacle blocking)

### Implementation Notes

**Mathematical Foundation**:

The BattleBot Universe is mathematically defined as:
- **Topological Space**: 2-dimensional Euclidean space R²
- **Metric**: Standard Euclidean metric d(p,q) = √((x₂−x₁)² + (y₂−y₁)²)
- **Coordinate Chart**: Cartesian coordinates φ: R² → R² where φ(p) = (x, y)
- **Note on Boundaries**: Specific bounded regions are defined in game mode ADRs (e.g., ADR-0011 for 1v1 battles). This ADR defines the universal topological space; game modes define specific arenas within that space.

**Numeric Value Refinement**:

Default arena size (100×100 units) is marked TBD and will be refined through playtesting. See ADR-0011 for specific arena configuration and tuning details.

1. Playtesting to tune arena size for engagement frequency (per game mode)
2. Visualization testing for rendering clarity
3. Timeout scenario frequency analysis (tuned with game mode timeout values)
4. Equipment balance testing to ensure stat-based equipment choices remain meaningful

**Key Design Insights**:
- 2D + R^n continuous + Cartesian create mutually reinforcing topological framework
- Mathematical rigor provides clear foundation for implementation independent of game modes
- Simplicity enables focus on core battle mechanics and accessibility
- Complete spatial framework enables users to implement sophisticated pathfinding and AI solutions
- Extensibility allows future enhancements (3D, new game modes, new physics properties) without disrupting foundation
- Separation of universal topology (ADR-0005) from game mode configuration (e.g., ADR-0011) enables flexibility

**Future Considerations**:
- **Multiple Game Modes**: Different battle types can use different arena configurations (team battles, tournaments, special events)
- **3D Spatial System**: If aerial combat or vertical positioning becomes strategically important, extend to R³ with (x, y, z) coordinates
- **Variable Physics**: Different game modes may apply different physics properties (gravity, friction models, collision mechanics)
- **Obstacles and Terrain**: Spatial system designed to support static/dynamic obstacles and variable terrain effects
- **Polar Coordinate Option**: For specific mechanics (turret rotation, sensor sweeps) polar coordinates could supplement Cartesian system
- **Non-Euclidean Space**: Theoretical extensibility to hyperbolic or spherical geometry for future game modes

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
