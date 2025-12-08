---
title: "[0008] Equipment and Loadout System"
description: >
    Bot customization through weapons and armor enabling strategic diversity
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

## Decision Outcome

The equipment system consists of two equipment categories (Weapons, Armor) with the following initial equipment options:
- **Weapons**: Rifle, Shotgun
- **Armor**: Light Armor, Medium Armor, Heavy Armor

Each bot equips a loadout with 1 weapon and 1 armor. Equipment modifies bot characteristics (ADR-0007) and determines available actions (ADR-0009), creating distinct tactical options for different builds.

### Consequences

* Good, because equipment loadout system creates meaningful build diversity (DPS, Tank, Balanced, Mobile/Skirmisher)
* Good, because pre-battle equipment selection enables strategic planning and counter-play
* Good, because stat modifications are clear and calculable (simple additive formula)
* Good, because equipment-action coupling creates distinct tactical options for different builds
* Good, because Defense vs. Speed tradeoffs in armor create natural power-mobility choices
* Good, because loadout constraints force meaningful choices
* Good, because modular structure (each equipment type and option in dedicated sections) makes equipment additions straightforward
* Good, because initial equipment set provides foundation for expansion
* Good, because example builds demonstrate viable diversity and counter-play options
* Neutral, because stat modification values (all TBD) require extensive playtesting
* Neutral, because loadout slot counts (1 weapon, 1 armor) need validation through testing
* Neutral, because equipment balance requires tuning to ensure no dominant loadout
* Bad, because equipment system adds complexity for bot developers to understand
* Bad, because stat modifications and action requirements create dependencies to track
* Bad, because equipment balance is critical - single dominant loadout would eliminate diversity

### Confirmation

The decision will be confirmed through:

1. Implementation of equipment system in game server with stat modification calculations
2. Bot SDK exposing loadout configuration and equipment selection interface
3. Playtesting with all four example builds (DPS, Tank, Balanced, Stealth) to validate viability
4. Competitive balance analysis ensuring no single loadout dominates all matchups
5. Counter-play validation confirming viable counter-builds exist for each archetype
6. Developer feedback on equipment configuration clarity and accessibility
7. Protocol integration testing for equipment selection and validation

## Weapons

Weapons enable combat actions and determine offensive capabilities. Each weapon provides a unique attack action with different energy costs, damage patterns, and tactical applications. All weapons contribute to bot Mass (ADR-0007), with heavier weapons providing greater firepower at the cost of mobility.

### Rifle

Standard precision weapon enabling reliable ranged attacks.

**Stat Effects**:
- No modifications (baseline weapon)
- Mass Contribution: TBD (baseline weapon mass)

**Tactical Profile**:
- Reliable ranged damage output
- Versatile baseline option suitable for any playstyle
- No stat penalties, maintains mobility
- Effective at medium to long range

### Shotgun

Close-range weapon enabling devastating burst damage with damage falloff based on distance.

**Stat Effects**:
- -1 Speed (weight penalty) (TBD)
- -1 Range (close-range weapon) (TBD)
- Mass Contribution: TBD (higher than Rifle)

**Tactical Profile**:
- High close-range burst damage
- Requires positioning to maximize effectiveness
- Weight penalty reduces mobility
- Ineffective at long range due to damage falloff

## Armor

Armor provides defensive bonuses and damage mitigation. Armor directly modifies Defense and Speed characteristics (ADR-0007), creating tradeoffs between survivability and mobility. All armor contributes to bot Mass, with heavier armor providing greater protection at the cost of reduced Speed.

### Light Armor

Minimal protection that maintains mobility for evasive playstyles.

**Stat Effects**:
- +1 Defense (TBD)
- +0 Speed (no speed penalty) (TBD)
- Mass Contribution: TBD (minimal)

**Tactical Profile**:
- Minimal defense bonus maintains baseline survivability
- No speed penalty preserves mobility
- Optimal for evasion-based and high-mobility builds
- Relies on Speed rather than damage absorption

### Medium Armor

Balanced protection with moderate defensive bonus and minor speed penalty.

**Stat Effects**:
- +2 Defense (TBD)
- -1 Speed (TBD)
- Mass Contribution: TBD (moderate)

**Tactical Profile**:
- Reasonable defense without severe mobility cost
- Versatile option for balanced builds
- Moderate survivability increase with manageable speed reduction
- Suitable for all-around playstyles

### Heavy Armor

Maximum protection with significant defensive bonus and substantial speed penalty.

**Stat Effects**:
- +3 Defense (TBD)
- -2 Speed (TBD)
- Mass Contribution: TBD (high)

**Tactical Profile**:
- Maximum damage reduction for tank builds
- Significant speed penalty limits mobility
- Enables prolonged engagements and damage absorption
- Requires positional awareness due to low mobility

## Example Loadouts

The following example loadouts demonstrate the range of viable bot configurations and playstyle diversity using the equipment options defined above:

**DPS (Damage Per Second) Build**

*Philosophy*: Maximize offensive capability and damage output. Accept lower survivability for high burst damage potential.

*Equipment*:
- **Weapon**: Shotgun (high damage spray at close range)
- **Armor**: Light Armor (minimal defense, maintain mobility)

*Stat Profile* (TBD values):
- Attack: 12 (high)
- Defense: 6 (low)
- Speed: 10 (good)
- Energy: 10 (average)

*Strategy*: Close distance quickly, deliver devastating shotgun blasts at close range. High risk, high reward playstyle focused on burst damage.

*Weaknesses*: Vulnerable to sustained fire, limited survivability if caught in poor position, ineffective at long range.

---

**Tank Build**

*Philosophy*: Maximum survivability and staying power. Control space through defensive presence and outlast opponents.

*Equipment*:
- **Weapon**: Rifle (reliable baseline offense)
- **Armor**: Heavy Armor (maximum damage reduction)

*Stat Profile* (TBD values):
- Attack: 8 (moderate)
- Defense: 13 (very high)
- Speed: 8 (low)
- Energy: 10 (average)

*Strategy*: Hold key positions, absorb damage with heavy armor, maintain reliable ranged offense with rifle. Win through attrition rather than burst damage.

*Weaknesses*: Low mobility makes positioning critical, vulnerable to kiting strategies, limited offensive pressure.

---

**Balanced Build**

*Philosophy*: Versatile configuration capable of adapting to various situations. No extreme weaknesses, no extreme strengths.

*Equipment*:
- **Weapon**: Rifle (reliable ranged damage)
- **Armor**: Medium Armor (reasonable defense without severe speed penalty)

*Stat Profile* (TBD values):
- Attack: 10 (above average)
- Defense: 10 (above average)
- Speed: 9 (average)
- Energy: 10 (average)

*Strategy*: Maintain optimal engagement distance with rifle, rely on well-rounded stats to handle unexpected situations. Adaptable to opponent strategies.

*Weaknesses*: Lacks specialization, may be outperformed by specialized builds in their areas of strength.

---

**Mobile/Skirmisher Build**

*Philosophy*: Control battlefield through mobility and positioning rather than direct combat superiority.

*Equipment*:
- **Weapon**: Rifle (precise long-range attacks)
- **Armor**: Light Armor (maintain mobility)

*Stat Profile* (TBD values):
- Attack: 8 (moderate)
- Defense: 6 (low)
- Speed: 12 (very high)
- Energy: 10 (average)

*Strategy*: Strike from optimal angles with rifle shots, rely on superior mobility to control engagement distance and disengage when needed. Win through tactical advantage and superior positioning rather than sustained combat.

*Weaknesses*: Ineffective in forced direct combat, vulnerable if caught in poor position, relies on maintaining optimal range.

## More Information

### Related Documentation

- **[ADR-0007: Bot Characteristics System](0007-bot-characteristics-system.md)**: Stats that equipment modifies (Health, Speed, Defense, Mass)

- **[ADR-0009: Bot Actions and Resource Management](0009-bot-actions-resource-management.md)**: Actions that equipment enables or requires

- **[Equipment System Analysis](../analysis/game-mechanics/equipment/)**: Detailed technical specifications for equipment types and loadouts

- **[ADR-0005: 1v1 Battle Orchestration](0005-1v1-battle-orchestration.md)**: High-level battle flow using these equipment loadouts

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for equipment selection and validation

### Implementation Notes

**Modular Structure and Extensibility**:

This ADR defines the initial equipment options for each category (2 weapons, 3 armor types). The modular structure—with each equipment type and option in dedicated sections—enables straightforward expansion:

- **Adding New Equipment Options**: New equipment can be added by simply adding new subsections under the appropriate category (Weapons or Armor). For example, "Laser Rifle" could be added as a new subsection under Weapons.
- **Adding New Equipment Categories**: New equipment categories (e.g., "Chassis Types" or "Power Cores") can be introduced by adding new top-level sections to the Equipment System Specification.
- **Modifying Existing Equipment**: Specific equipment subsections can be superseded to rebalance or redesign individual items without affecting other equipment.

**Stat Value Refinement**:

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Playtesting with all four example builds to validate competitive viability
2. Balance modeling to ensure no single loadout dominates all matchups
3. Counter-play analysis confirming viable counter-builds exist for each archetype
4. Stat tuning to create meaningful tradeoffs (Defense vs. Speed, Power vs. Mobility)
5. Weapon balance to ensure Rifle and Shotgun are situationally viable
6. Armor balance to ensure Light/Medium/Heavy each have optimal use cases

**Key Design Insights**:
- Equipment-action coupling ensures loadout choices directly affect tactical options
- Defense vs. Speed tradeoffs in armor create natural tank vs. mobile playstyle spectrum
- Example builds demonstrate diversity while providing optimization starting points
- Modular structure (dedicated sections per equipment) simplifies expansion and balance changes

### Design Principles

The equipment system follows these principles:
- **Tradeoffs over Power**: All equipment involves costs and benefits
- **Diversity over Dominance**: Multiple loadouts should be competitively viable
- **Coupling over Independence**: Equipment determines both stats and actions
- **Constraints over Freedom**: Slot limits force meaningful choices
- **Extensibility over Finality**: Framework supports future equipment additions
