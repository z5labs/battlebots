---
title: "1v1 Battles"
description: >
  Complete mechanics guide for 1v1 battles - arena, bot characteristics, equipment, and actions
type: docs
weight: 10
---

In 1v1 battles, two bots face off in direct combat within a **Battle Arena**. The arena is a configured instance of the BattleBot Universe with specific properties: terrain type (biome), boundary dimensions, visibility rules, starting positions, and win conditions. Victory goes to the bot that satisfies any of the win conditions first: reducing opponent's health to zero, having more health when the timeout expires, or preventing opponent disconnection recovery.

## What You'll Find Here

This section contains complete mechanics documentation for 1v1 battles:

- **[Getting Started](getting-started/)** - Quick start guide with a minimal bot example
- **[Arena](arena/)** - The 2D battle space, coordinates, boundaries, collision, and physics
- **[Bot Characteristics](bot-characteristics/)** - Health, Defense, and Mass stats that define your bot
- **[Equipment](equipment/)** - Weapons and armor that customize your bot's capabilities
- **[Actions](actions/)** - All available actions, energy costs, and cooldowns

## Key Concepts

### Real-Time Gameplay

1v1 battles operate in real-time with a continuous tick-based game loop. Your bot receives state updates each tick and can submit actions to perform.

### Win Conditions

A battle concludes when any of the following conditions is met (checked in order):

1. **Elimination** (Primary): Your bot defeats opponent by reducing their Health to 0 - immediate victory
2. **Timeout** (Fallback): Battle reaches 5-minute time limit - higher Health wins
3. **Disconnect** (Safety): Opponent disconnects and fails to reconnect within 30-second grace period - your bot wins

See **[ADR-0011: 1v1 Battles](../../../research_and_development/adrs/0011-1v1-battles.md)** for complete technical specification of arena properties and win conditions.

### Bot Customization

Before battle, you configure your bot's equipment loadout (1 weapon + 1 armor). During battle, your bot performs actions that consume energy and have cooldown periods.

### Programming Challenge

This documentation explains the mechanics and rules. Developing effective bot logic, decision-making algorithms, and winning strategies is your challenge as a programmer!

## Get Started

New to 1v1 battles? Start with the [Getting Started](getting-started/) guide to create your first bot, then explore the mechanics pages to understand the full system.
