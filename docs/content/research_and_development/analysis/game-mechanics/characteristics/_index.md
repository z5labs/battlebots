---
title: "Bot Characteristics"
description: "Core bot attributes and their gameplay effects"
type: docs
weight: 2
date: 2025-12-05
---

## Overview

Bot characteristics are the foundational attributes that define a bot's capabilities in battle. The game uses a five-stat system that determines how bots interact with the environment and each other:

- **Health**: Survivability and damage capacity
- **Energy**: Resource pool for performing actions
- **Speed**: Movement capability and positioning advantage
- **Power**: Offensive damage output
- **Defense**: Damage mitigation and resistance

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

## Energy

Energy is the fuel resource that powers all bot actions. Every action in combat consumes Energy, making it a critical limiting factor in bot performance.

**Key Properties**:
- **Energy Pool**: Maximum stored Energy available (TBD: placeholder 100-1000 range)
- **Regeneration Rate**: Energy restored per game tick (TBD)
- **Action Cost**: All actions (movement, attacks, abilities) consume Energy
- **Depletion Impact**: Bots cannot perform actions when Energy is insufficient

**Gameplay Impact**:
- Determines action frequency and sustained combat capability
- High Energy pools enable aggressive, action-intensive strategies
- Regeneration rate affects recovery time between action bursts
- Energy management becomes a core tactical consideration
- Low Energy bots must be more selective about action timing

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

## Power

Power is the damage multiplier that determines a bot's offensive output. This stat scales the base damage of attacks and offensive abilities.

**Key Properties**:
- **Damage Multiplier**: Scaling factor applied to attack damage (TBD: placeholder 1-10)
- **Offensive Scaling**: Higher Power means more damage per hit
- **No Energy Cost Reduction**: Power affects damage, not action efficiency
- **Applies to All Damage**: Affects basic attacks and damage-based abilities

**Gameplay Impact**:
- High Power enables burst damage strategies and quick eliminations
- Critical for glass cannon builds (high damage, low survivability)
- Must be balanced against Energy costs for sustained damage output
- Low Power bots may focus on attrition or support strategies
- Power differential determines time-to-kill in direct engagements

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
- Counters high Power strategies through damage nullification

## Stat Interactions

Bot characteristics don't operate in isolation - they create complex interactions that define combat dynamics:

### Effective Durability
Health and Defense combine to determine true survivability:
- **Effective HP** = Health × (1 + Defense modifier)
- High Defense multiplies the value of each Health point
- Balanced allocation is more efficient than single-stat stacking

### Damage Output Over Time
Power and Energy determine sustained damage capability:
- High Power with low Energy = burst damage, limited sustain
- Low Power with high Energy = sustained damage over time
- Action costs create natural balance between power and frequency

### Combat Range Control
Speed and Energy enable positioning strategies:
- Speed determines range closing/opening capability
- Energy enables sustained movement for kiting or pursuit
- Speed without Energy limits repositioning options

### Time-to-Kill vs. Time-to-Die
The relationship between offensive and defensive stats determines engagement outcomes:
- **Attacker TTK** = Defender Health / (Attacker Power - Defender Defense)
- **Defender TTK** = Attacker Health / (Defender Power - Attacker Defense)
- Favorable TTK ratios determine engagement winners
- Speed can modify these ratios by controlling engagement duration

### Resource Efficiency
All stats compete for allocation during bot design:
- Specialization (high single stat) vs. Generalization (balanced stats)
- Each stat point has opportunity cost in other areas
- Optimal builds depend on game mechanics and opponent meta
- No single dominant strategy (intended design goal)

---

**Note**: All numeric values and ranges in this document are preliminary placeholders marked TBD. Final values will be determined through gameplay testing and balance iteration.
