-- .header on
.mode col


-- First, list all the tables.
select name, type FROM sqlite_master
;

-- Quick snapshot of the main table and columns.
select *, count() from wire_pow_of_two
group by campaign_id
;

    -- Table of Campaigns.
create view if not exists Camp as
select distinct decay, proglen, campaign_id from wire_pow_of_two ;

-- Make a table of all the unique powers of two for each campaign.
create view UniquePo2 as
with
A as ( select distinct campaign_id, log(2, value) lg2val  from wire_pow_of_two ),
B as ( select *, floor(lg2val)=lg2val isPo2 from A ),
C as ( select * from B where isPo2 = true),
D as ( select C.*, P.decay, P.proglen from C join Camp P using (campaign_id))
select * from D -- unique powers of two for each campaign
;

-- Count unique powers of two per campaign For i64 values there are
-- only 63 possible unique values! If we get 54 unique values in a
-- single campaign, that's pretty good?
select campaign_id, decay, proglen, count() total_unique
from UniquePo2
group by campaign_id 
;

    -- Which specific powers of two are uniquely present in which campaign?
with
A as (
select *, count() over (partition by lg2val) as cnt
from UniquePo2
)
select * from A
where cnt = 2
;

-- Let's look at the reward over time. Does it look like the fancy
-- wiring approaches find powers of two earlier and then peter out?
-- x=time y=logrew c=logrew row=decay
select *, log(2, reward) as logrew
from wire_pow_of_two
where reward > 0
;

-- Plot cumulative log-reward over time. There is a clear advantage
-- to the wiring schemes that prefer the recent program lines!
-- x=time y=sum_reward c=_decay
with
A as (
select *, log(2, reward) as logrew, cast(decay as text) _decay
from wire_pow_of_two
where reward > 0
)
select *, sum(reward) over (partition by campaign_id order by time) "sum_reward"
from A
order by _decay
;

-- Reward Cumulative Distribution
-- Interesting how you don't see a massive difference in the plots or
-- in the total count of unique powers of two. Try toggling `row` to
-- see how the reward distributions perfectly overlap.
-- x=logrew y=bucket row=campaign_id c=campaign_id
with
S as (
    select *, log(2, reward + pow(2,-11)) as logrew
    from wire_pow_of_two
),
A as (
    select *, ntile(500) over (order by logrew) bucket
    from S
    order by value
),
B as (
    select *, row_number() over () as rn 
    from A
    group by bucket
)
    select * from B
;


-- select only rows from campaigns with more than 1k rows
create view myview as
with A as (
    select campaign_id, count() cnt 
    from wire_pow_of_two 
    group by campaign_id
),
B as (
    select *,
        row_number() over (partition by campaign_id order by time) as rid
    from wire_pow_of_two join A using (campaign_id) 
    where cnt > 1000
)
select * from B 
;
