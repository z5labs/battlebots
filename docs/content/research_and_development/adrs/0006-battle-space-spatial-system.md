---
title: "[0006] Battle Space Spatial System"
description: >
    Spatial system governing bot positioning, movement, and interactions in the battle arena
type: docs
weight: 6
category: "strategic"
status: "accepted"
date: 2025-12-07
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

Battle Bots requires a spatial environment where bots can move, position themselves tactically, and interact with each other. We need to define the fundamental coordinate system, boundaries, collision mechanics, and movement physics that govern all bot interactions. This spatial system must support real-time continuous gameplay while integrating cleanly with the gRPC bidirectional streaming protocol (ADR-0004).

Without a well-defined spatial system, we cannot:
- Implement bot movement mechanics
- Calculate distances for ranged attacks
- Detect collisions between bots or with boundaries
- Provide consistent positioning for visualization
- Define tactical positioning strategies
- Support future terrain and obstacle features

## Decision Drivers

* **Real-time Continuous Gameplay** - Must support smooth, continuous movement (not turn-based or tile-based)
* **Tactical Depth** - Precise positioning should enable strategic gameplay and tactical maneuvering
* **Protocol Integration** - Must integrate with gRPC streaming for position updates (ADR-0004)
* **Computational Efficiency** - Collision detection and position calculations must perform well in real-time
* **Predictability** - Movement physics should be deterministic and understandable to bot developers
* **Extensibility** - Should support future additions like obstacles, terrain effects, and variable arena sizes
* **Visualization Clarity** - Positions must map cleanly to visual representation for battle viewers

## Considered Options

* **Option 1: Grid-Based Tile System** - Discrete tiles, grid-based movement, simplified collision detection
* **Option 2: Continuous 2D Cartesian Coordinate System** - Smooth continuous movement, floating-point precision, circle-based collision
* **Option 3: Hexagonal Grid System** - Hex tiles for uniform distance, common in strategy games

## Decision Outcome

Chosen option: "**Option 2: Continuous 2D Cartesian Coordinate System**", because it enables smooth real-time movement that integrates naturally with gRPC streaming, provides precise tactical positioning depth, supports predictable physics calculations, and allows for future enhancements like variable terrain without fundamental redesign.

### Spatial System Specification

#### Coordinate System

The battle space uses a **Cartesian coordinate system** with the following properties:

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

1. **Speed Limits**: Bots have maximum movement speeds (defined in ADR-0007: Bot Characteristics)
2. **No Teleportation**: Bots cannot instantly move from one position to another; all movement follows continuous paths
3. **Elastic Wall Collisions**: When a bot collides with a wall during movement, it stops at the wall position without bouncing
4. **No Wall Damage**: Wall collisions do not inherently cause damage to bots

### Consequences

* Good, because continuous coordinate system enables smooth, realistic movement and precise tactical positioning
* Good, because centered origin (0,0) simplifies distance and angle calculations for bot AI developers
* Good, because floating-point precision allows sub-unit accuracy without limiting strategic options
* Good, because circle-based collision detection is computationally efficient and easy to understand
* Good, because friction system creates tactical depth through movement physics and enables future terrain variety
* Good, because extensible design supports future additions (obstacles, dynamic boundaries, variable terrain)
* Good, because integrates naturally with gRPC streaming protocol for real-time position updates
* Neutral, because arena size (100x100) and bot radius (2 units) are placeholders requiring playtesting
* Neutral, because friction coefficient (0.1) requires tuning to balance movement fluidity vs. stopping distances
* Neutral, because line of sight rules are simple but may need enhancement for fog of war features
* Bad, because floating-point coordinates introduce potential precision issues that grid-based systems avoid
* Bad, because continuous collision detection is more complex than tile-based collision
* Bad, because friction physics adds computational overhead compared to instant-stop movement

### Confirmation

The decision will be confirmed through:

1. Implementation of game server spatial system with continuous coordinates and circle collision detection
2. Bot SDK position and movement API that exposes Cartesian coordinates and velocity vectors
3. Visualization system successfully rendering bot positions and movements in the continuous space
4. Playtesting demonstrating that tactical positioning creates meaningful strategic choices
5. Performance testing confirming collision detection scales to required tick rates
6. Balance tuning of arena size, bot radius, and friction values through competitive gameplay

## Pros and Cons of the Options

### Option 1: Grid-Based Tile System

Discrete tiles where bots occupy grid cells, movement is tile-to-tile, collision is grid-based.

* Good, because simpler collision detection (just check tile occupancy)
* Good, because eliminates floating-point precision issues
* Good, because easier to visualize and debug (clear tile boundaries)
* Good, because path-finding algorithms are well-established for grids
* Neutral, because tactical depth depends on grid resolution (finer grids approach continuous)
* Bad, because discrete movement feels less fluid and realistic
* Bad, because limits precision of tactical positioning (can't be "between" tiles)
* Bad, because diagonal movement requires special handling for distance fairness
* Bad, because poor integration with real-time streaming (position changes are discrete jumps)

### Option 2: Continuous 2D Cartesian Coordinate System

Smooth continuous space with floating-point coordinates, circle-based collision (CHOSEN).

* Good, because smooth, realistic movement that feels natural and fluid
* Good, because precise tactical positioning enables fine-grained strategy
* Good, because integrates naturally with real-time gRPC streaming
* Good, because friction and physics systems are straightforward to implement
* Good, because extensible to obstacles, terrain effects, and variable arenas
* Good, because distance and angle calculations use standard Cartesian math
* Neutral, because requires careful tuning of movement speeds and collision parameters
* Bad, because floating-point precision can introduce edge cases
* Bad, because collision detection is more computationally expensive than grid-based
* Bad, because path-finding is more complex than grid-based approaches

### Option 3: Hexagonal Grid System

Hex tiles for movement and positioning, common in strategy games.

* Good, because uniform distance to all six neighbors (solves diagonal distance problem)
* Good, because provides tactical depth through grid-based positioning
* Good, because deterministic collision and movement (no floating-point issues)
* Good, because well-suited for turn-based tactical games
* Neutral, because offers middle ground between grid and continuous systems
* Bad, because more complex coordinate math (offset, cube, or axial coordinates)
* Bad, because discrete movement still feels less fluid than continuous
* Bad, because hex grids are less intuitive for many developers than Cartesian coordinates
* Bad, because poor fit for real-time continuous gameplay model
* Bad, because visualization and UI are more complex than rectangular grids or continuous space

## More Information

### Related Documentation

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol that transmits position updates and movement commands

- **[ADR-0007: Bot Characteristics System](0007-bot-characteristics-system.md)**: Defines Speed stat that determines movement rate in this spatial system

- **[Battle Space Analysis](../analysis/game-mechanics/battle-space/)**: Detailed technical specifications for the spatial system

- **[ADR-0005: 1v1 Battle Orchestration](0005-1v1-battle-orchestration.md)**: High-level battle flow that operates within this spatial system

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Performance testing to determine optimal tick rate and collision detection overhead
2. Playtesting to tune arena size for engagement frequency and tactical maneuvering space
3. Balance analysis of bot radius to ensure appropriate bot density and collision frequency
4. Friction coefficient tuning to balance movement fluidity with stopping control
5. Visualization testing to ensure positions render clearly and movement appears smooth

**Future Enhancements**: The spatial system is designed to support future additions without requiring fundamental redesign:
- **Obstacles**: Static or dynamic obstacles within the arena that block movement and line of sight
- **Advanced Terrain Effects**: Damage-over-time zones, healing zones, vision-reducing fog
- **Dynamic Boundaries**: Shrinking play areas or moving walls to force engagement
- **Variable Arena Sizes**: Different battle modes with different dimensions
- **Vertical Dimension**: Elevation or z-axis for flying bots or multi-level arenas

### Design Principles

The spatial system follows these principles:
- **Smooth over Discrete**: Continuous movement creates more engaging, fluid gameplay
- **Simple Physics**: Friction and collision mechanics are understandable and predictable
- **Tactical Positioning**: Precise coordinates enable meaningful strategic choices
- **Performance-Conscious**: Collision detection optimized for real-time gameplay
- **Extensible Foundation**: Core system supports future spatial features without redesign
