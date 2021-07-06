# raft-practice
- [raft-practice](#raft-practice)
  - [Overview](#overview)
    - [Leader Election](#leader-election)
    - [Log Replication](#log-replication)
  - [Dive into Leader Election](#dive-into-leader-election)
  - [Dive into Log Replication](#dive-into-log-replication)
  - [References](#references)


所有節點初始狀態都是 `follower`

若 `follower` 沒有聽到來自 `leader` 心跳，則可以變成 `candidate`

## Overview
### Leader Election
candidate 可以發起選舉 (vote)，要求其他人投票給它；其他人會回應選票。
- 如果 candidate 獲得大多數的票，則會成為 leader

### Log Replication
leader 收到一筆改動訊息，會先存在 log，接著會將訊息丟給其他人，其他人先收著資料並設定成 uncommitted 的狀態並回應
當收到大多數的人的回應時，leader commit message，並且告訴其他人它 committed 了
其他人收到 committed 訊息，也將自己的資料 committed

## Dive into Leader Election

two timeout settings 控制選舉流程
- `election timeout` - 是每一個 follower 等待成為 candidate 的時間
  - 會隨機設定: 150ms~300ms
  - 在時間到達時，會轉變自己的狀態成為 candidate，並發起一場選舉要求大家投給它 (election term)
  - 發起選舉的方式是 send out `Request Vote` to other nodes
  - 如果收到 `Request Vote` 的 node 在這一輪還沒有投過票，則會將票投給發出此 request 的 candidate
    - 並且 reset 他們的 election timeout
  - candidate 此時
    - 若收到多數票(自己會投給自己)
      - 則改變狀態為 leader
      - leader 傳送 `Append Entries` messages 給其他 followers
        - 此訊息會指定一個 `heartbeat timeout`
      - followers 回應 `Append Entries` messages (同時 reset election timeout)，與 leader 用此方式保持連線
      - 此一來一往，保持著 election term 的狀態，直到有某個 follower 不再收到 `Append Entries` messages、達到了 election timeout 為止
    - 若沒有收到多數票
      - 則不做事，等待新的一場選舉發生
- how to calculate majority of votes
  - 在初始化的時候應得知總共幾台 server

## Dive into Log Replication

## References
- main
  - http://thesecretlivesofdata.com/raft/
  - https://yuanchieh.page/posts/2020/2020-11-03_raft-%E6%BC%94%E7%AE%97%E6%B3%95%E4%BB%8B%E7%B4%B9/
- supplementary
  - https://iter01.com/590533.html
  - https://www.geeksforgeeks.org/raft-consensus-algorithm/