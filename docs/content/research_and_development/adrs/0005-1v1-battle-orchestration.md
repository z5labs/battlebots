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

Battle Bots requires high-level orchestration rules for 1v1 battles that tie together the spatial system, characteristics, equipment, and actions into a complete battle experience. We need to define spatial configuration (2D vs. 3D, bounded vs. unbounded arena), visibility rules (fog of war vs. full battlefield vision), battle pacing mechanisms (time limits, engagement encouragement), and win conditions (how battles conclude). These orchestration rules must create engaging gameplay while providing a foundation for detailed battle mechanics.

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
* **Integration** - Must work seamlessly with spatial mechanics, bot characteristics, equipment systems, and action mechanics
* **Observability** - Battle state and outcomes must be clear for visualization and analysis
* **Time Constraints** - Battles must conclude in reasonable timeframes

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
- Cartesian coordinate system naturally supports 2D battle space
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
- Bounded arena provides clear spatial constraints for the battle system

**Alternative Considered: Option 4.2 - Bounded and Shrinking** would add additional engagement pressure through environmental forcing function, but adds implementation complexity (dynamic boundaries, bot notification of shrinking, damage from being outside boundaries) that is unnecessary given the timeout mechanism already prevents stalemates. Could be reconsidered as an optional game mode for more aggressive pacing.

**Alternative Rejected: Option 4.3 - Unbounded** would allow indefinite evasion, potentially creating stalemate scenarios even with timeout (bots could avoid each other for entire battle). Unbounded space also complicates pathfinding (no clear boundaries for navigation) and visualization (camera must follow bots across potentially large distances).

### Consequences

**Visibility Decision (Full Visibility)**:

* Good, because simplest possible implementation (no detection mechanics, line-of-sight, or visibility protocol)
* Good, because encourages users to implement sophisticated pathfinding solutions with complete battlefield information
* Good, because enables advanced AI strategies - bots have complete information for tactical decisions
* Good, because easier to debug and visualize (all state always visible)
* Good, because reduces initial POC scope - allows focus on core battle mechanics
* Good, because integrates seamlessly with spatial system without additional visibility calculations
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
* Good, because bounded arena definition provides clear spatial constraints
* Neutral, because could be enhanced with shrinking boundaries in future game modes
* Bad, because fixed boundaries may enable defensive corner camping strategies
* Bad, because lacks additional engagement pressure that shrinking boundaries would provide

**Overall Integration**:

* Good, because all four properties integrate seamlessly with the spatial system, bot characteristics, equipment systems, and action mechanics
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
8. Bounded arena boundary enforcement testing to ensure bots cannot escape or clip through boundaries
9. 2D pathfinding and collision detection testing to confirm adequate performance
10. Visualization testing to confirm 2D display provides clear battle state representation
11. Future evaluation of fog of war as optional game mode once core mechanics are validated
12. Future consideration of 3D spatial system if aerial combat becomes strategically important
13. Future consideration of shrinking boundaries as optional game mode for aggressive pacing

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
* Neutral, because requires line-of-sight calculations
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

**Win Determination**: Bot whose opponent reaches `health <= 0` wins. No other win conditions exist.

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

**Win Determination**:
- Bot whose opponent reaches `health <= 0` wins
- If timeout: bot with higher health wins; equal health results in draw

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

**Win Determination**:
- Bot whose opponent reaches `health <= 0` wins
- If timeout: bot with higher health wins; equal health results in draw
- If disconnect: connected bot wins after grace period expires; bot that reconnects within grace period continues battle

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
* Good, because bounded arena definition provides clear spatial constraints
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

- **[ADR-0004: Bot to Battle Server Interface](0004-bot-battle-server-interface.md)**: gRPC protocol for battle state synchronization and visibility information

- **[Game Mechanics Analysis](../analysis/game-mechanics/)**: Detailed framework documentation including:
  - [Battle Space](../analysis/game-mechanics/battle-space/): Spatial system details
  - [Bot Characteristics](../analysis/game-mechanics/characteristics/): Stat system details
  - [Equipment System](../analysis/game-mechanics/equipment/): Loadout mechanics
  - [Bot Actions](../analysis/game-mechanics/actions/): Action catalog
  - [Win Conditions](../analysis/game-mechanics/win-conditions/): Battle resolution details

- **[POC User Journey](../user-journeys/0001-poc.md)**: Proof of concept implementation using these battle orchestration rules

**Future ADRs**: This orchestration ADR will inform future detailed ADRs for spatial systems, bot characteristics, equipment systems, and action mechanics.

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
- **Integration**: Seamlessly combines spatial mechanics, bot characteristics, equipment systems, and action mechanics
- **Future Extensibility**: All properties (visibility, spatial dimensions, boundaries) can be enhanced later without disrupting core mechanics
