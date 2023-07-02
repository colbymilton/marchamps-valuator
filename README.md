# Marvel Champions Pack Valuator

This tool is meant to help determine which Marvel Champions packs provide the best value for building up your collection!

## Setup

To be written.

## How Cards Are Valued

| Evaluated | Points | Example | Implemented |
| --- | --- | --- | --- |
| Already Owned | x1 or x0 | Cards you already own are worth 0 points. | Yes |
| Popularity in Eligable Decks | 0 -> 200 | If a card included in 25% of same aspect decks since its release, it will be worth 50 points. | Yes |
| How Many Heroes Match Trait | 0 -> 50 | If a card is trait-locked ("can only be played by XYZ Identities") and you own 50% of XYZ heroes, it will be worth 25 points. | No |