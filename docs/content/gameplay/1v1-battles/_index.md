---
title: "1v1 Battles"
description: >
  Complete mechanics guide for 1v1 battles - arena, bot characteristics, equipment, and actions
type: docs
weight: 10
---

## 1v1 Battle Mechanics

In 1v1 battles, two bots face off in direct combat within a bounded 2D arena. Victory goes to the bot that reduces its opponent's health to zero or has more health when the timeout expires.

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

- **Elimination**: Reduce opponent's Health to 0
- **Timeout**: Have more Health than opponent when time expires

### Bot Customization

Before battle, you configure your bot's equipment loadout (1 weapon + 1 armor). During battle, your bot performs actions that consume energy and have cooldown periods.

### Programming Challenge

This documentation explains the mechanics and rules. Developing effective bot logic, decision-making algorithms, and winning strategies is your challenge as a programmer!

## Get Started

New to 1v1 battles? Start with the [Getting Started](getting-started/) guide to create your first bot, then explore the mechanics pages to understand the full system.
