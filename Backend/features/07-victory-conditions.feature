# 勝利條件檢查（系統自動觸發）
Feature: 勝利條件

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

  Rule: 後置 - 潛伏戰線收集 3 張或以上紅色情報時全體獲勝

    Example: 潛伏戰線玩家達成勝利條件
      Given (UID="$Player001.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction  | redCount |
        | $Game001.gameId | $Player001.id | 潛伏戰線 | 2        |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1        |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | status | winner   |
        | $Game001.gameId | ENDED  | 潛伏戰線 |

    Example: 潛伏戰線玩家超過 3 張紅色情報
      Given (UID="$Player001.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction  | redCount |
        | $Game001.gameId | $Player001.id | 潛伏戰線 | 3        |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 2        |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | status | winner   |
        | $Game001.gameId | ENDED  | 潛伏戰線 |

  Rule: 後置 - 軍情處收集 3 張或以上藍色情報時全體獲勝

    Example: 軍情處玩家達成勝利條件
      Given (UID="$Player002.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction | blueCount |
        | $Game001.gameId | $Player002.id | 軍情處  | 2         |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | blueCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1         |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | status | winner  |
        | $Game001.gameId | ENDED  | 軍情處  |

  Rule: 後置 - 打醬油的跟隨獲勝陣營但無法單獨觸發勝利

    Example: 打醬油的玩家收集情報
      Given (UID="$Player003.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction    | redCount |
        | $Game001.gameId | $Player003.id | 打醬油的   | 2        |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1        |
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
        | gameId          | status  |
        | $Game001.gameId | PLAYING |

  Rule: 後置 - 未達成勝利條件時遊戲繼續

    Example: 紅色情報數量不足 3 張
      Given (UID="$Player001.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction  | redCount |
        | $Game001.gameId | $Player001.id | 潛伏戰線 | 1        |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | redCount |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | 1        |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player002.id") 接收情報, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 應該存在一個遊戲, with table:
        | gameId          | status  |
        | $Game001.gameId | PLAYING |
