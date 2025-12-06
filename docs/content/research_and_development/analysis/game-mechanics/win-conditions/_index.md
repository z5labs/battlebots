---
title: "Win Conditions"
description: "How battles conclude and determine winners"
type: docs
weight: 5
date: 2025-12-05
---

## Overview

Battle resolution defines how a 1v1 battle concludes and determines the winner. Every battle must have a definitive outcome: victory, defeat, or draw. This document specifies the conditions under which each outcome occurs and how edge cases are handled.

The battle engine must continuously monitor the game state to detect when any of these conditions are met and immediately conclude the battle with the appropriate outcome.

## Victory Conditions

A bot achieves victory when any of the following conditions are met:

### Destruction

The enemy bot's health reaches 0 or below. This is the primary victory condition and occurs when a bot successfully depletes its opponent's health through combat.

**Trigger**: `enemy.health <= 0`

**Outcome**: Immediate victory for the surviving bot

### Forfeit

The enemy bot disconnects from the battle server or fails to respond within the allowed timeout period. This represents an abandonment of the battle.

**Trigger**: Enemy bot connection lost or timeout exceeded

**Outcome**: Immediate victory for the connected bot

### Timeout

The battle time limit is reached and the bot has more health remaining than its opponent. This represents a victory by attrition.

**Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health > enemy.health`

**Outcome**: Victory for the bot with higher health

**Note**: The specific time limit value is TBD (placeholder: 5 minutes)

## Defeat Conditions

A bot is defeated when any of the following conditions are met:

### Destruction

The bot's own health reaches 0 or below due to enemy actions or environmental damage.

**Trigger**: `bot.health <= 0`

**Outcome**: Immediate defeat

### Forfeit

The bot disconnects from the battle server or fails to respond within the allowed timeout period.

**Trigger**: Bot connection lost or timeout exceeded

**Outcome**: Immediate defeat

### Timeout

The battle time limit is reached and the bot has less health remaining than its opponent.

**Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health < enemy.health`

**Outcome**: Defeat (opponent wins by timeout)

## Draw Conditions

A battle results in a draw when neither bot achieves a clear victory:

### Mutual Destruction

Both bots reach 0 health simultaneously in the same game tick. This can occur when:
- Both bots deal fatal damage to each other in the same action resolution phase
- Both bots are destroyed by simultaneous environmental effects
- A collision or explosion affects both bots fatally at the same moment

**Trigger**: `bot.health <= 0 && enemy.health <= 0` (same tick)

**Outcome**: Draw

### Equal Health at Timeout

The battle time limit is reached and both bots have exactly the same health remaining.

**Trigger**: `battle.elapsed_time >= TIME_LIMIT && bot.health == enemy.health`

**Outcome**: Draw

## Time Limits

Battle duration is capped to prevent indefinitely long battles and ensure matches conclude in a reasonable timeframe.

**Time Limit**: TBD (placeholder: 5 minutes / 300 seconds)

When the time limit is reached:
1. The battle engine stops accepting new actions
2. All pending actions in the current tick are resolved
3. Final health values are compared
4. The outcome is determined based on health comparison

**Design Considerations**:
- Time limit should be long enough to allow for strategic gameplay
- Time limit should be short enough to maintain engagement
- Time limit may need to be configurable for different battle modes or tournaments

## Edge Cases

### Simultaneous Actions

When both bots perform actions in the same game tick that could affect the battle outcome:

**Action Resolution Order**:
1. All actions for the current tick are collected
2. Actions are resolved in a deterministic order (e.g., by bot ID, action type priority)
3. Game state is updated after all actions are processed
4. Win conditions are checked after state update

**Example**: If both bots fire projectiles that would destroy each other in the same tick, both hits are processed, and if both bots reach 0 health, the result is a draw (mutual destruction).

### Disconnect Handling

When a bot disconnects or becomes unresponsive:

**Grace Period**: TBD - A brief timeout window allows the bot to reconnect or respond before forfeit is declared

**Reconnection**: If a bot reconnects within the grace period, the battle continues from the current state

**No Response**: If the grace period expires without response, the battle immediately ends with a forfeit victory for the opponent

**Both Disconnect**: If both bots disconnect simultaneously or within the grace period, the outcome is determined by:
- If one reconnects within grace period: that bot wins by forfeit
- If neither reconnects: battle ends as a draw

### Timeout and Destruction in Same Tick

If the time limit is reached in the same tick that a bot is destroyed:

**Priority**: Destruction takes precedence over timeout

**Rationale**: Destruction is a more definitive outcome and should be the primary victory condition

**Example**: If the time limit expires at tick 1000, and in that same tick a bot is destroyed, the outcome is victory by destruction (not timeout).

### Negative Health Values

Bots can temporarily have negative health values if damage exceeds remaining health:

**Example**: Bot has 5 health, takes 20 damage → health becomes -15

**Win Condition Check**: Any health value <= 0 triggers destruction

**Display**: For UI/visualization purposes, negative health should be displayed as 0

### Zero Health at Battle Start

If a bot somehow begins a battle with 0 or negative health (due to a bug or configuration error):

**Behavior**: Battle immediately ends with defeat for the bot with invalid health

**Prevention**: Battle engine should validate initial state before battle begins
