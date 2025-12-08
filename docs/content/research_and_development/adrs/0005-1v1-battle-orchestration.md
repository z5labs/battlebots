---
title: "[0005] 1v1 Battle Orchestration"
description: >
   High-level orchestration of 1v1 battles including visibility rules, battle pacing, and win conditions
type: docs
weight: 5
category: "strategic"
status: "accepted"
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

Battle Bots requires high-level orchestration rules for 1v1 battles that tie together the spatial system, characteristics, equipment, and actions into a complete battle experience. We need to define spatial configuration (2D vs. 3D, bounded vs. unbounded arena), visibility rules (fog of war vs. full battlefield vision), battle pacing mechanisms (time limits, engagement encouragement), and win conditions (how battles conclude). These orchestration rules must create engaging gameplay while integrating the detailed mechanics defined in ADR-0006 through ADR-0009.

Without well-defined battle orchestration, we cannot:
- Define the spatial dimensions and arena boundaries for battles (spatial configuration)
- Determine what information bots can see during battle (visibility)
- Prevent stalemate and ensure battles conclude in reasonable time (pacing)
- Definitively determine battle outcomes (win conditions)
- Balance information advantage vs. direct combat capability
- Create tension and encourage engagement between bots
- Design the battle flow from start to conclusion

## Decision Drivers

* **Implementation Complexity** - Spatial system should balance strategic depth with development simplicity
* **Accessibility** - Bot developers should have accessible pathfinding and movement algorithms
* **Computational Efficiency** - Spatial calculations (pathfinding, collision, physics) should be performant
* **Engagement Guarantee** - Arena boundaries should ensure bots cannot avoid combat indefinitely
* **Information Strategy Depth** - Visibility rules should create meaningful choices (sensor investment vs. stealth)
* **Battle Pacing** - Mechanisms should prevent passive stalemates and force interaction
* **Definitive Outcomes** - Win conditions must handle all scenarios including edge cases
* **Fairness** - All rules should be deterministic, predictable, and balanced
* **Integration** - Must work seamlessly with mechanics from ADR-0006 through ADR-0009
* **Observability** - Battle state and outcomes must be clear for visualization and analysis
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

### Property 3: Spatial Dimensionality

* **Option 3.1: 2D** - Battles occur in two-dimensional space (x, y coordinates)
* **Option 3.2: 3D** - Battles occur in three-dimensional space (x, y, z coordinates)

### Property 4: Spatial Manifold

* **Option 4.1: Bounded** - Fixed arena boundaries that remain constant throughout the battle
* **Option 4.2: Bounded and Shrinking** - Arena boundaries that progressively shrink over time to force engagement
* **Option 4.3: Unbounded** - No arena boundaries (infinite or very large playable area)

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

### Property 3: Spatial Dimensionality

**Chosen: Option 3.1 - 2D**

Rationale:
- Simpler implementation for initial POC - 2D spatial system reduces complexity in pathfinding, collision detection, and visualization
- Lower barrier to entry for bot developers - 2D pathfinding and movement algorithms are more accessible and well-documented
- Easier visualization - 2D battles can be displayed directly on screen without camera controls or 3D rendering complexity
- Sufficient strategic depth - 2D space provides adequate complexity for positioning, flanking, and tactical movement
- Consistent with ADR-0006 spatial system design - the cartesian coordinate system is designed for 2D battle space
- Reduces computational requirements - 2D physics and collision detection are significantly less computationally expensive than 3D

**Alternative Rejected: Option 3.2 - 3D** would add significant complexity to pathfinding (3D A*), collision detection (3D physics), and visualization (3D rendering, camera controls) without proportional gameplay benefit for 1v1 battles. Can be reconsidered for future game modes if aerial combat or vertical positioning becomes strategically important.

### Property 4: Spatial Manifold

**Chosen: Option 4.1 - Bounded**

Rationale:
- Guarantees engagement - Fixed boundaries ensure bots cannot escape indefinitely, promoting interaction
- Predictable arena - Bots can rely on constant arena size for pathfinding and strategic planning
- Simpler implementation - No dynamic boundary updates or shrinking mechanics to track
- Consistent with timeout mechanism - Bounded arena complements the timeout system (ADR Property 2) for ensuring battle conclusion
- Fair and balanced - Both bots have equal access to the full arena throughout the battle
- Enables positioning strategy - Fixed boundaries create meaningful positional choices (corner control, center dominance, edge play)
- Integrates with ADR-0006 - The bounded arena defined in the spatial system ADR

**Alternative Considered: Option 4.2 - Bounded and Shrinking** would add additional engagement pressure through environmental forcing function, but adds implementation complexity (dynamic boundaries, bot notification of shrinking, damage from being outside boundaries) that is unnecessary given the timeout mechanism already prevents stalemates. Could be reconsidered as an optional game mode for more aggressive pacing.

**Alternative Rejected: Option 4.3 - Unbounded** would allow indefinite evasion, potentially creating stalemate scenarios even with timeout (bots could avoid each other for entire battle). Unbounded space also complicates pathfinding (no clear boundaries for navigation) and visualization (camera must follow bots across potentially large distances).

### Battle Orchestration Specification

#### Spatial System Rules

**2D Spatial Dimensionality**:

All battles occur in two-dimensional space with cartesian coordinates.

**Coordinate System**:
- **X-axis**: Horizontal position (left-right movement)
- **Y-axis**: Vertical position (up-down movement)
- **Position Format**: (x, y) coordinate pairs
- **No Z-axis**: No vertical elevation or aerial movement

**Movement Implications**:
- Bots move in 2D plane only (no jumping, flying, or vertical positioning)
- Pathfinding algorithms operate in 2D space (A*, navigation meshes, etc.)
- Collision detection uses 2D bounding boxes or circles
- Line-of-sight calculations use 2D raycasting (if fog of war implemented in future)

**Strategic Considerations**:
- **Positioning**: Tactical positioning focuses on x,y coordinates relative to enemy and obstacles
- **Flanking**: Bots can flank by moving around obstacles in 2D space
- **Distance Management**: Range calculations use 2D distance formulas (Euclidean or Manhattan)
- **No Vertical Advantage**: No height-based tactical advantages (high ground, aerial bombardment)

**Integration with ADR-0006**: The 2D cartesian coordinate system aligns with the spatial system defined in ADR-0006 (Battle Space Spatial System).

**Bounded Arena Manifold**:

All battles occur within a fixed rectangular arena with constant boundaries.

**Arena Boundaries**:
- **Fixed Size**: Arena dimensions remain constant throughout the battle
- **Rectangular Shape**: Arena defined by minimum and maximum x,y coordinates
- **Boundary Enforcement**: Bots cannot move outside arena boundaries (collision/clamping)
- **No Shrinking**: Boundaries do not change or shrink during battle

**Boundary Behavior**:
- **Collision**: Bots that attempt to move beyond boundaries are stopped at the edge
- **Pathfinding**: Bot pathfinding must respect arena boundaries as impassable walls
- **Visibility**: Arena boundaries are visible to all bots (provided as part of battlefield state)

**Arena Information Provided**:
- **Arena Width**: Maximum x-coordinate minus minimum x-coordinate
- **Arena Height**: Maximum y-coordinate minus minimum y-coordinate
- **Arena Bounds**: Minimum x, minimum y, maximum x, maximum y coordinates

**Strategic Implications**:
- **Corner Control**: Corners provide defensive positioning but limit escape options
- **Center Dominance**: Center positioning maximizes movement options
- **Edge Play**: Bots can use boundaries as strategic barriers (preventing flanking from one side)
- **Bounded Engagement**: Bots cannot indefinitely flee; boundaries ensure eventual engagement

**Design Rationale**: Bounded arena guarantees engagement by preventing indefinite evasion, provides predictable space for pathfinding algorithms, and creates meaningful positional strategy through corner control and center dominance.

#### Battlefield Visibility Rules

**Full Visibility System**:

All bots have complete battlefield visibility at all times. Each bot receives full information about all entities and state on the battlefield.

**Information Provided to All Bots**:
- **Enemy Bot Position**: Complete 2D (x, y) coordinates of all enemy bots
- **Enemy Bot Health**: Current HP value of all enemy bots
- **Enemy Bot Orientation**: Facing direction of all enemy bots (for predicting movement)
- **Enemy Bot Equipment**: Loadout information for all enemy bots (visible equipment from ADR-0008)
- **Battlefield Obstacles**: All obstacle positions and dimensions in 2D space (from ADR-0006)
- **Arena Boundaries**: Complete bounded arena dimensions (min x, min y, max x, max y coordinates)
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

**Spatial Dimensionality Decision (2D)**:

* Good, because simplest spatial implementation (2D coordinates, pathfinding, collision detection)
* Good, because lowers barrier to entry for bot developers (accessible 2D algorithms)
* Good, because easier visualization (direct 2D display without 3D rendering complexity)
* Good, because reduces computational requirements (2D physics less expensive than 3D)
* Good, because sufficient strategic depth for positioning and tactical movement
* Good, because integrates seamlessly with ADR-0006 cartesian coordinate system
* Neutral, because may limit future aerial or vertical combat mechanics
* Bad, because eliminates vertical positioning as a strategic dimension
* Bad, because may feel limiting compared to 3D games if users expect aerial combat

**Spatial Manifold Decision (Bounded)**:

* Good, because guarantees engagement (bots cannot escape indefinitely)
* Good, because predictable arena size enables reliable pathfinding and strategy
* Good, because simplest manifold implementation (fixed boundaries)
* Good, because complements timeout mechanism for battle conclusion
* Good, because fair and balanced (equal arena access for both bots)
* Good, because enables meaningful positional strategy (corner control, center dominance)
* Good, because integrates with ADR-0006 bounded arena definition
* Neutral, because could be enhanced with shrinking boundaries in future game modes
* Bad, because fixed boundaries may enable defensive corner camping strategies
* Bad, because lacks additional engagement pressure that shrinking boundaries would provide

**Overall Integration**:

* Good, because all four properties integrate seamlessly with ADR-0006 (spatial system), ADR-0007 (characteristics), ADR-0008 (equipment), and ADR-0009 (actions)
* Good, because property-based decision structure allows independent tuning and future modifications
* Good, because each property addresses distinct battle orchestration concerns
* Good, because 2D bounded arena aligns with full visibility decision (simpler implementation, easier visualization)

### Confirmation

The decision will be confirmed through:

1. Implementation of full visibility system with complete battlefield state synchronization
2. Validation that all bots receive complete and accurate battlefield information
3. Playtesting to tune time limits and ensure engagement without stalemate
4. Win condition testing to verify all scenarios (destruction, forfeit, timeout, draw) resolve correctly
5. Edge case validation for simultaneous actions, disconnects, and timeout priorities
6. User testing to confirm pathfinding and AI strategy opportunities with complete information
7. Timeout scenario analysis to confirm aggressive play is rewarded over passive stalemate
8. 2D spatial system implementation validation with ADR-0006 cartesian coordinates
9. Bounded arena boundary enforcement testing to ensure bots cannot escape or clip through boundaries
10. 2D pathfinding and collision detection testing to confirm adequate performance
11. Visualization testing to confirm 2D display provides clear battle state representation
12. Future evaluation of fog of war as optional game mode once core mechanics are validated
13. Future consideration of 3D spatial system if aerial combat becomes strategically important
14. Future consideration of shrinking boundaries as optional game mode for aggressive pacing

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

### Property 3: Spatial Dimensionality

#### Option 3.1: 2D (CHOSEN)

Battles occur in two-dimensional space (x, y coordinates).

* Good, because simplest spatial implementation (2D pathfinding, collision detection, physics)
* Good, because lower computational requirements compared to 3D
* Good, because easier visualization (direct 2D display, no camera controls needed)
* Good, because lower barrier to entry for bot developers (2D algorithms more accessible)
* Good, because well-documented algorithms (A*, 2D vector math, 2D physics)
* Good, because sufficient strategic depth for positioning and movement tactics
* Good, because integrates seamlessly with ADR-0006 cartesian coordinate system
* Neutral, because appropriate for ground-based combat scenarios
* Neutral, because may be extended to 3D in future if needed
* Bad, because eliminates vertical positioning as strategic dimension
* Bad, because no aerial combat or flying units
* Bad, because may feel limiting if users expect 3D movement

#### Option 3.2: 3D

Battles occur in three-dimensional space (x, y, z coordinates).

* Good, because enables vertical positioning strategy (height advantage)
* Good, because supports aerial combat and flying units
* Good, because additional strategic dimension (above/below positioning)
* Good, because familiar from many modern games
* Neutral, because enables jump mechanics or flight equipment
* Neutral, because may be more engaging for some users
* Bad, because significantly more complex implementation (3D pathfinding, collision, physics)
* Bad, because much higher computational requirements (3D calculations expensive)
* Bad, because complex visualization (3D rendering, camera controls, depth perception)
* Bad, because higher barrier to entry (3D algorithms more complex)
* Bad, because more difficult to debug and visualize battles
* Bad, because may be unnecessary complexity for 1v1 ground combat

### Property 4: Spatial Manifold

#### Option 4.1: Bounded (CHOSEN)

Fixed arena boundaries that remain constant throughout the battle.

* Good, because guarantees engagement (bots cannot escape indefinitely)
* Good, because predictable arena enables reliable pathfinding
* Good, because simplest implementation (fixed boundary checks)
* Good, because fair and balanced (equal arena access)
* Good, because enables positional strategy (corner control, center dominance)
* Good, because complements timeout mechanism for battle conclusion
* Good, because integrates with ADR-0006 bounded arena definition
* Good, because clear visualization (fixed arena visible throughout)
* Neutral, because could be enhanced with shrinking in future modes
* Neutral, because requires boundary collision detection
* Bad, because may enable corner camping defensive strategies
* Bad, because lacks additional engagement pressure
* Bad, because fixed size may feel static compared to dynamic boundaries

#### Option 4.2: Bounded and Shrinking

Arena boundaries that progressively shrink over time to force engagement.

* Good, because provides additional engagement pressure (environmental forcing function)
* Good, because prevents passive stalemate through space reduction
* Good, because familiar mechanic from battle royale games
* Good, because creates urgency and encourages aggressive play
* Good, because adds dynamic element to battles
* Neutral, because requires tuning shrink rate and damage mechanics
* Neutral, because complements timeout mechanism (double engagement pressure)
* Bad, because significantly more complex implementation (dynamic boundaries)
* Bad, because requires bot notification system for boundary changes
* Bad, because requires out-of-bounds damage or elimination mechanic
* Bad, because may feel artificial or gimmicky
* Bad, because adds state tracking complexity (shrink timing, positions)
* Bad, because unnecessary given timeout already prevents stalemate
* Bad, because may disadvantage mobility-focused builds unfairly

#### Option 4.3: Unbounded

No arena boundaries (infinite or very large playable area).

* Good, because enables unlimited strategic space
* Good, because no artificial boundaries constraining movement
* Good, because may feel more "realistic" for some scenarios
* Neutral, because appropriate for exploration-focused gameplay
* Neutral, because eliminates boundary collision checks
* Bad, because allows indefinite evasion (bots can run away forever)
* Bad, because enables stalemate even with timeout (avoid combat entire battle)
* Bad, because complicates pathfinding (no clear navigation boundaries)
* Bad, because difficult visualization (camera must track across large distances)
* Bad, because may create boring gameplay (chasing fleeing opponents)
* Bad, because timeout becomes only victory condition (no engagement guarantee)
* Bad, because eliminates positional strategy (no boundaries to control)
* Bad, because fundamentally incompatible with engagement-focused 1v1 gameplay

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
- 2D spatial system reduces implementation complexity and lowers barrier to entry for bot developers
- Bounded arena guarantees engagement and complements timeout mechanism for battle conclusion
- Time limit with HP-based timeout resolution encourages aggressive play over passive stalemate
- Destruction takes precedence over timeout to prioritize definitive outcomes
- Grace period for disconnects balances technical issues with competitive fairness
- Sensor Array and Stealth Module equipment provide stat bonuses without visibility effects in this mode

**Future Considerations**:
- **Variable Time Limits**: Different battle modes (quick match vs. tournament) may have different time limits
- **Fog of War Game Modes**: Constant or revealed fog of war could be added as optional game modes to create information strategy layer
- **3D Spatial System**: If aerial combat or vertical positioning becomes strategically important, 3D space could be added as optional game mode
- **Shrinking Arena**: Bounded and shrinking boundaries could be added as optional game mode for more aggressive pacing
- **Overtime Mechanics**: If timeout occurs with close HP values, brief overtime period could be added
- **Spectator Mode**: Already full visibility for all participants; spectators share same view
- **Replay System**: Battle recordings with complete visibility throughout
- **Partial Information Modes**: Once core mechanics proven, visibility could become configurable property per game mode

### Design Principles

The battle orchestration follows these principles:
- **Simplicity First**: Full visibility and 2D bounded arena reduce initial complexity and enable focus on core mechanics
- **Complete Information Strategy**: Bots have full battlefield knowledge to implement sophisticated AI and pathfinding
- **Accessibility**: 2D spatial system lowers barrier to entry for bot developers with well-documented algorithms
- **Guaranteed Engagement**: Bounded arena and timeout mechanism ensure battles conclude with interaction
- **Engagement Encouragement**: Timeout resolution rewards aggressive play over passive stalemate
- **Definitive Outcomes**: All scenarios have clear win/loss/draw determination
- **Fairness**: Deterministic, predictable resolution rules with complete information and equal arena access
- **Integration**: Seamlessly combines mechanics from ADR-0006 through ADR-0009
- **Future Extensibility**: All properties (visibility, spatial dimensions, boundaries) can be enhanced later without disrupting core mechanics
