# [POST] /api/v1/games/{gameId}/actions/draw, [POST] /api/v1/games/{gameId}/actions/play-card, [POST] /api/v1/games/{gameId}/actions/pass
Feature: 行動階段

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
    And (UID="$Player001.id") 建立遊戲房, call table:
      | >Game001.gameId | hostPlayerId  | maxPlayers |
      | <gameId         | $Player001.id | 9          |
    And (UID="$Player002.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player002.id |
    And (UID="$Player003.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player003.id |
    And (UID="$Player001.id") 開始遊戲, call table:
      | gameId          | hostPlayerId  |
      | $Game001.gameId | $Player001.id |

  Rule: 後置 - 行動階段開始時當前玩家抽 2 張牌

    Example: 當前玩家抽 2 張牌
      Given 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId | phase  | >initialHandCount |
        | $Game001.gameId | $Player001.id   | ACTION | <handCardCount    |
      When (UID="$Player001.id") 行動階段抽牌, call table:
        | gameId          | playerId      | drawCount |
        | $Game001.gameId | $Player001.id | 2         |
      Then 回應, with table:
        | gameId          | playerId      | drawnCards |
        | $Game001.gameId | $Player001.id | 2          |
      And 應該存在一個遊戲, with table:
        | gameId          | playerId      | handCardCount                     |
        | $Game001.gameId | $Player001.id | &add($initialHandCount, 2)        |

  Rule: 後置 - 牌堆不足時只抽取剩餘的牌

    Example: 牌堆只剩 1 張牌
      Given (UID="$Player001.id") 設定牌堆狀態, call table:
        | gameId          | deckCardCount |
        | $Game001.gameId | 1             |
      And 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId | phase  | >initialHandCount |
        | $Game001.gameId | $Player001.id   | ACTION | <handCardCount    |
      When (UID="$Player001.id") 行動階段抽牌, call table:
        | gameId          | playerId      | drawCount |
        | $Game001.gameId | $Player001.id | 2         |
      Then 回應, with table:
        | gameId          | playerId      | drawnCards |
        | $Game001.gameId | $Player001.id | 1          |
      And 應該存在一個遊戲, with table:
        | gameId          | deckCardCount |
        | $Game001.gameId | 0             |

  Rule: 後置 - 玩家可以出功能牌（第一版無實際效果）

    Example: 玩家出功能牌
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 威脅     |
      When (UID="$Player001.id") 出功能牌, call table:
        | gameId          | playerId      | cardId        |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      Then 回應, with table:
        | gameId          | playerId      | message    |
        | $Game001.gameId | $Player001.id | 卡牌已使用 |
      And 應該不存在該手牌, with table:
        | playerId      | cardId          |
        | $Player001.id | $Card001.cardId |

  Rule: 前置 - 當前行動玩家不能出掉最後一張手牌

    Example: 當前行動玩家嘗試出掉最後一張手牌
      Given (UID="$Player001.id") 設定手牌數量, call table:
        | gameId          | playerId      | handCardCount |
        | $Game001.gameId | $Player001.id | 1             |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 威脅     |
      When (UID="$Player001.id") 出功能牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      Then 操作失敗

  Rule: 前置 - 卡牌必須在手牌中才能出牌

    Example: 嘗試出不在手牌中的卡牌
      When (UID="$Player001.id") 出功能牌, call table:
        | gameId          | playerId      | cardId                 |
        | $Game001.gameId | $Player001.id | non-existent-card-id   |
      Then 操作失敗

  Rule: 後置 - 所有玩家 pass 後進入情報階段

    Example: 所有玩家 pass
      Given (UID="$Player001.id") 選擇 pass, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player001.id |
      And (UID="$Player002.id") 選擇 pass, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      When (UID="$Player003.id") 選擇 pass, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | phase        |
        | $Game001.gameId | INTELLIGENCE |
