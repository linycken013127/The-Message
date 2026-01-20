# 死亡條件檢查（系統自動觸發）
Feature: 死亡條件

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

  Rule: 後置 - 玩家獲得 3 張或以上黑色情報時死亡

    Example: 玩家收集 3 張黑色情報死亡
      Given (UID="$Player001.id") 設定玩家情報數, call table:
        | gameId          | playerId      | blackCount |
        | $Game001.gameId | $Player001.id | 2          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1          |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | false |

    Example: 玩家收集超過 3 張黑色情報死亡
      Given (UID="$Player001.id") 設定玩家情報數, call table:
        | gameId          | playerId      | blackCount |
        | $Game001.gameId | $Player001.id | 2          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 4          |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | false |

  Rule: 後置 - 死亡玩家在回合中被跳過

    Example: 死亡玩家被跳過
      Given (UID="$Player002.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | false |
      And 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId |
        | $Game001.gameId | $Player001.id   |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player003.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | currentPlayerId |
        | $Game001.gameId | $Player003.id   |

  Rule: 後置 - 情報自動跳過死亡玩家

    Example: 情報跳過死亡玩家
      Given (UID="$Player003.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player003.id | false |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | SECRET   |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 拒絕接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個情報傳遞, with table:
        | gameId          | currentTargetPlayerId |
        | $Game001.gameId | $Player004.id         |

  Rule: 後置 - 黑色情報少於 3 張時玩家保持存活

    Example: 黑色情報數量不足 3 張
      Given (UID="$Player001.id") 設定玩家情報數, call table:
        | gameId          | playerId      | blackCount |
        | $Game001.gameId | $Player001.id | 1          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1          |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | true  |

  Rule: 後置 - 死亡條件優先於勝利條件

    Example: 玩家同時達成勝利和死亡條件
      Given (UID="$Player001.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction  | redCount | blackCount |
        | $Game001.gameId | $Player001.id | 潛伏戰線 | 2        | 2          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1        | 1          |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲玩家, with table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | false |
      And 應該存在一個遊戲, with table:
        | gameId          | status  |
        | $Game001.gameId | PLAYING |

  Rule: 後置 - 所有玩家死亡時遊戲結束為平局

    Example: 所有玩家都死亡
      Given (UID="$Player001.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player001.id | false |
      And (UID="$Player002.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player002.id | false |
      And (UID="$Player004.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player004.id | false |
      And (UID="$Player005.id") 設定玩家狀態, call table:
        | gameId          | playerId      | alive |
        | $Game001.gameId | $Player005.id | false |
      And (UID="$Player003.id") 設定玩家情報數, call table:
        | gameId          | playerId      | blackCount |
        | $Game001.gameId | $Player003.id | 2          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | blackCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1          |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player003.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | status | winner |
        | $Game001.gameId | ENDED  | &isNull |
