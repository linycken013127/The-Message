# [POST] /api/v1/games/{gameId}/intelligence/accept-card, [POST] /api/v1/games/{gameId}/intelligence/reject-card
Feature: 情報接收

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

  Rule: 後置 - 玩家可以接收傳遞來的情報

    Example: 玩家接收情報
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount | blueCount | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 3        | 3         | 1          |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個玩家情報, with table:
        | playerId      | redCount | blueCount | blackCount |
        | $Player002.id | 3        | 3         | 1          |
      And 應該不存在傳遞中的情報, with table:
        | gameId          |
        | $Game001.gameId |

  Rule: 後置 - 玩家可以拒絕接收情報並傳遞給下一位玩家

    Example: 玩家拒絕接收密電情報
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | SECRET   |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個情報傳遞, with table:
        | gameId          | currentTargetPlayerId |
        | $Game001.gameId | $Player003.id         |

    Example: 玩家拒絕接收文本情報
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 退回     | TEXT     |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個情報傳遞, with table:
        | gameId          | currentTargetPlayerId |
        | $Game001.gameId | $Player003.id         |

  Rule: 後置 - 密電/文本傳回發送者時系統自動接收

    Example: 密電情報回到發送者自動接收
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType | redCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | SECRET   | 2        |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      And (UID="$Player002.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      And (UID="$Player004.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player004.id |
      When (UID="$Player005.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player005.id |
      Then 應該存在一個玩家情報, with table:
        | playerId      | redCount |
        | $Player001.id | 2        |
      And 應該不存在傳遞中的情報, with table:
        | gameId          |
        | $Game001.gameId |
      And 應該存在一個遊戲, with table:
        | gameId          | phase        |
        | $Game001.gameId | ACTION       |

  Rule: 後置 - 直達牌沒人接收時系統自動歸屬發送者

    Example: 直達情報被拒絕後自動歸屬發送者
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType | blueCount |
        | $Game001.gameId | $Player001.id | <cardId         | 截獲     | DIRECT   | 3         |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          | targetPlayerId |
        | $Game001.gameId | $Player001.id | $Card001.cardId | $Player004.id  |
      When (UID="$Player004.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player004.id |
      Then 應該存在一個玩家情報, with table:
        | playerId      | blueCount |
        | $Player001.id | 3         |
      And 應該不存在傳遞中的情報, with table:
        | gameId          |
        | $Game001.gameId |

  Rule: 前置 - 只有當前目標玩家可以接收或拒絕情報

    Example: 非目標玩家嘗試接收情報
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player003.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      Then 操作失敗

  Rule: 後置 - 情報階段結束後進入下一位玩家

    Example: 進入下一位玩家的回合
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId | phase  |
        | $Game001.gameId | $Player002.id   | ACTION |
