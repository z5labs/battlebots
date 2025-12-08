---
title: "[0009] Bot Actions and Resource Management"
description: >
    Action system governing bot behavior and creating tactical decision-making in real-time battles
type: docs
weight: 9
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

Battle Bots requires an action system that governs how bots behave in real-time battles. We need to define what actions bots can perform, how actions are constrained (costs, cooldowns, equipment requirements), and how the resource management system prevents action spam while maintaining fluid gameplay. The action system must integrate with the real-time continuous gameplay model and support the equipment-based customization system (ADR-0007).

Without a well-defined action system, we cannot:
- Design the game server tick processing and action resolution
- Create bot SDK action interfaces and methods
- Define protocol messages for action submission and validation
- Implement resource management (energy) and timing constraints (cooldowns)
- Balance tactical decision-making and action variety
- Prevent action spam or create action starvation scenarios

## Decision Drivers

* **Real-time Gameplay Support** - Must work with continuous tick-based game loop (not turn-based)
* **Resource Management Depth** - Should create strategic decisions about when to use actions
* **Tactical Timing Decisions** - Cooldowns should create timing choices beyond raw resource availability
* **Equipment Integration** - Actions must respect equipment requirements (ADR-0008)
* **Balance Flexibility** - Energy costs and cooldowns provide tuning knobs for game balance
* **Action Spam Prevention** - System must prevent trivial spamming of powerful actions
* **Gameplay Fluidity** - Basic actions (Move) should remain fluid and responsive

## Considered Options

* **Option 1: Cooldown-Only System** - Actions limited by time-based cooldowns only (no energy)
* **Option 2: Energy-Only System** - Actions limited by energy pool regeneration only (no cooldowns)
* **Option 3: Dual-Constraint System** - Both energy costs AND cooldowns required
* **Option 4: Action Points Per Turn** - Traditional turn-based action economy

## Decision Outcome

Chosen option: "**Option 3: Dual-Constraint System (energy costs AND cooldowns)**", because it creates both resource management decisions (energy) and tactical timing decisions (cooldowns), prevents action spam while maintaining fluidity for basic actions, enables diverse action costs for balance tuning, and integrates naturally with equipment system and real-time gameplay.

### Bot Actions and Resource Management Specification

#### Resource Management System

**Energy Pool** (from ADR-0007: Bot Characteristics):
- All actions consume energy from a limited pool
- Energy regenerates over time at a fixed rate (TBD)
- Insufficient energy prevents action execution
- Energy creates strategic resource management (choosing when to spend energy)
- Equipment can modify energy capacity (e.g., Boost Engine: -1 Energy Capacity)

**Cooldown System**:
- Actions have cooldown periods measured in game ticks
- Cooldowns prevent immediate re-use of actions
- Cooldowns create tactical timing decisions (when to use powerful abilities)
- Independent of energy (both must be satisfied to perform action)

**Dual Constraints**: Both energy availability AND cooldown status must be satisfied to perform actions:
- Basic actions (Move) have low energy cost and no cooldown for fluid movement
- Powerful actions (Shield, ShotgunBlast) have higher energy costs and longer cooldowns for balance
- This dual system prevents both energy-based spam and cooldown-based spam

#### Action Categories

Actions are organized into four categories based on their primary purpose:

**Movement Actions**: Navigate the 2D battle space (ADR-0006)

**Combat Actions**: Deal damage to opponent bots

**Defensive Actions**: Mitigate incoming damage or avoid attacks

**Utility Actions**: Provide information or modify bot state

#### Action Catalog

**Movement Actions**

*Move*: Moves the bot in a specified direction within the battle space

- **Energy Cost**: 5 (TBD)
- **Cooldown**: 0 ticks (TBD)
- **Equipment Required**: None (universal action)
- **Parameters**: Direction vector or target coordinates
- **Constraints**: Cannot move through obstacles or other bots (ADR-0006 collision rules)
- **Tactical Use**: Positioning, pursuit, evasion, range control

---

**Combat Actions**

*RifleShot*: Single-shot, precise ranged attack

- **Energy Cost**: 15 (TBD)
- **Cooldown**: 1 tick (TBD)
- **Equipment Required**: Rifle (ADR-0008)
- **Damage**: Moderate (TBD)
- **Range**: Long (TBD)
- **Constraints**: Requires line of sight (ADR-0006), rifle must be equipped
- **Tactical Use**: Consistent ranged damage, reliable baseline offense

*ShotgunBlast*: Spray of projectiles effective at close range

- **Energy Cost**: 20 (TBD)
- **Cooldown**: 2 ticks (TBD)
- **Equipment Required**: Shotgun (ADR-0008)
- **Damage**: High at close range with falloff based on distance (TBD)
- **Range**: Short to medium (TBD)
- **Constraints**: Shotgun must be equipped, damage decreases with distance
- **Tactical Use**: High burst damage at close range, requires positioning

---

**Defensive Actions**

*Block*: Reduces damage from incoming attacks

- **Energy Cost**: 10 (TBD)
- **Cooldown**: 2 ticks (TBD)
- **Equipment Required**: None (universal action)
- **Effect**: Reduces damage by percentage (TBD)
- **Duration**: Active for current tick only
- **Constraints**: Must be activated before damage is received
- **Tactical Use**: Mitigate predicted incoming damage, reduce burst damage

*Evade*: Attempts to dodge incoming attacks

- **Energy Cost**: 15 (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: None (universal action)
- **Effect**: Chance to completely avoid attack (TBD)
- **Duration**: Active for current tick only
- **Constraints**: Success rate may depend on bot Speed stat (ADR-0007)
- **Tactical Use**: Avoid high-damage attacks, risky but high-reward defense

*Shield*: Activates energy shield that absorbs damage over multiple ticks

- **Energy Cost**: 20 initial + ongoing drain (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: None (universal action)
- **Effect**: Absorbs damage up to threshold (TBD)
- **Duration**: Sustained (multiple ticks with ongoing energy drain - TBD)
- **Constraints**: Ongoing energy cost while active may limit other actions
- **Tactical Use**: Extended damage mitigation, control space through invulnerability window

---

**Utility Actions**

*Scan*: Gathers information about environment and nearby bots

- **Energy Cost**: 5 (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: Sensor Array module recommended (ADR-0008)
- **Effect**: Returns information about nearby entities (positions, health, etc.)
- **Range**: Limited detection radius, enhanced by Sensor Array (TBD)
- **Constraints**: May only reveal information within range
- **Tactical Use**: Information gathering, enemy tracking, battlefield awareness

*Charge*: Increases energy regeneration rate temporarily

- **Energy Cost**: Variable (may cost initial energy or pause regen - TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: None (universal action)
- **Effect**: Boosts energy regeneration for duration (TBD)
- **Duration**: Multiple ticks (TBD)
- **Constraints**: Bot may be vulnerable while charging (cannot perform other actions - TBD)
- **Tactical Use**: Energy economy management, prepare for energy-intensive action sequences

---

**Equipment-Dependent Actions** (from ADR-0008)

*Boost*: Temporary speed increase for repositioning

- **Energy Cost**: Variable (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: Boost Engine module
- **Effect**: Increases Speed temporarily (TBD)
- **Tactical Use**: Close distance for shotgun attacks, escape dangerous positions

*Repair*: Limited self-repair during combat

- **Energy Cost**: Variable (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: Repair Kit module
- **Effect**: Restores HP (TBD amount, usage limits)
- **Tactical Use**: Extend combat effectiveness, sustain tank builds

*Cloak*: Reduces detection range by enemies

- **Energy Cost**: Variable (TBD)
- **Cooldown**: Variable (TBD)
- **Equipment Required**: Stealth Module
- **Effect**: Reduces enemy detection range temporarily
- **Tactical Use**: Avoid detection, enable surprise attacks, stealth positioning

#### Action Execution and Processing

**Action Submission**:
- Bots submit actions via the bot-to-server communication protocol
- Actions are queued and processed during the next game tick
- Multiple actions may be submitted if resources permit (parallel actions TBD)

**Action Validation**:
- Server validates energy availability before execution
- Server validates cooldown status before execution
- Server validates equipment requirements before execution
- Invalid actions are rejected with error feedback to bot

**Action Resolution**:
- Actions are resolved during the game tick cycle
- Action outcomes are broadcast to all bots via the communication protocol
- State changes (damage, position, energy) are updated and synchronized

### Consequences

* Good, because dual-constraint system (energy + cooldowns) creates both resource management and tactical timing layers
* Good, because basic actions (Move) have low cost and no cooldown for fluid gameplay
* Good, because powerful actions (Shield, ShotgunBlast) have high costs and cooldowns preventing spam
* Good, because equipment-dependent actions integrate naturally with loadout system (ADR-0008)
* Good, because energy costs and cooldowns provide independent tuning knobs for balance
* Good, because action variety (Movement, Combat, Defensive, Utility) enables diverse playstyles
* Good, because universal actions (Move, Block, Evade) are always available regardless of equipment
* Good, because real-time action submission enables responsive gameplay
* Neutral, because all energy costs and cooldown values are TBD requiring extensive playtesting
* Neutral, because dual resource system (energy + cooldowns) is more complex than single-constraint systems
* Neutral, because energy regeneration rate requires careful tuning to avoid starvation or abundance
* Bad, because dual-constraint system adds complexity for bot developers to manage
* Bad, because energy management adds cognitive load beyond simple cooldown tracking
* Bad, because equipment-dependent actions require validation overhead in protocol
* Bad, because action spam prevention must be tuned carefully to avoid frustrating resource starvation

### Confirmation

The decision will be confirmed through:

1. Implementation of action system in game server with energy and cooldown tracking
2. Bot SDK exposing action methods with energy/cooldown visibility
3. Playtesting action economy to validate energy costs and cooldowns feel balanced
4. Action spam testing to confirm dual constraints prevent trivial ability spam
5. Equipment-action integration testing to validate equipment requirements work correctly
6. Balance tuning of energy costs, cooldowns, and regeneration rates through competitive gameplay
7. Protocol integration testing for action submission, validation, and feedback

## Pros and Cons of the Options

### Option 1: Cooldown-Only System

Actions limited by time-based cooldowns only, no energy resource.

* Good, because simpler than dual-constraint system
* Good, because eliminates energy tracking and management complexity
* Good, because easier for developers to understand (just wait for cooldown)
* Good, because predictable action availability based on time
* Neutral, because may be sufficient for simple action economies
* Bad, because lacks resource management strategic layer
* Bad, because no cost for spamming available actions (just wait for cooldowns)
* Bad, because difficult to balance rapid low-cooldown actions vs. slow high-impact actions
* Bad, because no energy scarcity creates no opportunity cost for action use
* Bad, because equipment cannot modify energy capacity for strategic effect

### Option 2: Energy-Only System

Actions limited by energy pool regeneration only, no cooldowns.

* Good, because creates resource management decisions
* Good, because energy scarcity forces strategic action choices
* Good, because simpler than dual-constraint system
* Good, because equipment can modify energy capacity for build diversity
* Neutral, because single resource may be sufficient for action economy
* Bad, because powerful actions can be spammed rapidly if energy is available
* Bad, because no timing constraints beyond energy availability
* Bad, because burst damage strategies may dominate (spend all energy immediately)
* Bad, because difficult to prevent action spam without introducing effective cooldowns anyway
* Bad, because regeneration rate becomes sole balancing mechanism

### Option 3: Dual-Constraint System (energy + cooldowns)

Both energy costs AND cooldowns required (CHOSEN).

* Good, because creates resource management (energy) and tactical timing (cooldowns)
* Good, because prevents action spam through dual constraints
* Good, because basic actions can have low cost and no cooldown for fluidity
* Good, because powerful actions can have high cost and long cooldown for balance
* Good, because provides two independent tuning knobs for balance adjustments
* Good, because equipment can modify energy capacity for strategic effect
* Good, because cooldowns prevent burst-spam even with high energy
* Good, because energy scarcity creates opportunity cost for action use
* Neutral, because requires careful tuning to avoid frustrating starvation or trivial abundance
* Neutral, because dual constraints require understanding both systems
* Bad, because more complex than single-constraint systems
* Bad, because adds cognitive load for bot developers to track both resources
* Bad, because requires tuning two systems instead of one

### Option 4: Action Points Per Turn

Traditional turn-based action point economy.

* Good, because familiar system from many turn-based games
* Good, because deterministic action availability per turn
* Good, because easy to balance (X actions per turn)
* Good, because clear constraints on action usage
* Neutral, because works well for turn-based gameplay
* Bad, because poor fit for real-time continuous gameplay model
* Bad, because turn boundaries create artificial rhythm contrary to real-time design
* Bad, because gRPC protocol designed for continuous action submission, not turns
* Bad, because eliminates timing decisions within turns
* Bad, because rigid action economy may feel restrictive in real-time battles

## More Information

### Related Documentation

- **[ADR-0007: Bot Characteristics System](0007-bot-characteristics-system.md)**: Energy pool characteristic that actions consume

- **[ADR-0008: Equipment and Loadout System](0008-equipment-loadout-system.md)**: Equipment that enables/disables actions and modifies energy

- **[ADR-0006: Battle Space Spatial System](0006-battle-space-spatial-system.md)**: Spatial environment where movement and combat actions occur

- **[Bot Actions Analysis](../analysis/game-mechanics/actions/)**: Detailed technical specifications for action catalog

- **[ADR-0004: Bot to Battle Server Communication Protocol](0004-bot-battle-server-interface.md)**: Communication protocol choice (gRPC) for bot-to-server interface

- **[ADR-0005: 1v1 Battle Orchestration](0005-1v1-battle-orchestration.md)**: High-level battle flow using these actions

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Energy economy modeling to determine regeneration rates and pool sizes
2. Playtesting action costs to validate energy expenditure feels balanced
3. Cooldown tuning to prevent action spam without creating frustrating delays
4. Balance analysis of powerful vs. basic actions (Shield vs. Move)
5. Equipment-action integration testing to validate equipment requirements
6. Competitive gameplay to identify dominant action patterns and adjust accordingly
7. Developer feedback on action system complexity and usability

**Key Design Insights**:
- Basic actions (Move: 5 energy, 0 cooldown) enable fluid gameplay
- Powerful actions (Shield: 20 energy, variable cooldown) are constrained to prevent spam
- Equipment-dependent actions create loadout-specific tactical options
- Universal actions ensure all bots have baseline capabilities

**Energy System Tuning Considerations**:
- Energy regeneration rate determines action frequency baseline
- Energy pool size determines burst action capability
- Energy costs should scale with action power (Move: 5, RifleShot: 15, Shield: 20)
- Equipment energy modifiers enable build optimization (Boost Engine: -1 capacity)

**Cooldown System Tuning Considerations**:
- Cooldowns prevent rapid-fire use of powerful abilities
- Cooldown duration should reflect action impact (Shield longer than RifleShot)
- Zero-cooldown actions (Move) enable continuous use if energy permits
- Cooldowns create timing windows for counter-play

**Future Enhancements**:
- Additional action types (area-of-effect, status effects, debuffs)
- Combo actions that trigger based on action sequences
- Interrupt actions that counter opponent actions
- Charged actions that increase in power with longer preparation
- Resource conversion actions (health for energy, etc.)

### Design Principles

The action and resource system follows these principles:
- **Dual Constraints**: Energy and cooldowns create layered decision-making
- **Fluidity for Basics**: Low-cost, no-cooldown basic actions maintain gameplay flow
- **Constraints for Power**: High-cost, long-cooldown powerful actions prevent spam
- **Equipment Integration**: Actions respect equipment requirements and modifications
- **Real-time Fit**: System designed for continuous gameplay, not turn-based
