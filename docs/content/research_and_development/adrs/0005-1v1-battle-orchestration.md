---
title: "[0005] 1v1 Battle Orchestration"
description: >
   High-level orchestration of 1v1 battles including visibility rules, battle pacing, and win conditions
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

Battle Bots requires high-level orchestration rules for 1v1 battles that tie together the spatial system, characteristics, equipment, and actions into a complete battle experience. We need to define visibility rules (fog of war vs. full battlefield vision), battle pacing mechanisms (time limits, engagement encouragement), and win conditions (how battles conclude). These orchestration rules must create engaging gameplay while integrating the detailed mechanics defined in ADR-0006 through ADR-0009.

Without well-defined battle orchestration, we cannot:
- Determine what information bots can see during battle (visibility)
- Prevent stalemate and ensure battles conclude in reasonable time (pacing)
- Definitively determine battle outcomes (win conditions)
- Balance information advantage vs. direct combat capability
- Create tension and encourage engagement between bots
- Design the battle flow from start to conclusion

## Decision Drivers

* **Information Strategy Depth** - Visibility rules should create meaningful choices (sensor investment vs. stealth)
* **Engagement Encouragement** - Battle pacing should prevent passive stalemates and force interaction
* **Definitive Outcomes** - Win conditions must handle all scenarios including edge cases
* **Fairness** - Resolution rules should be deterministic and predictable
* **Integration** - Must work seamlessly with mechanics from ADR-0006 through ADR-0009
* **Observability** - Battle outcomes must be clear for visualization and analysis
* **Time Constraints** - Battles must conclude in reasonable timeframes

## Considered Options

* **Option 1: Full Visibility + Fixed Time Limit** - Complete battlefield vision, 5-minute hard cutoff
* **Option 2: Fog of War + Dynamic Arena** - Limited detection ranges, shrinking play area over time
* **Option 3: Fog of War + Fixed Time Limit** - Limited detection ranges with sensor/stealth mechanics, fixed time limit

## Decision Outcome

Chosen option: "**Option 3: Fog of War + Fixed Time Limit**", because it creates strategic depth through information advantage (sensor/stealth equipment becomes valuable), encourages tactical positioning (must get in range to detect), maintains simplicity (no dynamic arena complexity), and provides definitive conclusion through time limits. This option maximizes the value of Sensor Array and Stealth Module equipment (ADR-0008).

### Battle Orchestration Specification

#### Battlefield Visibility Rules

**Fog of War System**:

Bots do not have automatic full battlefield visibility. Instead, detection is range-based:

**Detection Mechanics**:
- **Base Detection Range**: Each bot has a detection radius (TBD: placeholder 30 units)
- **Sensor Array Bonus**: +2 Detection Range (TBD) when equipped (ADR-0008)
- **Stealth Module Effect**: -2 Enemy Detection Range (TBD) when equipped (ADR-0008)
- **Line of Sight Required**: Detection requires unobstructed line of sight (ADR-0006)

**Information Provided Within Detection Range**:
- Enemy bot position (x, y coordinates)
- Enemy bot health (current HP value)
- Enemy bot orientation/facing (for predicting movement)
- Enemy actions (visible when performed within range)

**Information Hidden Outside Detection Range**:
- Enemy bot position unknown (not visible on battlefield)
- Enemy bot health unknown
- Enemy actions unknown

**Tactical Implications**:
- **Sensor Array Equipment**: Increases detection range, providing information advantage (ADR-0008)
- **Stealth Module Equipment**: Reduces enemy detection range, enables surprise attacks (ADR-0008)
- **Scan Action**: Provides temporary detection boost within area (ADR-0009)
- **Positioning Strategy**: Bots must position to detect enemies or avoid detection
- **Equipment Tradeoffs**: Sensor/Stealth builds become viable strategic choices

**Alternative Considered**: Full Visibility (all bots always visible) would simplify implementation but eliminate information strategy layer and reduce value of sensor/stealth equipment.

#### Battle Pacing and Time Management

**Time Limit System**:

All 1v1 battles have a maximum duration to prevent indefinitely long battles and ensure matches conclude:

**Time Limit**: TBD (placeholder: 5 minutes / 300 seconds)

**Time Limit Enforcement**:
1. Battle engine tracks elapsed time from battle start
2. When time limit is reached, battle engine stops accepting new actions
3. All pending actions in the current tick are resolved
4. Final health values are compared to determine outcome
5. Battle concludes with timeout victory/defeat/draw

**Timeout Resolution**:
- **Higher Health Wins**: Bot with more HP at timeout achieves timeout victory
- **Equal Health Draws**: Bots with equal HP at timeout result in draw
- **Encourages Aggression**: Passive play risks timeout loss if opponent deals any damage

**Design Considerations**:
- Time limit should be long enough to allow for strategic gameplay and tactical maneuvering
- Time limit should be short enough to maintain engagement and prevent stalemate
- Time limit may need to be configurable for different battle modes or tournaments
- Future consideration: Shrinking arena over time to force engagement (not in current design)

**Engagement Encouragement Mechanisms**:
- **Timeout Favors Aggressor**: Dealing damage creates HP advantage for timeout scenario
- **No Passive Victory**: Cannot win by hiding; must either destroy enemy or have HP advantage
- **Fog of War**: Limited visibility encourages active scouting and engagement
- **Energy Regeneration**: Bots have resources to maintain active gameplay (ADR-0009)

#### Win Conditions

Battle resolution defines how a 1v1 battle concludes and determines the winner. Every battle must have a definitive outcome: victory, defeat, or draw.

**Victory Conditions** - A bot achieves victory when any of the following occurs:

*Destruction*: Enemy bot's health reaches 0 or below

- **Trigger**: `enemy.health <= 0`
- **Outcome**: Immediate victory for the surviving bot
- **Primary Condition**: Most definitive victory type

*Forfeit*: Enemy bot disconnects or fails to respond

- **Trigger**: Enemy bot connection lost or timeout exceeded
- **Grace Period**: TBD (brief window allows reconnection before forfeit declared)
- **Outcome**: Immediate victory for the connected bot

*Timeout*: Battle time limit reached with higher health

- **Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health > enemy.health`
- **Outcome**: Victory by health advantage
- **Encourages Aggression**: Rewards damage dealt over passive play

**Defeat Conditions** - A bot is defeated when any of the following occurs:

*Destruction*: Bot's own health reaches 0 or below

- **Trigger**: `bot.health <= 0`
- **Outcome**: Immediate defeat

*Forfeit*: Bot disconnects or fails to respond

- **Trigger**: Bot connection lost or timeout exceeded
- **Outcome**: Immediate defeat

*Timeout*: Battle time limit reached with lower health

- **Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health < enemy.health`
- **Outcome**: Defeat (opponent wins by health advantage)

**Draw Conditions** - A battle results in a draw when neither bot achieves a clear victory:

*Mutual Destruction*: Both bots reach 0 health simultaneously

- **Trigger**: `bot.health <= 0 && enemy.health <= 0` (same tick)
- **Scenarios**: Both deal fatal damage in same tick, simultaneous environmental effects
- **Outcome**: Draw

*Equal Health at Timeout*: Time limit reached with identical health

- **Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health == enemy.health`
- **Outcome**: Draw

#### Edge Cases and Resolution Rules

**Simultaneous Actions**:

When both bots perform actions in the same game tick that affect battle outcome:

1. All actions for the current tick are collected
2. Actions are resolved in a deterministic order (e.g., by bot ID, action type priority)
3. Game state is updated after all actions are processed
4. Win conditions are checked after state update

Example: If both bots fire projectiles that would destroy each other in the same tick, both hits are processed, and if both bots reach 0 health, the result is a draw (mutual destruction).

**Disconnect Handling**:

When a bot disconnects or becomes unresponsive:

- **Grace Period**: TBD timeout window allows bot to reconnect/respond
- **Reconnection**: If bot reconnects within grace period, battle continues from current state
- **No Response**: If grace period expires, battle immediately ends with forfeit victory for opponent
- **Both Disconnect**: If both bots disconnect, outcome determined by who reconnects first within grace period; if neither reconnects, battle ends as draw

**Timeout and Destruction Priority**:

If time limit is reached in the same tick that a bot is destroyed:

- **Priority**: Destruction takes precedence over timeout
- **Rationale**: Destruction is more definitive and should be the primary victory condition
- **Example**: Time limit expires at tick 1000, bot destroyed in same tick → victory by destruction (not timeout)

**Negative Health Values**:

Bots can temporarily have negative health if damage exceeds remaining health:

- **Example**: Bot has 5 health, takes 20 damage → health becomes -15
- **Win Condition Check**: Any health value <= 0 triggers destruction
- **Display**: For visualization, negative health displayed as 0

**Zero Health at Battle Start**:

If a bot begins battle with 0 or negative health (configuration error):

- **Behavior**: Battle immediately ends with defeat for the bot with invalid health
- **Prevention**: Battle engine should validate initial state before battle begins

### Consequences

* Good, because fog of war creates strategic depth through information advantage and detection mechanics
* Good, because sensor/stealth equipment becomes meaningfully valuable (not just stat modifiers)
* Good, because time limits ensure battles conclude in reasonable timeframes
* Good, because timeout resolution encourages aggressive play (damage creates HP advantage)
* Good, because win conditions handle all scenarios including edge cases definitively
* Good, because deterministic resolution rules ensure fairness and predictability
* Good, because integrates seamlessly with spatial system (ADR-0006), characteristics (ADR-0007), equipment (ADR-0008), and actions (ADR-0009)
* Good, because draw conditions are clear and handle simultaneous outcomes
* Good, because disconnect handling with grace period balances technical issues with competitive integrity
* Neutral, because detection range values (base 30 units, +2 sensor, -2 stealth) require playtesting
* Neutral, because time limit (5 minutes) requires validation through competitive gameplay
* Neutral, because fog of war adds complexity vs. full visibility but creates strategic value
* Bad, because fog of war implementation is more complex than full visibility
* Bad, because detection range mechanics require careful tuning to balance sensor/stealth value
* Bad, because time limit may feel artificial if tuned incorrectly (too short or too long)
* Bad, because timeout resolution may favor defensive play if time limit is too generous

### Confirmation

The decision will be confirmed through:

1. Implementation of fog of war system with detection ranges and line of sight
2. Equipment integration testing to validate sensor/stealth effects on detection
3. Playtesting to tune time limits and ensure engagement without stalemate
4. Win condition testing to verify all scenarios (destruction, forfeit, timeout, draw) resolve correctly
5. Edge case validation for simultaneous actions, disconnects, and timeout priorities
6. Competitive gameplay to ensure sensor/stealth builds are viable strategic choices
7. Timeout scenario analysis to confirm aggressive play is rewarded over passive stalemate

## Pros and Cons of the Options

### Option 1: Full Visibility + Fixed Time Limit

Complete battlefield vision, all bots always visible, 5-minute hard cutoff.

* Good, because simplest implementation (no detection mechanics)
* Good, because eliminates fog of war complexity
* Good, because bots always have complete information for decision-making
* Good, because easier to debug and visualize (everything always visible)
* Neutral, because may be sufficient for initial POC
* Bad, because eliminates information strategy layer
* Bad, because Sensor Array and Stealth Module equipment have no visibility effect (only stat modifiers)
* Bad, because no reward for scouting or positioning for detection
* Bad, because reduces strategic depth (no information advantage mechanic)

### Option 2: Fog of War + Dynamic Arena

Limited detection ranges with sensor/stealth mechanics, shrinking play area forces engagement.

* Good, because fog of war creates information strategy depth
* Good, because dynamic arena (shrinking boundaries) guarantees eventual engagement
* Good, because eliminates timeout stalemate possibility (arena forces contact)
* Good, because sensor/stealth equipment becomes highly valuable
* Neutral, because dynamic arena adds urgency and tension
* Bad, because dynamic arena significantly increases implementation complexity
* Bad, because shrinking boundaries may feel artificial or gimmicky
* Bad, because requires careful tuning of shrink rate to avoid frustration
* Bad, because arena size changes affect all spatial calculations over time
* Bad, because may disadvantage slower bots unfairly (cannot escape shrinking boundary)

### Option 3: Fog of War + Fixed Time Limit

Limited detection ranges with sensor/stealth mechanics, fixed 5-minute time limit (CHOSEN).

* Good, because fog of war creates information strategy depth
* Good, because sensor/stealth equipment becomes meaningfully valuable
* Good, because simpler than dynamic arena (static boundaries from ADR-0006)
* Good, because detection mechanics encourage tactical positioning and scouting
* Good, because time limit ensures definitive conclusion without arena complexity
* Good, because timeout resolution (HP comparison) encourages aggressive play
* Neutral, because requires tuning detection ranges for balance
* Neutral, because time limit needs validation through playtesting
* Bad, because fog of war implementation more complex than full visibility
* Bad, because passive play is technically possible until timeout (relies on timeout resolution to discourage)
* Bad, because detection range mechanics add protocol and server complexity

## More Information

### Related Documentation

- **[ADR-0006: Battle Space Spatial System](0006-battle-space-spatial-system.md)**: 2D arena where battles occur, line of sight calculations for detection

- **[ADR-0007: Bot Characteristics System](0007-bot-characteristics-system.md)**: Health stat that determines timeout and destruction outcomes

- **[ADR-0008: Equipment and Loadout System](0008-equipment-loadout-system.md)**: Sensor Array (+2 detection) and Stealth Module (-2 enemy detection) equipment

- **[ADR-0009: Bot Actions and Resource Management](0009-bot-actions-resource-management.md)**: Scan action for temporary detection boost, actions that occur within battles

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for battle state synchronization and visibility information

- **[Game Mechanics Analysis](../analysis/game-mechanics/)**: Detailed framework documentation including:
  - [Battle Space](../analysis/game-mechanics/battle-space/): Spatial system details
  - [Bot Characteristics](../analysis/game-mechanics/characteristics/): Stat system details
  - [Equipment System](../analysis/game-mechanics/equipment/): Loadout mechanics
  - [Bot Actions](../analysis/game-mechanics/actions/): Action catalog
  - [Win Conditions](../analysis/game-mechanics/win-conditions/): Battle resolution details

- **[POC User Journey](../user-journeys/0001-poc.md)**: Proof of concept implementation using these battle orchestration rules

### Implementation Notes

All numeric values in this ADR are marked TBD (To Be Determined) and serve as placeholder values to establish the framework structure. These values will be refined through:

1. Detection range playtesting to balance sensor/stealth value (base: 30 units, sensor: +2, stealth: -2)
2. Time limit tuning to ensure engagement without stalemate (placeholder: 5 minutes)
3. Grace period validation for disconnect handling (balance technical issues vs. competitive integrity)
4. Win condition edge case testing (simultaneous destruction, timeout priorities, etc.)
5. Fog of war implementation testing to validate detection mechanics work correctly
6. Competitive gameplay analysis to ensure sensor/stealth builds are viable
7. Timeout scenario frequency analysis to tune time limit appropriately

**Key Design Insights**:
- Fog of war makes Sensor Array and Stealth Module equipment strategically valuable beyond stat effects
- Time limit with HP-based timeout resolution encourages aggressive play over passive stalemate
- Destruction takes precedence over timeout to prioritize definitive outcomes
- Grace period for disconnects balances technical issues with competitive fairness

**Future Considerations**:
- **Variable Time Limits**: Different battle modes (quick match vs. tournament) may have different time limits
- **Shrinking Arena**: Dynamic boundaries could be added for specific game modes if passive play becomes problematic
- **Overtime Mechanics**: If timeout occurs with close HP values, brief overtime period could be added
- **Spectator Mode**: Full visibility for spectators even with fog of war for participants
- **Replay System**: Battle recordings with fog of war visualization

### Design Principles

The battle orchestration follows these principles:
- **Information Strategy**: Fog of war creates meaningful sensor/stealth choices
- **Engagement Encouragement**: Timeout resolution rewards aggressive play
- **Definitive Outcomes**: All scenarios have clear win/loss/draw determination
- **Fairness**: Deterministic, predictable resolution rules
- **Integration**: Seamlessly combines mechanics from ADR-0006 through ADR-0009
