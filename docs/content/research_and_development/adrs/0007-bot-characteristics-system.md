---
title: "[0007] Bot Characteristics System"
description: >
    Attribute system defining bot capabilities and creating strategic differentiation
type: docs
weight: 7
category: "strategic"
status: "proposed"
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

## Considered Options

* **Option 1: Minimal Stats (2-3 attributes)** - Health + Speed only, possibly Attack
* **Option 2: Four-Stat System** - Health, Speed, Defense, Mass (equipment-derived)
* **Option 3: Complex Stats (6+ attributes)** - Health, Speed, Defense, Attack, Energy, Accuracy, Evasion, etc.

## Decision Outcome

Chosen option: "**Option 2: Four-Stat System (Health, Speed, Defense, Mass)**", because it creates strategic depth through meaningful stat interactions (Effective HP, Effective Speed), enables equipment-driven tradeoffs via Mass, balances complexity with accessibility, and provides diverse optimization paths without overwhelming developers. The inclusion of equipment-derived Mass creates natural mobility-power tradeoffs that emerge from loadout choices.

### Bot Characteristics Specification

#### Health

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
- Low Health bots must rely on Speed and tactical positioning to survive

#### Speed

Speed determines how quickly a bot can move through the 2D battle space (ADR-0006). This stat directly affects positioning, engagement control, and evasion capabilities.

**Key Properties**:
- **Movement Rate**: Distance covered per game tick (TBD: placeholder 1-10 units/tick)
- **Positioning Advantage**: Faster bots control engagement range
- **Pursuit/Evasion**: Enables chasing or retreating from combat
- **No Direct Damage Scaling**: Speed affects tactics, not raw damage output

**Gameplay Impact**:
- High Speed enables kiting strategies (attacking while maintaining distance)
- Critical for evasion-based playstyles that avoid damage
- Allows control of engagement timing (when to fight vs. when to disengage)
- Low Speed bots must rely on durability or zone control
- Speed differences create natural predator/prey dynamics between bot types

#### Defense

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
- Low Defense bots must rely on Speed for damage avoidance
- Defense vs. Health allocation creates build optimization choices

#### Mass

Mass represents the total physical weight of a bot, determined by the equipment and components it carries. Unlike other characteristics, **Mass is not directly allocated** but is the cumulative result of loadout choices.

**Key Properties**:
- **Equipment-Derived**: Mass is calculated from the sum of all equipped items
- **Dynamic Value**: Changes based on equipped weapons, armor, and modules
- **Movement Impact**: Higher Mass reduces effective Speed
- **Momentum Effects**: Mass affects collision physics and knockback resistance (TBD)
- **No Direct Damage Scaling**: Mass affects mobility, not offensive capability

**Gameplay Impact**:
- Heavy equipment loadouts reduce mobility through increased Mass
- Creates natural tradeoff between firepower/protection and Speed
- Light bots sacrifice durability for superior maneuverability
- Mass cannot be optimized independently - it's a consequence of equipment choices
- Forces strategic decisions between powerful equipment and tactical mobility
- May affect collision mechanics and position displacement (future mechanics)

**Equipment Examples** (TBD):
- Weapons: Heavy weapons (high Mass) vs. light weapons (low Mass)
- Armor: Heavy plating (high Mass) vs. light armor (low Mass)
- Modules: Power cores, sensors, and systems each contribute Mass
- Loadout variety creates diverse Mass profiles across bot builds

#### Stat Interactions

Bot characteristics don't operate in isolation - they create complex interactions that define combat dynamics:

**Effective Durability**: Health and Defense combine to determine true survivability:
- **Effective HP Formula**: Health × (1 + Defense modifier)
- High Defense multiplies the value of each Health point
- Balanced allocation is more efficient than single-stat stacking
- Example: 100 Health + 50% Defense = 150 Effective HP

**Mass and Mobility**: Mass directly impacts effective movement capability:
- **Effective Speed Formula**: Base Speed / Mass modifier
- Heavy equipment reduces Speed, creating mobility-firepower tradeoffs
- Light loadouts maximize Speed at the cost of offensive/defensive power
- Equipment choices fundamentally alter tactical capabilities

**Combat Positioning**: Speed enables tactical positioning and engagement control:
- High Speed allows kiting, pursuit, and disengagement
- Speed advantage determines range control in combat
- Mass penalties from equipment reduce positioning flexibility
- Lightweight builds gain tactical mobility at the cost of durability

**Survivability Tradeoffs**: Defensive investment creates complex build choices:
- High Health + Low Defense = vulnerable to sustained damage
- Low Health + High Defense = vulnerable to burst damage
- Mass from defensive equipment reduces evasion capability
- Optimal defensive strategy depends on threat profile

**Loadout Optimization**: Equipment choices (ADR-0008) create cascading effects across all characteristics:
- Heavy weapons increase offensive capability but reduce Speed through Mass
- Armor improves Defense but adds Mass that limits mobility
- Lightweight builds sacrifice protection for superior positioning
- No equipment configuration dominates all scenarios (intended design goal)

### Consequences

* Good, because four-stat system creates strategic depth without overwhelming developers
* Good, because stat interactions (Effective HP, Effective Speed) enable emergent complexity from simple rules
* Good, because equipment-derived Mass creates natural mobility-firepower tradeoffs
* Good, because multiple playstyles are viable (tank, DPS, mobile, balanced) through different stat profiles
* Good, because Defense × Health interaction rewards balanced allocation over single-stat stacking
* Good, because Speed enables tactical gameplay through positioning and engagement control
* Good, because Mass cannot be optimized independently, forcing meaningful equipment tradeoffs
* Good, because stats map cleanly to combat calculations and are easy to understand
* Neutral, because stat values (Health range, Speed range, Defense values) require extensive playtesting
* Neutral, because Mass modifier formulas need tuning to balance equipment weight penalties
* Neutral, because four stats hit a middle ground - simpler than complex systems but more nuanced than minimal
* Bad, because more stats to track and balance than minimal 2-3 stat systems
* Bad, because stat interactions (especially Effective HP) add complexity to build optimization
* Bad, because equipment-derived Mass means loadout choices have cascading effects that may confuse new developers

### Confirmation

The decision will be confirmed through:

1. Implementation of characteristic system in game server with stat calculation formulas
2. Equipment system implementation (ADR-0008) that modifies stats and contributes Mass
3. Playtesting with diverse bot builds (tank, DPS, balanced, mobile) to validate viability
4. Balance analysis ensuring no single stat profile dominates all scenarios
5. Developer feedback on stat system accessibility and understandability
6. Combat simulation confirming stat interactions create meaningful tactical choices

## Pros and Cons of the Options

### Option 1: Minimal Stats (2-3 attributes)

Health + Speed only, possibly Attack. Simplest possible system.

* Good, because extremely accessible and easy to understand
* Good, because minimal complexity for bot developers to manage
* Good, because fast iteration on balance with fewer variables
* Good, because straightforward calculations (no complex stat interactions)
* Neutral, because may be sufficient for initial gameplay proof-of-concept
* Bad, because limited strategic depth and build diversity
* Bad, because shallow optimization (just maximize both stats)
* Bad, because difficult to create distinct playstyles (tank vs. DPS vs. mobile)
* Bad, because equipment customization has minimal impact (just adds to 2-3 stats)
* Bad, because no emergent complexity from stat interactions

### Option 2: Four-Stat System (Health, Speed, Defense, Mass)

Health, Speed, Defense as direct stats, Mass as equipment-derived (CHOSEN).

* Good, because strategic depth through stat interactions (Effective HP, Effective Speed)
* Good, because equipment-derived Mass creates natural mobility-firepower tradeoffs
* Good, because enables distinct playstyles (tank, DPS, mobile, balanced)
* Good, because accessible complexity - four stats are learnable without being overwhelming
* Good, because Defense × Health interaction rewards balanced builds
* Good, because Speed enables tactical positioning gameplay
* Good, because Mass as consequence of equipment forces meaningful loadout choices
* Neutral, because requires balance tuning but not excessively complex
* Neutral, because four stats hit sweet spot between simple and overwhelming
* Bad, because more complex than minimal stats to implement and explain
* Bad, because stat interactions add cognitive load for optimization
* Bad, because equipment Mass penalties require careful tuning to avoid punishing heavy loadouts

### Option 3: Complex Stats (6+ attributes)

Health, Speed, Defense, Attack, Energy, Accuracy, Evasion, Stamina, etc.

* Good, because maximum strategic depth and optimization potential
* Good, because extensive differentiation between bot builds
* Good, because appeals to players who enjoy complex optimization
* Good, because many tuning knobs for balance adjustments
* Neutral, because may enable very deep competitive meta
* Bad, because high learning curve may deter casual bot developers
* Bad, because complexity makes game balance extremely difficult
* Bad, because numerous stats create decision paralysis during bot design
* Bad, because difficult to understand why one build beats another
* Bad, because implementation complexity significantly higher
* Bad, because many stats may be redundant or minimally impactful

## More Information

### Related Documentation

- **[ADR-0006: Battle Space Spatial System](0006-battle-space-spatial-system.md)**: Spatial environment where Speed determines movement

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
5. Speed balancing to ensure mobility advantages are significant but not overwhelming
6. Competitive meta analysis to identify dominant builds and adjust accordingly

**Key Design Insight**: Mass is equipment-derived, not directly allocated. This creates emergent tradeoffs where powerful equipment inherently reduces mobility, forcing strategic loadout decisions without requiring explicit stat allocation.

**Future Considerations**:
- Additional derived stats (e.g., Effective HP, Effective Speed) may be exposed to developers
- Attack stat may be added if weapon damage needs per-bot customization
- Energy stat may become a direct characteristic if resource management complexity increases
- Evasion or Accuracy stats may be added if hit-chance mechanics are introduced

### Design Principles

The characteristic system follows these principles:
- **Interactions over Isolation**: Stats combine to create emergent complexity
- **Tradeoffs over Power**: Equipment choices involve costs and benefits through Mass
- **Diversity over Dominance**: Multiple stat profiles should be competitively viable
- **Clarity over Complexity**: Four stats balance depth with accessibility
- **Equipment Integration**: Stats are modified by loadout choices (ADR-0008)
