# [POST] /api/v1/games/{gameId}/start
Feature: 遊戲初始化

  Background:
    Given 準備一個玩家, with table:
      | >Player001.id | playerName |
      | <playerId     | Alice      |
    And 準備一個玩家, with table:
      | >Player002.id | playerName |
      | <playerId     | Bob        |
    And 準備一個玩家, with table:
      | >Player003.id | playerName |
      | <playerId     | Charlie    |

  Rule: 後置 - 根據人數分配身份牌

    Scenario Outline: 不同人數的身份配置
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers   |
        | <gameId         | $Player001.id | <playerCount> |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      When (UID="$Player001.id") 開始遊戲, call table:
        | gameId          | hostPlayerId  |
        | $Game001.gameId | $Player001.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | 潛伏戰線Count | 軍情處Count | 打醬油的Count |
        | $Game001.gameId | <red>         | <blue>      | <green>       |

      Examples:
        | playerCount | red | blue | green |
        | 3           | 1   | 1    | 1     |
        | 5           | 2   | 2    | 1     |

  Rule: 後置 - 每位玩家獲得 3 張初始手牌

    Example: 發放初始手牌
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      When (UID="$Player001.id") 開始遊戲, call table:
        | gameId          | hostPlayerId  |
        | $Game001.gameId | $Player001.id |
      Then 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | handCardCount |
        | $Game001.gameId | $Player001.id | 3             |
      And 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | handCardCount |
        | $Game001.gameId | $Player002.id | 3             |
      And 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | handCardCount |
        | $Game001.gameId | $Player003.id | 3             |

  Rule: 後置 - 初始化包含 45 張卡牌的牌堆

    Example: 建立完整牌堆
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      When (UID="$Player001.id") 開始遊戲, call table:
        | gameId          | hostPlayerId  |
        | $Game001.gameId | $Player001.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | deckCardCount | shuffled |
        | $Game001.gameId | 45            | true     |

  Rule: 後置 - 第一位進房的玩家成為首位行動玩家

    Example: 設定第一位行動玩家
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      When (UID="$Player001.id") 開始遊戲, call table:
        | gameId          | hostPlayerId  |
        | $Game001.gameId | $Player001.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId | phase  |
        | $Game001.gameId | $Player001.id   | ACTION |
