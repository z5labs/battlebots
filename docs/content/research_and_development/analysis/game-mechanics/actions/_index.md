---
title: "Bot Actions"
description: "Catalog of all actions bots can perform in battle"
type: docs
weight: 4
date: 2025-12-05
---

## Overview

Actions are the primary means by which bots interact with the battle environment. Each action has associated costs and constraints that govern when and how it can be used:

- **Energy Cost**: Amount of energy required to perform the action (deducted from the bot's energy pool)
- **Cooldown**: Number of ticks that must pass before the action can be used again
- **Equipment Requirements**: Some actions require specific equipment to be equipped
- **Execution**: Actions are submitted via gRPC and processed during the tick cycle

All energy values reference the energy pool system defined in [Characteristics](../characteristics/).

## Movement Actions

Movement actions allow bots to navigate the 2D battle space.

### Move

Moves the bot in a specified direction within the battle grid.

- **Energy Cost**: 5 (TBD)
- **Cooldown**: 0 ticks (TBD)
- **Parameters**: Direction vector or target coordinates
- **Constraints**: Cannot move through obstacles or other bots

### Rotate

Rotates the bot to face a new direction.

- **Energy Cost**: 2 (TBD)
- **Cooldown**: 0 ticks (TBD)
- **Parameters**: Target angle or relative rotation
- **Constraints**: Instantaneous rotation

### Dash

Performs a rapid movement burst, covering more distance than a standard move.

- **Energy Cost**: 15 (TBD)
- **Cooldown**: 3 ticks (TBD)
- **Parameters**: Direction vector
- **Constraints**: Higher energy cost than Move but covers more ground quickly

## Combat Actions

Combat actions deal damage to opponent bots.

### BasicAttack

A standard melee or short-range attack.

- **Energy Cost**: 10 (TBD)
- **Cooldown**: 1 tick (TBD)
- **Damage**: Variable based on bot stats (TBD)
- **Range**: Short (TBD)
- **Constraints**: Target must be within range

### LaserShot

A ranged energy weapon attack requiring laser equipment.

- **Energy Cost**: 25 (TBD)
- **Cooldown**: 2 ticks (TBD)
- **Equipment Required**: Laser
- **Damage**: Higher than BasicAttack (TBD)
- **Range**: Long (TBD)
- **Constraints**: Requires line of sight, equipment must be equipped

### HeavyBlow

A powerful melee attack with high damage and high cost.

- **Energy Cost**: 30 (TBD)
- **Cooldown**: 4 ticks (TBD)
- **Damage**: High (TBD)
- **Range**: Melee (TBD)
- **Constraints**: Longer cooldown balances high damage output

## Defensive Actions

Defensive actions mitigate incoming damage or avoid attacks.

### Block

Reduces damage from incoming attacks.

- **Energy Cost**: 10 (TBD)
- **Cooldown**: 2 ticks (TBD)
- **Effect**: Reduces damage by percentage (TBD)
- **Duration**: Active for current tick only
- **Constraints**: Must be activated before damage is received

### Evade

Attempts to dodge incoming attacks.

- **Energy Cost**: 15 (TBD)
- **Cooldown**: Variable (TBD)
- **Effect**: Chance to completely avoid attack (TBD)
- **Duration**: Active for current tick only
- **Constraints**: Success rate may depend on bot stats

### Shield

Activates an energy shield that absorbs damage over multiple ticks.

- **Energy Cost**: 20 (TBD)
- **Cooldown**: Variable (TBD)
- **Effect**: Absorbs damage up to threshold (TBD)
- **Duration**: Sustained (multiple ticks, ongoing energy drain - TBD)
- **Constraints**: May have ongoing energy cost while active

## Utility Actions

Utility actions provide information or modify bot state without direct combat effects.

### Scan

Gathers information about the environment and nearby bots.

- **Energy Cost**: 5 (TBD)
- **Cooldown**: Variable (TBD)
- **Effect**: Returns information about nearby entities (positions, health, etc.)
- **Range**: Limited detection radius (TBD)
- **Constraints**: May only reveal information within range

### Charge

Increases energy regeneration rate temporarily.

- **Energy Cost**: Variable (may cost initial energy or pause regen - TBD)
- **Cooldown**: Variable (TBD)
- **Effect**: Boosts energy regeneration for duration (TBD)
- **Duration**: Multiple ticks (TBD)
- **Constraints**: Bot may be vulnerable while charging (cannot perform other actions - TBD)

## gRPC Mapping

Actions are submitted from bots to the battle server using the gRPC protocol defined in [ADR-0004: Bot to Battle Server Interface](/battlebots/research_and_development/adrs/0004-bot-to-battle-server-interface/).

### BotAction Message Structure

Each action is encoded in a `BotAction` message that includes:

- **Action Type**: Enum identifying which action to perform (e.g., `ACTION_MOVE`, `ACTION_BASIC_ATTACK`)
- **Parameters**: Action-specific parameters (direction vectors, target IDs, etc.)
- **Tick Number**: The tick for which this action is intended

### Example Action Mapping

```protobuf
message BotAction {
  ActionType type = 1;
  map<string, string> parameters = 2;
  uint64 tick = 3;
}

enum ActionType {
  ACTION_UNKNOWN = 0;
  ACTION_MOVE = 1;
  ACTION_ROTATE = 2;
  ACTION_DASH = 3;
  ACTION_BASIC_ATTACK = 4;
  ACTION_LASER_SHOT = 5;
  ACTION_HEAVY_BLOW = 6;
  ACTION_BLOCK = 7;
  ACTION_EVADE = 8;
  ACTION_SHIELD = 9;
  ACTION_SCAN = 10;
  ACTION_CHARGE = 11;
}
```

### Parameter Encoding

Each action type has specific parameter requirements:

- **Movement actions**: Direction vector (x, y) or angle
- **Combat actions**: Target ID or direction
- **Defensive actions**: Duration or activation flag
- **Utility actions**: Action-specific parameters (scan radius, charge duration, etc.)

The exact parameter encoding and validation rules are defined in the battle server implementation.

### Action Validation

The battle server validates each submitted action to ensure:

1. Bot has sufficient energy to perform the action
2. Action is not on cooldown
3. Required equipment is equipped (if applicable)
4. Parameters are valid and within acceptable ranges
5. Action is legal given current game state

Invalid actions are rejected, and an error response is sent back to the bot via the gRPC stream.
