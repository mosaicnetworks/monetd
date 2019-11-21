# New Voting Test
3 initial validators and live node. 

- Node 4 self nominates
    - Node 1 votes Yes
    - Node 1 tries to vote again - should fail
    - Node 4 attempt to self-nominate - should fail
    - Node 2 votes No
- Node 4 self nominates 
    - Node 1 votes yes - this is the current point of failure
    - Node 2 votes yes
    - Node 3 votes yes - should be approved
    
- Node 1 nominates Node 4 for eviction
    - Node 1 votes to evict Node 4
    - Node 1 attempts to re-nominate node 4 for eviction - should fail
    - Node 2 votes No

- Node 1 nominates Node 4 for eviction
    - Node 1 votes yes - this is the current point of failure
    - Node 2 votes yes
    - Node 3 votes yes - should be approved
