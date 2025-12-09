---
title: "[0009] Equipment Type System"
description: >
    Define equipment type categories that enable bot customization through stat modification and tactical capabilities
type: docs
weight: 9
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

The BattleBot equipment system requires a clearly defined set of **equipment type categories** that bots can equip to customize their capabilities. ADR-0008 (Bot Characteristics System) established that equipment modifies bot stats—particularly Mass, Health, and Defense—and creates the foundation for equipment-based strategic differentiation. However, we have not yet decided which **types** of equipment will be supported in the initial implementation.

Equipment types define broad categories of items with similar mechanical purposes. Each type represents a distinct customization axis that creates different strategic choices and playstyles. The choice of which equipment types to support impacts:

- **Build diversity**: How many distinct bot archetypes are possible
- **Tactical complexity**: The depth of strategic decision-making in loadout selection
- **Implementation scope**: Development effort and balancing complexity
- **Future extensibility**: Ability to expand equipment options based on gameplay validation

Without defined equipment types, we cannot:
- Design specific equipment items (weapons, armor pieces, utility modules)
- Create equipment loadout constraints and balancing rules
- Implement equipment-based stat modifications in the game server
- Enable strategic build customization for bot developers
- Support future game modes that depend on specialized equipment (team battles with communication gear, fog of war with sensor equipment)

This ADR establishes which equipment type categories are supported for the **initial implementation**, with explicit acknowledgment that future ADRs may expand the equipment type system as new game mechanics and modes are developed and validated through playtesting.

## Decision Drivers

* **Strategic Depth** - Equipment types should create meaningful build choices and optimization paths without overwhelming bot developers
* **Implementation Feasibility** - Initial release should focus on core equipment types that deliver maximum gameplay value relative to development effort
* **Stat System Integration** - Equipment must integrate cleanly with Health, Defense, and Mass characteristics from ADR-0008
* **Physics Integration** - Equipment Mass must correctly affect movement mechanics (acceleration A = F/M from ADR-0007) and friction (F = μ × M × |v| from ADR-0006)
* **Incremental Complexity** - Start simple with fundamental types, expand based on gameplay validation and balance testing
* **Future Game Modes** - Consider which equipment types are essential now versus valuable for future features (team battles, fog of war, competitive modes)
* **Balance Testability** - Limit initial types to enable thorough balance testing and validation before expansion
* **Developer Experience** - Equipment categories should be understandable and implementable by bot developers with clear strategic roles

## Considered Options

* **Option 1: Weapons** - Equipment that enables bots to attack and damage other bots
* **Option 2: Armor** - Equipment that increases Defense stat (damage mitigation) at the cost of increased Mass
* **Option 3: Load Bearing Gear** - Equipment that increases equipment carrying capacity at significant Mass penalty
* **Option 4: Communication** - Equipment for bot-to-bot communication (useful in team battles)
* **Option 5: Sensor** - Equipment that detects nearby enemies (useful in fog of war battles)
* **Option 6: Boost** - Equipment that temporarily enhances bot characteristics

## Decision Outcome

Chosen options: **Weapons and Armor only** for the initial implementation.

The equipment system will initially support two equipment type categories:

1. **Weapons** - Enable offensive capabilities through damage output modification
2. **Armor** - Enable defensive capabilities through Defense stat modification

Both types modify bot Mass (affecting mobility per ADR-0007 and friction per ADR-0006), creating the fundamental **mobility-firepower-protection tradeoff** that defines strategic build diversity. This two-type system creates clear strategic depth through force balancing while remaining implementable within project scope.

**Rationale:**

1. **Core Combat Loop** - Weapons (offense) and Armor (defense) define the fundamental combat interaction between bots; together they form the essential combat mechanics
2. **Complete Stat Coverage** - Combined with base characteristics (Health, Defense, Mass from ADR-0008), these types create complete strategic build diversity enabling distinct archetypes (tank, DPS, mobile, balanced)
3. **Physics Integration Validation** - Both types modify Mass, enabling thorough validation of movement mechanics (ADR-0007) where acceleration depends on Mass (A = F/M) and friction depends on Mass (F = μ × M × |v|)
4. **Testable Combat Triangle** - Creates offensive builds, defensive builds, and balanced builds for comprehensive balance testing without single dominant strategy
5. **No System Dependencies** - Function independently without requiring additional game mechanics, protocols, or systems
6. **Maximum Value per Complexity** - Deliver core PvP gameplay experience with minimal implementation scope, enabling focus on balance and polish

## Pros and Cons of the Options

### Option 1: Weapons

Equipment that enables bots to attack and damage other bots through various weapon types. Weapons modify offensive capability (damage output) and bot Mass.

**Examples**: Rifles, shotguns, laser weapons, missile launchers, plasma cannons, melee weapons

* Good, because enables core offensive gameplay and damage dealing mechanics
* Good, because creates weapon variety opportunities (high damage vs high rate of fire, close range vs long range)
* Good, because Mass penalties create meaningful mobility tradeoffs (heavy weapons reduce acceleration and terminal velocity)
* Good, because integrates with future damage calculation system and combat flow
* Good, because familiar concept from existing games (accessible to bot developers)
* Good, because enables DPS and burst damage playstyles through different weapon types
* Good, because supports tactical positioning based on weapon effective ranges
* Neutral, because requires future ADRs to define specific weapons and damage mechanics
* Neutral, because requires weapon balance tuning through playtesting

### Option 2: Armor

Equipment that increases Defense stat (damage mitigation) at the cost of increased bot Mass. Armor improves survivability through damage reduction.

**Examples**: Light plating, medium armor, heavy plating, reactive armor, ablative armor, shields

* Good, because enables defensive gameplay and tank strategies
* Good, because creates Defense stat variation beyond base values (builds on ADR-0008)
* Good, because Mass penalties create protection vs mobility tradeoff
* Good, because interacts multiplicatively with Health for Effective HP (Effective HP = Health × (1 + Defense modifier) from ADR-0008)
* Good, because enables distinct tank archetype viable against offensive builds
* Good, because supports diverse survivability strategies (high Defense, medium Defense, low Defense approaches)
* Neutral, because requires future ADRs to define specific armor pieces and Defense values
* Neutral, because requires armor balance tuning and testing with weapons

### Option 3: Load Bearing Gear

Equipment that increases equipment carrying capacity (equipment slots) at significant Mass penalty. Enables bots to equip more items.

**Examples**: Backpacks, cargo frames, modular mounts, storage systems

* Good, because enables "heavy weapons platform" builds with multiple weapons or extra capacity
* Good, because creates emergent strategies for specialized loadouts
* Good, because expands equipment slot system for future equipment types
* Neutral, because requires equipment slot system design (not yet established)
* Neutral, because provides build customization options for advanced players
* Bad, because adds complexity without immediate gameplay value (weapons and armor needed first)
* Bad, because requires loadout constraint system beyond initial scope
* Bad, because creates management complexity for bot developers (slot optimization)
* Bad, because should be deferred until weapons and armor balance is established

### Option 4: Communication

Equipment for bot-to-bot communication within team battles. Enables bots to share information and coordinate actions.

**Examples**: Radio transmitters, encrypted communication systems, signal boosters, relay networks

* Good, because enables team coordination in future team battle modes
* Good, because adds strategic depth to multiplayer scenarios
* Good, because creates support role opportunities (communication specialist)
* Good, because can enable unique team strategies and synergies
* Neutral, because minimal Mass cost (communication gear is typically light)
* Bad, because only valuable in team battles (not in initial 1v1 PvP focus)
* Bad, because requires bot-to-bot communication protocol implementation (complex)
* Bad, because should be deferred until team battle modes are developed
* Bad, because adds scope without supporting initial single-player game modes
* Bad, because requires separate ADR for team battle architecture and team communication protocol

### Option 5: Sensor

Equipment that detects nearby enemies, useful in fog of war battles where bots have limited visibility. Increases detection range and may provide tactical information.

**Examples**: Radar systems, lidar, motion sensors, thermal scanners, sonar, detection modules

* Good, because enables information gathering in fog of war scenarios
* Good, because creates scout and reconnaissance roles
* Good, because adds strategic depth to asymmetric information scenarios
* Good, because supports recon builds with high sensor range
* Neutral, because minimal Mass cost (sensors are typically light)
* Bad, because only valuable with fog of war game mechanics (not in open-map initial design)
* Bad, because requires visibility system implementation (line of sight, detection radius, fog of war rendering)
* Bad, because should be deferred until fog of war game mode is developed
* Bad, because adds scope without supporting initial visibility model
* Bad, because requires separate ADR for visibility and fog of war architecture

### Option 6: Boost

Equipment that temporarily enhances bot characteristics for tactical advantage. Provides time-limited stat bonuses or special effects.

**Examples**: Repair modules, shield generators, thrust boosters, damage amplifiers, defense boosters, speed enhancers

* Good, because enables tactical burst capabilities and reactive gameplay
* Good, because adds resource management complexity (cooldowns, charges, activation timing)
* Good, because creates high-skill plays through optimal activation timing
* Good, because supports diverse playstyles (burst damage, defensive boost, speed boost)
* Neutral, because requires temporal effects and state management system
* Neutral, because adds complexity that may appeal to advanced players
* Bad, because adds significant implementation complexity (activation mechanics, cooldown system, temporal state)
* Bad, because requires resource management system (energy/charges/cooldowns)
* Bad, because should be deferred until core combat loop is validated and balanced
* Bad, because adds complexity that may overwhelm bot developers before core mechanics are learned
* Bad, because may require separate ADR for active abilities and temporal effects system
* Bad, because balance interactions between boost effects and base characteristics need careful analysis

## Consequences

### Overall Equipment System

* Good, because two-type system creates clear strategic depth through offense-defense-mobility triangle
* Good, because Weapons and Armor integrate perfectly with existing characteristics (Health, Defense, Mass from ADR-0008)
* Good, because Mass-based tradeoffs create natural balance between firepower, protection, and movement
* Good, because focused scope enables thorough balance testing and polish before expansion
* Good, because creates solid foundation for future equipment type additions without overcommitting implementation
* Good, because bot developers have understandable, clear equipment categories with distinct strategic roles
* Good, because supports diverse playstyles (tank, DPS, mobile, balanced) through equipment choices
* Neutral, because limits initial build diversity to combat-focused archetypes (no support/utility roles initially)
* Bad, because team battle features (Communication equipment) and fog of war (Sensor equipment) must wait for future implementation
* Bad, because specialized playstyles requiring Boost or Load Bearing Gear must wait for system expansion

### Weapons Equipment Type

* Good, because enables core offensive gameplay and weapon variety
* Good, because creates meaningful weapon differentiation (damage output, rate of fire, range)
* Good, because Mass penalties create observable mobility tradeoffs
* Good, because supports tactical positioning and range management gameplay
* Good, because integrates with future damage calculation and weapon mechanics systems
* Good, because provides clear offensive archetype for bot builders
* Neutral, because requires future ADRs to define specific weapon types and mechanics
* Neutral, because requires extensive balance tuning through playtesting

### Armor Equipment Type

* Good, because enables tank strategies and defensive gameplay
* Good, because creates Defense stat variation beyond base values
* Good, because Mass penalties create meaningful protection vs mobility tradeoff
* Good, because interacts multiplicatively with Health for Effective HP optimization
* Good, because enables distinct defensive archetype viable against offensive builds
* Good, because supports diverse survivability strategies
* Good, because provides clear defensive archetype for bot builders
* Neutral, because requires future ADRs to define specific armor and Defense values
* Neutral, because requires balance tuning with weapons

### Deferred Equipment Types

* Good, because focusing on core types (Weapons, Armor) ensures quality implementation and thorough balance testing
* Good, because incremental expansion allows learning from initial gameplay before adding complexity
* Good, because Load Bearing Gear, Communication, Sensor, and Boost remain available for future ADRs
* Good, because deferred types can be added based on community feedback and gameplay insights
* Neutral, because team battles and fog of war modes will require equipment expansion before competitive launch
* Neutral, because future types can be designed with knowledge from weapon/armor balance insights
* Bad, because limits initial game mode variety (1v1 combat only until expansion)
* Bad, because team-oriented players must wait for team battle features and Communication equipment

## Confirmation

The decision will be confirmed through:

1. **Equipment Implementation Validation**:
   - Weapon equipment items successfully implement damage output mechanics
   - Weapon equipment items correctly contribute to bot Mass
   - Armor equipment items successfully modify Defense stat
   - Armor equipment items correctly contribute to bot Mass
   - Equipment loadout system integrates with bot SDK

2. **Mass Integration Testing**:
   - Heavy weapons increase bot Mass and reduce acceleration (A = F/M verification from ADR-0007)
   - Heavy weapons reduce terminal velocity (v = F/(μ × M) verification from ADR-0006)
   - Heavy armor increases bot Mass and reduces acceleration
   - Heavy armor reduces terminal velocity
   - Light equipment enables high-mobility builds with superior agility

3. **Build Diversity Validation**:
   - Offensive builds (heavy weapons, light armor) are strategically viable
   - Defensive builds (light weapons, heavy armor) are strategically viable
   - Balanced builds (medium weapons, medium armor) are strategically viable
   - No single build dominates all scenarios or creates non-interactive gameplay

4. **Physics Integration Confirmation**:
   - Acceleration formula (A = F/M) correctly applies with equipment-derived Mass
   - Friction formula (F = μ × M × |v|) correctly applies with equipment-derived Mass
   - Terminal velocity varies appropriately with equipment Mass
   - Equipment Mass changes are observable in movement behavior

5. **Balance Validation**:
   - Weapon damage scales appropriately with Mass penalty (heavy weapons balanced)
   - Armor Defense bonus scales appropriately with Mass penalty (heavy armor balanced)
   - Mobility penalties from Mass are significant and observable but not punishing
   - Combat encounters have varied outcomes based on loadout choices
   - No equipment combinations create trivial or unwinnable matchups

6. **Developer Experience Testing**:
   - Bot developers understand weapon and armor choices and tradeoffs
   - Equipment selection is intuitive in bot SDK
   - Build optimization is accessible without requiring extensive analysis
   - Documentation clearly explains Mass impact on movement

7. **Playtesting Feedback**:
   - Weapon/armor choices feel meaningful and impactful
   - Equipment tradeoffs between firepower, protection, and mobility are clear
   - Equipment variety creates distinct playstyles
   - No dominant meta build emerges that eliminates alternatives

8. **Future Extensibility Validation**:
   - Equipment type system can accommodate future additions (Load Bearing Gear, Communication, Sensor, Boost)
   - Adding new types does not break existing weapon/armor balance
   - Type system is extensible through ADR without architectural changes required
   - Future equipment items can modify Mass and other stats consistently

## More Information

### Related Documentation

- **[ADR-0005: BattleBot Universe Topological Properties](0005-battlebot-universe-topological-properties.md)**: Defines 2D Euclidean continuous space and rectangular boundaries where equipped bots battle

- **[ADR-0006: BattleBot Universe Physics Laws](0006-battlebot-universe-physics-laws.md)**: Defines Mass as universal property affecting friction (F = μ × M × |v|) and collision physics; equipment contributes to Mass

- **[ADR-0007: Bot Movement Mechanics](0007-bot-movement-mechanics.md)**: Defines thrust-based movement where acceleration depends on Mass (A = F/M); equipment-derived Mass directly affects acceleration and terminal velocity

- **[ADR-0008: Bot Characteristics System](0008-bot-characteristics-system.md)**: Defines Health, Defense, and Mass characteristics; establishes equipment-derived Mass formula (M_total = M_base + Σ(M_equipment)) and stat interactions (Effective HP = Health × (1 + Defense modifier))

### Future Architecture Decision Records

This ADR establishes the equipment type framework for initial implementation. Future ADRs will define:

- **ADR-0010 (or subsequent)**: Specific Weapon Items - Define individual weapons (rifle, shotgun, laser, etc.) with concrete damage values, projectile mechanics, and Mass contributions
- **ADR-0011 (or subsequent)**: Specific Armor Items - Define individual armor pieces (light plating, medium armor, heavy plating, etc.) with concrete Defense bonuses and Mass values
- **Future ADR**: Load Bearing Gear Equipment Type - Define when equipment slot expansion is needed for advanced multi-item builds
- **Future ADR**: Communication Equipment Type - Define bot-to-bot communication protocols and mechanics for team battle modes
- **Future ADR**: Sensor Equipment Type - Define fog of war visibility and detection mechanics for visibility-limited game modes
- **Future ADR**: Boost Equipment Type - Define active abilities, temporal effects, cooldown system, and burst capability mechanics

### Implementation Notes

#### Equipment Type Properties

Each equipment type has these defining characteristics:

- **Stat Modifications**: Which bot characteristics (Health, Defense, Mass) are affected
- **Mass Contribution**: All equipment contributes to bot Mass (varies by type and item)
- **Strategic Purpose**: The gameplay role and build archetype enabled
- **Dependencies**: Which game systems or modes require this type

#### Weapons Type Specification

- **Primary Function**: Enable offensive capabilities through damage output
- **Stat Modifications**: Damage output (weapon-specific property, not bot characteristic); Mass contribution
- **Mass Range** (TBD): Light weapons 5-10 Mass units, Heavy weapons 20-30 Mass units
- **Damage-Mass Tradeoff**: Higher damage weapons have higher Mass, reducing acceleration and terminal velocity
- **Integration Points**: Future damage calculation system, combat resolution system
- **Examples** (to be defined in future ADRs): Rifle, Shotgun, Laser, Missile Launcher, Plasma Cannon, Melee Weapon

#### Armor Type Specification

- **Primary Function**: Enable defensive capabilities through Defense stat modification
- **Stat Modifications**: Defense bonus (damage reduction percentage or flat reduction); Mass contribution
- **Mass Range** (TBD): Light armor 3-8 Mass units, Heavy armor 15-25 Mass units
- **Defense-Mass Tradeoff**: Higher Defense armor has higher Mass, reducing acceleration and terminal velocity
- **Integration Points**: Damage calculation system (applies Defense multiplier), stat interaction system
- **Examples** (to be defined in future ADRs): Light Plating, Medium Armor, Heavy Plating, Reactive Armor, Ablative Armor

#### Mass Integration with Equipment

From ADR-0008: `M_total = M_base + Σ(M_equipment)`

**Formula expansion with equipment types**:
```
M_total = M_base + M_weapon + M_armor
```

**Movement Impact**:
- Acceleration: `A = F_thrust / M_total`
- Terminal Velocity: `v_terminal = F_thrust / (μ(position) × M_total)`
- Friction Effect: `F_friction = μ(position) × M_total × |v|`

**Implications**:
- Weapons contribute Mass, reducing acceleration and terminal velocity
- Armor contributes Mass, reducing acceleration and terminal velocity
- Heavy loadouts (heavy weapon + heavy armor) dramatically reduce mobility
- Light loadouts (light weapon + light armor) enable superior acceleration and speed
- Creates natural mobility-firepower-protection tradeoff emerging from equipment choices

#### Equipment Slot System (Future)

Current implementation assumes all bots can equip:
- One weapon (future Load Bearing Gear may enable multiple weapons)
- One armor piece (future may enable mixed armor or multiple pieces)

This simplifies initial implementation while leaving room for future equipment slot expansion through Load Bearing Gear equipment type.

#### Integration with Bot Characteristics System

From ADR-0008, equipment interacts with characteristics as follows:

**Health Modification**:
- Some equipment types may increase Health pool (future Boost equipment)
- Initial Weapons and Armor do not modify Health
- Equipment can be designed to increase survivability through Defense instead of Health

**Defense Modification**:
- Armor equipment directly increases Defense stat
- Enables Effective HP optimization: `Effective HP = Health × (1 + Defense modifier)`
- Supports tank playstyle through Defense stacking

**Mass Modification**:
- All equipment contributes to Mass
- Creates universal mobility penalty for equipment load
- Enables balanced gameplay through physics-based tradeoffs

### Design Principles

The equipment type system follows these principles:

- **Incremental Expansion**: Start with core types (Weapons, Armor) that enable core gameplay, expand based on gameplay validation and community feedback
- **Stat Integration**: All equipment types must integrate with the characteristic system (ADR-0008)
- **Mass-Based Tradeoffs**: All equipment contributes to Mass, creating universal mobility tradeoff across all types
- **No Dead Types**: Only include equipment types that provide immediate gameplay value and create meaningful strategic decisions
- **Future-Proof Architecture**: Type system designed to accommodate expansion without breaking existing balance or requiring architectural changes
- **Strategic Clarity**: Each type should enable distinct, understandable build archetype (tank, DPS, mobile, balanced)
- **Physics Consistency**: All equipment affects movement through Mass impact on physics formulas from ADR-0006 and ADR-0007
- **Developer Accessibility**: Equipment categories should be clear and implementable by bot developers at various skill levels

### Future Considerations

**Load Bearing Gear**:
- May be added to enable "weapons platform" builds with multiple weapons or extra equipment
- Would expand equipment slot system (currently assumed single weapon, single armor)
- Useful for specialized loadouts and advanced build variety
- Should be added after weapon and armor balance is established

**Communication Equipment**:
- Will be added when team battle modes are implemented
- Requires bot-to-bot communication protocol (separate ADR)
- Enables team coordination and support roles
- Should be prioritized for competitive/team play support

**Sensor Equipment**:
- Will be added when fog of war visibility mechanics are implemented
- Requires visibility system and detection radius mechanics (separate ADR)
- Enables scout and reconnaissance roles
- Should be prioritized for game mode variety

**Boost Equipment**:
- May be added after core combat loop is balanced and validated
- Would require active ability system, temporal effects, cooldowns
- Enables burst gameplay and reactive tactics
- Should be subject of separate "Active Abilities" ADR
- May require resource management system (energy, charges)

**Mobility Equipment** (potential future type):
- Thrust enhancers, boosters, or hover systems
- Would modify movement characteristics without direct damage/defense
- Could enable mobility-focused builds as alternative to light equipment
- Requires separate ADR for movement enhancement mechanics

**Utility Equipment** (potential future type):
- Non-combat equipment for specialized purposes
- Could include repair systems, decoys, or terrain-altering equipment
- Enables support and control playstyles
- Requires separate ADR for utility mechanics framework

### Design Insights

- **Focused Scope**: Weapons and Armor provide maximum strategic depth relative to implementation complexity, enabling quality release over feature-packed chaos
- **Physics Foundation**: Equipment-derived Mass creates emergent gameplay through physics formulas; no need for explicit mechanics to balance heavy equipment (physics does the work)
- **Stat Integration**: Equipment modifies existing characteristics (Health, Defense, Mass) rather than adding new stats, keeping bot developer mental model simple
- **Future Extensibility**: Adding new equipment types is straightforward—define stat modifications, Mass contribution, and strategic purpose; physics handles the rest
- **Balance Considerations**: All equipment types share universal Mass contribution, creating consistent physics-based balance across types and enabling predictable interaction analysis

