# [POST] /api/v1/games, [POST] /api/v1/games:join, [POST] /api/v1/games:start
Feature: 遊戲房管理

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

  Rule: 前置 - 建立遊戲房需要認證

    Example: 成功建立遊戲房
      When (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      Then 回應, with table:
        | gameId          | status  | hostPlayerId  |
        | $Game001.gameId | WAITING | $Player001.id |
      And 應該存在一個遊戲, with table:
        | gameId          | hostPlayerId  | status  | currentPlayers |
        | $Game001.gameId | $Player001.id | WAITING | 1              |

    Example: 未認證無法建立遊戲房
      When (No Actor) 建立遊戲房, call table:
        | hostPlayerId  | maxPlayers |
        | $Player001.id | 9          |
      Then 操作失敗

  Rule: 後置 - 玩家可以加入等待中的遊戲房

    Example: 成功加入遊戲房
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      When (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 回應, with table:
        | gameId          | playerId      | message  |
        | $Game001.gameId | $Player002.id | 加入成功 |
      And 應該存在一個遊戲, with table:
        | gameId          | currentPlayers |
        | $Game001.gameId | 2              |

  Rule: 前置 - 只有房主可以開始遊戲

    Example: 房主可以開始遊戲
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
      Then 回應, with table:
        | gameId          | status  |
        | $Game001.gameId | PLAYING |
      And 應該存在一個遊戲, with table:
        | gameId          | status  |
        | $Game001.gameId | PLAYING |

    Example: 非房主無法開始遊戲
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      And (UID="$Player003.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player003.id |
      When (UID="$Player002.id") 開始遊戲, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 操作失敗

  Rule: 前置 - 遊戲房人數限制為 3-9 人

    Example: 人數少於 3 人時無法開始
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      When (UID="$Player001.id") 開始遊戲, call table:
        | gameId          | hostPlayerId  |
        | $Game001.gameId | $Player001.id |
      Then 操作失敗

    Example: 第 10 位玩家無法加入
      Given 準備一個玩家, with table:
        | >Player004.id | playerName |
        | <playerId     | David      |
      And 準備一個玩家, with table:
        | >Player005.id | playerName |
        | <playerId     | Emma       |
      And 準備一個玩家, with table:
        | >Player006.id | playerName |
        | <playerId     | Frank      |
      And 準備一個玩家, with table:
        | >Player007.id | playerName |
        | <playerId     | Grace      |
      And 準備一個玩家, with table:
        | >Player008.id | playerName |
        | <playerId     | Henry      |
      And 準備一個玩家, with table:
        | >Player009.id | playerName |
        | <playerId     | Ivy        |
      And 準備一個玩家, with table:
        | >Player010.id | playerName |
        | <playerId     | Jack       |
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
      And (UID="$Player006.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player006.id |
      And (UID="$Player007.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player007.id |
      And (UID="$Player008.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player008.id |
      And (UID="$Player009.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player009.id |
      When (UID="$Player010.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player010.id |
      Then 操作失敗

  Rule: 後置 - 玩家不能重複加入同一遊戲房

    Example: 嘗試重複加入同一遊戲房
      Given (UID="$Player001.id") 建立遊戲房, call table:
        | >Game001.gameId | hostPlayerId  | maxPlayers |
        | <gameId         | $Player001.id | 9          |
      And (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      When (UID="$Player002.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player002.id |
      Then 操作失敗

  Rule: 前置 - 無法加入進行中的遊戲

    Example: 嘗試加入進行中的遊戲
      Given 準備一個玩家, with table:
        | >Player004.id | playerName |
        | <playerId     | Latecomer  |
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
      When (UID="$Player004.id") 加入遊戲房, call table:
        | gameId          | playerId      |
        | $Game001.gameId | $Player004.id |
      Then 操作失敗
