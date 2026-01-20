# [POST] /api/v1/games/{gameId}/intelligence/pass-card, [POST] /api/v1/games/{gameId}/intelligence/receive-card
Feature: 情報階段

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
    And 準備一個玩家, with table:
      | >Player004.id | playerName |
      | <playerId     | David      |
    And 準備一個玩家, with table:
      | >Player005.id | playerName |
      | <playerId     | Emma       |
    And (UID="$Player001.id") 建立遊戲房, call table:
      | >Game001.gameId | hostPlayerId  | maxPlayers |
      | <gameId         | $Player001.id | 9          |
    And (UID="$Player002.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player002.id |
    And (UID="$Player003.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player003.id |
    And (UID="$Player004.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player004.id |
    And (UID="$Player005.id") 加入遊戲房, call table:
      | gameId          | playerId      |
      | $Game001.gameId | $Player005.id |
    And (UID="$Player001.id") 開始遊戲, call table:
      | gameId          | hostPlayerId  |
      | $Game001.gameId | $Player001.id |
    And (UID="$Player001.id") 進入情報階段, call table:
      | gameId          |
      | $Game001.gameId |

  Rule: 後置 - 當前行動玩家選擇手牌傳遞

    Example: 選擇手牌作為情報
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      Then 回應, with table:
        | gameId          | status            | message          |
        | $Game001.gameId | INTELLIGENCE_SENT | 情報牌已送出     |
      And 應該存在一個情報傳遞, with table:
        | gameId          | cardId          | senderPlayerId | status     |
        | $Game001.gameId | $Card001.cardId | $Player001.id  | IN_TRANSIT |

  Rule: 後置 - 密電卡牌蓋牌向右傳遞

    Example: 傳遞密電卡牌
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | SECRET   |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      Then 回應, with table:
        | gameId          | status            | message      |
        | $Game001.gameId | INTELLIGENCE_SENT | 情報牌已送出 |
      And 應該存在一個情報傳遞, with table:
        | gameId          | cardId          | faceUp | currentTargetPlayerId |
        | $Game001.gameId | $Card001.cardId | false  | $Player002.id         |

  Rule: 後置 - 文本卡牌明牌向右傳遞

    Example: 傳遞文本卡牌
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 退回     | TEXT     |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      Then 回應, with table:
        | gameId          | status            | message      |
        | $Game001.gameId | INTELLIGENCE_SENT | 情報牌已送出 |
      And 應該存在一個情報傳遞, with table:
        | gameId          | cardId          | faceUp | currentTargetPlayerId |
        | $Game001.gameId | $Card001.cardId | true   | $Player002.id         |

  Rule: 後置 - 直達卡牌蓋牌指定玩家傳遞

    Example: 傳遞直達卡牌給指定玩家
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 截獲     | DIRECT   |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          | targetPlayerId |
        | $Game001.gameId | $Player001.id | $Card001.cardId | $Player004.id  |
      Then 回應, with table:
        | gameId          | status            | message      |
        | $Game001.gameId | INTELLIGENCE_SENT | 情報牌已送出 |
      And 應該存在一個情報傳遞, with table:
        | gameId          | cardId          | faceUp | currentTargetPlayerId | originalTargetPlayerId |
        | $Game001.gameId | $Card001.cardId | false  | $Player004.id         | $Player004.id          |

  Rule: 前置 - 直達卡牌不能指定自己

    Example: 嘗試將直達卡牌傳給自己
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 截獲     | DIRECT   |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          | targetPlayerId |
        | $Game001.gameId | $Player001.id | $Card001.cardId | $Player001.id  |
      Then 操作失敗

  Rule: 前置 - 直達卡牌只能指定存活玩家

    Example: 嘗試將直達卡牌傳給死亡玩家
      Given (UID="$Player003.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player003.id | false |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 截獲     | DIRECT   |
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          | targetPlayerId |
        | $Game001.gameId | $Player001.id | $Card001.cardId | $Player003.id  |
      Then 操作失敗

  Rule: 前置 - 只有當前行動玩家可以傳遞情報

    Example: 非當前玩家嘗試傳遞情報
      Given (UID="$Player002.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player002.id | <cardId         | 破譯     |
      When (UID="$Player002.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player002.id | $Card001.cardId |
      Then 操作失敗

  Rule: 前置 - 卡牌必須在手牌中才能傳遞

    Example: 嘗試傳遞不在手牌中的卡牌
      When (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId                 |
        | $Game001.gameId | $Player001.id | non-existent-card-id   |
      Then 操作失敗
