---
title: "[0005] 1v1 Battle Game Mechanics"
description: >
    Core game mechanics definition for 1v1 battles including bot actions, characteristics, equipment, and win conditions
type: docs
weight: 5
category: "strategic"
status: "proposed"
date: 2025-12-05
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

Battle Bots requires a clear definition of the core game mechanics for 1v1 battles. We need to specify what actions bots can perform, what attributes they have, how equipment affects gameplay, and how battles conclude. These mechanics must support the real-time continuous gameplay model while providing enough complexity for strategic depth.

Without well-defined game mechanics, we cannot:
- Design the game server implementation
- Create bot SDK interfaces
- Define the protocol for bot-server communication beyond basic structure
- Develop the battle visualization system
- Begin playtesting and balance iteration

## Decision Drivers

* **Bot Developer Experience** - Mechanics must be clear and understandable for developers implementing autonomous bots
* **Game Balance** - Multiple viable strategies should exist with no single dominant approach
* **Implementation Feasibility** - Mechanics must be implementable in a real-time continuous gameplay system
* **Observability** - Game events must be trackable for visualization and debugging
* **Customization Depth** - Moderate complexity that allows personalization without overwhelming new developers
* **Protocol Integration** - Mechanics must map cleanly to the gRPC bidirectional streaming interface (ADR-0004)

## Considered Options

* **Option 1: Simple Arcade-Style** - Move + shoot only, minimal stats (health, speed), no equipment customization
* **Option 2: Moderate with Customization** - 4 core stats (with equipment-derived Mass), equipment loadouts, varied action catalog
* **Option 3: Deep Simulation-Style** - 10+ stats, complex resource management, many abilities and systems

## Decision Outcome

Chosen option: "**Option 2: Moderate with Customization**", because it provides strategic depth and customization options while remaining accessible to bot developers. This approach balances complexity with clarity and supports diverse playstyles.

### Core Mechanics Framework

**Bot Characteristics** (4 core attributes):
- **Health**: HP pool (100-500 TBD), determines survivability, bot destroyed at 0
- **Speed**: Movement rate (1-10 units/tick TBD), affects positioning and tactical options
- **Defense**: Damage reduction (1-10 TBD), mitigates incoming damage
- **Mass**: Derived from equipped items, reduces effective Speed (equipment weight impacts mobility)

**Bot Actions** (organized by category):
- **Movement**: Move (0 cooldown), Rotate (0 cooldown), Dash (3 tick cooldown) - all TBD
- **Combat**: BasicAttack (1 tick cooldown), LaserShot (2 tick cooldown, requires Laser), HeavyBlow (4 tick cooldown) - all TBD
- **Defensive**: Block (2 tick cooldown), Evade (3 tick cooldown), Shield (5 tick cooldown, duration-based) - all TBD
- **Utility**: Scan (2 tick cooldown), RepairKit (10 tick cooldown, requires Repair Kit module) - all TBD

**Equipment System** (3 types in fixed slots):
- **Weapons**: Enable combat actions (Laser, Missile, Plasma, EMP), each has Mass value affecting mobility
- **Armor**: Provide Defense bonuses (Light, Medium, Heavy, Reactive), heavier armor increases Mass
- **Modules**: Enable utility functions (Shield Generator, Boost Engine, Repair Kit, Sensors, Stealth), each contributes Mass
- **Loadout**: 1 weapon slot, 1 armor slot, 2 module slots (all TBD)
- **Mass Calculation**: Total Mass = Weapon Mass + Armor Mass + Module Mass (sum), affects effective Speed

**Battle Space**:
- 2D rectangular arena with Cartesian coordinates (100x100 units TBD)
- Center origin (0,0), continuous coordinate system
- Bot radius 2 units (TBD), circle-based collision detection
- Line of sight for ranged attacks, movement boundaries

**Win Conditions**:
- **Victory**: Destruction (enemy at 0 health), Forfeit (enemy disconnect), Timeout (higher health at time limit)
- **Defeat**: Destruction (bot at 0 health), Forfeit (bot disconnect), Timeout (lower health)
- **Draw**: Mutual destruction, equal health at timeout
- **Time Limit**: 5 minutes (TBD)

**Cooldown System**:
- Actions have cooldown periods measured in game ticks to prevent spam
- Cooldowns create tactical timing decisions (when to use powerful abilities)
- Basic actions (Move, Rotate) have no cooldown for fluid movement
- Powerful actions (Dash, HeavyBlow, Shield) have longer cooldowns for balance

### Consequences

* Good, because the framework provides clear structure for bot developers to understand what their bots can do
* Good, because equipment customization enables diverse strategies (DPS, Tank, Balanced, Utility builds)
* Good, because real-time continuous gameplay integrates naturally with gRPC bidirectional streaming (ADR-0004)
* Good, because moderate complexity is accessible while still supporting strategic depth
* Good, because all game events (actions, damage, state changes) are observable for visualization and debugging
* Good, because cooldown-only system simplifies bot logic (no resource pool management)
* Good, because Mass as derived stat creates natural tradeoffs between equipment power and mobility
* Neutral, because placeholder values require extensive playtesting and balance iteration before finalization
* Neutral, because the 4-stat system with equipment-based Mass is more complex than arcade-style but simpler than deep simulation
* Neutral, because cooldown-only system may allow action spam in some scenarios (requires careful cooldown tuning)
* Bad, because moderate complexity has a higher learning curve than simple arcade-style mechanics
* Bad, because equipment dependencies create validation requirements in the protocol and game server

### Confirmation

The decision will be confirmed through:
1. Implementation of game server mechanics following this framework
2. Creation of bot SDK that exposes these actions and characteristics
3. Playtesting with sample bots demonstrating different equipment loadouts
4. Balance iteration and adjustment of TBD placeholder values
5. Successful integration with the gRPC protocol defined in ADR-0004

## Pros and Cons of the Options

### Option 1: Simple Arcade-Style

Move + shoot only, minimal stats (health, speed), no equipment customization.

* Good, because extremely simple to understand and implement
* Good, because low learning curve for bot developers
* Good, because fast iteration on game balance with fewer variables
* Neutral, because limited complexity may be sufficient for initial POC
* Bad, because lack of customization limits strategic diversity
* Bad, because shallow gameplay may not sustain long-term engagement
* Bad, because limited differentiation between bots beyond raw implementation skill

### Option 2: Moderate with Customization

4 core stats (with equipment-derived Mass), equipment loadouts, varied action catalog (CHOSEN).

* Good, because strategic depth through stat allocation and equipment choices
* Good, because multiple viable builds and playstyles (DPS, Tank, Utility, Balanced)
* Good, because customization creates interesting pre-battle decisions
* Good, because complexity is manageable and learnable
* Good, because equipment enables action variety without protocol bloat
* Neutral, because requires balance tuning but not excessively complex
* Neutral, because 4 stats hit a sweet spot between simple and overwhelming
* Bad, because more complex than arcade-style to implement and explain
* Bad, because equipment validation adds implementation overhead

### Option 3: Deep Simulation-Style

10+ stats, complex resource management, many abilities and systems.

* Good, because maximum strategic depth and optimization potential
* Good, because extensive customization options and build diversity
* Good, because appeals to players who enjoy complexity and optimization
* Neutral, because may enable very deep competitive meta
* Bad, because high learning curve may deter casual bot developers
* Bad, because complexity makes game balance extremely difficult
* Bad, because numerous stats and abilities create protocol bloat
* Bad, because difficult to visualize and understand battle outcomes
* Bad, because implementation complexity significantly higher

## More Information

### Related Documentation

- **[Game Mechanics Analysis](../analysis/game-mechanics/)**: Detailed framework documentation including:
  - [Battle Space](../analysis/game-mechanics/battle-space/): 2D arena definition
  - [Bot Characteristics](../analysis/game-mechanics/characteristics/): 4-stat system details
  - [Equipment System](../analysis/game-mechanics/equipment/): Weapons, armor, modules
  - [Bot Actions](../analysis/game-mechanics/actions/): Complete action catalog
  - [Win Conditions](../analysis/game-mechanics/win-conditions/): Battle resolution rules

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol that these actions map to

- **[POC User Journey](../user-journeys/0001-poc.md)**: Proof of concept implementation using these mechanics

- **ADR-NNNN: Game Runtime Architecture** (Pending): Future ADR covering tick rate, game loop implementation, and state processing details

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Mathematical modeling and simulation
2. Playtesting with real bot implementations
3. Balance analysis and competitive meta observation
4. Performance testing and optimization requirements

The framework separates **WHAT** mechanics exist (this ADR) from **HOW** they are implemented (future Game Runtime Architecture ADR). This allows independent iteration on game balance and implementation details.

### Design Principles

The mechanics follow these principles:
- **Clarity over Complexity**: Prefer understandable mechanics to clever systems
- **Tradeoffs over Power**: Equipment choices involve costs and benefits
- **Diversity over Dominance**: Multiple strategies should be viable
- **Observable over Opaque**: All game state changes should be trackable
- **Extensible over Final**: Framework supports future additions without redesign
