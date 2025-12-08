---
title: "[0008] Equipment and Loadout System"
description: >
    Bot customization through weapons, armor, and modules enabling strategic diversity
type: docs
weight: 8
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

Battle Bots requires a customization system that enables players to differentiate their bots before battle begins. We need to determine how bots customize their capabilities, what equipment types exist, how equipment affects performance, and how loadout constraints create strategic tradeoffs. The equipment system must create meaningful build diversity while integrating with the characteristics system (ADR-0007) and action system (ADR-0009).

Without a well-defined equipment system, we cannot:
- Create distinct bot builds and playstyles
- Enable pre-battle strategic decisions
- Provide stat customization and optimization paths
- Define which actions are available to bots
- Balance power vs. mobility tradeoffs
- Support future equipment additions and expansions

## Decision Drivers

* **Build Diversity** - Multiple viable loadouts should exist with distinct strengths and weaknesses
* **Pre-Battle Strategic Decisions** - Equipment choices should matter and create differentiation
* **Stat Modification Clarity** - Equipment effects on characteristics (ADR-0007) should be understandable
* **Equipment-Action Coupling** - Equipment should enable/disable specific actions (ADR-0009)
* **Tradeoff Mechanics** - Equipment should involve costs and benefits (no strictly superior choices)
* **Extensibility** - System should support future equipment additions without redesign
* **Protocol Integration** - Equipment selection must map to gRPC protocol (ADR-0004)
* **Developer Accessibility** - Loadout configuration should be straightforward for bot developers

## Considered Options

* **Option 1: No Customization** - All bots identical, differentiation through AI only
* **Option 2: Skill/Ability Selection** - Choose abilities from pool, no stat modification
* **Option 3: Equipment Loadout System** - Weapons, armor, modules with stat modifications and action enablement
* **Option 4: Class-Based System** - Predefined bot classes with fixed equipment

## Decision Outcome

Chosen option: "**Option 3: Equipment Loadout System**", because it enables meaningful pre-battle strategic decisions, creates stat tradeoffs through equipment modifications, enables diverse builds (DPS, Tank, Balanced, Stealth), provides extensibility for future equipment, and maps cleanly to the protocol while maintaining developer accessibility.

### Equipment and Loadout Specification

#### Equipment Categories

Equipment falls into three primary categories, each serving a distinct role in bot customization:

**Weapons**: Enable combat actions and determine offensive capabilities

- **Rifle**: Enables single-shot, precise attacks with moderate damage
  - Stat Effects: No modifications (baseline weapon)
  - Enabled Action: RifleShot (15 energy, 1 tick cooldown)
  - Tactical Profile: Reliable ranged damage, versatile baseline option

- **Shotgun**: Enables spray of projectiles with damage falloff based on distance
  - Stat Effects: -1 Speed (weight penalty), -1 Range (close-range weapon) (TBD)
  - Enabled Action: ShotgunBlast (20 energy, 2 tick cooldown)
  - Tactical Profile: High close-range burst damage, requires positioning

**Armor**: Provides defensive bonuses and damage mitigation

- **Light Armor**: Minimal defense bonus, no speed penalty
  - Stat Effects: +1 Defense, +0 Speed (TBD)
  - Tactical Profile: Maintains mobility for evasive playstyles

- **Medium Armor**: Balanced defense bonus with minor speed penalty
  - Stat Effects: +2 Defense, -1 Speed (TBD)
  - Tactical Profile: Versatile option for balanced builds

- **Heavy Armor**: Significant defense bonus with substantial speed penalty
  - Stat Effects: +3 Defense, -2 Speed (TBD)
  - Tactical Profile: Maximizes survivability for tank builds

**Modules**: Provide utility functions, special abilities, and tactical advantages

- **Boost Engine**: Enables temporary speed increases
  - Stat Effects: +1 Max Speed, -1 Energy Capacity (TBD)
  - Enabled Action: Boost (variable energy cost and duration)
  - Tactical Profile: Repositioning and engagement control

- **Repair Kit**: Allows limited self-repair during combat
  - Stat Effects: No continuous stat modifications
  - Enabled Action: Repair (restore HP, usage limits TBD)
  - Tactical Profile: Extended combat effectiveness through self-sustain

- **Sensor Array**: Increases detection range and provides tactical information
  - Stat Effects: +2 Detection Range (TBD)
  - Enabled Action: Scan (5 energy, variable cooldown)
  - Tactical Profile: Information advantage and awareness

- **Stealth Module**: Reduces detection range by enemies
  - Stat Effects: -2 Enemy Detection Range, -1 Defense (exposed systems) (TBD)
  - Enabled Action: Cloak (variable energy cost and duration)
  - Tactical Profile: Avoid detection and enable surprise attacks

#### Loadout Constraints

Each bot has limited loadout capacity to prevent overpowered configurations and maintain balance:

- **1 Weapon Slot** (TBD - may expand in future iterations)
- **1 Armor Slot** (TBD - single armor type per bot)
- **2 Module Slots** (TBD - count may vary based on testing)

These slot limitations force meaningful choices during bot configuration. Players cannot equip all available equipment and must prioritize based on strategy.

**Alternative Configurations Under Consideration**:
- Multiple weapon slots with restrictions (e.g., one primary, one secondary)
- Armor layering systems with weight limits
- Variable module slots based on bot chassis type
- Equipment weight/point systems instead of fixed slots

#### Stat Modification Mechanics

Equipment modifies bot characteristics (ADR-0007), creating different performance profiles:

**Stat Calculation Formula**:
```
Final Stat = Base Stat + Weapon Modifier + Armor Modifier + Module Modifiers
```

**Example Calculation** (all values TBD):
```
Bot Base Speed: 10
Heavy Armor Speed Penalty: -2
Boost Engine Speed Bonus: +1
Final Speed: 10 - 2 + 1 = 9
```

**Mass Contribution**: All equipment contributes to bot Mass (ADR-0007), creating natural mobility-firepower tradeoffs:
- Heavy weapons and armor increase Mass, reducing Effective Speed
- Light loadouts minimize Mass for maximum mobility
- Equipment choices create cascading effects on movement capabilities

#### Action Requirements

Equipment directly determines which actions (ADR-0009) are available during combat:

**Weapon-Dependent Actions**:
- **RifleShot**: Requires Rifle equipped
- **ShotgunBlast**: Requires Shotgun equipped

Without the appropriate weapon, these actions are unavailable in the bot's action set.

**Module-Dependent Actions**:
- **Boost**: Requires Boost Engine module
- **Repair**: Requires Repair Kit module
- **Scan**: Requires Sensor Array module
- **Cloak**: Requires Stealth Module module

Module-dependent actions provide tactical options beyond direct combat.

**Universal Actions** (always available):
- **Move**: Movement in the 2D battle space (ADR-0006)
- **Evade**: Defensive positioning or dodge action
- **Block**: Damage reduction stance
- **Shield**: Energy-based damage absorption

#### Example Loadouts

The following example loadouts demonstrate the range of viable bot configurations and playstyle diversity:

**DPS (Damage Per Second) Build**

*Philosophy*: Maximize offensive capability and damage output. Accept lower survivability for high burst damage potential.

*Equipment*:
- **Weapon**: Shotgun (high damage spray at close range)
- **Armor**: Light Armor (minimal defense, maintain mobility)
- **Module 1**: Boost Engine (enable repositioning for close-range attacks)
- **Module 2**: Sensor Array (track enemies for optimal engagement range)

*Stat Profile* (TBD values):
- Attack: 12 (high)
- Defense: 6 (low)
- Speed: 10 (good)
- Energy: 10 (average)

*Strategy*: Close distance quickly using Boost, deliver devastating shotgun blasts at close range, use sensors to track enemy positions for optimal engagement. High risk, high reward playstyle.

*Weaknesses*: Vulnerable to sustained fire, limited survivability if caught in poor position, ineffective at long range.

---

**Tank Build**

*Philosophy*: Maximum survivability and staying power. Control space through defensive presence and outlast opponents.

*Equipment*:
- **Weapon**: Rifle (reliable baseline offense)
- **Armor**: Heavy Armor (maximum damage reduction)
- **Module 1**: Repair Kit (self-sustain during extended engagements)
- **Module 2**: Sensor Array (maintain awareness despite low mobility)

*Stat Profile* (TBD values):
- Attack: 8 (moderate)
- Defense: 13 (very high)
- Speed: 8 (low)
- Energy: 10 (average)

*Strategy*: Hold key positions, absorb damage with heavy armor, use Repair to extend combat effectiveness, rely on sensors to track enemy movements. Win through attrition rather than burst damage.

*Weaknesses*: Low mobility makes positioning critical, vulnerable to kiting strategies, limited offensive pressure.

---

**Balanced Build**

*Philosophy*: Versatile configuration capable of adapting to various situations. No extreme weaknesses, no extreme strengths.

*Equipment*:
- **Weapon**: Rifle (reliable ranged damage)
- **Armor**: Medium Armor (reasonable defense without severe speed penalty)
- **Module 1**: Sensor Array (tactical awareness)
- **Module 2**: Repair Kit (survivability boost during combat)

*Stat Profile* (TBD values):
- Attack: 10 (above average)
- Defense: 10 (above average)
- Speed: 9 (average)
- Energy: 10 (average)

*Strategy*: Use sensors for tactical awareness, maintain optimal engagement distance with rifle, use repair to extend combat effectiveness, rely on well-rounded stats to handle unexpected situations. Adaptable to opponent strategies.

*Weaknesses*: Lacks specialization, may be outperformed by specialized builds in their areas of strength.

---

**Stealth/Scout Build**

*Philosophy*: Control battlefield through information advantage and mobility rather than direct damage.

*Equipment*:
- **Weapon**: Rifle (precise long-range attacks)
- **Armor**: Light Armor (maintain mobility)
- **Module 1**: Sensor Array (information gathering)
- **Module 2**: Stealth Module (avoid detection)

*Stat Profile* (TBD values):
- Attack: 8 (moderate)
- Defense: 6 (low)
- Speed: 12 (very high)
- Energy: 10 (average)

*Strategy*: Scout with sensors, avoid detection with stealth, strike from unexpected angles with rifle shots, rely on mobility to disengage. Win through tactical advantage and superior positioning rather than sustained combat.

*Weaknesses*: Ineffective in forced direct combat, relies heavily on stealth and positioning mechanics working as intended, vulnerable if detected in poor position.

### Consequences

* Good, because equipment loadout system creates meaningful build diversity (DPS, Tank, Balanced, Stealth)
* Good, because pre-battle equipment selection enables strategic planning and counter-play
* Good, because stat modifications are clear and calculable (simple additive formula)
* Good, because equipment-action coupling creates distinct tactical options for different builds
* Good, because Defense vs. Speed tradeoffs in armor create natural power-mobility choices
* Good, because module variety enables utility and tactical flexibility
* Good, because loadout constraints force meaningful choices (can't equip everything)
* Good, because extensible design supports future equipment additions
* Good, because example builds demonstrate viable diversity and counter-play options
* Neutral, because stat modification values (all TBD) require extensive playtesting
* Neutral, because loadout slot counts (1 weapon, 1 armor, 2 modules) need validation through testing
* Neutral, because equipment balance requires tuning to ensure no dominant loadout
* Bad, because equipment system adds complexity for bot developers to understand
* Bad, because stat modifications and action requirements create dependencies to track
* Bad, because equipment balance is critical - single dominant loadout would eliminate diversity
* Bad, because module-dependent actions require protocol validation of equipment

### Confirmation

The decision will be confirmed through:

1. Implementation of equipment system in game server with stat modification calculations
2. Bot SDK exposing loadout configuration and equipment selection interface
3. Playtesting with all four example builds (DPS, Tank, Balanced, Stealth) to validate viability
4. Competitive balance analysis ensuring no single loadout dominates all matchups
5. Counter-play validation confirming viable counter-builds exist for each archetype
6. Developer feedback on equipment configuration clarity and accessibility
7. Protocol integration testing for equipment selection and validation

## Pros and Cons of the Options

### Option 1: No Customization

All bots identical, differentiation through AI implementation only.

* Good, because simplest possible implementation
* Good, because eliminates equipment balance concerns
* Good, because focuses all differentiation on AI quality
* Good, because easier for new developers to start (no loadout decisions)
* Neutral, because may be sufficient for initial proof-of-concept
* Bad, because eliminates pre-battle strategic decisions
* Bad, because no build diversity or customization depth
* Bad, because limits long-term engagement and replayability
* Bad, because reduces strategic gameplay to AI implementation skill only

### Option 2: Skill/Ability Selection

Choose abilities from pool, no stat modification system.

* Good, because creates action diversity without stat complexity
* Good, because simpler than full equipment system
* Good, because focuses customization on tactical options (abilities)
* Neutral, because enables some build diversity through ability choices
* Bad, because lacks stat customization and optimization depth
* Bad, because no power vs. mobility tradeoffs
* Bad, because abilities may be imbalanced without stat constraints
* Bad, because equipment flavor and thematic coherence is lost

### Option 3: Equipment Loadout System

Weapons, armor, modules with stat modifications and action enablement (CHOSEN).

* Good, because creates strategic depth through stat customization
* Good, because equipment-action coupling creates diverse tactical options
* Good, because Defense vs. Speed tradeoffs enable distinct playstyles
* Good, because enables multiple viable builds (DPS, Tank, Balanced, Stealth)
* Good, because extensible to future equipment additions
* Good, because equipment choices create cascading effects via Mass (ADR-0007)
* Good, because loadout constraints force meaningful tradeoffs
* Good, because thematically coherent (weapons, armor, modules)
* Neutral, because requires extensive balance tuning
* Neutral, because stat modification complexity is moderate (simple additive formula)
* Bad, because more complex than ability-only or no customization
* Bad, because equipment validation adds protocol overhead
* Bad, because balance is critical - single dominant loadout kills diversity

### Option 4: Class-Based System

Predefined bot classes (Tank, DPS, Scout) with fixed equipment.

* Good, because guarantees build diversity through class design
* Good, because simpler for developers (just pick class)
* Good, because easier to balance (fewer total configurations)
* Neutral, because provides some customization through class choice
* Bad, because eliminates loadout customization and optimization
* Bad, because rigid class definitions may not match all desired playstyles
* Bad, because limits strategic depth (no stat allocation decisions)
* Bad, because less extensible (new equipment requires new classes or class modifications)
* Bad, because reduces pre-battle strategic planning to class selection

## More Information

### Related Documentation

- **[ADR-0007: Bot Characteristics System](0007-bot-characteristics-system.md)**: Stats that equipment modifies (Health, Speed, Defense, Mass)

- **[ADR-0009: Bot Actions and Resource Management](0009-bot-actions-resource-management.md)**: Actions that equipment enables or requires

- **[Equipment System Analysis](../analysis/game-mechanics/equipment/)**: Detailed technical specifications for equipment types and loadouts

- **[ADR-0005: 1v1 Battle Orchestration](0005-1v1-battle-orchestration.md)**: High-level battle flow using these equipment loadouts

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for equipment selection and validation

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Playtesting with all four example builds to validate competitive viability
2. Balance modeling to ensure no single loadout dominates all matchups
3. Counter-play analysis confirming viable counter-builds exist for each archetype
4. Stat tuning to create meaningful tradeoffs (Defense vs. Speed, Power vs. Mobility)
5. Module effectiveness testing to ensure utility value justifies module slots
6. Weapon balance to ensure Rifle and Shotgun are situationally viable
7. Armor balance to ensure Light/Medium/Heavy each have optimal use cases

**Key Design Insights**:
- Equipment-action coupling ensures loadout choices directly affect tactical options
- Defense vs. Speed tradeoffs in armor create natural tank vs. mobile playstyle spectrum
- Module slots enable tactical customization beyond raw combat stats
- Example builds demonstrate diversity while providing optimization starting points

**Future Enhancements**:
- Equipment rarity/tier systems for progression
- Equipment upgrade mechanics for long-term advancement
- Synergy bonuses for specific equipment combinations
- Dynamic equipment swapping during combat (risky but strategic)
- Equipment durability and damage mechanics
- Additional weapon types (laser, plasma, melee)
- Additional armor types (reactive, ablative, energy shields)
- Additional modules (hacking, jamming, ECM, decoys)

### Design Principles

The equipment system follows these principles:
- **Tradeoffs over Power**: All equipment involves costs and benefits
- **Diversity over Dominance**: Multiple loadouts should be competitively viable
- **Coupling over Independence**: Equipment determines both stats and actions
- **Constraints over Freedom**: Slot limits force meaningful choices
- **Extensibility over Finality**: Framework supports future equipment additions
