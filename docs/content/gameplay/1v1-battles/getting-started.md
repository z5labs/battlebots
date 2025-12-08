---
title: "Getting Started"
description: >
  Quick start guide for creating your first 1v1 battle bot
type: docs
weight: 10
---

This guide will help you create your first bot for 1v1 battles. We'll cover the basics of bot structure, equipment configuration, and connecting to the battle server.

## Prerequisites

- Programming experience in any language
- Battle Bots SDK installed for your chosen language
- Development environment set up

## Bot Structure

A Battle Bots bot consists of:

1. **Equipment Configuration** - Choose your weapon and armor before battle
2. **State Processing** - Receive and process battle state updates each tick
3. **Action Submission** - Decide which actions to perform based on current state
4. **Communication** - Connect to the battle server via gRPC

## Minimal Bot Example

Here's a minimal bot structure to get you started:

### Step 1: Configure Equipment

Before battle, configure your bot's loadout with 1 weapon and 1 armor:

```python
# Example: Simple balanced loadout
loadout = {
    "weapon": "Rifle",      # Baseline precision weapon
    "armor": "Light Armor"  # Minimal defense, maintains mobility
}
```

See the [Equipment](equipment/) page for all available weapons and armor options.

### Step 2: Connect to Battle Server

Your bot connects to the battle server using the SDK:

```python
from battlebots_sdk import BattleBot

bot = BattleBot(loadout=loadout)
bot.connect()
```

### Step 3: Process State Updates

Each tick, your bot receives a state update with:

- Your bot's position, health, energy
- Opponent bot's position (if visible)
- Arena boundaries
- Available actions

```python
def process_state(state):
    my_position = state.my_bot.position  # (x, y) coordinates
    my_health = state.my_bot.health
    my_energy = state.my_bot.energy

    opponent = state.opponent_bot  # May be None if out of sight
    if opponent:
        opponent_position = opponent.position
        opponent_health = opponent.health

    return decide_action(state)
```

### Step 4: Submit Actions

Based on the state, your bot decides which action to perform:

```python
def decide_action(state):
    # Example: Simple logic
    if state.opponent_bot:
        # Opponent visible - check if in range
        distance = calculate_distance(state.my_bot.position,
                                      state.opponent_bot.position)

        if distance < RIFLE_RANGE and state.my_bot.energy >= 15:
            return Action.RifleShot(target=state.opponent_bot.position)

    # Move toward center if opponent not visible
    return Action.Move(target=(0, 0))
```

### Step 5: Run Your Bot

```python
# Main bot loop
while bot.is_battle_active():
    state = bot.get_current_state()
    action = process_state(state)
    bot.submit_action(action)
```

## Next Steps

Now that you have a basic bot structure, learn about the game mechanics:

1. **[Arena](arena/)** - Understand the battle space, coordinates, and collision rules
2. **[Bot Characteristics](bot-characteristics/)** - Learn about Health, Defense, and Mass stats
3. **[Equipment](equipment/)** - Explore different weapons and armor options
4. **[Actions](actions/)** - Discover all available actions, their costs, and constraints

## Important Notes

- **SDK Abstracts Complexity**: Your SDK handles gRPC communication, coordinate calculations, and state management
- **Tick-Based**: Battles run in real-time ticks; your bot should process state and respond quickly
- **Energy Management**: Monitor your energy pool; actions cost energy to perform
- **Cooldowns**: Some actions have cooldown periods before they can be used again

Start simple, test your bot in battles, and iterate on your approach. The mechanics documentation will help you understand the full system, but figuring out effective bot logic is your challenge!
