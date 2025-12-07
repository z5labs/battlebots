---
title: "Equipment System"
description: "Bot customization through weapons, armor, and modules"
type: docs
weight: 3
date: 2025-12-05
---

## Overview

The equipment system provides the primary mechanism for bot customization in Battle Bots. Before entering battle, players configure their bot's loadout by selecting weapons, armor, and utility modules. These equipment choices fundamentally shape how a bot performs in combat, creating distinct playstyles through stat modifications and enabling specific combat actions.

Equipment-based customization creates strategic depth by forcing players to make tradeoffs. A heavily armored bot gains survivability but sacrifices speed. A bot loaded with offensive modules gains firepower but may lack defensive options. The equipment system transforms identical base bots into diverse combat units, each optimized for different tactical approaches.

This document defines the equipment types, loadout constraints, stat modification mechanics, and example configurations that demonstrate the range of viable bot builds.

## Equipment Types

Equipment falls into three primary categories, each serving a distinct role in bot customization:

### Weapons

Weapons enable combat actions and determine a bot's offensive capabilities. Each weapon unlocks specific attack actions with unique characteristics:

- **Rifle**: Enables single-shot, precise attacks with moderate damage (TBD exact values)
- **Shotgun**: Enables spray of projectiles attack with damage falloff based on distance (TBD)

Weapon selection directly determines which offensive actions are available to a bot during combat. A bot without a weapon equipped cannot perform any attack actions.

### Armor

Armor provides defensive bonuses and damage mitigation. Different armor types offer varying levels of protection with corresponding mobility costs:

- **Light Armor**: Minimal defense bonus, no speed penalty - suitable for fast, evasive builds (TBD values)
- **Medium Armor**: Balanced defense bonus with minor speed penalty - versatile option for most builds (TBD values)
- **Heavy Armor**: Significant defense bonus with substantial speed penalty - maximizes survivability (TBD values)

Armor affects how much damage a bot takes when hit and may influence movement speed and energy consumption.

### Modules

Modules provide utility functions, special abilities, and tactical advantages beyond direct combat:

- **Boost Engine**: Enables temporary speed increases for repositioning or escape (TBD duration and cooldown)
- **Repair Kit**: Allows limited self-repair during combat (TBD healing amount and usage limits)
- **Sensor Array**: Increases detection range and provides tactical information (TBD range and intel bonuses)
- **Stealth Module**: Reduces detection range by enemies (TBD effectiveness and duration)

Modules enable tactical flexibility and allow bots to adapt to different combat scenarios. Module selection determines which utility actions are available during battle.

## Equipment Slots

Each bot has a limited loadout capacity to prevent overpowered configurations and maintain game balance. The current proposed loadout structure is:

- **1 Weapon Slot** (TBD - may expand to multiple weapon slots in future iterations)
- **1 Armor Slot** (TBD - single armor type per bot)
- **2 Module Slots** (TBD - count may vary based on testing)

These slot limitations force meaningful choices during bot configuration. Players cannot equip all available equipment types and must prioritize based on their intended strategy.

**Note**: Slot counts and restrictions are subject to change during playtesting and balance iterations. Alternative configurations under consideration include:

- Multiple weapon slots with restrictions (e.g., one primary, one secondary)
- Armor layering systems with weight limits
- Variable module slots based on bot chassis type
- Equipment weight/point systems instead of fixed slots

## Stat Modifications

Equipment modifies a bot's base characteristics, creating different performance profiles. All values below are TBD and will be refined through playtesting and balance analysis.

### Weapon Stat Effects

Weapons primarily enable actions but may also modify stats:

- **Rifle**: No stat modifications (baseline weapon)
- **Shotgun**: -1 Speed (weight penalty), -1 Range (close-range weapon) (TBD)

### Armor Stat Effects

Armor provides defense at the cost of mobility:

- **Light Armor**: +1 Defense, +0 Speed (TBD)
- **Medium Armor**: +2 Defense, -1 Speed (TBD)
- **Heavy Armor**: +3 Defense, -2 Speed (TBD)

### Module Stat Effects

Modules provide diverse effects beyond simple stat changes:

- **Boost Engine**: +1 Max Speed, -1 Energy Capacity (TBD)
- **Repair Kit**: +0 continuous stats, Action: Restore 20 HP (TBD usage limits)
- **Sensor Array**: +2 Detection Range, +0 other stats (TBD)
- **Stealth Module**: -2 Enemy Detection Range, -1 Defense (exposed systems) (TBD)

### Stat Calculation

Final bot stats are calculated as:

```
Final Stat = Base Stat + Weapon Modifier + Armor Modifier + Module Modifiers
```

Example calculation (all values TBD):
```
Bot Base Speed: 10
Heavy Armor Speed Penalty: -2
Boost Engine Speed Bonus: +1
Final Speed: 10 - 2 + 1 = 9
```

## Action Requirements

Certain actions require specific equipment to be used. This creates a direct link between loadout choices and tactical options during combat.

### Weapon-Dependent Actions

- **RifleShot**: Requires Rifle equipped
- **ShotgunBlast**: Requires Shotgun equipped

Without the appropriate weapon, these actions are unavailable in the bot's action set.

### Module-Dependent Actions

- **Boost**: Requires Boost Engine module
- **Repair**: Requires Repair Kit module
- **Scan**: Requires Sensor Array module
- **Cloak**: Requires Stealth Module module

Module-dependent actions provide tactical options beyond direct combat, enabling diverse strategies.

### Universal Actions

Some actions are always available regardless of equipment:

- **Move**: Movement in the 2D battle space (TBD mechanics)
- **Evade**: Defensive positioning or dodge action (TBD mechanics)
- **Wait**: Skip turn or charge energy (TBD mechanics)

## Example Loadouts

The following example loadouts demonstrate the range of viable bot configurations. All stat values are TBD and serve as illustrations of the customization system's strategic depth.

### DPS (Damage Per Second) Build

**Philosophy**: Maximize offensive capability and damage output. Accept lower survivability in exchange for high burst damage potential.

**Equipment**:
- **Weapon**: Shotgun (high damage spray at close range)
- **Armor**: Light Armor (minimal defense, maintain mobility)
- **Module 1**: Boost Engine (enable repositioning for close-range attacks)
- **Module 2**: Sensor Array (track enemies for optimal engagement range)

**Stat Profile** (TBD values):
- Attack: 12 (high)
- Defense: 6 (low)
- Speed: 10 (good)
- Energy: 10 (average)

**Strategy**: Close distance quickly using Boost, deliver devastating shotgun blasts at close range, use sensors to track enemy positions for optimal engagement. High risk, high reward playstyle.

**Weaknesses**: Vulnerable to sustained fire, limited survivability if caught in poor position, ineffective at long range.

### Tank Build

**Philosophy**: Maximum survivability and staying power. Control space through defensive presence and outlast opponents.

**Equipment**:
- **Weapon**: Rifle (reliable baseline offense)
- **Armor**: Heavy Armor (maximum damage reduction)
- **Module 1**: Repair Kit (self-sustain during extended engagements)
- **Module 2**: Sensor Array (maintain awareness despite low mobility)

**Stat Profile** (TBD values):
- Attack: 8 (moderate)
- Defense: 13 (very high)
- Speed: 8 (low)
- Energy: 10 (average)

**Strategy**: Hold key positions, absorb damage with heavy armor, use Repair to extend combat effectiveness, rely on sensors to track enemy movements. Win through attrition rather than burst damage.

**Weaknesses**: Low mobility makes positioning critical, vulnerable to kiting strategies, limited offensive pressure.

### Balanced Build

**Philosophy**: Versatile configuration capable of adapting to various situations. No extreme weaknesses, no extreme strengths.

**Equipment**:
- **Weapon**: Rifle (reliable ranged damage)
- **Armor**: Medium Armor (reasonable defense without severe speed penalty)
- **Module 1**: Sensor Array (tactical awareness)
- **Module 2**: Repair Kit (survivability boost during combat)

**Stat Profile** (TBD values):
- Attack: 10 (above average)
- Defense: 10 (above average)
- Speed: 9 (average)
- Energy: 10 (average)

**Strategy**: Use sensors for tactical awareness, maintain optimal engagement distance with rifle, use repair to extend combat effectiveness, rely on well-rounded stats to handle unexpected situations. Adaptable to opponent strategies.

**Weaknesses**: Lacks specialization, may be outperformed by specialized builds in their areas of strength.

### Stealth/Scout Build

**Philosophy**: Control the battlefield through information advantage and mobility rather than direct damage.

**Equipment**:
- **Weapon**: Rifle (precise long-range attacks)
- **Armor**: Light Armor (maintain mobility)
- **Module 1**: Sensor Array (information gathering)
- **Module 2**: Stealth Module (avoid detection)

**Stat Profile** (TBD values):
- Attack: 8 (moderate)
- Defense: 6 (low)
- Speed: 12 (very high)
- Energy: 10 (average)

**Strategy**: Scout with sensors, avoid detection with stealth, strike from unexpected angles with rifle shots, rely on mobility to disengage. Win through tactical advantage and superior positioning rather than sustained combat.

**Weaknesses**: Ineffective in forced direct combat, relies heavily on stealth and positioning mechanics working as intended, vulnerable if detected in poor position.

## Design Considerations

The equipment system must balance several competing concerns:

1. **Meaningful Choices**: Equipment selections must create genuinely different playstyles, not merely cosmetic variations
2. **Balance**: No single loadout should dominate all situations; viable counter-builds should exist for any strategy
3. **Clarity**: Players should understand how equipment affects their bot's capabilities without consulting extensive documentation
4. **Extensibility**: The system should accommodate future equipment additions without requiring fundamental redesign

All values and mechanics in this document are TBD and subject to revision based on:
- Playtesting results and player feedback
- Mathematical balance modeling
- Implementation complexity and performance considerations
- Overall game design direction decisions

Future iterations may introduce:
- Equipment rarity/tier systems
- Equipment upgrade mechanics
- Synergy bonuses for specific equipment combinations
- Dynamic equipment swapping during combat
- Equipment durability and damage mechanics
