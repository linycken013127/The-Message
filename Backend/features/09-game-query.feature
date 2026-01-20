# [GET] /api/v1/games/{gameId}
Feature: 遊戲資訊查詢

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

  Rule: 後置 - 玩家可以查詢自己參與的遊戲資訊

    Example: 成功查詢遊戲資訊
      When (UID="$Player001.id") 查詢遊戲資訊, call table:
        | gameId          |
        | $Game001.gameId |
      Then 回應, with table:
        | gameId          | status  | currentPlayerId |
        | $Game001.gameId | PLAYING | $Player001.id   |

  Rule: 後置 - 查詢結果只包含玩家可見的資訊

    Example: 查詢結果包含可見資訊
      Given (UID="$Player001.id") 設定玩家身份與情報, call table:
        | gameId          | playerId      | faction  | redCount | blueCount | blackCount |
        | $Game001.gameId | $Player001.id | 潛伏戰線 | 1        | 0         | 1          |
      And (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     |
      When (UID="$Player001.id") 查詢遊戲資訊, call table:
        | gameId          |
        | $Game001.gameId |
      Then 回應, with table:
        | gameId          | myFaction | myRedCount | myBlueCount | myBlackCount | myHandCardCount |
        | $Game001.gameId | 潛伏戰線  | 1          | 0           | 1            | 1               |

  Rule: 前置 - 未認證玩家無法查詢

    Example: 未提供認證
      When (No Actor) 查詢遊戲資訊, call table:
        | gameId          |
        | $Game001.gameId |
      Then 操作失敗

  Rule: 後置 - 查詢結果包含遊戲狀態

    Example: 查詢包含完整遊戲狀態
      Given (UID="$Player001.id") 準備手牌, call table:
        | gameId          | playerId      | >Card001.cardId | cardName | cardType |
        | $Game001.gameId | $Player001.id | <cardId         | 破譯     | SECRET   |
      And (UID="$Player001.id") 進入情報階段, call table:
        | gameId          |
        | $Game001.gameId |
      And (UID="$Player001.id") 傳遞情報牌, call table:
        | gameId          | playerId      | cardId          |
        | $Game001.gameId | $Player001.id | $Card001.cardId |
      When (UID="$Player001.id") 查詢遊戲資訊, call table:
        | gameId          |
        | $Game001.gameId |
      Then 回應, with table:
        | gameId          | status  | currentPlayerId | phase        | hasIntelInTransit |
        | $Game001.gameId | PLAYING | $Player001.id   | INTELLIGENCE | true              |
      And 回應包含所有玩家公開資訊, with table:
        | gameId          | playerCount |
        | $Game001.gameId | 3           |
