# Magnum

Magnum is a program that checks if the list of specified light novels has any updates and notes the release dates of any new entries.

## Supported Publishers

- JNovel Club
- Yen Press
- Seven Seas Entertainment (Uses Wikipedia)

## Light Novels to Account for

Daily Life of the Immortal King - Novel Updates?
Eighth Son - Novel Updates
Moonlit Fantasy - Novel Updates?
The rising of the shield hero - Open Peace Books?  https://www.onepeacebooks.com/jt/ShieldHeroLNV.html
```json
    {
      "name": "The Rising of the Shield Hero",
      "total_volumes": 0,
      "latest_volume": "",
      "unreleased_volumes": null,
      "slug_override": "List_of_The_Rising_of_the_Shield_Hero_volumes",
      "type": "LN",
      "publisher": "OnePeaceBooks",
      "status": "O"
    }
```

## Features to Add

- Get-info should be able to determine if the latest release for a series already came out and handle that accordingly
- Add more unit tests and validation for commands and parsing logic to make sure it works as intended and is easier to refactor down the road since breaking changes should be easier to catch
