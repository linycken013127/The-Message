# [POST] /api/v1/players
Feature: 玩家管理

  Rule: 後置 - 玩家名稱不可重複

    Example: 玩家名稱沒有重複
      When (No Actor) 新增玩家, call table:
        | >Player001.id | playerName |
        | <playerId     | Alice      |
      Then 回應, with table:
        | playerId       | playerName |
        | $Player001.id  | Alice      |
      And 應該存在一個玩家, with table:
        | playerId      | playerName |
        | $Player001.id | Alice      |

  Rule: 後置 - 玩家名稱重複

    Example: 玩家無法使用重複名稱
      Given (No Actor) 新增玩家, call table:
        | >Player001.id | playerName |
        | <playerId     | Alice      |
      When (No Actor) 新增玩家, call table:
        | playerName |
        | Alice      |
      Then 操作失敗

  Rule: 前置 - 玩家名稱最大長度為 10 字元

    Scenario Outline: 玩家使用無效的名稱建立失敗
      When (No Actor) 新增玩家, call table:
        | playerName   |
        | <playerName> |
      Then 操作失敗

      Examples:
        | playerName     |
        | 12345678901    |
        | player00000    |
