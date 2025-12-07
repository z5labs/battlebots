---
title: "Battle Space"
description: "2D arena definition including coordinate system, boundaries, and spatial rules"
type: docs
weight: 1
date: 2025-12-05
---

## Overview

The Battle Space is a 2D rectangular arena where bots engage in combat. This bounded playing field provides a constrained environment that encourages tactical positioning and strategic decision-making. The battle space uses a continuous coordinate system (as opposed to a tile-based grid), allowing for smooth movement and precise positioning.

All bots exist within this shared arena, and their positions, movements, and interactions are governed by the spatial rules defined here. The battle space serves as the fundamental playing field for all battle scenarios.

## Coordinate System

The battle space uses a **Cartesian coordinate system** with the following properties:

- **Origin (0, 0)**: Located at the center of the arena
- **X-axis**: Horizontal axis, with positive values extending to the right and negative values to the left
- **Y-axis**: Vertical axis, with positive values extending upward and negative values downward
- **Units**: Abstract spatial units (not meters, pixels, or other real-world measurements)
- **Precision**: Floating-point coordinates allow for sub-unit positioning accuracy

This centered origin simplifies calculations for distance, angle, and relative positioning between bots. It also provides symmetry for balanced starting positions in various battle configurations.

## Boundaries

The battle space is defined by rectangular boundaries:

- **Arena Size**: **100 x 100 units** (TBD - subject to tuning based on playtesting)
- **X-axis Range**: -50 to +50 units
- **Y-axis Range**: -50 to +50 units

### Out-of-Bounds Handling

Bots cannot move outside the arena boundaries. The following rules apply:

1. **Movement Blocking**: Any movement command that would place a bot outside the boundaries is clamped to the nearest valid position at the boundary edge
2. **No Wrapping**: Coordinates do not wrap around (i.e., exiting the right side does not place a bot on the left side)
3. **Boundary Contact**: Bots may be positioned exactly on the boundary line
4. **Force Effects**: External forces (knockback, explosions, etc.) that would push a bot out-of-bounds will stop at the boundary

This approach prevents bots from leaving the arena while maintaining predictable physics and movement behavior.

## Bot Positioning

Each bot occupies a circular area within the battle space:

- **Bot Radius**: **2 units** (TBD - subject to balance tuning)
- **Center Position**: Each bot's coordinates (x, y) represent the center of its circular footprint
- **Minimum Separation**: Bots cannot overlap; their circular areas must not intersect

### Collision Detection

Collision detection uses simple circle-to-circle distance calculations:

1. **Bot-to-Bot Collision**: Two bots collide when the distance between their centers is less than the sum of their radii (2r = 4 units for identical bots)
2. **Bot-to-Wall Collision**: A bot collides with a wall when its center position plus its radius exceeds the boundary
   - Left wall: `x - radius < -50`
   - Right wall: `x + radius > 50`
   - Bottom wall: `y - radius < -50`
   - Top wall: `y + radius > 50`

3. **Collision Resolution**: When a collision is detected, the movement is adjusted to place the bot in contact with the obstacle without overlapping

## Friction and Movement Physics

The battle space applies friction to all moving bots, which affects their velocity and movement behavior. Friction provides realistic physics that require bots to continuously apply force to maintain movement, and it enables the possibility of variable terrain types with different surface properties.

### Friction Mechanics

1. **Friction Force**: Opposes the direction of bot movement, proportional to velocity
2. **Friction Coefficient (μ)**: Determines the strength of friction applied to a bot
   - **Default Coefficient**: **0.1** (TBD - subject to balance tuning)
   - **Range**: 0.0 (frictionless) to 1.0 (maximum friction)
3. **Velocity Decay**: Each update tick, a bot's velocity is reduced by the friction force
4. **Natural Stopping**: Without continuous thrust, a bot will gradually slow to a stop due to friction

### Friction Calculation

The friction force applied to a moving bot is calculated as:

```
friction_force = -μ × velocity
new_velocity = velocity + friction_force
```

Where:
- `μ` is the friction coefficient at the bot's current position
- `velocity` is the bot's current velocity vector
- The negative sign indicates friction opposes the direction of movement

### Variable Friction Zones

The battle space supports different friction coefficients across different areas, enabling terrain variety:

1. **Uniform Friction**: By default, the entire battle space has a uniform friction coefficient
2. **Friction Zones**: Specific rectangular or circular areas may define different friction values
   - **Low Friction** (μ < 0.1): "Slippery" surfaces where bots slide more easily
   - **Standard Friction** (μ = 0.1): Normal battle space surface
   - **High Friction** (μ > 0.1): "Rough" surfaces that slow bot movement more quickly

3. **Zone Priority**: When friction zones overlap, the highest friction coefficient applies
4. **Transition Behavior**: Moving between friction zones immediately applies the new coefficient (no gradual transition)

### Tactical Implications

Friction creates several tactical considerations:

- **Movement Planning**: Bots must account for deceleration when planning movements
- **Pursuit and Evasion**: Understanding friction helps predict opponent stopping distances
- **Zone Control**: High-friction zones can limit mobility, while low-friction zones enable faster repositioning
- **Energy Management**: Continuous thrust is required to maintain velocity, affecting energy/turn economy

This friction system provides a foundation for diverse terrain types and strategic positioning while maintaining simple, predictable physics.

## Line of Sight

Line of sight determines whether one bot can "see" or target another bot, which is essential for ranged attacks and targeting systems.

### Line of Sight Rules

1. **Direct Path**: A bot has line of sight to another bot if an unobstructed straight line can be drawn between their center positions
2. **Obstacle-Free**: Currently, the only obstacles are other bots. A bot blocks line of sight between two other bots if the line passes through its circular area
3. **Boundary Walls**: Walls do not block line of sight (bots can see through walls but cannot shoot through them - weapon-specific rules apply)

### Line of Sight Calculation

To determine if Bot A has line of sight to Bot B:

1. Draw a line segment from A's center to B's center
2. For each other bot C in the arena:
   - Calculate the perpendicular distance from C's center to the line segment
   - If this distance is less than C's radius, line of sight is blocked

This calculation may be optimized by only checking bots that fall within a bounding box around the line segment.

## Spatial Rules

The following spatial rules govern bot behavior within the battle space:

### Movement Constraints

1. **Speed Limits**: Bots have maximum movement speeds (TBD - defined in bot mechanics)
2. **Acceleration**: Bots may have acceleration/deceleration constraints (TBD - affects turn-by-turn movement)
3. **No Teleportation**: Bots cannot instantly move from one position to another; all movement follows continuous paths

### Wall Collisions

1. **Elastic Collisions**: When a bot collides with a wall during movement, it stops at the wall position
2. **No Damage**: Wall collisions do not inherently cause damage to bots (unless specific game mechanics introduce this)
3. **No Bouncing**: Bots do not bounce off walls; they simply stop at the boundary

### Bot-to-Bot Collisions

1. **Movement Blocking**: Bots cannot move through each other
2. **Collision Physics**: When two bots attempt to occupy overlapping space:
   - The moving bot stops at the point of contact
   - Both bots remain stationary in their final positions (no pushing or displacement)
3. **Damage**: Bot-to-bot collisions do not inherently cause damage (unless specific collision damage mechanics are introduced)

### Future Considerations

The following spatial features may be added in future iterations:

- **Obstacles**: Static or dynamic obstacles within the arena that block movement and line of sight
- **Advanced Terrain Effects**: Beyond friction, areas could have additional effects (e.g., damage-over-time zones, healing zones, vision-reducing fog)
- **Complex Friction Patterns**: Non-rectangular friction zones with gradual transitions between coefficients
- **Vertical Dimension**: Elevation or z-axis for flying bots or multi-level arenas
- **Dynamic Boundaries**: Shrinking play areas or moving walls to force engagement
- **Variable Arena Sizes**: Different battle modes with different arena dimensions

All numeric values and specific mechanics in this document are marked as TBD and will be refined through playtesting and balance tuning.
