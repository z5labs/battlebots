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

## Combat Actions

Combat actions deal damage to opponent bots. All combat actions require specific weapon equipment to be available.

### RifleShot

A single-shot, precise ranged attack.

- **Energy Cost**: 15 (TBD)
- **Cooldown**: 1 tick (TBD)
- **Equipment Required**: Rifle
- **Damage**: Moderate (TBD)
- **Range**: Long (TBD)
- **Constraints**: Requires line of sight, rifle must be equipped

### ShotgunBlast

A spray of projectiles attack effective at close range.

- **Energy Cost**: 20 (TBD)
- **Cooldown**: 2 ticks (TBD)
- **Equipment Required**: Shotgun
- **Damage**: High at close range with damage falloff based on distance (TBD)
- **Range**: Short to medium (TBD)
- **Constraints**: Shotgun must be equipped, damage decreases with distance

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
