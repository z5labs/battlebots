---
title: "[0010] Biome Definition Framework"
description: >
    Meta-definition establishing which properties define biomes for battle arenas
type: docs
weight: 10
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

BattleBots requires diverse battle arenas with distinct tactical characteristics to provide meaningful gameplay variety. Currently, the game physics system (ADR-0006) includes variable friction with position-dependent coefficients, explicitly anticipating terrain-based mechanics. To operationalize this capability, we need a consistent framework for defining what makes one arena distinct from another.

A "Biome" is a meta-definition establishing which properties differentiate battle arenas. This ADR defines WHAT properties define biomes (not specific biomes themselves—those will follow in future ADRs). The framework must integrate seamlessly with existing spatial (ADR-0005) and physics (ADR-0006) systems.

ADR-0005 established the 2D Cartesian coordinate system with rectangular boundaries ([-50, 50] × [-50, 50]). ADR-0006 established variable friction with the formula `F_friction = μ(position) × M × |v|`, explicitly mentioning biomes as a future mechanic: "Variable friction zones create interesting terrain effects (ice zones, mud zones)" and "enables future 'biome' mechanics with different terrain types." Biomes operationalize this vision by defining how geography provides the μ(position) function.

We must choose which properties will define biomes from five candidates: Geography, Climate, Vegetation, Animal Life, and Ecosystem. Future ADRs will define specific concrete biomes (Desert, Arctic, Forest) using the framework established here.

## Decision Drivers

* **Physics System Integration** - Biome properties must integrate with the variable friction μ(position) function from ADR-0006
* **Implementation Complexity** - Properties must be implementable within reasonable development scope
* **Gameplay Impact** - Properties should create meaningful tactical differences between arenas
* **Visual Clarity** - Properties must be easily visualized and understandable to players
* **Determinism** - Must work with deterministic physics simulation requirements
* **Extensibility** - Framework must allow future property additions without invalidating existing biomes
* **Computational Efficiency** - Friction queries must be real-time calculable at game tick rates
* **Developer Predictability** - Bot developers must understand and predict terrain effects
* **Spatial Integration** - Must work seamlessly with 2D Euclidean space and arena boundaries
* **Arena Balance** - Properties must enable balanced competitive gameplay

## Considered Options

* Geography (topology and textures mapping to friction)
* Climate (weather, wind, temperature effects)
* Vegetation (obstacles, destructible plants, line-of-sight blocking)
* Animal Life (NPC entities, hazards, environmental challenges)
* Ecosystem (integrated climate, vegetation, and animal systems)

## Decision Outcome

Chosen option: **Geography**, because it perfectly integrates with the variable friction system already built into ADR-0006, requires minimal implementation complexity, creates immediate tactical gameplay impact, provides visual clarity, and alone provides sufficient arena differentiation.

### Geography Definition

Geography defines battle arena terrain through two complementary sub-properties:

**Topology** - The position-dependent friction structure:
* Mathematical definition: Position-dependent friction coefficient function `μ(position) : R² → R⁺`
* Implementation: Arena divided into friction zones with distinct coefficient values
* Zone types: Uniform (single friction across arena), Zoned (distinct regions with different coefficients), Gradient (smooth friction transitions)
* Data structure: Friction zone maps associating region coordinates with coefficient values

**Textures** - Visual representation communicating friction properties:
* Visual appearance tied to friction zones (ice, mud, sand, grass, rock, dirt, etc.)
* Texture-friction mapping: Low friction materials (ice → μ = 0.2), medium friction (grass → μ = 0.5), high friction (mud → μ = 1.2)
* Purely visual representation; friction coefficient is the source of truth
* Enables intuitive player understanding without memorizing numerical coefficients

**Rationale for Geography Selection:**

1. **Perfect Physics Integration** - Directly maps to μ(position) from ADR-0006; no new physics systems required
2. **Minimal Implementation** - Purely data-driven; friction maps and texture assets are sufficient
3. **Immediate Gameplay Impact** - Creates tactical positioning decisions, movement trade-offs, and strategic escape routes
4. **Visual Clarity** - Terrain is naturally intuitive; players understand friction effects from visual appearance
5. **Complete Differentiation** - Geography alone creates sufficient arena variety for meaningful gameplay

**Example Biome Definitions (from future ADRs):**

* **Desert**: Sand regions (μ = 0.6), rocky outcrops (μ = 0.4), gravel paths (μ = 1.0)
* **Arctic**: Ice plains (μ = 0.2), snow banks (μ = 0.5), frozen rock (μ = 0.4)
* **Forest**: Grass clearings (μ = 0.5), muddy swamps (μ = 1.2), dirt paths (μ = 0.6)

### Rejected Properties

All rejected properties share a core reason: **implementation complexity exceeds current project scope**. They are deferred to future ADRs after geography-based biomes are validated through gameplay.

**Climate** - Rejected (deferred):
* Requirements: Dynamic weather simulation, wind force vectors, temperature state management, temporal tracking, dynamic friction modifiers
* Complexity: New physics systems (wind effects on movement), bot state expansion (temperature), significant development effort
* Status: Deferred until geography-based biomes validate biome concept
* Future ADR: "Add Climate Property to Biome Framework"

**Vegetation** - Rejected (deferred):
* Requirements: Obstacle collision system, pathfinding around obstacles, destructible entity system, line-of-sight mechanics
* Complexity: Collision system expansion, visual rendering complexity, AI pathfinding integration
* Status: Deferred; static obstacles must prove tactically valuable through playtesting
* Future ADR: "Add Vegetation Obstacles to Biome Framework"

**Animal Life** - Rejected (deferred indefinitely):
* Requirements: NPC entity system with AI, animal movement, collision detection, damage system expansion, visual rendering
* Complexity: Extremely high implementation cost; minimal tactical impact on bot-versus-bot combat
* Status: Deferred indefinitely unless specific game mode (e.g., "Survival Mode") explicitly requires environmental hazards
* Unlikely path: Niche use case with separate development prioritization

**Ecosystem** - Rejected (deferred indefinitely):
* Requirements: Climate + Vegetation + Animal Life as dependencies, dynamic interaction rules between systems
* Complexity: Exponential implementation cost; no clear mapping to arena combat mechanics
* Status: Deferred indefinitely; more appropriate for environmental simulation than arena combat games
* Design fit: Better suited for open-world or survival games than competitive PvP

**Design Principle**: "Start with the simplest property that creates meaningful biome differentiation, validate through gameplay, then consider expansion."

### Consequences

**Overall Geography Decision:**
* Good, because it perfectly aligns with the variable friction system ADR-0006 already implements
* Good, because data-driven approach minimizes implementation complexity
* Good, because friction-based mechanics create clear tactical gameplay (positioning, escape routes, speed zones)
* Good, because visual textures provide intuitive player understanding
* Neutral, because friction coefficients require careful tuning per biome for balance
* Bad, because initial biomes are limited to terrain-based tactics (defers advanced environmental mechanics)

**Rejected Properties:**
* Good, because reduced scope enables high-quality Geography implementation
* Good, because clear extensibility path allows future property additions without breaking existing biomes
* Bad, because defers dynamic environmental mechanics (weather, obstacles, hazards)
* Bad, because limits initial biome variety to terrain-only differentiation
* Bad, because some properties (climate, vegetation) offer rich gameplay potential (deferred, not eliminated)

## Confirmation

Validation approach for ADR-0010 and subsequent biome implementations:

1. **Geography Framework Implementation** - Friction coefficient maps query correctly based on bot position
2. **Integration Testing** - ADR-0005 coordinates and ADR-0006 friction formula integrate seamlessly
3. **Concrete Biome ADRs** - Future biome definitions successfully instantiate framework (Desert, Arctic, Forest)
4. **Playtesting** - Confirms tactical differentiation between biomes; friction effects create meaningful positioning decisions
5. **Performance Validation** - Real-time friction queries meet tick rate requirements (verify O(1) or O(log n) lookup)
6. **Developer Experience** - Bot developers can predict and understand terrain effects on movement

## Pros and Cons of the Options

### Geography (Chosen)

Primary mechanism: Position-dependent friction coefficients providing terrain differentiation.

* Good, because direct integration with ADR-0006 variable friction formula `F_friction = μ(position) × M × |v|`
* Good, because minimal complexity—purely data-driven friction maps and texture assets
* Good, because creates immediate tactical gameplay: positioning matters, escape routes vary, movement trade-offs
* Good, because visual clarity enables intuitive player understanding
* Good, because sufficient arena differentiation (Desert vs. Arctic vs. Forest feel distinctly different)
* Good, because completely non-breaking for future property additions
* Neutral, because friction coefficient tuning requires careful balance per biome
* Bad, because arena topology mapping requires design effort per biome
* Bad, because limited to terrain-based tactics initially (defers weather, obstacles, hazards)

### Climate

Primary mechanism: Dynamic weather affecting movement, visibility, and damage.

* Good, because weather variety creates dynamic, changing tactical environments
* Good, because intuitive to players (ice is slippery, wind affects trajectory)
* Good, because rich gameplay potential (weather systems, temporal effects)
* Bad, because requires temporal state system (time progression, weather changes over match duration)
* Bad, because requires new physics: wind force vectors affecting movement trajectory
* Bad, because requires bot state expansion: temperature tracking, heat/cold effects
* Bad, because dynamic friction modifiers complicate ADR-0006 physics (time-dependent μ)
* Bad, because significant development effort for uncertain gameplay value
* **Rejected**: Implementation complexity too high for initial release

### Vegetation

Primary mechanism: Static and destructible obstacles providing cover and tactical positioning.

* Good, because tactical obstacles create cover-based positioning strategies
* Good, because line-of-sight blocking enables stealth and ambush tactics
* Good, because destructible obstacles add dynamic environmental interaction
* Bad, because requires obstacle collision detection system (new physics component)
* Bad, because requires pathfinding enhancements for bot AI
* Bad, because requires destructible entity system and damage integration
* Bad, because visual complexity in rendering (obstacle placement, destruction states)
* Bad, because significant development effort before gameplay validation
* **Rejected**: Implementation complexity too high for initial release

### Animal Life

Primary mechanism: NPC entities with AI providing hazards and environmental challenges.

* Good, because environmental hazards create immersive atmosphere
* Good, because NPC behavior adds unpredictability and engagement
* Neutral, because tactical impact on bot-versus-bot combat is marginal
* Bad, because requires NPC entity system with independent AI
* Bad, because requires animal movement, collision, and rendering
* Bad, because requires damage system expansion (animals attack bots)
* Bad, because extremely high implementation cost relative to gameplay impact
* Bad, because uncertain gameplay value in competitive PvP setting
* **Rejected**: Marginal tactical impact does not justify complexity cost

### Ecosystem

Primary mechanism: Integrated systems where climate, vegetation, and animals interact dynamically.

* Good, because emergent environmental interactions create rich simulations
* Good, because integrated systems feel cohesive (weather affects vegetation, animals react to weather)
* Bad, because requires Climate + Vegetation + Animal Life as implementation prerequisites
* Bad, because exponential complexity with interdependent systems
* Bad, because no clear mapping to arena combat mechanics (design is for simulation, not competition)
* Bad, because prohibitive implementation cost
* **Rejected**: More suitable for open-world simulation games than competitive PvP

## More Information

### Related Documentation

* **ADR-0005: BattleBot Universe Topological Properties** - Defines 2D Cartesian coordinate system with [-50, 50]² boundaries
* **ADR-0006: BattleBot Universe Physics Laws** - Defines variable friction with μ(position) function and explicitly anticipates biomes
* **ADR-0007: Bot Movement Mechanics** - Demonstrates gameplay impact of friction through terminal velocity formula

### Future Concrete Biome ADRs

* **ADR-00XX: Desert Biome** - First concrete biome implementation with friction topology and texture definition
* **ADR-00YY: Arctic Biome** - Second biome demonstrating framework reusability
* **ADR-00ZZ: Forest Biome** - Third biome validating biome differentiation

### Future Deferred Properties

**Climate Property (Deferred)**
* Implementation prerequisites: Temporal state system, visual weather effects, dynamic friction modifiers, bot temperature state
* Validation gate: Geography-only biomes must prove successful through playtesting
* Potential ADR: "Add Climate Property to Biome Framework" (depends on gameplay feedback)

**Vegetation Property (Deferred)**
* Implementation prerequisites: Obstacle collision system, pathfinding enhancement, line-of-sight system, destructible entities
* Validation gate: Static obstacles must prove tactically valuable; requires prior work on collision systems
* Potential ADR: "Add Vegetation Obstacles to Biome Framework" (depends on terrain tactics validation)

**Animal Life Property (Deferred Indefinitely)**
* Implementation prerequisites: NPC entity system, AI framework, entity collision/damage expansion
* Unlikely activation: Only if specific game mode (e.g., "Survival Mode") explicitly requires environmental hazards
* Non-standard track: Would require separate design, playtesting, and development prioritization

**Ecosystem Property (Deferred Indefinitely)**
* Implementation prerequisites: All three deferred properties implemented first, complex interaction rules
* Design fit: More appropriate for open-world or survival simulation games than competitive PvP
* Architectural note: Would require fundamental redesign of arena mechanics away from pure bot-versus-bot focus

### Extensibility Design

Non-breaking addition of properties through future ADRs:
* Existing biomes remain valid with Geography-only definitions
* New biomes created after property additions can leverage additional properties
* Property additions require new ADRs and validation before adoption
* Each biome ADR explicitly documents which properties it uses

### Design Principles

1. **Simplicity First** - Start with Geography, validate through gameplay, then expand properties
2. **Leverage Existing Systems** - Use ADR-0006 variable friction; don't invent new physics
3. **Minimize Initial Complexity** - Defer properties until proven necessary by gameplay
4. **Validate Before Expanding** - Demonstrate need through competitive playtesting before adding properties
