# Marvel Champions Pack Valuator

## What is Marvel Champions?

If you've somehow come across this repo without knowing what Marvel Champions is, here's a brief description:

> "Marvel Champions is a cooperative living card game from Fantasy Flight Games, with the Marvel license. You take control of a hero from the Marvel universe and attempt to defeat a scenario represented by a major villain, such as Rhino in the Core Set. Every hero has its own signature cards that must be included when building a deck for that hero, while the remainder of the deck is made up of whichever aspect is chosen: Justice, Leadership, Aggression, or Protection. In this way, you can play your favorite hero in whatever role best complements the team." - teamcovenant.com

For more information, see the game's official page: https://www.fantasyflightgames.com/en/products/marvel-champions-the-card-game/

## So what is this tool?

The Marvel Champions Pack Valuator is a tool that is meant to help players decide which pack they should buy next when attempting to build up their collection. It's worth noting that I believe players should primarily focus on acquiring the heroes (and villains) that interest them - but if you've already bought your favorite heroes, this tool might help you decide what to get next!

Another note: this tool is focused on improving your collection for deck building. That means packs are not evaluated based on the hero that they include; they are instead evaluated based on the *other* cards in the pack. 

For example: Hulk is widely considered to be a weak hero but his pre-constructed deck includes some fantastic aggression cards. Due to these great included cards, Hulk's pack may end up being valued quite highly by this tool despite the fact Hulk himself is typically not.

## How Cards Are Valued

Every card is valued at 100 points by default. There are then multipliers added to each card based on certain criteria, detailed below.

| Evaluated | Mod | Example | Implemented? |
| --- | --- | --- | --- |
| Already Owned | ×1 or ×0 | Cards you already own are worth 0 points. | Yes |
| Popularity in Eligable Decks* | ×1 -> ×2 | If a card is included in 25% of all eligible decks, it will have a ×1.25 modifier. | Yes |
| How Many Heroes Match Trait | ×0 -> ×1 | If a card is trait-locked** and 6 out of your 10 owned heroes have that trait, it will have a ×0.6 modifier. | Yes |
| Aspect Weights | ×0 -> ×1 | If the user specifies a 0.5 weight for leadership cards, then all leadership cards will have a ×0.5 modifier. | Yes |

\* An eligible deck is defined as "a deck that could feasibly include the card":
- The deck must be running the appropriate aspect (a protection card is not eligible in an aggression deck).
- The deck must have been updated since the release of the card (a card from Wolverine's pack which released in 2022 is not eligible in a deck from 2020).
- For trait-locked cards, the deck must be for a hero that has or can reasonably acquire the specifed trait (Dive Bomb can only be played if your identity has the aerial trait and thus is not eligible in Captain America decks, but is eligible with Spectrum, Dr. Strange, Nova, etc.)

\*\* A trait-locked card is a card that can only be played if your identity has a specific trait. Dive Bomb requires aerial, The Sorcerer Supreme requires Mystic, etc.