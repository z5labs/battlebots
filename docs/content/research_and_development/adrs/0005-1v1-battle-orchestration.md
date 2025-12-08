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

Battle orchestration is decomposed into independent properties, each with distinct options:

### Property 1: Visibility

* **Option 1.1: Full Visibility** - All bots can see everything on the battlefield at all times
* **Option 1.2: Constant Fog of War** - Each bot can only see what's within a limited radius around themselves (similar to arcade-style games)
* **Option 1.3: Revealed Fog of War** - Each bot can see within a radius around themselves, but as they move they reveal the map and can see any items/enemies within the revealed area (similar to Age of Empires and Civilization)

### Property 2: Battle Termination

* **Option 2.1: Health-Only Termination** - Battle ends when `enemy.health <= 0`
* **Option 2.2: Health or Timeout Termination** - Battle ends when `enemy.health <= 0` OR max time reached (prevents indefinite stalemates)
* **Option 2.3: Health, Timeout, or Disconnect Termination** - Battle ends when `enemy.health <= 0` OR max time reached OR bot disconnected (forfeit victory)

## Decision Outcome

### Property 1: Visibility

**Chosen: Option 1.1 - Full Visibility**

Rationale:
- Simpler implementation for initial POC - no detection mechanics, line-of-sight calculations, or visibility protocol complexity
- Encourages users to implement full pathfinding solutions - with complete battlefield information, bots can calculate optimal paths considering enemy positions and obstacles
- Enables sophisticated AI strategies - bots have complete information to make informed tactical decisions
- Easier to debug and visualize - all battlefield state always visible
- Reduces initial implementation scope - allows focus on core battle mechanics before adding fog of war complexity
- Future enhancement path - fog of war can be added in later iterations once core mechanics are validated

**Alternative Considered: Option 1.2 - Constant Fog of War** would create strategic depth through information advantage but adds significant implementation complexity (detection ranges, line-of-sight, sensor/stealth mechanics) that is not necessary for initial POC. Can be reconsidered for future game modes once core mechanics are proven.

**Alternative Rejected: Option 1.3 - Revealed Fog of War** adds even more complexity (persistent revealed state tracking) and is not appropriate for initial implementation.

### Property 2: Battle Termination

**Chosen: Option 2.3 - Health, Timeout, or Disconnect Termination**

Rationale:
- **Health termination** (`enemy.health <= 0`) provides the primary definitive victory condition
- **Timeout termination** (max time reached with HP comparison) prevents indefinite stalemates from equally matched bots or passive/buggy strategies, encourages aggressive play by rewarding damage dealt
- **Disconnect termination** (forfeit on disconnect/unresponsive) handles technical failures gracefully with grace period for reconnection, ensures battles conclude even with connection issues
- All three termination conditions are necessary to handle the complete set of battle scenarios

**Time Limit Value**: TBD (placeholder: 5 minutes / 300 seconds) - requires playtesting to balance engagement encouragement vs. strategic gameplay depth

**Grace Period Value**: TBD - requires testing to balance reconnection allowance vs. competitive integrity

**Timeout Resolution**: Bot with higher health at timeout wins; equal health results in draw. This rewards aggressive play and damage dealing over passive stalemate.

### Battle Orchestration Specification

#### Battlefield Visibility Rules

**Full Visibility System**:

All bots have complete battlefield visibility at all times. Each bot receives full information about all entities and state on the battlefield.

**Information Provided to All Bots**:
- **Enemy Bot Position**: Complete x, y coordinates of all enemy bots
- **Enemy Bot Health**: Current HP value of all enemy bots
- **Enemy Bot Orientation**: Facing direction of all enemy bots (for predicting movement)
- **Enemy Bot Equipment**: Loadout information for all enemy bots (visible equipment from ADR-0008)
- **Battlefield Obstacles**: All obstacle positions and dimensions (from ADR-0006)
- **Arena Boundaries**: Complete arena size and boundary information (from ADR-0006)
- **Enemy Actions**: All actions performed by enemy bots are visible when executed

**No Hidden Information**:
- All battlefield state is visible to all bots at all times
- No detection ranges or line-of-sight restrictions
- No visibility-based equipment effects (Sensor Array and Stealth Module provide stat bonuses only, not visibility effects)

**Tactical Implications**:
- **Complete Information Gameplay**: Bots can make fully informed decisions based on complete battlefield state
- **Pathfinding Opportunities**: Bots can calculate optimal paths considering enemy positions, obstacles, and arena boundaries
- **Predictive AI**: Bots can implement sophisticated prediction algorithms knowing enemy positions and past actions
- **Strategic Positioning**: Position-based tactics still valuable (flanking, distance management, obstacle usage)
- **Equipment Focus**: Equipment choices focus on combat stats rather than detection/stealth capabilities

**Design Rationale**: Full visibility simplifies initial implementation, enables sophisticated pathfinding and AI strategies, and allows focus on core battle mechanics. Fog of war mechanics can be added in future iterations once core gameplay is validated.

**Future Consideration**: Fog of War (constant or revealed) could be added as optional game modes once core mechanics are proven, providing information strategy depth and making sensor/stealth equipment affect visibility.

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

**Visibility Decision (Full Visibility)**:

* Good, because simplest possible implementation (no detection mechanics, line-of-sight, or visibility protocol)
* Good, because encourages users to implement sophisticated pathfinding solutions with complete battlefield information
* Good, because enables advanced AI strategies - bots have complete information for tactical decisions
* Good, because easier to debug and visualize (all state always visible)
* Good, because reduces initial POC scope - allows focus on core battle mechanics
* Good, because integrates seamlessly with spatial system (ADR-0006) without additional visibility calculations
* Good, because provides clear future enhancement path - fog of war can be added later
* Neutral, because sensor/stealth equipment provides stat bonuses only (no visibility effects in this mode)
* Neutral, because eliminates information strategy layer - all gameplay is complete information
* Bad, because may reduce strategic depth compared to fog of war (no information advantage mechanic)
* Bad, because sensor/stealth builds less distinctive (equipment affects stats but not visibility)

**Battle Termination Decision (Health, Timeout, or Disconnect)**:

* Good, because handles all battle termination scenarios comprehensively
* Good, because timeout ensures battles conclude in reasonable timeframes (prevents infinite stalemates)
* Good, because timeout resolution encourages aggressive play (damage creates HP advantage)
* Good, because disconnect handling with grace period balances technical issues with competitive integrity
* Good, because provides definitive outcomes for all scenarios including edge cases
* Good, because deterministic resolution rules ensure fairness and predictability
* Good, because draw conditions are clear (mutual destruction, equal health at timeout)
* Neutral, because time limit (5 minutes) requires validation through competitive gameplay
* Neutral, because grace period duration requires tuning (balance fairness vs. delays)
* Neutral, because timeout may feel artificial if tuned incorrectly (too short or too long)
* Bad, because most complex termination implementation (three conditions to track and handle)
* Bad, because requires robust disconnect detection and reconnection infrastructure
* Bad, because timeout resolution may favor defensive play if time limit is too generous

**Overall Integration**:

* Good, because decisions integrate seamlessly with ADR-0006 (spatial system), ADR-0007 (characteristics), ADR-0008 (equipment), and ADR-0009 (actions)
* Good, because property-based decision structure allows independent tuning and future modifications
* Good, because each property addresses distinct battle orchestration concerns

### Confirmation

The decision will be confirmed through:

1. Implementation of full visibility system with complete battlefield state synchronization
2. Validation that all bots receive complete and accurate battlefield information
3. Playtesting to tune time limits and ensure engagement without stalemate
4. Win condition testing to verify all scenarios (destruction, forfeit, timeout, draw) resolve correctly
5. Edge case validation for simultaneous actions, disconnects, and timeout priorities
6. User testing to confirm pathfinding and AI strategy opportunities with complete information
7. Timeout scenario analysis to confirm aggressive play is rewarded over passive stalemate
8. Future evaluation of fog of war as optional game mode once core mechanics are validated

## Pros and Cons of the Options

### Property 1: Visibility

#### Option 1.1: Full Visibility (CHOSEN)

All bots can see everything on the battlefield at all times.

* Good, because simplest implementation (no detection mechanics required)
* Good, because eliminates fog of war complexity
* Good, because bots always have complete information for decision-making
* Good, because easier to debug and visualize (everything always visible)
* Good, because no line-of-sight calculations needed
* Neutral, because may be sufficient for initial POC or simple game modes
* Bad, because eliminates information strategy layer entirely
* Bad, because Sensor Array and Stealth Module equipment have no visibility effect (reduced to stat modifiers only)
* Bad, because no reward for scouting, positioning, or detection tactics
* Bad, because reduces strategic depth (no information advantage mechanic)
* Bad, because all builds become combat-focused with no sensor/stealth viability

#### Option 1.2: Constant Fog of War

Each bot can only see what's within a limited radius around themselves.

* Good, because creates information strategy depth through limited visibility
* Good, because sensor/stealth equipment becomes meaningfully valuable (affects detection ranges)
* Good, because encourages tactical positioning and maneuvering
* Good, because enables viable sensor/stealth builds as strategic alternatives
* Good, because balances complexity vs. value - simpler than revealed fog of war
* Good, because detection mechanics add depth without persistent state tracking
* Neutral, because requires tuning detection ranges for balance
* Neutral, because requires line-of-sight calculations (already needed for ADR-0006)
* Bad, because more complex implementation than full visibility
* Bad, because detection range mechanics add protocol and server complexity
* Bad, because requires careful balancing of base detection, sensor bonus, and stealth penalty

#### Option 1.3: Revealed Fog of War

Bots reveal map areas as they move, maintaining visibility of revealed areas.

* Good, because creates strong information strategy depth
* Good, because rewards exploration and map control
* Good, because enables "scouting" strategies and territorial gameplay
* Good, because familiar mechanic from RTS games (Age of Empires, Civilization)
* Neutral, because sensor/stealth equipment still valuable for initial detection
* Bad, because significant implementation complexity (persistent revealed state per bot)
* Bad, because requires tracking revealed areas across battle duration
* Bad, because arena state synchronization becomes more complex
* Bad, because memory overhead for revealed area maps
* Bad, because may favor mobility-focused builds disproportionately
* Bad, because complexity may not justify strategic benefit for 1v1 battles

### Property 2: Battle Termination

#### Option 2.1: Health-Only Termination

Battle ends only when `enemy.health <= 0`.

* Good, because simplest termination logic (single victory condition)
* Good, because most intuitive outcome (destroy enemy to win)
* Good, because no timeout mechanics needed
* Neutral, because appropriate for games where engagement is guaranteed
* Bad, because battles could run indefinitely with equally matched bots
* Bad, because passive/defensive strategies could create stalemates
* Bad, because buggy bot logic could cause infinite battles
* Bad, because no mechanism to conclude matches that reach equilibrium
* Bad, because requires external intervention to end problematic matches

#### Option 2.2: Health or Timeout Termination

Battle ends when `enemy.health <= 0` OR max time reached.

* Good, because ensures battles conclude in reasonable timeframes
* Good, because prevents stalemates from equally matched bots
* Good, because handles buggy/passive bot strategies gracefully
* Good, because timeout resolution (HP comparison) encourages aggressive play
* Good, because rewards damage dealing over pure defense
* Neutral, because requires tuning time limit for balance
* Neutral, because timeout may feel artificial if tuned incorrectly
* Bad, because doesn't handle disconnection/unresponsive bots
* Bad, because disconnect scenarios require separate handling or manual intervention
* Bad, because incomplete solution for all battle termination scenarios

#### Option 2.3: Health, Timeout, or Disconnect Termination (CHOSEN)

Battle ends when `enemy.health <= 0` OR max time reached OR bot disconnected.

* Good, because handles all battle termination scenarios comprehensively
* Good, because timeout prevents indefinite stalemates
* Good, because disconnect handling ensures battles conclude gracefully
* Good, because grace period for reconnection balances fairness with technical issues
* Good, because rewards aggressive play through timeout HP comparison
* Good, because provides definitive outcomes for all edge cases
* Neutral, because requires tuning both time limit and grace period
* Neutral, because disconnect logic adds complexity to battle flow
* Bad, because most complex termination implementation (three conditions to track)
* Bad, because grace period tuning difficult (too short = unfair, too long = delays)
* Bad, because requires robust disconnect detection and handling infrastructure

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

1. Time limit tuning to ensure engagement without stalemate (placeholder: 5 minutes)
2. Grace period validation for disconnect handling (balance technical issues vs. competitive integrity)
3. Win condition edge case testing (simultaneous destruction, timeout priorities, etc.)
4. Full visibility state synchronization testing to validate complete battlefield information delivery
5. Pathfinding and AI strategy validation to confirm complete information enables sophisticated bot implementations
6. Timeout scenario frequency analysis to tune time limit appropriately
7. Equipment balance testing to ensure stat-based equipment choices remain meaningful

**Key Design Insights**:
- Full visibility simplifies initial implementation and focuses on core battle mechanics
- Complete battlefield information enables users to implement sophisticated pathfinding and AI solutions
- Time limit with HP-based timeout resolution encourages aggressive play over passive stalemate
- Destruction takes precedence over timeout to prioritize definitive outcomes
- Grace period for disconnects balances technical issues with competitive fairness
- Sensor Array and Stealth Module equipment provide stat bonuses without visibility effects in this mode

**Future Considerations**:
- **Variable Time Limits**: Different battle modes (quick match vs. tournament) may have different time limits
- **Fog of War Game Modes**: Constant or revealed fog of war could be added as optional game modes to create information strategy layer
- **Shrinking Arena**: Dynamic boundaries could be added for specific game modes if passive play becomes problematic
- **Overtime Mechanics**: If timeout occurs with close HP values, brief overtime period could be added
- **Spectator Mode**: Already full visibility for all participants; spectators share same view
- **Replay System**: Battle recordings with complete visibility throughout
- **Partial Information Modes**: Once core mechanics proven, visibility could become configurable property per game mode

### Design Principles

The battle orchestration follows these principles:
- **Simplicity First**: Full visibility reduces initial complexity and enables focus on core mechanics
- **Complete Information Strategy**: Bots have full battlefield knowledge to implement sophisticated AI and pathfinding
- **Engagement Encouragement**: Timeout resolution rewards aggressive play
- **Definitive Outcomes**: All scenarios have clear win/loss/draw determination
- **Fairness**: Deterministic, predictable resolution rules with complete information
- **Integration**: Seamlessly combines mechanics from ADR-0006 through ADR-0009
- **Future Extensibility**: Visibility property can be enhanced later (fog of war modes) without disrupting core mechanics
