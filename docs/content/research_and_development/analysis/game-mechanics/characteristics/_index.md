---
title: "Bot Characteristics"
description: "Core bot attributes and their gameplay effects"
type: docs
weight: 2
date: 2025-12-05
---

## Overview

Bot characteristics are the foundational attributes that define a bot's capabilities in battle. The game uses a four-stat system that determines how bots interact with the environment and each other:

- **Health**: Survivability and damage capacity
- **Speed**: Movement capability and positioning advantage
- **Defense**: Damage mitigation and resistance
- **Mass**: Physical weight determined by equipped items

Each characteristic directly impacts gameplay mechanics and creates strategic tradeoffs during bot design. Bots must balance their stat allocation to match their intended combat strategy.

## Health

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

## Speed

Speed determines how quickly a bot can move through the 2D battle space. This stat directly affects positioning, engagement control, and evasion capabilities.

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

## Defense

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
- Reduces effective damage from enemy attacks

## Mass

Mass represents the total physical weight of a bot, determined by the equipment and components it carries. Unlike other characteristics, Mass is not directly allocated but is the cumulative result of loadout choices.

**Key Properties**:
- **Equipment-Derived**: Mass is calculated from the sum of all equipped items
- **Dynamic Value**: Changes based on equipped weapons, armor, and components
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
- Components: Power cores, sensors, and systems each contribute Mass
- Loadout variety creates diverse Mass profiles across bot builds

## Stat Interactions

Bot characteristics don't operate in isolation - they create complex interactions that define combat dynamics:

### Effective Durability
Health and Defense combine to determine true survivability:
- **Effective HP** = Health × (1 + Defense modifier)
- High Defense multiplies the value of each Health point
- Balanced allocation is more efficient than single-stat stacking

### Mass and Mobility
Mass directly impacts effective movement capability:
- **Effective Speed** = Base Speed / Mass modifier
- Heavy equipment reduces Speed, creating mobility-firepower tradeoffs
- Light loadouts maximize Speed at the cost of offensive/defensive power
- Equipment choices fundamentally alter tactical capabilities

### Combat Positioning
Speed enables tactical positioning and engagement control:
- High Speed allows kiting, pursuit, and disengagement
- Speed advantage determines range control in combat
- Mass penalties from equipment reduce positioning flexibility
- Lightweight builds gain tactical mobility at the cost of durability

### Survivability Tradeoffs
Defensive investment creates complex build choices:
- High Health + Low Defense = vulnerable to sustained damage
- Low Health + High Defense = vulnerable to burst damage
- Mass from defensive equipment reduces evasion capability
- Optimal defensive strategy depends on threat profile

### Loadout Optimization
Equipment choices create cascading effects across all characteristics:
- Heavy weapons increase offensive capability but reduce Speed through Mass
- Armor improves Defense but adds Mass that limits mobility
- Lightweight builds sacrifice protection for superior positioning
- No equipment configuration dominates all scenarios (intended design goal)

---

**Note**: All numeric values and ranges in this document are preliminary placeholders marked TBD. Final values will be determined through gameplay testing and balance iteration.
