---
title: "[0007] Bot Characteristics System"
description: >
    Attribute system defining bot capabilities and creating strategic differentiation
type: docs
weight: 7
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

Battle Bots requires an attribute system that defines bot capabilities and creates strategic differentiation between different bot builds. We need to determine how many stats, what they represent, how they interact, and how they integrate with equipment customization. The characteristic system must support diverse playstyles while remaining accessible and understandable to bot developers.

Without a well-defined characteristic system, we cannot:
- Define bot capabilities and limitations
- Create meaningful equipment modifications
- Balance different bot builds and strategies
- Design combat calculations and damage systems
- Enable strategic depth through stat optimization
- Provide clear feedback to bot developers about their bot's strengths

## Decision Drivers

* **Strategic Depth** - System should create meaningful build choices and optimization paths
* **Equipment Integration** - Stats must be modifiable by equipment to enable customization
* **Build Diversity** - Multiple viable stat allocations should exist (no single dominant build)
* **Developer Accessibility** - Stats should be understandable without being overwhelming
* **Calculation Simplicity** - Stat interactions should be calculable and predictable
* **Playstyle Support** - Should enable distinct archetypes (tank, DPS, mobile, balanced)
* **Spatial Integration** - Must work with continuous 2D movement system (ADR-0006)

## Decision Outcome

The bot characteristics system consists of three core attributes that create strategic depth through meaningful interactions while remaining accessible to bot developers:

1. **Health** - Bot's survivability pool; total damage a bot can sustain before destruction
2. **Defense** - Damage mitigation capability; reduces effective damage from enemy attacks
3. **Mass** - Equipment-derived weight; calculated from equipped items and impacts effective thrust-to-movement conversion

This three-attribute system creates strategic depth through stat interactions (Effective HP, thrust-based mobility), enables equipment-driven tradeoffs via Mass, balances complexity with accessibility, and provides diverse optimization paths without overwhelming developers. Movement is governed by a thrust-based system (ADR-0009) where bots apply continuous thrust force to overcome friction (ADR-0006), with Mass affecting how efficiently thrust translates into velocity. The inclusion of equipment-derived Mass creates natural mobility-power tradeoffs that emerge from loadout choices.

## Bot Characteristics Specification

### Health

Health (HP) represents a bot's survivability in combat. This is the primary resource that determines whether a bot remains operational in battle.

**Key Properties**:
- **HP Pool**: Total damage a bot can sustain before destruction (TBD: placeholder 100-500 range)
- **Damage Resistance**: Works in conjunction with Defense stat to reduce incoming damage
- **Destruction Threshold**: Bot is eliminated when Health reaches 0
- **No Regeneration**: Health does not regenerate during battle (current design)

**Gameplay Impact**:
- Higher Health allows bots to sustain longer engagements
- Critical for aggressive playstyles that prioritize direct confrontation
- Must be balanced against offensive capabilities to ensure threat viability
- Low Health bots must rely on mobility (via thrust actions) and tactical positioning to survive

### Defense

Defense represents a bot's ability to mitigate incoming damage. This stat reduces the effective damage from enemy attacks.

**Key Properties**:
- **Damage Reduction**: Percentage or flat reduction applied to incoming damage (TBD: placeholder 1-10)
- **All Damage Types**: Applies to all incoming damage sources (current design)
- **Multiplicative with Health**: Effective HP = Health × Defense multiplier
- **No Evasion Component**: Defense reduces damage taken, not hit chance

**Gameplay Impact**:
- High Defense enables tank strategies and prolonged engagements
- Multiplicatively increases effective Health pool
- Critical for front-line and damage-absorbing playstyles
- Low Defense bots must rely on mobility (via thrust actions) for damage avoidance
- Defense vs. Health allocation creates build optimization choices

### Mass

Mass represents the total physical weight of a bot, determined by the equipment and components it carries. Unlike other characteristics, **Mass is not directly allocated** but is the cumulative result of loadout choices.

**Key Properties**:
- **Equipment-Derived**: Mass is calculated from the sum of all equipped items
- **Dynamic Value**: Changes based on equipped weapons, armor, and modules
- **Movement Impact**: Higher Mass reduces acceleration from thrust actions (more force needed to overcome inertia and friction)
- **Momentum Effects**: Mass affects collision physics and knockback resistance (TBD)
- **No Direct Damage Scaling**: Mass affects mobility, not offensive capability

**Gameplay Impact**:
- Heavy equipment loadouts reduce mobility through increased Mass (requiring more thrust to achieve same velocity)
- Creates natural tradeoff between firepower/protection and maneuverability
- Light bots sacrifice durability for superior acceleration and agility
- Mass cannot be optimized independently - it's a consequence of equipment choices
- Forces strategic decisions between powerful equipment and tactical mobility
- Higher Mass requires sustained thrust to overcome friction and maintain velocity (ADR-0006)
- May affect collision mechanics and position displacement (future mechanics)

**Equipment Examples** (TBD):
- Weapons: Heavy weapons (high Mass) vs. light weapons (low Mass)
- Armor: Heavy plating (high Mass) vs. light armor (low Mass)
- Modules: Power cores, sensors, and systems each contribute Mass
- Loadout variety creates diverse Mass profiles across bot builds

### Stat Interactions

Bot characteristics don't operate in isolation - they create complex interactions that define combat dynamics:

**Effective Durability**: Health and Defense combine to determine true survivability:
- **Effective HP Formula**: Health × (1 + Defense modifier)
- High Defense multiplies the value of each Health point
- Balanced allocation is more efficient than single-stat stacking
- Example: 100 Health + 50% Defense = 150 Effective HP

**Mass and Mobility**: Mass directly impacts effective movement capability:
- **Thrust-to-Velocity Relationship**: Acceleration = Thrust Force / Mass (influenced by friction from ADR-0006)
- Heavy equipment increases Mass, reducing acceleration from thrust actions and requiring more sustained thrust to overcome friction
- Light loadouts maximize agility and responsiveness at the cost of offensive/defensive power
- Equipment choices fundamentally alter tactical capabilities through Mass-based mobility tradeoffs

**Combat Positioning**: Thrust actions (ADR-0009) enable tactical positioning and engagement control:
- High thrust capacity allows kiting, pursuit, and disengagement
- Low-Mass bots have positioning advantage through superior acceleration
- Mass penalties from heavy equipment reduce positioning flexibility and increase thrust requirements
- Lightweight builds gain tactical mobility at the cost of durability

**Survivability Tradeoffs**: Defensive investment creates complex build choices:
- High Health + Low Defense = vulnerable to sustained damage
- Low Health + High Defense = vulnerable to burst damage
- Mass from defensive equipment reduces mobility and increases thrust requirements for evasion
- Optimal defensive strategy depends on threat profile

**Loadout Optimization**: Equipment choices (ADR-0008) create cascading effects across all characteristics:
- Heavy weapons increase offensive capability but increase Mass, reducing acceleration and requiring more thrust
- Armor improves Defense but adds Mass that limits mobility and increases friction effects
- Lightweight builds sacrifice protection for superior acceleration and lower thrust requirements
- No equipment configuration dominates all scenarios (intended design goal)

## Consequences

* Good, because three-stat system creates strategic depth without overwhelming developers
* Good, because stat interactions (Effective HP, thrust-based mobility) enable emergent complexity from simple rules
* Good, because equipment-derived Mass creates natural mobility-firepower tradeoffs through thrust mechanics
* Good, because multiple playstyles are viable (tank, DPS, mobile, balanced) through different stat profiles
* Good, because Defense × Health interaction rewards balanced allocation over single-stat stacking
* Good, because thrust-based movement (ADR-0009) with Mass and friction (ADR-0006) enables tactical gameplay through physics
* Good, because Mass cannot be optimized independently, forcing meaningful equipment tradeoffs
* Good, because stats map cleanly to combat calculations and physics-based movement
* Good, because movement physics create natural skill expression through thrust management
* Neutral, because stat values (Health range, Defense values) require extensive playtesting
* Neutral, because Mass modifier formulas and thrust-to-acceleration ratios need tuning to balance equipment weight penalties
* Neutral, because three stats create a minimal but sufficient foundation for strategic depth
* Bad, because thrust-based movement adds complexity compared to simple speed-based systems
* Bad, because stat interactions (especially Effective HP and thrust-Mass-friction relationships) add complexity to build optimization
* Bad, because equipment-derived Mass means loadout choices have cascading effects on movement that may confuse new developers

## Confirmation

The decision will be confirmed through:

1. Implementation of characteristic system in game server with stat calculation formulas
2. Equipment system implementation (ADR-0008) that modifies stats and contributes Mass
3. Playtesting with diverse bot builds (tank, DPS, balanced, mobile) to validate viability
4. Balance analysis ensuring no single stat profile dominates all scenarios
5. Developer feedback on stat system accessibility and understandability
6. Combat simulation confirming stat interactions create meaningful tactical choices

## More Information

### Related Documentation

- **[ADR-0006: Battle Space Spatial System](0006-battle-space-spatial-system.md)**: Spatial environment with friction mechanics that govern movement physics

- **[ADR-0008: Equipment and Loadout System](0008-equipment-loadout-system.md)**: Equipment that modifies stats and contributes to Mass

- **[ADR-0009: Bot Actions and Resource Management](0009-bot-actions-resource-management.md)**: Actions that consume resources and leverage bot characteristics

- **[Bot Characteristics Analysis](../analysis/game-mechanics/characteristics/)**: Detailed technical specifications for the stat system

- **[ADR-0005: 1v1 Battle Orchestration](0005-1v1-battle-orchestration.md)**: High-level battle flow that uses these characteristics

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Combat simulation to model stat interactions and balance implications
2. Playtesting with diverse bot builds across different stat profiles
3. Equipment balance analysis to ensure Mass penalties are meaningful but not punishing
4. Health and Defense tuning to create appropriate effective HP ranges
5. Thrust-to-Mass ratio balancing to ensure mobility advantages are significant but not overwhelming
6. Friction coefficient tuning (ADR-0006) to balance movement physics
7. Competitive meta analysis to identify dominant builds and adjust accordingly

**Key Design Insight**: Mass is equipment-derived, not directly allocated. This creates emergent tradeoffs where powerful equipment inherently reduces mobility (through increased thrust requirements and friction effects), forcing strategic loadout decisions without requiring explicit stat allocation. Movement is governed by thrust actions (ADR-0009) that must overcome both Mass-based inertia and friction forces (ADR-0006).

**Future Considerations**:
- Additional derived stats (e.g., Effective HP, effective acceleration) may be exposed to developers
- Attack stat may be added if weapon damage needs per-bot customization
- Energy stat may become a direct characteristic if resource management complexity increases
- Evasion or Accuracy stats may be added if hit-chance mechanics are introduced
- Thrust capacity may become a direct characteristic if thrust-based movement requires per-bot customization beyond equipment

### Design Principles

The characteristic system follows these principles:
- **Interactions over Isolation**: Stats combine to create emergent complexity
- **Tradeoffs over Power**: Equipment choices involve costs and benefits through Mass
- **Diversity over Dominance**: Multiple stat profiles should be competitively viable
- **Clarity over Complexity**: Four stats balance depth with accessibility
- **Equipment Integration**: Stats are modified by loadout choices (ADR-0008)
